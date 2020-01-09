# Trimg 

[![Actions Status](https://github.com/esakat/trimg/workflows/Go/badge.svg)](https://github.com/esakat/trimg/actions)

Trimg is CLI to help using Kubernetes with AWS Elastic Container Registry.

Background
----------

If you use kubernetes on AWS with off-line VPC, you should push images into ECR.  
Also if you want to use existing manifest, e.g. kubernetes-dashboard, metric-server..,  
you should replace manifest about image path.

trimg help you about it!

Installation
------------

homebrew:

    brew install esakat/trimg/trimg
    
Usage
-----

trimg support 2 feature

1. transfer(pull images, create ecr repository, push into ecr)
2. replace manifest about image path

### transfer

you can get target image from manifest file

```bash 
$ trimg transfer -f testfiles/input/replicaset.yml --dry-run
following images will be transfer
gcr.io/google_samples/gb-frontend:v3 -> <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/gcr.io/google_samples/gb-frontend:v3

$ trimg transfer -f testfiles/input/replicaset.yml
[gcr.io/google_samples/gb-frontend:v3] [==============================================================================]  100 %
1: gcr.io/google_samples/gb-frontend:v3 transfer to <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/gcr.io/google_samples/gb-frontend:v3

// ecr repository is created
$ aws ecr list-images --repository-name gcr.io/google_samples/gb-frontend
{
    "imageIds": [
        {
            "imageDigest": "sha256:60049e8aa1bb97242ce1a5fc5f9d86478d3f3407c2643edb054c717ac12c14bb",
            "imageTag": "v3"
        }
    ]
}

```

also, you can specify by manual
```bash
$ trimg transfer nginx:latest redis golang:1.13.5 --dry-run
following images will be transfer
nginx:latest -> <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/nginx:latest
redis -> <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/redis
golang:1.13.5 -> <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/golang:1.13.5
```

### replace

```bash 
$ trimg replace testfiles/input/replicaset.yml > replacedManifest.yml


$ cat testfiles/input/replicaset.yml  | grep image:
        image: gcr.io/google_samples/gb-frontend:v3
$ cat replacedManifest.yml | grep image:
        image: <YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/gcr.io/google_samples/gb-frontend:v3
```

### Use with Kubernetes

```bash
$ trimg transfer -f testfiles/input/replicaset.yml
$ trimg replace testfiles/input/replicaset.yml > replacedManifest.yml

// kubernetes on EKS
$ kubectl cluster-info | grep master
Kubernetes master is running at https://......eks.amazonaws.com

// deploy replaced manifes
$ kubectl apply -f replacedManifest.yml
replicaset.apps/frontend created

// check
$ kubectl get replicaset
NAME       DESIRED   CURRENT   READY   AGE
frontend   3         3         3       97s

$ kubectl get replicaset frontend -ojson | jq .spec.template.spec.containers[0].image
"<YourAccountId>.dkr.ecr.<YourDefaultRegion>.amazonaws.com/gcr.io/google_samples/gb-frontend:v3"
```
