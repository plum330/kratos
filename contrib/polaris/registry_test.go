package polaris

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go"

	"github.com/go-kratos/kratos/v2/registry"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("183.47.111.80:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithTimeout(time.Second),
		WithHealthy(true),
		WithIsolate(false),
		WithRegistryNamespace("default"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	service, err := r.GetService(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
}

func TestDeregister(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithTimeout(time.Second),
		WithHealthy(true),
		WithIsolate(false),
		WithRegistryNamespace("default"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)
	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatch(t *testing.T) {
	name := fmt.Sprintf("test-ut-%d", time.Now().Unix())
	sdk, err := polaris.NewSDKContextByAddress("183.47.111.80:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithTimeout(time.Second),
		WithHealthy(false),
		WithIsolate(false),
		WithRegistryNamespace("default"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    name,
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	
	w, err := r.Watch(context.Background(), name)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	service, err := w.Next()
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Errorf("want 1, got %d", len(service))
	}

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    name,
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 2 {
		t.Errorf("want 1, got %d", len(service))
	}

	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    name,
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 1 {
		t.Errorf("want 1, got %d", len(service))
	}
	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    name,
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 0 {
		t.Errorf("want 0, got %d", len(service))
	}
}
