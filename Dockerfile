FROM        quay.io/prometheus/busybox:latest
LABEL maintainer="Timon Wong <timon86.wang@gmail.com>"

COPY ding2wechat /bin/ding2wechat

EXPOSE      8080
ENTRYPOINT  [ "/bin/ding2wechat" ]
