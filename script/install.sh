#!/bin/bash -e

# Color
ERROR='\033[0;31m'
INFO='\033[0;36m'
WARN='\033[0;33m'
SUCCESS='\033[0;32m'
NC='\033[0m'

# Check if user is root
if [ "$EUID" -ne 0 ]; then
    echo -e "${ERROR}Please run as root${NC}"
    exit 1
fi

# Check directory
if [ -d "/etc/uptime-client" ]; then
    echo -e "${WARN}uptime-client is already installed, please run update.sh${NC}"
    exit 1
else
    echo -e "${INFO}Creating work directory...${NC}"
    mkdir -p /etc/uptime-client/bin
fi

# Check package manager
function checkPackageManager() {
    if [ -x "$(command -v apt-get)" ]; then # Debian/Ubuntu
        PACKAGE_MANAGER="apt-get"
    elif [ -x "$(command -v yum)" ]; then # CentOS 7/RHEL
        PACKAGE_MANAGER="yum"
    elif [ -x "$(command -v dnf)" ]; then # Fedora/CentOS 8+/RHEL
        PACKAGE_MANAGER="dnf"
    elif [ -x "$(command -v zypper)" ]; then # OpenSUSE
        PACKAGE_MANAGER="zypper"
    elif [ -x "$(command -v apk)" ]; then # Alpine Linux
        PACKAGE_MANAGER="apk"
    elif [ -x "$(command -v pkg)" ]; then # FreeBSD
        PACKAGE_MANAGER="pkg"
    elif [ -x "$(command -v pacman)" ]; then # Arch Linux
        PACKAGE_MANAGER="pacman"
    elif [ -x "$(command -v brew)" ]; then # macOS
        PACKAGE_MANAGER="brew"
    else
        echo -e "${ERROR}Unsupported package manager${NC}"
        exit 1;
    fi
}

function checkDependencies() {
  # Check if curl is installed
    if [ -x "$(command -v curl)" ]; then
        echo -e "${SUCCESS}Curl is installed${NC}"
    else
        echo -e "${INFO}Curl is not installed, installing...${NC}"
        $PACKAGE_MANAGER install curl -y
    fi
  # Check connect to raw.githubusercontent.com
    echo -e "${INFO}Checking connection to raw.githubusercontent.com...${NC}"
    if curl -sL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/assets/v -w "%{http_code}\n" | grep -q "200"; then
        echo -e "${SUCCESS}Connection to raw.githubusercontent.com is OK${NC}"
    else
        echo -e "${ERROR}Cannot connect to raw.githubusercontent.com, please check your network.${NC}"
        exit 1
    fi
}

# Get latest version
function getLatestVersion() {
    echo -e "${INFO}Getting latest version...${NC}"
    LATEST_VERSION=$(curl -sL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/assets/v)
    echo -e "${SUCCESS}Latest version: v$LATEST_VERSION${NC}"
}

SUPPORTED_ARCHITECTURES=(
    "freebsd-amd64"
    "freebsd-arm64"
    "linux-amd64"
    "linux-arm64"
    "windows-amd64"
    "windows-arm64"
    "macos-amd64"
    "macos-arm64"
)

# Check architecture
function checkArchitecture() {
    echo -e "${INFO}Checking architecture...${NC}"
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    if [ "$OS" == "darwin" ]; then
        OS="macos"
    elif [ "$OS" == "freebsd" ]; then
        OS="freebsd"
    elif [ "$OS" == "linux" ]; then
        OS="linux"
    elif [ "$OS" == "cygwin" ] || [ "$OS" == "msys" ] || [ "$OS" == "win32" ]; then
        OS="windows"
    else
        echo -e "${ERROR}Unsupported OS: $OS and you can open an issue to request support.${NC}"
        exit 1
    fi
    if [ "$ARCH" == "x86_64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
        ARCH="arm64"
    else
        echo -e "${ERROR}Unsupported architecture: $ARCH and you can open an issue to request support.${NC}"
        exit 1
    fi
    ARCHITECTURE="$OS-$ARCH"
    echo "Architecture: $ARCHITECTURE"
    for SUPPORTED_ARCHITECTURE in "${SUPPORTED_ARCHITECTURES[@]}"; do
        if [ "$ARCHITECTURE" == "$SUPPORTED_ARCHITECTURE" ]; then
            echo -e "${SUCCESS}Architecture is supported: $ARCHITECTURE${NC}"
            return
        fi
    done
    echo -e "${ERROR}Architecture is not supported: $ARCHITECTURE and you can open an issue to request support.${NC}"
    exit 1
}

# Download binary
function downloadBinary() {
    echo -e "${INFO}Downloading binary...${NC}"
    curl -sSL "https://bin.xmsl.dev/uptime-client/main-$ARCHITECTURE" -o "/etc/uptime-client/bin/main"
    curl -sSL "https://raw.githubusercontent.com/Yuiinars/uptime-client/main/config.example.yaml" -o "/etc/uptime-client/config.example.yaml"

    # Change permission to executable
    chmod +x "/etc/uptime-client/bin/main"

    # Save version cache
    curl -sSL "https://raw.githubusercontent.com/Yuiinars/uptime-client/main/assets/v" -o "/etc/uptime-client/.v"
    chmod 600 "/etc/uptime-client/.v"

    # Save architecture cache
    echo $ARCHITECTURE > "/etc/uptime-client/.arch"
    chmod 400 "/etc/uptime-client/.arch"
}

# Register service
function registerService() {
    # If not use systemd, exit
    if [ ! -d "/etc/systemd/system" ]; then
        echo "Systemd is not supported, please register service manually"
        return
    fi
    # If service is already register, exit
    if [ -f "/etc/systemd/system/uptime-client.service" ]; then
        echo "Service is already registered, please run systemctl enable --now uptime-client to start service"
        return
    fi
    echo "Registering service..."
    curl -sSL "https://raw.githubusercontent.com/Yuiinars/uptime-client/main/assets/uptime-client.service" -o "/etc/systemd/system/uptime-client.service"
    systemctl daemon-reload
}

# Main function
function main {

    echo -e "\n\n"
    echo -e "     __  __      __  _                   _________            __ \n" \
    "   / / / /___  / /_(_)___ ___  ___     / ____/ (_)__  ____  / /_\n" \
    "  / / / / __ \\/ __/ / __ \`__ \\/ _ \\   / /   / / / _ \\/ __ \\/ __/\n" \
    " / /_/ / /_/ / /_/ / / / / / /  __/  / /___/ / /  __/ / / / /_  \n" \
    " \\____/ .___/\\__/_/_/ /_/ /_/\\___/   \\____/_/_/\\___/_/ /_/\\__/  \n" \
    "      /_/               \n\n"
    echo -e "${INFO}Installing uptime-client...${NC}"

    checkPackageManager
    checkDependencies
    getLatestVersion
    checkArchitecture
    downloadBinary

    echo -e "\n${SUCCESS}uptime-client is installed successfully!${NC}\n"
    echo -e "${INFO}Please edit${NC} /etc/uptime-client/config.example.yaml ${INFO}and rename it to config.yaml${NC}"
    echo -e "${INFO}Then run:${NC} systemctl enable --now uptime-client\n"
    echo -e "${WARN}Like it? Give a star to this project: ${NC}https://github.com/Yuiinars/uptime-client${NC}"
    echo -e "${WARN}Have any questions? Please open an issue: ${NC}https://github.com/Yuiinars/uptime-client/issues"
}

main