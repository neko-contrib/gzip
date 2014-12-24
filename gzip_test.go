package gzip

import (
	"bufio"
	"github.com/rocwong/neko"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	gzipTestString = "gzip test"
)

func Test_Gzip(t *testing.T) {

	Convey("No gzip", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		m := neko.New()
		m.Use(Gzip(DefaultCompression))
		m.ServeHTTP(w, req)

		_, ok := w.HeaderMap[HeaderContentEncoding]
		So(ok, ShouldBeFalse)

		ce := w.Header().Get(HeaderContentEncoding)
		So(strings.EqualFold(ce, "gzip"), ShouldBeFalse)
	})

	Convey("Gzip response content", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(HeaderAcceptEncoding, "gzip")

		m := neko.New()
		m.Use(Gzip(DefaultCompression))
		m.GET("/", func(ctx *neko.Context) {
			ctx.Text(gzipTestString)
		})
		m.ServeHTTP(w, req)

		_, ok := w.HeaderMap[HeaderContentEncoding]
		So(ok, ShouldBeTrue)

		ce := w.Header().Get(HeaderContentEncoding)
		So(strings.EqualFold(ce, "gzip"), ShouldBeTrue)
	})

	Convey("Invalid compression level", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(HeaderAcceptEncoding, "gzip")

		m := neko.New()
		m.Use(Gzip(99))
		m.GET("/", func(ctx *neko.Context) {
			ctx.Text(gzipTestString)
		})
		m.ServeHTTP(w, req)

		_, ok := w.HeaderMap[HeaderContentEncoding]
		So(ok, ShouldBeTrue)

		ce := w.Header().Get(HeaderContentEncoding)
		So(strings.EqualFold(ce, "gzip"), ShouldBeTrue)

		So(w.Body.String(), ShouldEqual, gzipTestString)
	})
}

type hijackableResponse struct {
	Hijacked bool
	header   http.Header
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{header: make(http.Header)}
}

func (h *hijackableResponse) Header() http.Header           { return h.header }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func Test_ResponseWriter_Hijack(t *testing.T) {
	Convey("Hijack gzip", t, func() {
		hijackable := newHijackableResponse()

		m := neko.New()
		m.Use(Gzip(DefaultCompression))
		m.Use(func(ctx *neko.Context) {
			hj, ok := ctx.Writer.(http.Hijacker)
			So(ok, ShouldBeTrue)
			hj.Hijack()
		})

		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		r.Header.Set(HeaderAcceptEncoding, "gzip")
		m.ServeHTTP(hijackable, r)

		So(hijackable.Hijacked, ShouldBeTrue)
	})
}
