# Prerequisites

* Kubernetes cluster with at least two nodes (anti-affinity requirement)
* kubectl
* helm
* sshuttle
* curl
* golang (dev)
* docker (dev)
* make (dev)

# Build

```
make all
```

# Run server

```
docker run --name server --rm kayrus/servclient:v1
```

# Run client

```
docker run --name client --rm -e APP_URL=http://$(docker inspect server --format '{{.NetworkSettings.IPAddress}}'):8080 kayrus/servclient:v1 /opt/siclient
# or
APP_URL=http://$(docker inspect sa --format '{{.NetworkSettings.IPAddress}}'):8080 ./siclient
```

# curl

```
curl localhost:8081
# or using docker
curl $(docker inspect client --format '{{.NetworkSettings.IPAddress}}'):8081
```

# Deploy

```
kubectl create -f sa.yaml
helm init --debug --service-account=tiller
helm install --name mychart mychart
```

# Upgarde

## Major

Change the [chart](mychart/Chart.yaml#L4) and [software version](mychart/values.yaml#L4) and run:

```
helm install --name mychart-v2 mychart
```

this will allow you to have two major software versions.

## Minor

Change the [software version](mychart/values.yaml#L4) and run:

```
helm upgrade mychart mychart
```

# Test

## Prepare

```
wget https://github.com/kayrus/kuttle/raw/master/kuttle
chmod +x kuttle
kubectl run kuttle --image=alpine:latest --restart=Never -- sh -c 'apk add python --update && exec tail -f /dev/null'
sshuttle -r kuttle -e ./kuttle $(kubectl get svc --no-headers --output=custom-columns=:spec.clusterIP | tr '\n' ' ')
```

## curl

```
curl $(kubectl get svc siclient-v1 --no-headers --output=custom-columns=:spec.clusterIP):8081
curl $(kubectl get svc siserver-v1 --no-headers --output=custom-columns=:spec.clusterIP):8080
```
