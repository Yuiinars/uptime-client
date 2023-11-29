package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"sync/atomic"

	"github.com/go-ping/ping"
	"github.com/miekg/dns"
	"github.com/quic-go/quic-go"
	"gopkg.in/yaml.v2"
)

var version = "Uptime Client 1.1.3"                       // app version
var repoUrl = "https://github.com/Yuiinars/uptime-client" // app repo url
var totalTcpRequests uint64                               // atomic

// totalIcmpRequests is a variable that represents the total number of ICMP requests made.
var totalIcmpRequests uint64 // atomic
var totalHttpRequests uint64 // atomic
var totalQuicRequests uint64 // atomic
var totalDnsRequests uint64  // atomic

// Config represents the configuration for the uptime client.

type Target struct {
	Token    string `yaml:"token"`
	Mode     string `yaml:"mode"`
	Name     string `yaml:"name,omitempty"`
	Timeout  int    `yaml:"timeout,omitempty"`  // seconds
	Interval int    `yaml:"interval,omitempty"` // seconds

	TcpTarget  string `yaml:"tcp_target,omitempty"`  // host:port
	IcmpTarget string `yaml:"icmp_target,omitempty"` // host/ip
	HttpTarget string `yaml:"http_target,omitempty"` // host:port

	DnsTarget     string `yaml:"dns_target,omitempty"`      // host / FQDN
	DnsServer     string `yaml:"dns_server,omitempty"`      // ip:port server
	DnsServerPort int    `yaml:"dns_server_port,omitempty"` // port only
	DnsType       uint16 `yaml:"dns_type,omitempty"`        // A:1, AAAA:28, CNAME:5, MX:15, NS:2, PTR:12, SOA:6, TXT:16
}

type Config struct {
	ApiMode   string   `yaml:"api_mode"`   // async, sync
	ApiDomain string   `yaml:"api_domain"` // domain only
	ApiScheme string   `yaml:"api_scheme"` // http, https
	ApiPort   int      `yaml:"api_port"`   // port only
	ApiPath   string   `yaml:"api_path"`   // example: /api/push
	Targets   []Target `yaml:"targets"`
}

func main() {
	config := readConfig()
	// @TODO: Add log file support

	switch config.ApiMode {
	case "async":
		runAsync(config)
	case "sync":
		runSync(config)
	default:
		fmt.Println("Invalid exec_mode in config.yaml")
	}
}

// queryDNSOverUDP sends a DNS query over UDP to the specified server and port.
func queryDNSOverUDP(domain string, server string, port int, queryType uint16) (*dns.Msg, error) {
	// Create a DNS message
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), queryType)
	msg.RecursionDesired = true

	// Create a DNS client
	client := new(dns.Client)

	// Send the query to DNS over UDP server
	serverAddress := fmt.Sprintf("%s:%d", server, port)
	resp, _, err := client.Exchange(msg, serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query DNS over UDP: %w", err)
	}

	return resp, nil
}

// queryDNSOverQUIC is a function that sends a DNS query over QUIC to a specified server and port.
// It takes a domain name, server address, port number, and query type as input parameters.
// It returns the DNS response message and an error if any.
func queryDNSOverQUIC(domain string, server string, port int, queryType uint16) (*dns.Msg, error) {
	// Create a DNS message
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), queryType)
	msg.RecursionDesired = true

	// Pack the DNS message into bytes
	msgBytes, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %w", err)
	}

	// Create a QUIC session
	quicConfig := &quic.Config{
		MaxIdleTimeout: time.Second * 10, // Timeout to 10 seconds
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"doq"},
	}

	ctx := context.Background()
	serverAddress := fmt.Sprintf("%s:%d", server, port)
	sess, err := quic.DialAddr(ctx, serverAddress, tlsConfig, quicConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to dial QUIC server: %w", err)
	}
	defer func(sess quic.Connection, code quic.ApplicationErrorCode, s string) {
		_ = sess.CloseWithError(code, s)
	}(sess, 0, "")

	// Open a QUIC stream
	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to open QUIC stream: %w", err)
	}
	defer func(stream quic.Stream) {
		err := stream.Close()
		if err != nil {
			fmt.Println("Error closing QUIC stream:", err, "[", serverAddress, "]")
			return
		}
	}(stream)

	// send the DNS message to QUIC server
	_, err = stream.Write(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write DNS message to QUIC stream: %w", err)
	}

	// read the response from QUIC server
	buf := make([]byte, 65536)
	n, err := io.ReadFull(stream, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read DNS response from QUIC stream: %w", err)
	}

	// unpack the response
	resp := new(dns.Msg)
	err = resp.Unpack(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack DNS response: %w", err)
	}

	return resp, nil
}

