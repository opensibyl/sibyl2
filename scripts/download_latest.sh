set -e

LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/opensibyl/sibyl2/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"v\([^"]*\)".*/\1/')
echo "latest version: ${LATEST_VERSION}"

DOWNLOAD_URL="https://github.com/opensibyl/sibyl2/releases/download/v${LATEST_VERSION}/sibyl2_${LATEST_VERSION}_linux_amd64"
echo "download url: ${DOWNLOAD_URL}"
wget ${DOWNLOAD_URL}
