FROM harbor.ceclouddyn.com/library/golang:alpine as builder
# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build
# 将代码复制到容器中
COPY ./cloud-monitor-region .

RUN go build -o app .


FROM  harbor.ceclouddyn.com/baseos/centos:8.4.2105-0.0.1
# 移动到用于存放生成的二进制文件的 /dist 目录
WORKDIR /dist
COPY --from=builder /build/app .
COPY --from=builder /build/cloud-monitor-region/*.yml .
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
    echo 'Asia/Shanghai' >/etc/timezone

# 声明服务端口
EXPOSE 80

# 启动容器时运行的命令
ENTRYPOINT ["./app"]
