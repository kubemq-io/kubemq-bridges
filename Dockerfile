FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-hub/kubemq-bridges
WORKDIR $GOPATH/github.com/kubemq-hub/kubemq-bridges
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags container -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-bridges-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER KubeMQ info@kubemq.io
LABEL name="KubeMQ Bridges Connectors" \
      maintainer="info@kubemq.io" \
      vendor="kubemq.io" \
      version="v1.0.0" \
      release="stable" \
      summary="KubeMQ Bridges bridge, replicate, aggregate, and transform messages between KubeMQ clusters no matter where they are, allowing to build a true cloud-native messaging single network running globally." \
      description=""
COPY licenses /licenses
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir /kubemq-connector
COPY --from=builder $GOPATH/github.com/kubemq-hub/kubemq-bridges/kubemq-bridges-run ./kubemq-connector
RUN chown -R 1001:root  /kubemq-connector && chmod g+rwX  /kubemq-connector
WORKDIR kubemq-connector
USER 1001
CMD ["./kubemq-bridges-run"]
