# KubeMQ Bridges - Bridge Example

In this example we demonstrate how to bridge query request from one cluster as source to multiple clusters as target.

![bridge-example](../../.github/assets/bridge-example.jpeg)

## Run

Run the following deployment

```bash
kubectl apply -f ./deploy.yaml
```
Where deploy.yaml:

```yaml
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-a
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30501
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-b
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30502
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqConnector
metadata:
  name: kubemq-bridges
  namespace: kubemq
spec:
  type: bridges
  replicas: 1
  image: kubemq/kubemq-bridges:latest
  config: |-
    bindings:
    - name: clusters-sources
      properties:
        log_level: "debug"
      sources:
        kind: source.queue
        connections:
          - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
            channel: "queue1"
      targets:
        kind: target.queue
        name: cluster-targets
        connections:
          - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
            channels: "queue1"


```

