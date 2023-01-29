package app

import (
	"log"
	"reflect"
)

type ServiceLocator struct {
	services map[string]any
}

func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		services: make(map[string]any),
	}
}

func (s *ServiceLocator) Set(alias string, service any) {
	log.Printf("[DI] Inject service: %s, with alias: %s", getType(service), alias)

	s.services[alias] = service
}

func (s *ServiceLocator) Get(alias string) any {
	return s.services[alias]
}

func (s *ServiceLocator) Clear() {
	s.services = make(map[string]any)
}

func getType(service interface{}) string {
	if t := reflect.TypeOf(service); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
