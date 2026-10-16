package notice

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const wantXML = "<?xml version=\"1.0\"?>\r\n<RootNode>\r\n\t<Table>\r\n\t\t<MessageText>Xbox LIVE is up and running</MessageText>\r\n\t\t<Severity> </Severity>\r\n\t\t<Category> </Category>\r\n\t</Table>\r\n</RootNode>"

func get(t *testing.T, h http.Handler, method, path string) *http.Response {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec.Result()
}

func TestNCSI(t *testing.T) {
	res := get(t, Handler("x"), "GET", "/XBOXNCSI.txt")
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 || string(body) != "Xbox NCSI" || len(body) != 9 {
		t.Fatalf("%d %q", res.StatusCode, body)
	}
	if got := res.Header.Get("Content-Type"); got != "text/plain" {
		t.Fatalf("content type %q", got)
	}
	if got := res.Header.Get("Cache-Control"); got != "no-cache, no-store" {
		t.Fatalf("cache control %q", got)
	}
	if got := res.Header.Get("Content-Length"); got != "9" {
		t.Fatalf("content length %q", got)
	}
}

func TestStatusXMLAnyLocale(t *testing.T) {
	h := Handler("Xbox LIVE is up and running")
	for _, loc := range []string{"en-US", "de-DE", "fr-FR", "ja-JP", "en", "pt-BR", "zh-Hant-TW"} {
		res := get(t, h, "GET", "/xocdata/xml/"+loc+"/servicestatusChannel.xml")
		body, _ := io.ReadAll(res.Body)
		if res.StatusCode != 200 || string(body) != wantXML {
			t.Errorf("%s: %d %q", loc, res.StatusCode, body)
		}
		if res.Header.Get("Content-Type") != "text/xml" {
			t.Errorf("%s: content type %q", loc, res.Header.Get("Content-Type"))
		}
	}
}

func TestStatusTextEscaped(t *testing.T) {
	res := get(t, Handler("A & <B>"), "GET", "/xocdata/xml/en-US/servicestatusChannel.xml")
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "<MessageText>A &amp; &lt;B&gt;</MessageText>") {
		t.Fatalf("%q", body)
	}
}

func TestNotFound(t *testing.T) {
	h := Handler("x")
	for _, c := range [][2]string{
		{"GET", "/"},
		{"GET", "/xboxncsi.txt"},
		{"GET", "/XBOXNCSI.txt/"},
		{"GET", "/xocdata/xml/servicestatusChannel.xml"},
		{"GET", "/xocdata/xml/en-US/other.xml"},
		{"GET", "/xocdata/xml/..%2f/servicestatusChannel.xml"},
		{"GET", "/xocdata/xml/not%20a%20locale/servicestatusChannel.xml"},
		{"GET", "/xocdata/xml/en-US/x/servicestatusChannel.xml"},
		{"GET", "/xocdata/xml/1234567890/servicestatusChannel.xml"},
	} {
		if res := get(t, h, c[0], c[1]); res.StatusCode != 404 {
			t.Errorf("%s %s: %d", c[0], c[1], res.StatusCode)
		}
	}
	for _, m := range []string{"POST", "PUT", "DELETE"} {
		if res := get(t, h, m, "/XBOXNCSI.txt"); res.StatusCode != 405 {
			t.Errorf("%s: %d", m, res.StatusCode)
		}
	}
}

func TestHead(t *testing.T) {
	res := get(t, Handler("x"), "HEAD", "/XBOXNCSI.txt")
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 || len(body) != 0 || res.Header.Get("Content-Length") != "9" {
		t.Fatalf("%d %q %q", res.StatusCode, body, res.Header.Get("Content-Length"))
	}
}

// The firmware's request, as the dashboard sends it: HTTP/1.0, Host, a User-Agent, nothing else.
func TestFirmwareRequestOverTCP(t *testing.T) {
	srv := httptest.NewServer(Handler("Xbox LIVE is up and running"))
	defer srv.Close()
	for _, c := range []struct{ host, path, want, ctype string }{
		{"xncsi.xboxlive.com", "/XBOXNCSI.txt", "Xbox NCSI", "text/plain"},
		{"notice.xbox.com", "/xocdata/xml/en-US/servicestatusChannel.xml", wantXML, "text/xml"},
	} {
		conn, err := net.Dial("tcp", srv.Listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(conn, "GET "+c.path+" HTTP/1.0\r\nHost: "+c.host+"\r\nUser-Agent: Xbox Downloader/2.0.17559.0\r\n\r\n")
		res, err := http.ReadResponse(bufio.NewReader(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		conn.Close()
		if res.StatusCode != 200 || string(body) != c.want || res.Header.Get("Content-Type") != c.ctype {
			t.Errorf("%s: %d %q %q", c.path, res.StatusCode, body, res.Header.Get("Content-Type"))
		}
		if res.ContentLength != int64(len(c.want)) {
			t.Errorf("%s: content length %d", c.path, res.ContentLength)
		}
	}
}

func TestHealthz(t *testing.T) {
	if res := get(t, Handler("x"), "GET", "/healthz"); res.StatusCode != 200 {
		t.Fatal(res.StatusCode)
	}
}
