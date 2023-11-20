# Uptime Client
A simple health check client for [Uptime Kuma](https://github.com/louislam/uptime-kuma).

## :eight_spoked_asterisk: Support API Provider
- [x] Uptime Kuma @ v1.0.1-alpha+
- [ ] Uptime Robot (Todo)

## :ballot_box_with_check: Support Protocol/Service List

| Protocol/Service       | Version           | Status                      |
|------------------------|-------------------|-----------------------------|
| ICMP                   | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| TCP/UDP                | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| HTTP(S) (GET)          | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| DNS over UDP/QUIC      | Beta              |   :ballot_box_with_check:   |
| Custom Command         | Todo              |   :x:                       |

## :arrow_down: Running

### :package: Binary

  1. Download binary to `/opt/uptime-client/` from `Github Releases`
  2. Run `./bin/main-*`

### :octopus: Compile from source code

```bash
mkdir /opt/uptime-client/ && cd /opt/uptime-client/
git clone https://github.com/Yuiinars/uptime-client .
go build -o ./bin/main .
```

### :hammer: Compile from source code (Cross-compile)

#### :penguin: Linux
- Linux x64
`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go`
- Linux arm64
`CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go`

#### :window: Windows

- Windows x64
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64.exe main.go
```

- Windows arm64
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/main-windows-arm64.exe main.go
```

#### macOS
- macOS x64
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
```

- macOS arm64
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/main-darwin-arm64 main.go
```

- :whale: Docker
  - RUN `@TODO`
  - Compose `@TODO`
  - Dockerfile `@TODO`

## :toolbox: Usage
1. Create and edit config file to binary directory
  - Edit `config.example.yaml`
  - Rename `config.example.yaml` to `config.yaml`
2. Run binary
  - `./opt/uptime-client/bin/main-*`
3. Init daemon
```bash
# Linux
cat <<EOF > /etc/systemd/system/uptime-client.service
[Unit]
Description=Uptime Client
After=network.target

[Service]
Type=simple
User=root # Root user required for ICMP
WorkingDirectory=/opt/uptime-client
ExecStart=/opt/uptime-client/bin/main
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable uptime-client
systemctl start uptime-client
```

## 📝Config
- Look at `config.example.yaml`

## 📄License
- [GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/)