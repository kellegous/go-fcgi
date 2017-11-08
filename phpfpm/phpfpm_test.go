package phpfpm

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kellegous/fcgi"
)

func mustGetWd() string {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return d
}

func TestSingleGet(t *testing.T) {
	p := MustStart(DefaultConfig)
	defer p.Shutdown()

	c, err := fcgi.Dial("tcp", p.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var bout, berr bytes.Buffer
	req, err := c.BeginRequest(map[string][]string{
		"SCRIPT_FILENAME": {filepath.Join(mustGetWd(), "hello.php")},
		"REQUEST_METHOD":  {"GET"},
		"CONTENT_LENGTH":  {"0"},
	}, nil, &bout, &berr)
	if err != nil {
		t.Fatal(err)
	}

	if err := req.Wait(); err != nil {
		t.Fatal(err)
	}

	res, err := readResponse(&bout)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", res.Body)
}

func TestServeHTTP(t *testing.T) {
	p := MustStart(DefaultConfig)
	defer p.Shutdown()

	c, err := fcgi.Dial("tcp", p.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := responseWriter{
		Head: map[string][]string{},
	}

	params := fcgi.ParamsFromRequest(req)
	params["SCRIPT_FILENAME"] = []string{filepath.Join(mustGetWd(), "hello.php")}

	c.ServeHTTP(params, &rw, req)
}
