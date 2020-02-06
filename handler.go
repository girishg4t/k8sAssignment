package main

import (
	"fmt"
	sc "strconv"

	"github.com/girishg4t/k8sAssignment/utils"

	apps_v1 "k8s.io/api/apps/v1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
}

// TestHandler is a sample implementation of Handler
type TestHandler struct{}

// Init handles any handler initialization
func (t *TestHandler) Init() error {
	fmt.Println("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *TestHandler) ObjectCreated(obj interface{}) {
	fmt.Println("TestHandler.ObjectCreated")
	// assert the type to a Pod object to pull out relevant data
	dep := obj.(*apps_v1.Deployment)
	fmt.Println("Deployment object created", dep.Name)

	name := dep.GetName()

	if name == "" {
		return
	}

	isServiceReq, _ := sc.ParseBool(dep.ObjectMeta.Annotations["auto-create-svc"])
	if !isServiceReq {
		fmt.Println("service not required")
		return
	}

	service, err := extractServiceInfoFromDeployment(dep)
	if err != nil {
		return
	}
	cs := utils.GetKubeHandle()
	exists := isServiceExists(service, cs)
	if exists {
		return
	}
	err = createService(service, cs)
	if err != nil {
		panic(err)
	}
}

// ObjectDeleted is called when an object is deleted
func (t *TestHandler) ObjectDeleted(obj interface{}) {
	dep := obj.(*apps_v1.Deployment)
	name := dep.GetName()
	if name == "" {
		return
	}
	serviceName := name + "-service"
	s := &serviceResource{ns: ns,
		name: serviceName}
	cs := utils.GetKubeHandle()
	exists := isServiceExists(s, cs)
	if !exists {
		return
	}
	err := deleteService(s, cs)
	if err != nil {
		panic(err)
	}
}
