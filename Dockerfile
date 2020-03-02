# Build the manager binary
FROM golang:1.13 as builder

# Versions
ARG OS=linux
ARG ARCH=amd64
ARG KUBEBUILDER_VERSION=2.3.0

# Install kubebuilder
WORKDIR /tmp
RUN curl -Lo kubebuilder.tar.gz https://go.kubebuilder.io/dl/${KUBEBUILDER_VERSION}/${OS}/${ARCH} \
      && tar -xzvf kubebuilder.tar.gz \
      && mkdir -p /usr/local/kubebuilder/bin \
      && mv kubebuilder_${KUBEBUILDER_VERSION}_${OS}_${ARCH}/bin/* /usr/local/kubebuilder/bin/ \
      && rm -rf kubebuilder*

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY config/ config/

# Build
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
