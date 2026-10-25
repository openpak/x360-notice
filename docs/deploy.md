# Deploy: x360-notice

Architecture: **a real Xbox 360 and the Wharf emulator hit the same endpoints.** OpenPak's DNS (or Wharf's
`nic_dns_hosts`) redirects `xncsi.xboxlive.com` and `notice.xbox.com` to OpenPak, the consoles connect on port 80,
and the reverse proxy in front forwards those two Host names to this service on `X360_NOTICE_ADDR` (default `:8080`).
The service does not look at the Host header: the paths tell the requests apart. There is no emulator-only path.

## Run

```sh
cp example.env .env
set -a; source .env; set +a
go run ./cmd/x360-notice          # listens on :8080
curl -i http://127.0.0.1:8080/XBOXNCSI.txt
```

```sh
go test ./...        # needs no network beyond loopback
go vet ./... && gofmt -l .
docker build -t x360-notice .
```

## Configuration

| Variable | Required | Default | Meaning |
|---|---|---|---|
| `X360_NOTICE_ADDR` | no | `:8080` | HTTP listen address (plain HTTP; real consoles arrive on :80 through the reverse proxy) |
| `X360_NOTICE_STATUS_TEXT` | no | `Xbox LIVE is up and running` | the `MessageText` of the service-status document, XML-escaped |
