FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-hub/kubemq-bridges
WORKDIR $GOPATH/github.com/kubemq-hub/kubemq-bridges
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-bridges-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER KubeMQ info@kubemq.io
LABEL name="KubeMQ Target Connectors" \
      maintainer="info@kubemq.io" \
      vendor="" \
      version="" \
      release="" \
      summary="" \
      description=""
COPY licenses /licenses
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir /kubemq-bridges
COPY --from=builder $GOPATH/github.com/kubemq-hub/kubemq-bridges/kubemq-bridges-run ./kubemq-bridges
RUN chown -R 1001:root  /kubemq-bridges && chmod g+rwX  /kubemq-bridges
WORKDIR kubemq-bridges
USER 1001
CMD ["./kubemq-bridges-run"]
