package main

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"

	"github.com/girishg4t/k8sAssignment/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type serviceResource struct {
	name     string
	stype    corev1.ServiceType
	port     int32
	selector map[string]string
	ns       string
}

type controller struct {
	cs kubernetes.Interface
	s  *serviceResource
}

const ns string = "mynamespace"

var c *controller

func main() {
	fmt.Println("Shared Informer app started")

	clientset := utils.GetKubeHandle()

	c = &controller{
		cs: clientset,
	}

	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset, 0,
		informers.WithNamespace(ns))
	informer := factory.Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		DeleteFunc: onDelete,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onAdd(obj interface{}) {
	dep := obj.(*v1.Deployment)
	name := dep.GetName()

	if name == "" {
		return
	}

	isServiceReq, _ := strconv.ParseBool(dep.ObjectMeta.Annotations["auto-create-svc"])
	if !isServiceReq {
		fmt.Println("service not required")
		return
	}

	service, err := extractServiceInfoFromDeployment(dep)
	if err != nil {
		return
	}
	c.s = service
	exists := isServiceExists(c)
	if exists {
		return
	}
	err = createService(c)
	if err != nil {
		panic(err)
	}
}

func onDelete(obj interface{}) {
	dep := obj.(*v1.Deployment)
	name := dep.GetName()
	if name == "" {
		return
	}
	serviceName := name + "-service"
	s := &serviceResource{ns: ns,
		name: serviceName}
	c.s = s
	exists := isServiceExists(c)
	if !exists {
		return
	}
	err := deleteService(c)
	if err != nil {
		panic(err)
	}
}
