# KubeMQ Bridges Events-Store Source

KubeMQ Bridges Events-Store source provides an RPC events-store subscriber for processing target commands.

## Prerequisites
The following are required to run the events-store source connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Events-Store source connector configuration properties:

| Properties Key             | Required | Description                            | Example                                              |
|:---------------------------|:---------|:---------------------------------------|:-----------------------------------------------------|
| address                    | yes      | kubemq server address (gRPC interface) | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id                  | no       | set client id                          | "client_id"                                          |
| auth_token                 | no       | set authentication token               | JWT token                                            |
| channel                    | yes      | set channel to subscribe               |                                                      |
| group                      | no       | set subscriber group                   |                                                      |
| auto_reconnect             | no       | set auto reconnect on lost connection  | "false", "true"                                      |
| reconnect_interval_seconds | no       | set reconnection seconds               | "5"                                                  |
| max_reconnects             | no       | set how many time to reconnect         | "0"                                                  |


Example:

```yaml
bindings:
  - name:  events-store-binding 
    properties: 
      log_level: error
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
      rate_per_second: 100
    sources:
      kind: source.events-store # Sources kind
      name: 3-clusters-source # sources name 
      connections: # Array of connections settings per each source kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-events-store-connection"
          auth_token: ""
          channel: "events-store"
          group: ""
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-events-store-connection"
          auth_token: ""
          channel: "events-store"
          group: ""
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-events-store-connection"
          auth_token: ""
          channel: "events-store"
          group: ""              
    targets:
    .....
```

