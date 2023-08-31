FROM golang:alpine3.18 as builder
ADD AuthProxy /workspace
WORKDIR /workspace

# 挂载构建缓存。
# GOPROXY防止下载失败。
RUN --mount=type=cache,target=/go \
  env GOPROXY=https://goproxy.cn,direct \
  go build -o /workspace/auth-proxy /workspace

FROM alpine:3.18.2
LABEL authors="Riley Long"

###############################################################################
#                                INSTALLATION
###############################################################################
# Set the time zone to East 8.
RUN apk update && apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV WORKDIR=/app
ADD AuthProxy/static                         $WORKDIR/static
COPY --from=builder /workspace/auth-proxy    $WORKDIR/auth-proxy
RUN chmod +x $WORKDIR/auth-proxy

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
EXPOSE 80 6060
CMD ./auth-proxy