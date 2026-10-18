# x360-notice

OpenPak's Xbox 360 connectivity-test web pages. During **Network Settings -> Test Xbox Live Connection** the
360 dashboard (firmware 17559) makes two plain-HTTP requests, and this service answers both:

| Hostname the console resolves | Request | Answer |
|---|---|---|
| `xncsi.xboxlive.com` | `GET /XBOXNCSI.txt` | `200`, `Content-Type: text/plain`, `Cache-Control: no-cache, no-store`, `Content-Length: 9`, body exactly `Xbox NCSI` |
| `notice.xbox.com` | `GET /xocdata/xml/<locale>/servicestatusChannel.xml` (`en-US` observed; any locale code is served the same document) | `200`, `Content-Type: text/xml`, the `RootNode`/`Table` service-status XML below |

Everything else is `404` (a wrong method on those paths is `405`). `GET /healthz` returns `ok` for the container
health check; no firmware path uses it.

```xml
<?xml version="1.0"?>
<RootNode>
	<Table>
		<MessageText>Xbox LIVE is up and running</MessageText>
		<Severity> </Severity>
		<Category> </Category>
	</Table>
</RootNode>
```

(CRLF line ends, tabs, no trailing newline; `MessageText` is `X360_NOTICE_STATUS_TEXT`.) Microsoft's real
`notice.xbox.com` answers the XML request with a 301 to `xnotify-production-origin.xboxlive.com` over HTTPS and serves
this document there; here the document is served directly, which is what the firmware was shown to accept.

Naming rule: every OpenPak Xbox 360 service is called `x360-<thing>`. Sibling: [`x360-xexds`](x360-xexds).

Architecture: **a real Xbox 360 and the Wharf emulator hit the same endpoints.** OpenPak's DNS (or Wharf's
`nic_dns_hosts`) redirects `xncsi.xboxlive.com` and `notice.xbox.com` to OpenPak, the consoles connect on port 80,
and the reverse proxy in front forwards those two Host names to this service on `X360_NOTICE_ADDR` (default `:8080`).
The service does not look at the Host header: the paths tell the requests apart. There is no emulator-only path.

Status, 2026-10-16: the firmware's Test Xbox Live Connection shows **Network: Connected** and, with the XEXDS probe
answered too, **Internet: Connected** against the stand-in equivalents of these two answers (`scratch/wharf/livecap`
`-ncsi`, which is what these responses were copied from). This Go server is verified by unit tests (including the
firmware's own HTTP/1.0 request shape over TCP) and by a raw HTTP/1.0 request against a running instance. It has **not**
yet been pointed at by the emulator itself, not tried by a real console, and not deployed.

The request shape (from [`docs/xbox360/wharf-base-plan.md`](../docs/xbox360/wharf-base-plan.md), "Connectivity test
traffic"): `GET <path> HTTP/1.0`, a `Host` header, user agent `Xbox Downloader/2.0.17559.0`, nothing else. The NCSI answer
Microsoft gave was `200 OK` with the 9-byte `text/plain` body `Xbox NCSI`. Protocol knowledge comes only from OpenPak's
own documents and the firmware's observed behaviour; see the [clean-room policy](../docs/clean-room-policy.md).

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

## Not yet

Every other `notice.xbox.com` and Xbox Live page. The locale is accepted but not interpreted: all locales get the same
English text.

## Layout

`cmd/x360-notice` main (request logging); `internal/config`; `internal/notice` (routes). License: AGPL-3.0-only, like the
rest of the OpenPak servers. Releases: tag `v*.*.*` publishes `ghcr.io/openpak/x360-notice` and deployed hosts auto-update
from it, so nothing is tagged until a deploy is wanted.
