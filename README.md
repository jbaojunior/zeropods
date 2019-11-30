# **Zeropods**

Zeropods is a project to scale (up or down) deployments and statefulset resources on kubernetes. This was created by using auto-scale feature on cluster to shutdown nodes that do not need stay on, like nodes that have only developement/test pods.

To auto-cluster-scale do the right job I created this app to scale down (0 pods) and scale up (the scale value before the down action) all of deployments and statefulset resources. With 0 pods the autoscaler can shutdown nodes, saving money =].

Zeropods is build using Golang :heart:

In the process of scale down I create a annotation (*zeropods/last-scale-number*) to save the size of replicas. With this annotation I can return to original scale size.

## **Options**
```
# zeropods -h
Usage of /usr/local/bin/zeropods:
  -action string
        Action to do. Possible values are "up" or "down"
  -conn string
        Connect method to cluster. Possible values are "cluster" and "config".
                 - "cluster" is to deploy inside a cluster, using a Service Account.
                 - "config" is to using outsite of cluster, with a kubeconfig (default "cluster")
  -kubeconfig string
        (optional) absolute path to the kubeconfig file. Usage when the parameter "connection" is "config" (default "/root/.kube/config")
  -n string
        Namespace to do the action
```

## **Build**
### Manually
Clone this repo and execute inside:
```
export GO111MODULE=on
go mod init
go get k8s.io/client-go@kubernetes-1.15.3
go build -o zeropods
```

### Docker
```
docker build -t zeropods .
```
PS.: If you want to use this image inside minikube is need execute a command before of build:
```
eval $(minikube docker-env)
docker build -t zeropods .
```

## **Executing**

Outsite of cluster we need a kubeconfig with a defined cluster set to run the commands:
```
kubectl config use-context minikube
docker run --rm --name zeropods -v ${HOME}:${HOME} zeropods "/usr/local/bin/zeropods -n dev -action down -conn config -kubeconfig ${HOME}/.kube/config"
```

Inside cluster we need specified a Service Account with permissions to scale deployments and statefulsets. Is a example of ClusterRole in the directory [kubernetes](./kubernetes).

## **Minikube**
I created example to be deployed inside a minikube cluster. First build the image using minikube environment. After apply the resource on the cluster:
```
kubectl apply -f kubernetes/*
```

Feel confortable to modify and use .