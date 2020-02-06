package main

import (
	"errors"
	"testing"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestCreateService(t *testing.T) {

	s := &serviceResource{
		ns:   "test-ns",
		name: "test-service",
		port: 8080,
		selector: map[string]string{
			"app": "demo-app",
		}}
	cs := testclient.NewSimpleClientset()
	err := createService(s, cs)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteService(t *testing.T) {

	s := &serviceResource{
		ns:   "test-ns",
		name: "test-service",
		port: 8080,
		selector: map[string]string{
			"app": "demo-app",
		}}
	cs := testclient.NewSimpleClientset()

	createService(s, cs)
	err := deleteService(s, cs)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteServiceNotFound(t *testing.T) {

	s := &serviceResource{ns: "test-ns",
		name: "test-service",
		port: 8080,
		selector: map[string]string{
			"app": "demo-app",
		}}
	cs := testclient.NewSimpleClientset()
	err := deleteService(s, cs)
	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestIsServiceExists(t *testing.T) {

	s := &serviceResource{ns: "test-ns",
		name: "test-service",
		port: 8080,
		selector: map[string]string{
			"app": "demo-app",
		}}
	cs := testclient.NewSimpleClientset()
	createService(s, cs)
	exists := isServiceExists(s, cs)
	if !exists {
		t.Fatal(errors.New("service should be present"))
	}
}

func TestIsServiceNotExists(t *testing.T) {

	s := &serviceResource{ns: "test-ns",
		name: "test-service",
		port: 8080,
		selector: map[string]string{
			"app": "demo-app",
		}}
	cs := testclient.NewSimpleClientset()
	exists := isServiceExists(s, cs)
	if exists {
		t.Fatal(errors.New("service should not be present"))
	}
}
