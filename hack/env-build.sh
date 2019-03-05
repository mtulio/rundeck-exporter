export APP_NAME=rundeck-exporter
export APP_VERSION=$(cat ./VERSION)
export APP_COMMIT=$(git rev-parse --short HEAD)
export APP_TAG=$(git describe --tags --always)
export APP_BUILD_USER=$(whoami)
export LDFLAGS="-X main.VersionCommit=${APP_COMMIT} -X main.VersionTag=${APP_TAG} LAGS -X main.VersionFull=${APP_VERSION}"