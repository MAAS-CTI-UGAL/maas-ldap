package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestToTargetRewritesURLHostAndBody(t *testing.T) {
	var gotMethod, gotPath, gotQuery, gotHost, gotBody string
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		gotHost = r.Host

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll(target body) returned error: %v", err)
		}
		gotBody = string(body)

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("proxied"))
	}))
	defer targetServer.Close()

	targetURL, err := url.Parse(targetServer.URL + "/target/login")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/original?next=%2FMAAS", strings.NewReader("original"))
	rr := httptest.NewRecorder()

	if err := ToTarget(rr, req, *targetURL, []byte("rewritten")); err != nil {
		t.Fatalf("ToTarget() returned error: %v", err)
	}

	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
	if rr.Body.String() != "proxied" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "proxied")
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("target method = %q, want %q", gotMethod, http.MethodPost)
	}
	if gotPath != "/target/login" {
		t.Fatalf("target path = %q, want %q", gotPath, "/target/login")
	}
	if gotQuery != "next=%2FMAAS" {
		t.Fatalf("target query = %q, want %q", gotQuery, "next=%2FMAAS")
	}
	if gotHost != targetURL.Host {
		t.Fatalf("target host = %q, want %q", gotHost, targetURL.Host)
	}
	if gotBody != "rewritten" {
		t.Fatalf("target body = %q, want %q", gotBody, "rewritten")
	}
}

func TestToTargetPreservesInboundBodyWhenBodyIsNil(t *testing.T) {
	var gotBody string
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll(target body) returned error: %v", err)
		}
		gotBody = string(body)
	}))
	defer targetServer.Close()

	targetURL, err := url.Parse(targetServer.URL + "/target/login")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/original", strings.NewReader("original"))
	rr := httptest.NewRecorder()

	if err := ToTarget(rr, req, *targetURL, nil); err != nil {
		t.Fatalf("ToTarget() returned error: %v", err)
	}

	if gotBody != "original" {
		t.Fatalf("target body = %q, want %q", gotBody, "original")
	}
}

func TestToTargetReturnsProxyError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() returned error: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}

	target := url.URL{Scheme: "http", Host: addr, Path: "/target/login"}
	req := httptest.NewRequest(http.MethodPost, "/original", strings.NewReader("original"))
	rr := httptest.NewRecorder()

	if err := ToTarget(rr, req, target, nil); err == nil {
		t.Fatal("ToTarget() returned nil error")
	}
}
