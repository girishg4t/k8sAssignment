### Installing M7 application through helm chart

#### Subchart :
Each subchart can be enabled/disabled based on the condition specified in values.yaml file e.g.
```
oracle:
   enabled: false
mongodb:
   enabled: false
tomcat:
   enabled: true
   ...
```

Values of each subchart can be overridden through parent helm chart

```
libreoffice:
   image:
      registry: get.rnd.metricstream.com
      pullPolicy: IfNotPresent

activemq:
   image:
      registry: get.rnd.metricstream.com
      pullPolicy: IfNotPresent

...

```

#### Parent chart creates below common resources :
1) service account creation is condition based and name can be taken from values.yaml, if not specified it will be autogenerated
2) Image pull secrets is created based on value specified in secrets.yaml file


To run the chart we can do :
```
helm dependency update .
cp secrets.sample.yaml secrets.yaml
# Edit secrets.yaml with the credentials for registry
helm upgrade --install ms-app . -f secrets.yaml --set oracle.oracle.username=metricstream --set oracle.oracle.password=password -n dev

```

#### Pipeline pre -requisite
1) k8s Cluster should should be up and running
2) helm should be installed in the master node
3) helm push plugin should be installed
```
helm plugin install https://github.com/chartmuseum/helm-push

```
4) Passwordless connection between runner and master node
------------------------------------------------------------------------
## Kubernetes Assignment 
### Problem Statement

Create a controller to automatically create a Service resource for every new Deployment
Consider the following constraints for creating the Service:  

1) Service ports and targetPorts should be taken from pod.spec.containers.ports  
2) If ports are not set for containers in the specs, consider port 8080 as a default port.  
3) There should be a way to disable auto-creation of service by setting annotations.  
e.g if "auto-create-svc: false" is set in the annotation of the parent resource (Deployment) spec, the controller should not create a service for that resource
4) Default service type should be ClusterIP, but should be configurable through annotations of the parent resource (i.e Deployment)
     e.g "auto-create-svc-type: NodePort"  
5) The Services should get automatically removed when the parent resource - Deployment - is deleted  
