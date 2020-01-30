package main

import (
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"

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

type k8sServiceAPI struct {
	cs kubernetes.Interface
	s  *serviceResource
}

const ns string = "mynamespace"

var api *k8sServiceAPI

func main() {
	api = &k8sServiceAPI{
		cs: getKubeHandle(),
	}

	watchlist := cache.NewListWatchFromClient(api.cs.AppsV1().RESTClient(), "deployments",
		ns, fields.Everything())
	_, controller := cache.NewInformer(watchlist, &v1.Deployment{}, time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAdd,
			DeleteFunc: onDelete})

	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func onAdd(obj interface{}) {
	dep := obj.(*v1.Deployment)
	name := dep.GetName()

	if name == "" {
		return
	}
	service, err := extractServiceInfoFromDeployment(dep)
	if err != nil {
		return
	}
	api.s = service
	exists := isServiceExists(api)
	if exists {
		return
	}
	err = createService(api)
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
	api.s = s
	exists := isServiceExists(api)
	if !exists {
		return
	}
	err := deleteService(api)
	if err != nil {
		panic(err)
	}
}
