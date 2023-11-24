FROM alpine:latest

COPY bot /go/bin/
RUN chmod +x /go/bin/bot
ENV DIARKIS_BOT_SEVER_MODE true
CMD ["/go/bin/bin"]
