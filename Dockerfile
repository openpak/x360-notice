# Multi-stage, cross-compiled: the build stage runs on the runner's native arch and Go
# emits the target arch, so linux/arm64 builds need no QEMU.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/x360-notice ./cmd/x360-notice

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
USER app
COPY --from=build /out/x360-notice /usr/local/bin/x360-notice
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/x360-notice"]
