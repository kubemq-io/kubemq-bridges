# Kubemq Queue Source

Kubemq Queue source provides a queue subscriber for processing messages from queues

## Prerequisites
The following are required to run queue source connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Queue source connector configuration properties:

| Properties Key | Required | Description                                            | Example     |
|:---------------|:---------|:-------------------------------------------------------|:------------|
| address                    | yes      | kubemq server address (gRPC interface) | kubemq-cluster:50000 |
| client_id      | no       | set client id                                          | "client_id" |
| auth_token     | no       | set authentication token                               | jwt token   |
| channel        | yes      | set channel to subscribe                               |             |
| sources        | no      | set how many concurrent sources to subscribe                               |    1        |
| batch_size     | no      | set how many messages to pull from queue | "1"         |
| wait_timeout   | no      | set how long to wait for messages to arrive in seconds | "5"        |


Example:

```yaml
bindings:
  - name:  queue-binding 
    properties: 
      log_level: error
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
      rate_per_second: 100
    sources:
      kind: source.queue-stream # Sources kind
      name: 3-clusters-source # sources name 
      connections: # Array of connections settings per each source kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-queue-connection"
          auth_token: ""
          channel: "queue"
          sources: 1
          batch_size: "1"
          wait_timeout: "3600"
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-queue-connection"
          auth_token: ""
          channel: "queue"
          sources: 1
          batch_size: 1
          wait_timeout: "3600"
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-queue-connection"
          auth_token: ""
          channel: "queue"
          sources: 1
          batch_size: 1
          wait_timeout: "3600"
    targets:
    .....
```
