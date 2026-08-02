package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	httpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/http"
	"github.com/gllm-dev/molinets/internal/adapters/handlers/http/healthzhdl"
)

func TestRoutes(t *testing.T) {
	s := mcp.NewServer(&mcp.Implementation{Name: "t", Version: "0"}, nil)
	ts := httptest.NewServer(httpServer.NewHandler(s, healthzhdl.New()))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), `"ok"`) {
		t.Fatalf("healthz: %d %s", res.StatusCode, body)
	}

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mcp", strings.NewReader(
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), `"serverInfo"`) {
		t.Fatalf("initialize: %d %s", res.StatusCode, body)
	}

	if res, _ := http.Get(ts.URL + "/nope"); res.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown route: %d", res.StatusCode)
	}
}

func TestLoadServerConfig(t *testing.T) {
	t.Setenv("PORT", "9090")
	cfg, err := httpServer.LoadServerConfig()
	if err != nil || cfg.Port != 9090 {
		t.Fatalf("%+v %v", cfg, err)
	}
	t.Setenv("PORT", "70000")
	if _, err := httpServer.LoadServerConfig(); err == nil {
		t.Fatal("port 70000 should be rejected")
	}
}
