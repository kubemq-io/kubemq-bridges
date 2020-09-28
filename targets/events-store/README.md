# KubeMQ Bridges Events-Store Target

KubeMQ Bridges Events-Store target provides an events-store sender for processing sources requests.

## Prerequisites
The following are required to run the events-store target connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Events-Store target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| channels | no       | set array of channels values to send the event                |  "events-store.a,events-store.b,events-store.c"                                                    |

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
    .....
    targets:
      kind: target.events-store # Sources kind
      name: 3-clusters-targets # targets name 
      connections: # Array of connections settings per each target kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-events-store-connection"
          auth_token: ""
          channels: "events-store.a,events-store.b,events-store.c"
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-events-store-connection"
          auth_token: ""
          channels: "events-store.a,events-store.b,events-store.c"
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-events-store-connection"
          auth_token: ""
          channels: "events-store.a,events-store.b,events-store.c"
```

