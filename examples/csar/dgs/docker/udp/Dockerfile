FROM alpine:latest

COPY udp /go/bin/
COPY health-check /go/bin/
RUN chmod +x /go/bin/udp && chmod +x /go/bin/health-check
CMD ["/go/bin/udp"]