// runAsync is a function that runs the targets asynchronously.
// It takes a Config parameter which contains the configuration settings.
// For each target in the configuration, it creates a goroutine that periodically processes the target.
// The processing is done by calling the processTarget function with the API domain and the target.
// The interval between each processing is determined by the target's Interval field.
// This function waits for all the goroutines to finish before returning.

func runAsync(config Config) {
	var wg sync.WaitGroup
	for _, target := range config.Targets {
		wg.Add(1)
		go func(target Target) {
			defer wg.Done()

			ticker := time.NewTicker(time.Duration(target.Interval) * time.Second)
			for range ticker.C {
				processTarget(config.ApiDomain, target)
			}
		}(target)
	}
	wg.Wait()
}

func runSync(config Config) {
	for _, target := range config.Targets {
		processTarget(config.ApiDomain, target)
		time.Sleep(time.Duration(target.Interval) * time.Second)
	}
}

// processTarget is a function that processes a target based on its mode.
// It performs different types of pings (TCP, ICMP, HTTP, DNS-over-QUIC, DNS) on the target
// and sends the results to an API endpoint.
// The function takes the API domain and the target as parameters.
// It updates the status and message based on the ping results,
// and prints the latency and status of the target.
// It also increments the total request counters for each type of ping.
func processTarget(apiDomain string, target Target) {
	var duration time.Duration
	var status = "up" // up, down
	var msg = "OK"    // OK, ...Custom message

	switch target.Mode {
	case "tcp":
		// TCP Ping Mode
		start := time.Now()
		dialer := net.Dialer{Timeout: time.Duration(target.Timeout) * time.Second}
		conn, err := dialer.Dial("tcp", target.TcpTarget)
		if err != nil {
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (TCP)", target.Name)
			fmt.Printf("Error connecting to TCP target (%s): %v\n", target.Name, err)
		} else {
			err := conn.Close()
			if err != nil {
				fmt.Printf("Error closing TCP connection (%s): %v\n", target.Name, err)
				return
			}
		}
		duration1 := time.Since(start)

		// Second TCP Ping
		start = time.Now()
		conn, err = dialer.Dial("tcp", target.TcpTarget)
		if err != nil {
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (TCP)", target.Name)
			fmt.Printf("Error connecting to TCP target (%s): %v\n", target.Name, err)
		} else {
			err := conn.Close()
			if err != nil {
				return
			}
		}
		duration2 := time.Since(start)

		// Choose the minimum duration
		if duration1 < duration2 {
			duration = duration1
		} else {
			duration = duration2
		}
		atomic.AddUint64(&totalTcpRequests, 2) // increment by 2

	case "icmp":
		// ICMP Ping Mode
		pinger, err := ping.NewPinger(target.IcmpTarget)
		if err != nil {
			fmt.Printf("Error creating ICMP pinger (%s): %v\n", target.Name, err)
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (ICMP)", target.Name)
		} else {
			// Pinger settings
			pinger.SetPrivileged(true)                                   // Required for ICMP
			pinger.Timeout = time.Duration(target.Timeout) * time.Second // Timeout in seconds
			pinger.Count = 3                                             // Send 3 packets

			err = pinger.Run()
			if err != nil {
				fmt.Printf("Error running ICMP pinger (%s): %v\n", target.Name, err)
			} else {
				stats := pinger.Statistics()
				duration = stats.MinRtt
				if stats.PacketsRecv == 0 {
					status = "down"
					msg = fmt.Sprintf("[Error] Cannot connect to %s (ICMP)", target.Name)
				}
				atomic.AddUint64(&totalIcmpRequests, 3) // increment by 3
			}
		}
	case "http", "https":
		// HTTP Ping Mode
		start := time.Now()
		client := &http.Client{
			Timeout: time.Duration(target.Timeout) * time.Second,
		}
		url := fmt.Sprintf("%s://%s", target.Mode, target.HttpTarget)
		resp, err := client.Get(url)
		if err != nil {
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (%s)", target.Name, target.Mode)
			fmt.Printf("Error connecting to %s target (%s): %v\n", target.Mode, target.Name, err)
		} else {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("Error closing HTTP response body (%s): %v\n", target.Name, err)
				return
			}
		}
		duration = time.Since(start)
		atomic.AddUint64(&totalHttpRequests, 1) // increment
	case "doq":
		// DNS-over-QUIC Ping Mode
		start := time.Now()
		resp, err := queryDNSOverQUIC(target.DnsTarget, target.DnsServer, target.DnsServerPort, target.DnsType)

		if err != nil {
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (DNS-over-QUIC)", target.Name)
			fmt.Printf("Error connecting to DNS-over-QUIC target (%s): %v\n", target.Name, err)
		} else {
			if len(resp.Answer) == 0 {
				status = "down"
				msg = fmt.Sprintf("[Error] Cannot connect to %s (DNS-over-QUIC)", target.Name)
			} else {
				msg = "OK"
			}
		}
		duration = time.Since(start)
		atomic.AddUint64(&totalQuicRequests, 1) // increment
	case "dns":
		// DNS Ping Mode
		start := time.Now()
		resp, err := queryDNSOverUDP(target.DnsTarget, target.DnsServer, target.DnsServerPort, target.DnsType)

		if err != nil {
			status = "down"
			msg = fmt.Sprintf("[Error] Cannot connect to %s (DNS)", target.Name)
			fmt.Printf("Error connecting to DNS target (%s): %v\n", target.Name, err)
		} else {
			if len(resp.Answer) == 0 {
				status = "down"
				msg = fmt.Sprintf("[Error] Cannot connect to %s (DNS)", target.Name)
			} else {
				msg = "OK"
			}
		}
		duration = time.Since(start)
		atomic.AddUint64(&totalDnsRequests, 1) // increment
	default:
		fmt.Println("Invalid mode in config.yaml")
		return
	}

	fmt.Printf("[%s] to [%s], latency: %d ms, status: %s\n",
		strings.ToUpper(target.Mode),
		target.Name,
		duration.Milliseconds(),
		strings.ToUpper(status),
	)

	sendData(apiDomain, target.Token, duration, status, msg)

	// Print total requests every 10 requests
	if totalIcmpRequests%10 == 0 && totalIcmpRequests > 0 {
		fmt.Printf("Total ICMP requests: %d\n", totalIcmpRequests)
	}
	if totalTcpRequests%10 == 0 && totalTcpRequests > 0 {
		fmt.Printf("Total TCP requests: %d\n", totalTcpRequests)
	}
	if totalHttpRequests%10 == 0 && totalHttpRequests > 0 {
		fmt.Printf("Total HTTP requests: %d\n", totalHttpRequests)
	}
}

