# Uptime Client
A simple health check client for [Uptime Kuma](https://github.com/louislam/uptime-kuma).

## ✳️ Support API Provider
- [x] Uptime Kuma @ v1.0.1-alpha+
- [ ] Uptime Robot (Todo)

## ✅ Support Protocol/Service List
- [x] ICMP @ v1.0.1-alpha+
- [x] TCP/UDP @ v1.0.1-alpha+
- [x] HTTP(S) (GET) @ v1.0.1-alpha+
- [x] DNS over UDP/QUIC @ v1.0.1-alpha+
- [ ] BIRD (Todo)

## ⏬Download
- 💲 Binary
  - Download binary to `/opt/uptime-client/` from `Github Releases`
  - Run `./bin/main-*`
- 🐙 Source
```bash
mkdir /opt/uptime-client/ && cd /opt/uptime-client/
git clone https://github.com/Yuiinars/uptime-client .
go build -o ./bin/main .
```
- ☁️ Cross platform
```bash
# Linux x64
## to 🐧Linux arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
## to 🪟Windows x64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64.exe main.go
## to 🪟Windows arm64
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/main-windows-arm64.exe main.go
## to macOS x64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
## to macOS arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/main-darwin-arm64 main.go
```

- 🐳 Docker
  - RUN `@TODO`
  - Compose `@TODO`
  - Dockerfile `@TODO`

## ☸️Usage
- Create and edit config file to binary directory
  - Edit `config.example.yaml`
  - `mv config.example.yaml config.yaml`
- Run binary
  - `./opt/uptime-client/bin/main-*`
- Init daemon
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