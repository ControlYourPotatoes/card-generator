package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Container defines the interface for dependency injection
type Container interface {
	RegisterSingleton(name string, factory interface{}) error
	RegisterTransient(name string, factory interface{}) error
	Resolve(name string) (interface{}, error)
	ResolveAs(name string, target interface{}) error
}

// Lifetime defines the lifetime of a service
type Lifetime int

const (
	Singleton Lifetime = iota
	Transient
)

// ServiceDescriptor holds service registration information
type ServiceDescriptor struct {
	Name     string
	Factory  interface{}
	Lifetime Lifetime
	Instance interface{}
}

// container implements the Container interface
type container struct {
	services map[string]*ServiceDescriptor
	mu       sync.RWMutex
}

// NewContainer creates a new dependency injection container
func NewContainer() Container {
	return &container{
		services: make(map[string]*ServiceDescriptor),
	}
}

// RegisterSingleton registers a service with singleton lifetime
func (c *container) RegisterSingleton(name string, factory interface{}) error {
	if err := c.validateFactory(factory); err != nil {
		return fmt.Errorf("invalid factory for %s: %w", name, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[name] = &ServiceDescriptor{
		Name:     name,
		Factory:  factory,
		Lifetime: Singleton,
	}

	return nil
}

// RegisterTransient registers a service with transient lifetime
func (c *container) RegisterTransient(name string, factory interface{}) error {
	if err := c.validateFactory(factory); err != nil {
		return fmt.Errorf("invalid factory for %s: %w", name, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[name] = &ServiceDescriptor{
		Name:     name,
		Factory:  factory,
		Lifetime: Transient,
	}

	return nil
}

// Resolve resolves a service by name
func (c *container) Resolve(name string) (interface{}, error) {
	c.mu.RLock()
	service, exists := c.services[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not registered", name)
	}

	// For singleton, return existing instance if available
	if service.Lifetime == Singleton && service.Instance != nil {
		return service.Instance, nil
	}

	// Create new instance
	instance, err := c.createInstance(service)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance of %s: %w", name, err)
	}

	// Store instance for singleton
	if service.Lifetime == Singleton {
		c.mu.Lock()
		service.Instance = instance
		c.mu.Unlock()
	}

	return instance, nil
}

// ResolveAs resolves a service and assigns it to the target pointer
func (c *container) ResolveAs(name string, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	if targetValue.IsNil() {
		return fmt.Errorf("target cannot be nil")
	}

	instance, err := c.Resolve(name)
	if err != nil {
		return err
	}

	instanceValue := reflect.ValueOf(instance)
	targetElement := targetValue.Elem()

	if !instanceValue.Type().AssignableTo(targetElement.Type()) {
		return fmt.Errorf("resolved instance type %s is not assignable to target type %s",
			instanceValue.Type(), targetElement.Type())
	}

	targetElement.Set(instanceValue)
	return nil
}

// validateFactory validates that the factory is a function
func (c *container) validateFactory(factory interface{}) error {
	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return fmt.Errorf("factory must be a function")
	}

	// Factory should return at least one value
	if factoryType.NumOut() == 0 {
		return fmt.Errorf("factory must return at least one value")
	}

	// If factory returns two values, second should be error
	if factoryType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !factoryType.Out(1).Implements(errorType) {
			return fmt.Errorf("second return value must be error")
		}
	}

	return nil
}

// createInstance creates a new instance using the factory function
func (c *container) createInstance(service *ServiceDescriptor) (interface{}, error) {
	factoryValue := reflect.ValueOf(service.Factory)
	factoryType := reflect.TypeOf(service.Factory)

	// Prepare arguments for factory function
	args := make([]reflect.Value, factoryType.NumIn())
	for i := 0; i < factoryType.NumIn(); i++ {
		paramType := factoryType.In(i)
		
		// Try to resolve parameter from container
		// For now, we'll pass zero values, but this can be enhanced
		// to support dependency resolution
		args[i] = reflect.Zero(paramType)
	}

	// Call factory function
	results := factoryValue.Call(args)

	// Check for error (if factory returns error as second value)
	if len(results) == 2 && !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results[0].Interface(), nil
}

// GetRegisteredServices returns the names of all registered services
func (c *container) GetRegisteredServices() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.services))
	for name := range c.services {
		names = append(names, name)
	}
	return names
}

// Clear removes all registered services
func (c *container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.services = make(map[string]*ServiceDescriptor)
} 