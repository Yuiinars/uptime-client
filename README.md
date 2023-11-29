<div style="display: block; margin: auto;" align="center">
  <img src="./assets/logo.svg" alt="Uptime Client" width="200px" height="200px" title="Logo Image">
</div>

# Uptime Client

![Powered by Golang](./assets/golang.svg)
![Build and Release](https://github.com/Yuiinars/uptime-client/actions/workflows/build.yml/badge.svg)

[English](README.md) | [简体中文](README_zh_CN.md)

A Network Health Checking tools built with Golang.

## :eight_spoked_asterisk: Support API Provider
- [x] [Uptime Kuma](https://github.com/louislam/uptime-kuma)
- [ ] Uptime Robot (Todo)

## :ballot_box_with_check: Support Protocol List

> [!TIP]
> If you want to add a new protocol, create a new PR, please.

| Protocol               | Version           | Support Status              |
|------------------------|-------------------|-----------------------------|
| ICMP                   | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| TCP/UDP                | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| HTTP(S) (GET)          | v1.0.1-alpha+     |   :ballot_box_with_check:   |
| DNS over UDP/QUIC      | Beta              |   :ballot_box_with_check:   |
| Custom Command         | Todo              |   :x:                       |

## :computer_mouse: One-click Deploy

> [!WARNING]
> **One-Click Deploy is not currently supported for Windows.**  
> **Please check the contents of the script before running it.**

### Install
```bash
curl -fsSL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/install.sh | bash
```

### Update
```bash
curl -fsSL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/update.sh | bash
```

## :arrow_down_small: Download

> [!WARNING]
> **DO NOT USE THE BINARY IF THE HASH CHECKSUM IS DIFFERENT.**

These binaries are built from the latest commit of the `main` branch by GitHub Action.

| Icon          | Platform | Architecture |        Download Link        |        Checksum File         |
|---------------|----------|--------------|:---------------------------:|:----------------------------:|
| :penguin:     | Linux    | AMD64 (x64)  |   [:airplane:][linux_x64]   |   [:lock:][linux_x64_hash]   |
| :penguin:     | Linux    | ARM64        |  [:airplane:][linux_arm64]  |  [:lock:][linux_arm64_hash]  |
| :smiling_imp: | FreeBSD  | AMD64 (x64)  |  [:airplane:][freebsd_x64]  |  [:lock:][freebsd_x64_hash]  |
| :smiling_imp: | FreeBSD  | ARM64        | [:airplane:][freebsd_arm64] | [:lock:][freebsd_arm64_hash] |
| :window:      | Windows  | AMD64 (x64)  |  [:airplane:][windows_x64]  |  [:lock:][windows_x64_hash]  |
| :window:      | Windows  | ARM64        | [:airplane:][windows_arm64] | [:lock:][windows_arm64_hash] |
| :apple:       | macOS    | AMD64 (Beta) |   [:airplane:][macos_x64]   |   [:lock:][macos_x64_hash]   |
| :apple:       | macOS    | ARM64 (Beta) |  [:airplane:][macos_arm64]  |  [:lock:][macos_arm64_hash]  |

[linux_x64]: https://bin.xmsl.dev/uptime-client/main-linux-amd64
[linux_x64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-linux-amd64.txt
[linux_arm64]: https://bin.xmsl.dev/uptime-client/main-linux-arm64
[linux_arm64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-linux-arm64.txt

[freebsd_x64]: https://bin.xmsl.dev/uptime-client/main-freebsd-amd64
[freebsd_x64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-freebsd-amd64.txt
[freebsd_arm64]: https://bin.xmsl.dev/uptime-client/main-freebsd-arm64
[freebsd_arm64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-freebsd-arm64.txt

[windows_x64]: https://bin.xmsl.dev/uptime-client/main-windows-amd64.exe
[windows_x64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-windows-amd64.txt
[windows_arm64]: https://bin.xmsl.dev/uptime-client/main-windows-arm64.exe
[windows_arm64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-windows-arm64.txt

[macos_x64]: https://bin.xmsl.dev/uptime-client/main-darwin-amd64
[macos_x64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-darwin-amd64.txt
[macos_arm64]: https://bin.xmsl.dev/uptime-client/main-darwin-arm64
[macos_arm64_hash]: https://bin.xmsl.dev/uptime-client/hash/main-darwin-arm64.txt


## :arrow_down: Running

### :package: Binary

  1. Download binary from [Download](#arrow_down_small-download)
  2. Run the binary

### :octopus: Compile from source code

```bash
mkdir /etc/uptime-client/ && cd /etc/uptime-client/
git clone https://github.com/Yuiinars/uptime-client .

go mod download
go mod verify
go mod tidy
go mod vendor

go build -o ...
```

### :hammer: Support Platform List

> [!TIP]
> **If you want to add a new platform, please create a new Pull Request.**

| Icon          | Platform | Architecture |         Support         | Note             |
|---------------|----------|--------------|:-----------------------:|------------------|
| :penguin:     | Linux    | AMD64 (x64)  | :ballot_box_with_check: | Tested           |
| :penguin:     | Linux    | ARM64        | :ballot_box_with_check: | Tested           |
| :smiling_imp: | FreeBSD  | AMD64 (x64)  |     :yellow_circle:     | Not Tested       |
| :smiling_imp: | FreeBSD  | ARM64        |     :yellow_circle:     | Not Tested       |
| :window:      | Windows  | AMD64 (x64)  | :ballot_box_with_check: | Tested           |
| :window:      | Windows  | ARM64        | :ballot_box_with_check: | Tested           |
| :apple:       | macOS    | AMD64 (x64)  |   :large_blue_circle:   | Beta, Not Tested |
| :apple:       | macOS    | ARM64        |   :large_blue_circle:   | Beta, Not Tested |

- :whale: Docker
  - RUN `@TODO`
  - Compose `@TODO`
  - Dockerfile `@TODO`

## :toolbox: Usage
1. Create and edit config file to binary directory
  - Edit `config.example.yaml`
  - Rename `config.example.yaml` to `config.yaml`

1. Run binary
  - `./etc/uptime-client/bin/main-*`

1. Init daemon
```bash
cat <<EOF > /etc/systemd/system/uptime-client.service
[Unit]
Description=Uptime Client
After=network.target

[Service]
Type=simple
User=root # Root user required for ICMP
WorkingDirectory=/etc/uptime-client
ExecStart=/etc/uptime-client/bin/main
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now uptime-client.service
```

## 📄License
- [GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/)