FROM golang:1.21-alpine

WORKDIR /app

# 安装 Git（用于 go mod）
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
