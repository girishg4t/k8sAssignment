package main

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createService(api *k8sServiceAPI) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      api.s.name,
			Namespace: api.s.ns,
			Labels: map[string]string{
				"app": "demo-app",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       api.s.port,
					TargetPort: intstr.FromInt(int(api.s.port)),
				},
			},
			Selector: api.s.selector,
			Type:     api.s.stype,
		},
	}

	// Create Service
	fmt.Println("Creating service...")
	result, err := api.cs.CoreV1().Services(api.s.ns).Create(service)
	if err != nil {
		return err
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func deleteService(api *k8sServiceAPI) error {
	fmt.Println("Deleting service...")
	err := api.cs.CoreV1().Services(api.s.ns).Delete(api.s.name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted service %q.\n", api.s.name)
	return nil
}

func isServiceExists(api *k8sServiceAPI) bool {
	_, err := api.cs.CoreV1().Services(api.s.ns).Get(api.s.name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Service with name %s in namespace %s not found \n",
			api.s.name, api.s.ns)
		return false
	}
	return true
}

func extractServiceInfoFromDeployment(dep *v1.Deployment) (*serviceResource, error) {
	metadata := dep.Spec.Template.ObjectMeta
	spec := dep.Spec.Template.Spec
	serviceName := dep.GetName() + "-service"

	antServiceType := metadata.Annotations["auto-create-svc-type"]
	serviceType := corev1.ServiceTypeClusterIP
	if antServiceType == "NodePort" {
		serviceType = corev1.ServiceTypeNodePort
	}
	port := int32(8080)
	if spec.Containers[0].Ports != nil && len(spec.Containers[0].Ports) > 0 {
		port = spec.Containers[0].Ports[0].ContainerPort
	}

	service := &serviceResource{ns: dep.GetNamespace(),
		name:  serviceName,
		port:  port,
		stype: serviceType,
		selector: map[string]string{
			"app": "demo-app",
		}}

	return service, nil
}

func int32Ptr(i int32) *int32 { return &i }
