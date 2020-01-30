## Kubernetes Assignment  
#### Problem Statement :  
Create a controller to automatically create a Service resource for every new Deployment  
  
Consider the following constraints for creating the Service:  
  
1) Service ports and targetPorts should be taken from pod.spec.containers.ports 
2) If ports are not set for containers in the specs, consider port 8080 as a default port.  
3) There should be a way to disable auto-creation of service by setting annotations.  
e.g if "auto-create-svc: false" is set in the annotation of the parent resource (Deployment) spec, the controller should not create a service for that resource  
4) Default service type should be ClusterIP, but should be configurable through annotations of the parent resource (i.e Deployment)  
     e.g "auto-create-svc-type: NodePort"  
5) The Services should get automatically removed when the parent resource - Deployment - is deleted
Instructions  
  
### Add integration/unit tests wherever applicable.  
Please try to following programming standards standards wherever applicable.  
  
  

Deliverables  
Github repo/Archive containing all source code and related artifacts(README, docker files etc)  
