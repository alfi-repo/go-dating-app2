#!/bin/sh

VERSION=${1}

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script.'
    exit 1
fi

if [ "$(uname -s)" != "Linux" ] ; then
    echo "(!) OS $(uname -s) unsupported."
    exit 1
fi

architecture="$(uname -m)"
if [ "${architecture}" != "amd64" ] && [ "${architecture}" != "x86_64" ] && [ "${architecture}" != "arm64" ] && [ "${architecture}" != "aarch64" ]; then
    echo "(!) Architecture $architecture unsupported"
    exit 1
fi

if ! mockery --version &> /dev/null ; then
    echo "Installing mockery v${VERSION}..."

    if [ "${architecture}" == "arm64" ] || [ "${architecture}" == "aarch64" ]; then
        arch="arm64"
    else
        arch="x86_64"
    fi

    tar_file="mockery_${VERSION}_Linux_${arch}.tar.gz"
    checksum_file="checksum.txt"

    curl -fsSLO --compressed "https://github.com/vektra/mockery/releases/download/v${VERSION}/${tar_file}"
    curl -fsSLO "https://github.com/vektra/mockery/releases/download/v${VERSION}/${checksum_file}"

    actual_checksum=$(sha256sum "${tar_file}" | awk '{ print $1 }')
    stored_checksum=$(grep "${tar_file}" "${checksum_file}" | awk '{ print $1 }')

    if [ "${actual_checksum}" != "${stored_checksum}" ]; then
        echo "(!) The tarball is NOT valid."
        exit 1
    fi

    bin_dir="/usr/local/bin"
    exe="${bin_dir}/mockery"
    if [ ! -d "${bin_dir}" ]; then
	    mkdir -p "${bin_dir}"
    fi

    tar -xzf "${tar_file}"
    mv mockery "${exe}"
    chmod +x "${exe}"
    rm -rf "${tar_file}" "${checksum_file}"
else
    echo "mockery already installed"
fi

echo "mockery installer done!"