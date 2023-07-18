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

var version = "Uptime Kuma Client 1.0.0 (Alpha)"           // app version
var repo_url = "https://github.com/Yuiinars/uptime-client" // app repo url
var totalTcpRequests uint64                                // atomic
var totalIcmpRequests uint64                               // atomic
var totalHttpRequests uint64                               // atomic
var totalQuicRequests uint64                               // atomic
var totalDnsRequests uint64                                // atomic

type Config struct {
	ApiDomain string   `yaml:"api_domain"`
	ExecMode  string   `yaml:"exec_mode"`
	Targets   []Target `yaml:"targets"`
}

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

func main() {
	config := readConfig()
	// @TODO: Add log file support

	switch config.ExecMode {
	case "async":
		runAsync(config)
	case "sync":
		runSync(config)
	default:
		fmt.Println("Invalid exec_mode in config.yaml")
	}
}

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
	defer sess.CloseWithError(0, "")

	// Open a QUIC stream
	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to open QUIC stream: %w", err)
	}
	defer stream.Close()

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
			conn.Close()
		}
		duration = time.Since(start)
		atomic.AddUint64(&totalTcpRequests, 1) // increment

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
			pinger.Count = 1                                             // Send 1 packet

			err = pinger.Run()
			if err != nil {
				fmt.Printf("Error running ICMP pinger (%s): %v\n", target.Name, err)
			} else {
				stats := pinger.Statistics()
				duration = stats.AvgRtt
				if stats.PacketsRecv == 0 {
					status = "down"
					msg = fmt.Sprintf("[Error] Cannot connect to %s (ICMP)", target.Name)
				}
				atomic.AddUint64(&totalIcmpRequests, 1) // increment
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
			resp.Body.Close()
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
				msg = ""
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
				msg = ""
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

func readConfig() Config {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("Error reading config.yaml")
	}

	var config Config
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		panic("Error parsing config.yaml")
	}

	fmt.Println("Config loaded successfully")
	logoText :=
		"    __  __      __  _                   _________            __ \n" +
			"   / / / /___  / /_(_)___ ___  ___     / ____/ (_)__  ____  / /_\n" +
			"  / / / / __ \\/ __/ / __ `__ \\/ _ \\   / /   / / / _ \\/ __ \\/ __/\n" +
			" / /_/ / /_/ / /_/ / / / / / /  __/  / /___/ / /  __/ / / / /_  \n" +
			" \\____/ .___/\\__/_/_/ /_/ /_/\\___/   \\____/_/_/\\___/_/ /_/\\__/  \n" +
			"      /_/               \n" + version + "\n" + repo_url + "\n"
	fmt.Println(logoText)

	return config
}

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Validate HTTP Status Code
		fmt.Println(
			"-------------------------------------------------------------------\n"+
				"API Error: HTTP Status Code:", resp.StatusCode, "[", apiDomain+token, "]",
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
