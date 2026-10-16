# Changelog: x360-notice

## Unreleased

- 2026-10-16: first cut. `GET /XBOXNCSI.txt` (`Xbox NCSI`, text/plain, no-cache) and
  `GET /xocdata/xml/<locale>/servicestatusChannel.xml` (RootNode/Table status XML, text/xml) for the Xbox 360
  dashboard's Test Xbox Live Connection; 404 for anything else. Tests with httptest. Not deployed, no tag, not yet
  tried by the emulator itself.
