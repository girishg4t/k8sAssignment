package main

import (
	"errors"
	"testing"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestCreateService(t *testing.T) {

	api := &k8sServiceAPI{
		s: &serviceResource{ns: "test-ns",
			name: "test-service",
			port: 8080,
			selector: map[string]string{
				"app": "demo-app",
			}},
		cs: testclient.NewSimpleClientset(),
	}

	err := createService(api)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteService(t *testing.T) {

	api := &k8sServiceAPI{
		s: &serviceResource{ns: "test-ns",
			name: "test-service",
			port: 8080,
			selector: map[string]string{
				"app": "demo-app",
			}},
		cs: testclient.NewSimpleClientset(),
	}
	createService(api)
	err := deleteService(api)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteServiceNotFound(t *testing.T) {

	api := &k8sServiceAPI{
		s: &serviceResource{ns: "test-ns",
			name: "test-service",
			port: 8080,
			selector: map[string]string{
				"app": "demo-app",
			}},
		cs: testclient.NewSimpleClientset(),
	}
	err := deleteService(api)
	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestIsServiceExists(t *testing.T) {

	api := &k8sServiceAPI{
		s: &serviceResource{ns: "test-ns",
			name: "test-service",
			port: 8080,
			selector: map[string]string{
				"app": "demo-app",
			}},
		cs: testclient.NewSimpleClientset(),
	}
	createService(api)
	exists := isServiceExists(api)
	if !exists {
		t.Fatal(errors.New("service should be present"))
	}
}

func TestIsServiceNotExists(t *testing.T) {

	api := &k8sServiceAPI{
		s: &serviceResource{ns: "test-ns",
			name: "test-service",
			port: 8080,
			selector: map[string]string{
				"app": "demo-app",
			}},
		cs: testclient.NewSimpleClientset(),
	}
	exists := isServiceExists(api)
	if exists {
		t.Fatal(errors.New("service should not be present"))
	}
}
