# KubeMQ Bridges Query Target

KubeMQ Bridges Query target provides an RPC query sender for processing sources requests.

## Prerequisites
The following are required to run the query target connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Query target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| default_channel | no       | set default channel to send request                |                                                      |
| timeout_seconds | no       | sets query request default timeout (600 seconds) |                                                      |


Example:

```yaml
bindings:
  - name:  query-binding 
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
      kind: target.query # Sources kind
      name: 3-clusters-targets # targets name 
      connections: # Array of connections settings per each target kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-query-connection"
          auth_token: ""
          channel: "query"
          group: ""
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-query-connection"
          auth_token: ""
          channel: "query"
          group: ""
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-query-connection"
          auth_token: ""
          channel: "query"
          group: ""              
```

