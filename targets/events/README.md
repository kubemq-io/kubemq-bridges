# KubeMQ Bridges Events Target

KubeMQ Bridges Events target provides an events sender for processing targets requests.

## Prerequisites
The following are required to run the events target connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Events target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| channels | no       | set array of channels values to send the event                |  "events.a,events.b,events.c"                                                    |

Example:

```yaml
bindings:
  - name:  events-binding 
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
      kind: target.events # Sources kind
      name: 3-clusters-targets # targets name 
      connections: # Array of connections settings per each target kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-events-connection"
          auth_token: ""
          channels: "events.a,events.b,events.c"
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-events-connection"
          auth_token: ""
          channels: "events.a,events.b,events.c"
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-events-connection"
          auth_token: ""
          channels: "events.a,events.b,events.c"
```

