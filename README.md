# Overview
Microservice for the reversing string and generating random number.
There are two main webservicess

1. random webservice (it includes redis cache) : exposes POST /api 
2. verso webservice : exposes POST /reverse 

The verso service api `reverse` is not exposed directly to end users.

# Runing webservices

```
$ git clone https://github.com/nitinpatil1992/microdigi.git
```

### Local execution

```
$ docker-compose up 

# wait for all services to up

$ curl -X POST -d '{"message":"abcdefg"}' localhost:9000/api

# on retry, you will get the same result along with random number as it is cached in redis 
$ curl -X POST -d '{"message":"abcdefg"}' localhost:9000/api
$ curl -X POST -d '{"message":"abcdefg"}' localhost:9000/api

```

### Minikube execution

```
$ minikube start

$ kubectl create ns diginex

$ helm init
# wait for tiller container to spin up 
$ kubectl get po -l'name=tiller' -n kube-system --watch

$ git clone https://github.com/nitinpatil1992/digihelm.git

$ cd digihelm/charts/verso

$ helm install --name=verso --namespace=diginex .

$ kubectl get all -n diginex

$ cd ../random

$ helm install --name=random --namespace=diginex .
# wait for random web service to spin up

$ kubectl get all -n diginex

$ curl -X POST -d '{"message":"0123456"}'  $(minikube ip):30501/api
```
