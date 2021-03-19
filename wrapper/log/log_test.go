package log

import (
	"context"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/memory"
	"testing"
)

var (
	// mock data
	testData = map[string][]*registry.Service{
		"foo": {
			{
				Name:    "foo",
				Version: "1.0.0",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.0-123",
						Address: "localhost:9999",
						Metadata: map[string]string{
							"protocol": "mucp",
						},
					},
					{
						Id:      "foo-1.0.0-321",
						Address: "localhost:9999",
						Metadata: map[string]string{
							"protocol": "mucp",
						},
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.1",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.1-321",
						Address: "localhost:6666",
						Metadata: map[string]string{
							"protocol": "mucp",
						},
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.3",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.3-345",
						Address: "localhost:8888",
						Metadata: map[string]string{
							"protocol": "mucp",
						},
					},
				},
			},
		},
	}
)

func newTestRegistry() registry.Registry {
	return memory.NewRegistry(memory.Services(testData))
}

func TestNewClientWrapper(t *testing.T) {
	var called bool
	service := "test.service"
	endpoint := "Test.Endpoint"
	address := "localhost:8080"

	r := newTestRegistry()
	c := NewClientWrapper(client.NewClient(
		client.Registry(r),
	))

	req := c.NewRequest(service, endpoint, nil)

	// test calling remote address
	if err := c.Call(context.Background(), req, nil, client.WithAddress(address)); err != nil {
		t.Fatal("call with address error", err)
	}

	if !called {
		t.Fatal("wrapper not called")
	}
}
