# syntax=docker/dockerfile:1

# Build on the host architecture and cross-compile, so multi-arch builds
# don't need QEMU emulation.
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /ilpost-podcast-feed .

FROM scratch
# CA certificates are needed to call the ilpost.it https API
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=build /ilpost-podcast-feed /app/
COPY static ./static
# fixture read at runtime by the /test endpoint
COPY pkg/endpoint/bordone.json ./pkg/endpoint/bordone.json
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/ilpost-podcast-feed"]
