FROM ubuntu:24.04

COPY dgs_bin /dgs_bin
COPY remote_bin/health-check /go/bin/

RUN chmod +x /dgs_bin/linux-server.x86_64 && chmod +x /go/bin/health-check

CMD ["/dgs_bin/linux-server.x86_64", "/configs/shared/log.json", "/configs/shared/mesh.json", "0.0.0.0", "7400", "127.0.0.1"]
