<div style="display: block; margin: auto;" align="center">
  <img src="./assets/logo.svg" alt="Uptime Client" width="200px" height="200px" title="Logo 图片">
</div>

# Uptime Client

![由 Golang 强力驱动](./assets/golang.svg)
![编译状态](https://github.com/Yuiinars/uptime-client/actions/workflows/build.yml/badge.svg)

[English](README.md) | [简体中文](README_zh_CN.md)

一个使用 Golang 编写的全能型网络健康检查工具。

## :eight_spoked_asterisk: 支持对接的 API 服务商
- [x] [Uptime Kuma](https://github.com/louislam/uptime-kuma)
- [ ] Uptime Robot (Todo)

## :ballot_box_with_check: 支持的协议列表

> [!TIP]
> 如果你想要添加一个新的协议，请创建一个新的 Pull Request，谢谢。

| 协议                | 支持的版本         | 支持状态                         |
|-------------------|---------------|------------------------------|
| ICMP              | v1.0.1-alpha+ | :ballot_box_with_check: 完全支持 |
| TCP/UDP           | v1.0.1-alpha+ | :ballot_box_with_check: 完全支持 |
| HTTP(s) (GET)     | v1.0.1-alpha+ | :ballot_box_with_check: 完全支持 |
| DNS over UDP/QUIC | Beta          | :ballot_box_with_check: 完全支持 |
| Custom Command    | Todo          | :x: 计划中                      |

## ：computer_mouse：一键部署

> [!WARNING]
> Windows 目前不支持一键部署。  
> **请在运行陌生脚本前先检查代码，这是一个好习惯。**

### :package: 安装
```bash
curl -fsSL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/script/install.sh | sudo bash
```

### :hammer: 更新
```bash
curl -fsSL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/script/update.sh | sudo bash
```

## :arrow_down_small: 二进制下载

> [!WARNING]
> **如果二进制文件的哈希校验不一致，请立即停止使用并保持警惕。**

这些二进制文件是 Github Action 自动从 `main` 分支的最新 commit 构建的。

| 图标            | 平台      | 架构           |            下载链接             |            哈希校验文件            |
|---------------|---------|--------------|:---------------------------:|:----------------------------:|
| :penguin:     | Linux   | AMD64 (x64)  |   [:airplane:][linux_x64]   |   [:lock:][linux_x64_hash]   |
| :penguin:     | Linux   | ARM64        |  [:airplane:][linux_arm64]  |  [:lock:][linux_arm64_hash]  |
| :smiling_imp: | FreeBSD | AMD64 (x64)  |  [:airplane:][freebsd_x64]  |  [:lock:][freebsd_x64_hash]  |
| :smiling_imp: | FreeBSD | ARM64        | [:airplane:][freebsd_arm64] | [:lock:][freebsd_arm64_hash] |
| :window:      | Windows | AMD64 (x64)  |  [:airplane:][windows_x64]  |  [:lock:][windows_x64_hash]  |
| :window:      | Windows | ARM64        | [:airplane:][windows_arm64] | [:lock:][windows_arm64_hash] |
| :apple:       | macOS   | AMD64 (Beta) |   [:airplane:][macos_x64]   |   [:lock:][macos_x64_hash]   |
| :apple:       | macOS   | ARM64 (Beta) |  [:airplane:][macos_arm64]  |  [:lock:][macos_arm64_hash]  |

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


## :arrow_down: 安装与使用

### :package: 二进制安装

  1. 从 [这个列表](#arrow_down_small-二进制下载) 下载二进制文件
  2. 运行二进制文件

### :octopus: 从源码编译安装

```bash
git clone https://github.com/Yuiinars/uptime-client
cd uptime-client

go mod download
go mod verify
go mod tidy
go mod vendor

go build -o bin/main
./bin/main
```

### :hammer: 支持的平台列表

> [!TIP]
> **如果你想要添加一个新的平台，请创建一个新的 Pull Request，谢谢。**

| Icon          | Platform | Architecture |         Support         | Note     |
|---------------|----------|--------------|:-----------------------:|----------|
| :penguin:     | Linux    | AMD64 (x64)  | :ballot_box_with_check: | 已测试      |
| :penguin:     | Linux    | ARM64        | :ballot_box_with_check: | 已测试      |
| :smiling_imp: | FreeBSD  | AMD64 (x64)  |     :yellow_circle:     | 未经测试     |
| :smiling_imp: | FreeBSD  | ARM64        |     :yellow_circle:     | 未经测试     |
| :window:      | Windows  | AMD64 (x64)  | :ballot_box_with_check: | 已测试      |
| :window:      | Windows  | ARM64        | :ballot_box_with_check: | 已测试      |
| :apple:       | macOS    | AMD64 (x64)  |   :large_blue_circle:   | 测试版，未经测试 |
| :apple:       | macOS    | ARM64        |   :large_blue_circle:   | 测试版，未经测试 |

- :whale: 从 Docker 运行
  - RUN `@TODO`
  - Compose `@TODO`
  - Dockerfile `@TODO`

## :toolbox: 使用
1. 编辑配置文件
  - 编辑 `config.example.yaml`
  - 将 `config.example.yaml` 重命名为 `config.yaml`

1. 运行二进制文件
  - `./bin/main`

1. 初始化 systemd 服务
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