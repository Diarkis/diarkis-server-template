FROM alpine:latest
ARG TARGETARCH

COPY --chmod=755 ${TARGETARCH}/http /go/bin/
COPY --chmod=755 ${TARGETARCH}/health-check /go/bin/
CMD ["/go/bin/http"]
