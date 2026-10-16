// Package notice serves the two plain-HTTP requests the Xbox 360 dashboard
// makes during Test Xbox Live Connection:
//
//	GET /XBOXNCSI.txt                                    (xncsi.xboxlive.com)
//	GET /xocdata/xml/<locale>/servicestatusChannel.xml   (notice.xbox.com)
//
// Both are documented in docs/xbox360/wharf-base-plan.md ("Connectivity test
// traffic"). The real firmware sends HTTP/1.0 with a Host header and a
// "Xbox Downloader/2.0.17559.0" User-Agent, nothing else. The hostnames are
// not looked at: the paths tell the two apart, and DNS decides who reaches us.
package notice

import (
	"encoding/xml"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NCSIBody is the exact 9-byte body of /XBOXNCSI.txt.
const NCSIBody = "Xbox NCSI"

var locale = regexp.MustCompile(`^[A-Za-z]{2,3}(-[A-Za-z0-9]{2,8}){0,2}$`)

// Handler returns the routes; statusText is the service-status message.
func Handler(statusText string) http.Handler {
	status := statusXML(statusText)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /XBOXNCSI.txt", func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Cache-Control", "no-cache, no-store")
		h.Set("Content-Type", "text/plain")
		h.Set("Content-Length", strconv.Itoa(len(NCSIBody)))
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			w.Write([]byte(NCSIBody))
		}
	})
	mux.HandleFunc("GET /xocdata/xml/{locale}/servicestatusChannel.xml", func(w http.ResponseWriter, r *http.Request) {
		if !locale.MatchString(r.PathValue("locale")) {
			http.NotFound(w, r)
			return
		}
		h := w.Header()
		h.Set("Content-Type", "text/xml")
		h.Set("Content-Length", strconv.Itoa(len(status)))
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			w.Write(status)
		}
	})
	// Not a firmware path: for the container/orchestrator only.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ok\n"))
	})
	return mux
}

// statusXML is the service-status document: one Table row with the message and
// two blank-space Severity and Category fields, CRLF line ends and tabs.
func statusXML(text string) []byte {
	var esc strings.Builder
	xml.EscapeText(&esc, []byte(text))
	return []byte("<?xml version=\"1.0\"?>\r\n<RootNode>\r\n\t<Table>\r\n\t\t<MessageText>" + esc.String() +
		"</MessageText>\r\n\t\t<Severity> </Severity>\r\n\t\t<Category> </Category>\r\n\t</Table>\r\n</RootNode>")
}

// New returns an http.Server for addr with conservative timeouts.
func New(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
