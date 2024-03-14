# 阶段一，用于构建
FROM golang:1.21.0 as build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o helloworld cmd/app/main.go

# 阶段二，用于运行
FROM golang:1.20-buster as runner

# 复制制品+配置文件
COPY --from=builder /usr/src/app/helloworld /opt/app/
RUN mkdir config
COPY --from=builder /usr/src/app/config/config.yml /opt/app/config/

#设置时区为上海
ENV TZ=Asia/Shanghai

# 执行分析程序
ENTRYPOINT ["/opt/app/helloworld"]