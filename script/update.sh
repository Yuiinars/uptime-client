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
  # Check binary exist
    if [ -f "/etc/uptime-client/bin/main" ]; then
        echo -e "${SUCCESS}Binary exist${NC}"
    else
        echo -e "${ERROR}Binary not exist, please run install.sh first.${NC}"
        exit 1
    fi

  # Check .v cache exist
    if [ -f "/etc/uptime-client/.v" ]; then
        echo -e "${SUCCESS}Version cache file exist${NC}"
    else
        echo -e "${WARN}Version cache file not exist, will force to update.${NC}"
        echo "0" > /etc/uptime-client/.v
        chmod 600 /etc/uptime-client/.v
    fi
}

# Get Version
function getVersion() {
    # Check local version
    LOCAL_VERSION=$(cat /etc/uptime-client/.v)
    echo -e "${INFO}Local version: v$LOCAL_VERSION${NC}"

    # Check remote version
    REMOTE_VERSION=$(curl -sL https://raw.githubusercontent.com/Yuiinars/uptime-client/main/assets/v)
    echo -e "${INFO}Remote version: v$REMOTE_VERSION${NC}"
}

# Verify version
function verifyVersion() {
    if [ "$LOCAL_VERSION" == "$REMOTE_VERSION" ]; then
        echo -e "${SUCCESS}You are using the latest version${NC}"
        exit 0
    else
        echo -e "${WARN}New version found: ${NC}v$REMOTE_VERSION ${WARN}Downloading...${NC}"
    fi
}

# Update
function update() {
    # Read Architecture
    ARCH=$(cat /etc/uptime-client/.arch)
    echo -e "${INFO}Download binary...${NC}"
    curl -sL "https://bin.xmsl.dev/uptime-client/$ARCH" -o /etc/uptime-client/bin/main
    chmod +x /etc/uptime-client/bin/main
    echo -e "${SUCCESS}Download binary success${NC}"

    # Update .v cache
    echo "$REMOTE_VERSION" > /etc/uptime-client/.v

    # restart service
    echo -e "${INFO}Restarting service...${NC}"
    if systemctl restart uptime-client; then
        echo -e "${SUCCESS}Service restarted${NC}"
    else
        echo -e "${ERROR}Service restart failed, please run ${NC}systemctl restart uptime-client${ERROR} manually.${NC}"
    fi
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
    echo -e "${INFO}Updating uptime-client...${NC}"

    checkDependencies
    getVersion
    verifyVersion
    update

    echo -e "\n${SUCCESS}uptime-client is updated to v$REMOTE_VERSION!${NC}\n"
    echo -e "${WARN}Like it? Give a star to this project: ${NC}https://github.com/Yuiinars/uptime-client${NC}"
    echo -e "${WARN}Have any questions? Please open an issue: ${NC}https://github.com/Yuiinars/uptime-client/issues"
}

main