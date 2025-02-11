# Golang后端服务Dockerfile
FROM golang:1.22-bookworm

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY backend/go.mod ./
COPY backend/go.sum ./

# 下载依赖
RUN go mod tidy
RUN go mod download

# 复制项目代码
COPY backend/ .

# 构建项目
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# 暴露端口
EXPOSE 8080

# 启动服务
CMD ["./server"]
