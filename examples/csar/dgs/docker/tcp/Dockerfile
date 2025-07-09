FROM alpine:latest

COPY tcp /go/bin/
COPY health-check /go/bin/
RUN chmod +x /go/bin/tcp && chmod +x /go/bin/health-check
CMD ["/go/bin/tcp"]
