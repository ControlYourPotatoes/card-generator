package di

import (
	"errors"
	"testing"
)

// Test types for DI testing
type TestService interface {
	GetName() string
}

type testServiceImpl struct {
	name string
}

func (t *testServiceImpl) GetName() string {
	return t.name
}

func newTestService() TestService {
	return &testServiceImpl{name: "test-service"}
}

func newTestServiceWithError() (TestService, error) {
	return &testServiceImpl{name: "test-service-with-error"}, nil
}

func newFailingTestService() (TestService, error) {
	return nil, errors.New("failed to create service")
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("NewContainer() returned nil")
	}

	// Verify container is empty initially
	services := container.GetRegisteredServices()
	if len(services) != 0 {
		t.Errorf("Expected empty container, got %d services", len(services))
	}
}

func TestRegisterSingleton(t *testing.T) {
	container := NewContainer()

	// Test successful registration
	err := container.RegisterSingleton("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register singleton: %v", err)
	}

	// Verify service is registered
	services := container.GetRegisteredServices()
	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}

	if services[0] != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", services[0])
	}
}

func TestRegisterTransient(t *testing.T) {
	container := NewContainer()

	// Test successful registration
	err := container.RegisterTransient("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register transient: %v", err)
	}

	// Verify service is registered
	services := container.GetRegisteredServices()
	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}
}

func TestRegisterInvalidFactory(t *testing.T) {
	container := NewContainer()

	tests := []struct {
		name    string
		factory interface{}
		wantErr bool
	}{
		{
			name:    "non-function factory",
			factory: "not a function",
			wantErr: true,
		},
		{
			name:    "function with no return values",
			factory: func() {},
			wantErr: true,
		},
		{
			name:    "function with invalid error return type",
			factory: func() (TestService, string) { return nil, "" },
			wantErr: true,
		},
		{
			name:    "valid function with one return",
			factory: newTestService,
			wantErr: false,
		},
		{
			name:    "valid function with error return",
			factory: newTestServiceWithError,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := container.RegisterSingleton(tt.name, tt.factory)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterSingleton() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	container := NewContainer()

	// Register a service
	err := container.RegisterSingleton("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	// Resolve the service
	instance, err := container.Resolve("test-service")
	if err != nil {
		t.Fatalf("Failed to resolve service: %v", err)
	}

	// Verify instance type
	service, ok := instance.(TestService)
	if !ok {
		t.Fatalf("Resolved instance is not TestService type")
	}

	if service.GetName() != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", service.GetName())
	}
}

func TestResolveNonExistentService(t *testing.T) {
	container := NewContainer()

	_, err := container.Resolve("non-existent")
	if err == nil {
		t.Fatal("Expected error when resolving non-existent service")
	}

	expectedErr := "service non-existent not registered"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestSingletonLifetime(t *testing.T) {
	container := NewContainer()

	// Register singleton
	err := container.RegisterSingleton("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register singleton: %v", err)
	}

	// Resolve twice
	instance1, err := container.Resolve("test-service")
	if err != nil {
		t.Fatalf("Failed to resolve service first time: %v", err)
	}

	instance2, err := container.Resolve("test-service")
	if err != nil {
		t.Fatalf("Failed to resolve service second time: %v", err)
	}

	// Verify same instance (pointer comparison)
	if instance1 != instance2 {
		t.Error("Singleton should return the same instance")
	}
}

func TestTransientLifetime(t *testing.T) {
	container := NewContainer()

	// Register transient
	err := container.RegisterTransient("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register transient: %v", err)
	}

	// Resolve twice
	instance1, err := container.Resolve("test-service")
	if err != nil {
		t.Fatalf("Failed to resolve service first time: %v", err)
	}

	instance2, err := container.Resolve("test-service")
	if err != nil {
		t.Fatalf("Failed to resolve service second time: %v", err)
	}

	// Verify different instances (pointer comparison)
	if instance1 == instance2 {
		t.Error("Transient should return different instances")
	}
}

func TestResolveAs(t *testing.T) {
	container := NewContainer()

	// Register service
	err := container.RegisterSingleton("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	// Resolve using ResolveAs
	var service TestService
	err = container.ResolveAs("test-service", &service)
	if err != nil {
		t.Fatalf("Failed to resolve service with ResolveAs: %v", err)
	}

	if service == nil {
		t.Fatal("ResolveAs returned nil service")
	}

	if service.GetName() != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", service.GetName())
	}
}

func TestResolveAsWithInvalidTarget(t *testing.T) {
	container := NewContainer()

	// Register service
	err := container.RegisterSingleton("test-service", newTestService)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	tests := []struct {
		name   string
		target interface{}
	}{
		{
			name:   "non-pointer target",
			target: "not a pointer",
		},
		{
			name:   "nil pointer",
			target: (*TestService)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := container.ResolveAs("test-service", tt.target)
			if err == nil {
				t.Error("Expected error with invalid target")
			}
		})
	}
}

func TestResolveWithFactoryError(t *testing.T) {
	container := NewContainer()

	// Register service that returns error
	err := container.RegisterSingleton("failing-service", newFailingTestService)
	if err != nil {
		t.Fatalf("Failed to register failing service: %v", err)
	}

	// Try to resolve
	_, err = container.Resolve("failing-service")
	if err == nil {
		t.Fatal("Expected error when factory fails")
	}

	if err.Error() != "failed to create instance of failing-service: failed to create service" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestClear(t *testing.T) {
	container := NewContainer()

	// Register some services
	container.RegisterSingleton("service1", newTestService)
	container.RegisterTransient("service2", newTestService)

	// Verify services are registered
	services := container.GetRegisteredServices()
	if len(services) != 2 {
		t.Errorf("Expected 2 services before clear, got %d", len(services))
	}

	// Clear container
	container.Clear()

	// Verify container is empty
	services = container.GetRegisteredServices()
	if len(services) != 0 {
		t.Errorf("Expected 0 services after clear, got %d", len(services))
	}
}

// Benchmarks
func BenchmarkRegisterSingleton(b *testing.B) {
	container := NewContainer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.RegisterSingleton("test-service", newTestService)
	}
}

func BenchmarkResolve(b *testing.B) {
	container := NewContainer()
	container.RegisterSingleton("test-service", newTestService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.Resolve("test-service")
	}
}

func BenchmarkResolveTransient(b *testing.B) {
	container := NewContainer()
	container.RegisterTransient("test-service", newTestService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = container.Resolve("test-service")
	}
}
