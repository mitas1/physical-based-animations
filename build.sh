COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)

building() {
    echo -e "\e[96m$1\e[0m"
}

success() {
    echo -e "\e[92m$1\e[0m"
}

failed() {
    echo -e "\e[91m$1\e[0m"
}

retval=0

if [ "${1}" == "linux" ] ; then
    if [ "${OSTYPE}" == "linux-gnu" ]; then
        building "Building for ${1}"
        GOOS="${1}" go build -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    else
        building "Cross-building for ${1}"
        GOOS="${1}" GOARCH="${5}" go build CGO_ENABLED=1 CC="${4}" -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    fi
fi

if [ "${1}" == "darwin" ] ; then
    if [ "${OSTYPE}" == "darwin" ]; then
        building "Building for ${1}"
        GOOS="$1" go build -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    else
        building "Cross-building for ${1}"
        GOOS="$1" GOARCH="${5}" CGO_ENABLED=1 CC="${4}" go build -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    fi
fi

if [ "${1}" == "windows" ]; then
    if [ "${OSTYPE}" == "win32" ] || [ "${OSTYPE}" == "msys" ] || [ "${OSTYPE}" == "cygwin" ]; then
        building "Building for ${1}"
        GOOS="${1}" go build -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    else
        building "Cross-building for ${1}"
        GOOS="${1}" GOARCH="${5}" CGO_ENABLED=1 CC="${4}" go build -ldflags "-X main.VERSION="${2}" -X main.COMMIT="${COMMIT}" -X main.BRANCH="${BRANCH}"" -o "${3}" . || retval=$?
    fi
fi

if [ $retval -eq 0 ]; then
    success "Successfully built. Output in ${3}"
else
    failed "Error building for $1"
fi
