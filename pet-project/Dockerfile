# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.20

# 设置工作目录
WORKDIR /app

# 复制项目代码到容器中
COPY . .

# 编译 Go 项目
RUN go build -o main

# 暴露应用程序端口
EXPOSE 8082

# 启动应用程序
CMD ["./main"]




