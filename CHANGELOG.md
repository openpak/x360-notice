# Changelog: x360-notice

## v0.1.1 (2026-11-03)

- Every request logs one stderr line, `http <method> <path> -> <status> <duration>ms` (x360-common v0.5.1 `httplog`): no query, headers, cookies, tokens or bodies; `/healthz` skipped. Replaces the slog `request` line (which also logged host and user agent).

## Unreleased

- 2026-10-16: first cut. `GET /XBOXNCSI.txt` (`Xbox NCSI`, text/plain, no-cache) and
  `GET /xocdata/xml/<locale>/servicestatusChannel.xml` (RootNode/Table status XML, text/xml) for the Xbox 360
  dashboard's Test Xbox Live Connection; 404 for anything else. Tests with httptest. Not deployed, no tag, not yet
  tried by the emulator itself.
