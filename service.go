package main

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createService(c *controller) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.s.name,
			Namespace: c.s.ns,
			Labels: map[string]string{
				"app": "demo-app",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       c.s.port,
					TargetPort: intstr.FromInt(int(c.s.port)),
				},
			},
			Selector: c.s.selector,
			Type:     c.s.stype,
		},
	}

	// Create Service
	fmt.Println("Creating service...")
	result, err := c.cs.CoreV1().Services(c.s.ns).Create(service)
	if err != nil {
		return err
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func deleteService(c *controller) error {
	fmt.Println("Deleting service...")
	err := c.cs.CoreV1().Services(c.s.ns).Delete(c.s.name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted service %q.\n", c.s.name)
	return nil
}

func isServiceExists(c *controller) bool {
	_, err := c.cs.CoreV1().Services(c.s.ns).Get(c.s.name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Service with name %s in namespace %s not found \n",
			c.s.name, c.s.ns)
		return false
	}
	return true
}

func extractServiceInfoFromDeployment(dep *v1.Deployment) (*serviceResource, error) {
	spec := dep.Spec.Template.Spec
	serviceName := dep.GetName() + "-service"

	antServiceType := dep.ObjectMeta.Annotations["auto-create-svc-type"]
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
