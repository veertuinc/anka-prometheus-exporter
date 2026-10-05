package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/veertuinc/anka-prometheus-exporter/src/state"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

func TestGetRegistryTemplatesData_TagFetchErrorKeepsTemplateNames(t *testing.T) {
	registryState := state.GetState()
	registryState.TemplatesMap = map[string]types.Template{
		"template-a": {
			UUID: "template-a",
			Name: "old-name",
			Size: 1,
			Tags: []types.TemplateTag{{Name: "v1", Size: 1}},
		},
	}
	registryState.TemplateTagsFetchTime = map[string]time.Time{}
	t.Cleanup(func() {
		registryState.TemplatesMap = map[string]types.Template{}
		registryState.TemplateTagsFetchTime = map[string]time.Time{}
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/registry/vm" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("id") != "" {
			_, _ = w.Write([]byte(`{"status":"FAIL","message":"tag fetch failed","body":{}}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"OK","message":"","body":[{"id":"template-a","name":"new-name","size":2},{"id":"template-b","name":"other-template","size":3}]}`))
	}))
	t.Cleanup(server.Close)

	comm := &Communicator{controllerAddress: server.URL}
	_, err := comm.GetRegistryTemplatesData()
	if err == nil {
		t.Fatal("expected tag fetch error")
	}

	templates := registryState.GetTemplatesMap()
	templateA, ok := templates["template-a"]
	if !ok {
		t.Fatal("template-a missing from template map after tag fetch error")
	}
	if templateA.Name != "new-name" {
		t.Errorf("template-a name = %q, want new-name", templateA.Name)
	}
	if len(templateA.Tags) != 1 || templateA.Tags[0].Name != "v1" {
		t.Errorf("template-a tags = %+v, want cached tag v1", templateA.Tags)
	}

	templateB, ok := templates["template-b"]
	if !ok {
		t.Fatal("template-b missing from template map after tag fetch error")
	}
	if templateB.Name != "other-template" {
		t.Errorf("template-b name = %q, want other-template", templateB.Name)
	}
}
