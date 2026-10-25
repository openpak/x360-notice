# Notes: x360-notice

Naming rule: every OpenPak Xbox 360 service is called `x360-<thing>`. Sibling: [`x360-xexds`](https://github.com/openpak/x360-xexds).

## Not yet

Every other `notice.xbox.com` and Xbox Live page. The locale is accepted but not interpreted: all locales get the same
English text.

## Layout

`cmd/x360-notice` main (request logging); `internal/config`; `internal/notice` (routes). License: AGPL-3.0-only, like the
rest of the OpenPak servers. Releases: tag `v*.*.*` publishes `ghcr.io/openpak/x360-notice` and deployed hosts auto-update
from it, so nothing is tagged until a deploy is wanted.