// readConfig reads the configuration from the "config.yaml" file and returns a Config struct.
// If there is an error reading or parsing the file, it panics.
func readConfig() Config {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("Error reading config.yaml, Is it exists?\n" + err.Error())
	}

	var config Config
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		panic("Error parsing config.yaml, Is it valid?\n" + err.Error())
	}

	config.ApiDomain = fmt.Sprintf(
		"%s://%s:%d/%s",
		config.ApiScheme,
		config.ApiDomain,
		config.ApiPort,
		config.ApiPath,
	)

	fmt.Println("Config loaded successfully")
	logoText :=
		"    __  __      __  _                   _________            __ \n" +
			"   / / / /___  / /_(_)___ ___  ___     / ____/ (_)__  ____  / /_\n" +
			"  / / / / __ \\/ __/ / __ `__ \\/ _ \\   / /   / / / _ \\/ __ \\/ __/\n" +
			" / /_/ / /_/ / /_/ / / / / / /  __/  / /___/ / /  __/ / / / /_  \n" +
			" \\____/ .___/\\__/_/_/ /_/ /_/\\___/   \\____/_/_/\\___/_/ /_/\\__/  \n" +
			"      /_/               \n" + version + "\n" + repoUrl
	fmt.Println(logoText)

	return config
}

// sendData sends data to the specified API endpoint.
// It takes the following parameters:
// - apiDomain: the domain of the API
// - token: the token used for authentication
// - duration: the duration of the request
// - status: the status of the request
// - msg: the message associated with the request
//
// It performs a GET request to the API endpoint with the provided parameters.
// If the request is successful, it checks the response for validity.
// If the response is not valid, it prints an error message.
// If there is an error during the request or response handling, it prints the corresponding error message.
func sendData(apiDomain, token string, duration time.Duration, status, msg string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", apiDomain, token), nil)
	if err != nil {
		fmt.Println("Cannot creating HTTP request:", err)
		return
	}

	q := req.URL.Query()
	q.Add("ping", strconv.FormatInt(duration.Milliseconds(), 10))
	q.Add("status", status)
	q.Add("msg", msg)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Cannot Request to API, Detail:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing HTTP response body:", err, "[", apiDomain, "/", token, "]")
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK { // Validate HTTP Status Code
		fmt.Println(
			"-------------------------------------------------------------------\n"+
				"API Error: HTTP Status Code:", resp.StatusCode, "[", apiDomain, "/", token, "]",
			"\n-------------------------------------------------------------------")
		return
	}

	var response struct {
		Ok bool `json:"ok"` // Validate API response
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Invalid API response:", err)
		return
	}

	if !response.Ok {
		fmt.Printf("Error: API response invalid: %s/%s\n", apiDomain, token)
	}
}
