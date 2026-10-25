# Protocol: x360-notice

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

The request shape (from internal research notes, "Connectivity test
traffic"): `GET <path> HTTP/1.0`, a `Host` header, user agent `Xbox Downloader/2.0.17559.0`, nothing else. The NCSI answer
Microsoft gave was `200 OK` with the 9-byte `text/plain` body `Xbox NCSI`. Protocol knowledge comes only from OpenPak's
own documents and the firmware's observed behaviour; see the [clean-room policy](https://docs.openpak.org/contributors/clean-room/).
