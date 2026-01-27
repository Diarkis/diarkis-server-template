FROM alpine:latest

COPY mars /go/bin/
COPY health-check /go/bin/
RUN chmod +x /go/bin/mars && chmod +x /go/bin/health-check

CMD ["/go/bin/mars"]
