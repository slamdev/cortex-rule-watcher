OS="$(go env GOOS)/$(go env GOARCH)"

if test -f "${OS}/bin/etcd"; then
    exit 0
fi

K8S_VERSION=1.21.2

curl -sSLo envtest-bins.tar.gz "https://go.kubebuilder.io/test-tools/${K8S_VERSION}/${OS}"
mkdir -p "${OS}"
tar -C "${OS}" --strip-components=1 -zvxf envtest-bins.tar.gz
rm -f envtest-bins.tar.gz
