FROM alpine:3.17.3 AS production

USER root
WORKDIR /

COPY osp .
COPY pkg/version/VERSION .

ENTRYPOINT ["/osp"]
