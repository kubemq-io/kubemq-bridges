module github.com/kubemq-hub/kubemq-bridges

go 1.15

require (
	github.com/fortytw2/leaktest v1.3.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty/v2 v2.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/json-iterator/go v1.1.10
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kubemq-hub/builder v0.6.2
	github.com/kubemq-io/kubemq-go v1.4.4
	github.com/labstack/echo/v4 v4.1.17
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.3 // indirect
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20200630154851-b2d8b0336632 // indirect
	google.golang.org/genproto v0.0.0-20200626011028-ee7919e894b5 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

//replace github.com/kubemq-hub/builder => ../builder
