FROM alpine:latest

COPY http /go/bin/
COPY health-check /go/bin/
RUN chmod +x /go/bin/http && chmod +x /go/bin/health-check
CMD ["/go/bin/http"]
