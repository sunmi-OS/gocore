## 容器化


DockerFile
```
FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/golang As builder

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

#step 1 build go cache
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

#step 2 build binary project
WORKDIR /project
ADD . .
RUN ls
RUN go build main.go

FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/centos:7.8.2003
#run binary project
WORKDIR /app
COPY --from=builder /project/main .

# your project shell [project] [arg1] [arg2] ...
CMD [ "/app/main","api","start"]
```

- 基础centos镜像：centos:7.8.2003
  - 依赖库打包缓存
  - 默认时区：东八区
  - nscd dns缓存避免core-dns高负载导致的dns解析异常
  - 基础工具：mtr net-tools telnet bind-utils wget
  - 联网检测脚本：网络通畅在启动程序
  - 二阶段构建减少最终包大小

