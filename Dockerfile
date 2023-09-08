# 设置基础镜像
FROM golang:1.21-alpine as builder

# 设置工作目录
WORKDIR /app

# 将应用程序代码复制到镜像中
COPY . .

# 复制依赖项清单并下载依赖项
RUN go env -w GO111MODULE=on \
&& go env -w GOPROXY=https://goproxy.cn,direct \
&& go mod tidy -go=1.21

# build 应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auction-backend .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auction-backend /app/
COPY --from=builder /app/conf/app.toml  /app/conf/
COPY --from=builder /app/conf/rbac_model.conf  /app/conf/

# 启动应用程序
CMD ["/app/auction-backend"]