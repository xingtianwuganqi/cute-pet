# 第一阶段：编译Go应用
FROM golang:1.20 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
# 复制源代码到容器中
COPY . .
COPY ./config /config
# 编译Go应用
#Alpine镜像基于musl libc而不是glibc，因此如果你的Go应用依赖于glibc，
#那么它可能无法在Alpine镜像中运行。这可能就是为什么你看到no such file or directory错误的原因。
#为了解决这个问题，你可以尝试在构建Go应用时设置CGO_ENABLED=0环境变量以禁用cgo。
#这将生成一个静态二进制文件，它不依赖于系统的C库（如glibc或musl libc），因此可以在任何Linux系统上运行。
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# 设置环境变量
#ARG ENVIRONMENT
#ENV ENVIRONMENT ${ENVIRONMENT}

# 第二阶段：运行Go应用
FROM alpine:latest
WORKDIR /app
# 从builder阶段复制编译好的应用到当前镜像
COPY --from=builder /app/main /app
COPY --from=builder /config /config
COPY --from=builder /app/locales /app/locales

## 设置环境变量
#ARG PORT
## 暴露端口
#EXPOSE ${PORT}
# Example command to print the values
#RUN echo "ENVIRONMENT is $ENVIRONMENT"
#RUN echo "PORT is $PORT"

# 运行应用
CMD ["/app/main"]


