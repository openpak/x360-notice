# Evidence: x360-notice

Status, 2026-10-16: the firmware's Test Xbox Live Connection shows **Network: Connected** and, with the XEXDS probe
answered too, **Internet: Connected** against the stand-in equivalents of these two answers (private notes, `-ncsi`,
which is what these responses were copied from). This Go server is verified by unit tests (including the
firmware's own HTTP/1.0 request shape over TCP) and by a raw HTTP/1.0 request against a running instance. It has **not**
yet been pointed at by the emulator itself, not tried by a real console, and not deployed.
