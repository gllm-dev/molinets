package mcp_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	httpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/http"
	"github.com/gllm-dev/molinets/internal/adapters/handlers/http/healthzhdl"
	mcpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp"
	conversionHandler "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/conversion"
	grinderHandler "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/grinder"
	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded/grindersrepo"
	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded/methodsrepo"
	conversionApp "github.com/gllm-dev/molinets/internal/applications/conversion"
	grinderApp "github.com/gllm-dev/molinets/internal/applications/grinder"
)

func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	grinders, err := grindersrepo.New()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := methodsrepo.New()
	if err != nil {
		t.Fatal(err)
	}
	gApp := grinderApp.New(grinders)
	server := mcpServer.NewServer(grinderHandler.New(gApp), conversionHandler.New(conversionApp.New(gApp, methods)))
	ts := httptest.NewServer(httpServer.NewHandler(server, healthzhdl.New()))
	t.Cleanup(ts.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	cs, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: ts.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func call(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) (map[string]any, bool, string) {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: protocol error: %v", name, err)
	}
	var text string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			text += tc.Text
		}
	}
	var out map[string]any
	if res.StructuredContent != nil {
		b, _ := json.Marshal(res.StructuredContent)
		_ = json.Unmarshal(b, &out)
	}
	return out, res.IsError, text
}

func TestToolsAreListedReadOnly(t *testing.T) {
	cs := connect(t)
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"list_grinders": false, "convert_grind_setting": false, "grind_setting_for_method": false, "microns_to_setting": false}
	for _, tool := range res.Tools {
		want[tool.Name] = true
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("%s should be read-only", tool.Name)
		}
		if tool.InputSchema == nil || tool.OutputSchema == nil {
			t.Errorf("%s is missing a schema", tool.Name)
		}
	}
	for name, ok := range want {
		if !ok {
			t.Errorf("tool %s missing", name)
		}
	}
}

func TestListGrinders(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "list_grinders", map[string]any{"type": "electric"})
	if isErr {
		t.Fatal(text)
	}
	gs := out["grinders"].([]any)
	if len(gs) == 0 || len(gs) >= 34 {
		t.Fatalf("electric filter returned %d grinders", len(gs))
	}
	for _, g := range gs {
		if g.(map[string]any)["type"] != "electric" {
			t.Fatalf("non-electric grinder in the list: %v", g)
		}
	}
}

func TestConvertMatchesBeean(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "convert_grind_setting", map[string]any{
		"from_grinder": "Comandante C40", "setting": "26", "to_grinders": []string{"Chestnut C3S Pro"},
	})
	if isErr {
		t.Fatal(text)
	}
	eq := out["equivalents"].([]any)[0].(map[string]any)
	if s := eq["setting"]; s != "18" && s != "19" {
		t.Fatalf("C40 26 -> C3S Pro %v, want 18 or 19", s)
	}
}

func TestConvertToEveryGrinder(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "convert_grind_setting", map[string]any{"from_grinder": "k6", "setting": "1.10"})
	if isErr {
		t.Fatal(text)
	}
	if n := len(out["equivalents"].([]any)); n != 33 {
		t.Fatalf("got %d equivalents, want 33", n)
	}
}

func TestBurrContactWarning(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "convert_grind_setting", map[string]any{"from_grinder": "comandante", "setting": "0", "to_grinders": []string{"k6"}})
	if isErr {
		t.Fatal(text)
	}
	if out["approx_microns"].(float64) != 0 || out["matching_methods"] == nil || !strings.Contains(out["note"].(string), "burr contact") {
		t.Fatalf("got %v", out)
	}
}

func TestToolErrorCodes(t *testing.T) {
	cs := connect(t)
	cases := []struct {
		tool string
		args map[string]any
		want string
	}{
		{"convert_grind_setting", map[string]any{"from_grinder": "hario", "setting": "3"}, "AMBIGUOUS_GRINDER"},
		{"convert_grind_setting", map[string]any{"from_grinder": "bialetti", "setting": "3"}, "GRINDER_NOT_FOUND"},
		{"convert_grind_setting", map[string]any{"from_grinder": "comandante", "setting": "60"}, "SETTING_OUT_OF_RANGE"},
		{"convert_grind_setting", map[string]any{"from_grinder": "j-max", "setting": "1.9.0"}, "INVALID_SETTING"},
		{"convert_grind_setting", map[string]any{"from_grinder": " ", "setting": "3"}, "INVALID_ARGUMENT"},
		{"grind_setting_for_method", map[string]any{"grinder": "skerton", "method": "turkish"}, "METHOD_UNREACHABLE"},
		{"grind_setting_for_method", map[string]any{"grinder": "c40", "method": "turkish delight"}, "METHOD_NOT_FOUND"},
		{"microns_to_setting", map[string]any{"grinder": "c40", "microns": -5}, "INVALID_ARGUMENT"},
		{"list_grinders", map[string]any{"type": "hybrid"}, "INVALID_ARGUMENT"},
	}
	for _, c := range cases {
		_, isErr, text := call(t, cs, c.tool, c.args)
		if !isErr || !strings.HasPrefix(text, c.want+":") {
			t.Errorf("%s %v: isError=%v text=%q, want %s", c.tool, c.args, isErr, text, c.want)
		}
	}
}

func TestMethodWithBrewerName(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "grind_setting_for_method", map[string]any{"grinder": "c3s pro", "method": "chemex"})
	if isErr {
		t.Fatal(text)
	}
	if out["method"] != "Pour over" {
		t.Fatalf("chemex resolved to %v", out["method"])
	}
}

func TestPartialMethodRange(t *testing.T) {
	cs := connect(t)
	out, isErr, text := call(t, cs, "grind_setting_for_method", map[string]any{"grinder": "smart grinder pro", "method": "french press"})
	if isErr {
		t.Fatal(text)
	}
	if !strings.Contains(out["note"].(string), "only part") {
		t.Fatalf("Sage SGP tops out at 820 µm, so French press is partial: %v", out["note"])
	}
}

func TestMissingRequiredArgument(t *testing.T) {
	cs := connect(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "convert_grind_setting", Arguments: map[string]any{"from_grinder": "c40"}})
	if err == nil && (res == nil || !res.IsError) {
		t.Fatal("a missing setting must fail")
	}
}
