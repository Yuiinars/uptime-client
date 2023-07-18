# Uptime Client
A simple health check client for [Uptime Kuma](https://github.com/louislam/uptime-kuma).

## ⏬Download
- 💲 Binary
  - Download binary to `/opt/uptime-client/` from `Github Releases`
- 🐙 Source
```bash
mkdir /opt/uptime-client/ && cd /opt/uptime-client/
git clone https://github.com/Yuiinars/uptime-client .
go build -o ./bin/main .
```
- 🐳 Docker
  - `@TODO`

## ☸️Usage
- Create and edit config file to binary directory
  - `nano config.example.yaml` / `vim config.example.yaml`
  - `mv config.example.yaml config.yaml`
- Run binary
  - `./opt/uptime-client/bin/main`
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
```

## 📝Config
- Look at `config.example.yaml`

## 📄License
- [GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/)