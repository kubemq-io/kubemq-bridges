# KubeMQ Bridges Queue Source

KubeMQ  Bridges Queue source provides an RPC queue subscriber for processing target commands.

## Prerequisites
The following are required to run the queue source connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Queue source configuration properties:

| Property     | Required | Description                                             | Possible Values                                      |
|:-------------|:---------|:--------------------------------------------------------|:-----------------------------------------------------|
| address      | yes      | kubemq server address (gRPC interface)                  | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id    | no       | sets client_id value for connection                     | "cluster-a-queue-connection"                         |
| auth_token   | no       | JWT auth token for connection authentication            | JWT token                                            |
| channel      | yes      | kubemq channel to pull queue messages                   | queue.a                                              |
| batch_size   | no       | sets how many messages the source will pull in one call | default - 1                                          |
| wait_timeout | no       | sets how many seconds to wait per each pull             | 60                                                   |


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
      kind: source.queue # Sources kind
      name: 3-clusters-source # sources name 
      connections: # Array of connections settings per each source kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-queue-connection"
          auth_token: ""
          channel: "queue"
          batch_size: 1
          wait_timeout: 60
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-queue-connection"
          auth_token: ""
          channel: "queue"
          batch_size: 1
          wait_timeout: 60
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-queue-connection"
          auth_token: ""
          channel: "queue"
          batch_size: 1
          wait_timeout: 60    
    targets:
    .....
```

