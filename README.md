# Kubernetes deployment example (with Minikube!) 

In this example we are going to deploy a simple app which consists of a service and a small redis cluster. The app gives back an Geolocation data based on the IP-address provided. It also caches requests so if you already requested geodata for the same IP-address - you're going to get the data from Redis instead of API-call to external data provider.

The code of the app is out of the focus today. Let's focus on the k8s objects and how we use them.
  

# Step-by-step  

## 1: Build app and redis container images

All you need is to run  

`docker build -t svc-ip-location-image .
`

and 

`docker build -t redis-local .
` 

in the respective folders. This way you add the images to your local "repository" which minikube can use to pull the images from. Notice how our redis config file which we introduce and move inside the container has a password set for the redis instance. This leads us to secret and config management of our cluster.

> If you get issues with deployments which cannot find your image - You may need to add the images to the minikube cache or even build the images using minikube's own container engine depending on your environment. 

  

## 2: Secrets and Configs
### 2a: Configmap
We need to create both **ConfigMap** and **Secrets** kubernetes objects here. The goal of the config object is to hold configuration of our applications. In this example it simply holds the cluster address of the redis instance. To add this object to a cluster you need to apply the existing *app-configmap.yaml* using

    kubectl apply -f .\k8s\configs\app-configmap.yaml
Notice the name of the ConfigMap inside the file. It's in the metadata.name: app-configmap. It's useful to know the name of the object so we can check it using kubectl like this

    kubectl describe configmap app-configmap
It'll give us all the data this ConfigMap holds so we can be sure it's been created properly.
### 2b: Secrets
Adding Secrets object is quite simple as well. All you need to do is to apply the *app-secret.yaml* file like this

    kubectl apply -f .\k8s\secrets\app-secret.yaml

 Please notice two values we're adding here: ***apikey*** and ***redis-passwd***. Instead of adding those to the ConfigMap we add them here as a secret. The main difference between configs and secrets is that Configs are usually do not hold sensitive data, whereas secrets are ideal for stuff like passwords and keys. Here we store them inside the kubernetes object but I strongly recommend getting familiar with secret-management solutions (this is outside of the scope for this example).
 
 Just like with ConfigMap object, we can check out our secrets using kubectl like this:
 

    kubectl describe secret app-secret
As you can see, this command does not really show you the secret values as-is, rather it obfuscates them in a way. These are secrets after all!

## 3: Deployments and Services
### 3a: Redis Deployment
Now to the fun part. Deployments, in a sense, are the most important objects for you as a developer. If you interact with K8S - you interact with deployments. So getting familiar with those are crucial for understanding the whole thing and how it ties everything together.

Let's start with redis deployment. It's quite simple really. All you need to keep track of: replicas, ports and image name. To add the deployment object just apply it using

    kubectl apply -f .\k8s\deployments\redis-deployment.yaml
The lifecycle of a deployment then begins! You can check the pods out using

    kubectl describe pods redis

Let's move on to the more sophisticated example - our application deployment.

### 3b: Application Deployment and Service
You can notice right away that it's quite a bit more complicated now with more things to keep track of. But it's actually quite simple. Key aspects of this **Deployment** are:

1. Health check and Readiness Check endpoints which your app provides as a way to tell kubernetes (it polls the endpoint) that your app is well and healthy. 
2. Environment variables with values which we take from ConfigMap.
3. Env. variables with values from Secrets.

These are quite intuitive. After you got familiar with the deployment object file just apply it to your cluster.

    kubectl apply -f .\k8s\deployments\app-deployment.yaml
To see the deployment stats please use 

    kubectl describe deployment ip-location

After we got our deployments running, let's talk about **Services**. Of you've been attentive to the file contents thus far you might've notice that we had Service description in our deployment files. And now's the good time to take a look at the Service we described for our application.

The main two things in there are:

1. Service type: in this case - **NodePort**, a simple way to provide access to the app in a cluster.
2. NodePort *binding*: we bind a port in a cluster to a port on a local machine. In this example - 30000.

But this in itself won't let you access the app just yet. You need to get the tunnel running!

    minikube service ip-location
This provides a tunnel which you can then use to access your app. This is not the only way btw, but it's something you need to learn on your own :)
