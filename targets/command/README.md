# KubeMQ Bridges Command Target

KubeMQ Bridges Command target provides an RPC command sender for processing sources requests.

## Prerequisites
The following are required to run the command target connector:

- kubemq cluster
- kubemq-bridges deployment


## Configuration

Command target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| default_channel | no       | set default channel to send request                |   "commands"                                                   |
| timeout_seconds | no       | sets command request default timeout (600 seconds) |                                                      |


Example:

```yaml
bindings:
  - name:  command-binding 
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
      kind: target.command # Sources kind
      name: 3-clusters-targets # targets name 
      connections: # Array of connections settings per each target kind
        - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-a-command-connection"
          auth_token: ""
          channel: "command"
          group: ""
        - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-b-command-connection"
          auth_token: ""
          channel: "command"
          group: ""
        - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
          client_id: "cluster-c-command-connection"
          auth_token: ""
          channel: "command"
          group: ""              
```

