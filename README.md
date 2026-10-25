# Xbox 360 Notice Server

OpenPak's Xbox 360 Notice Server.

| | |
| --- | --- |
| Network | Xbox 360 |
| Role | Plain-HTTP pages the Xbox 360 dashboard reads during Test Xbox Live Connection: the connectivity text and the service-status XML |
| Used by | Xbox 360 dashboard (firmware 17559), Wharf emulator |
| Image | `ghcr.io/openpak/x360-notice` |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| `GET /XBOXNCSI.txt` (`xncsi.xboxlive.com`) | Connectivity text, body `Xbox NCSI` |
| `GET /xocdata/xml/<locale>/servicestatusChannel.xml` (`notice.xbox.com`) | Service-status XML (`RootNode`/`Table`) |

## Run

```sh
go run ./cmd/x360-notice          # listens on :8080
curl -i http://127.0.0.1:8080/XBOXNCSI.txt
```

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
