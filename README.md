# gocore
SUNMI go开发封装核心库 

## 终端开发

CLI库(ETH等使用):[CLI库(ETH等使用)](github.com/urfave/cli)

CLI库(K8S,ETCD等使用):[CLI库(K8S,ETCD等使用)](github.com/spf13/cobra)

终端仪表盘:[终端仪表盘](github.com/gizak/termui)

终端文字美化输出各种色彩终端:[终端文字美化](github.com/fatih/color)

在终端上输出进度条:[终端上输出进度条](github.com/schollz/progressbar)


## 系统组件

DNS库:[DNS库](github.com/miekg/dns)

docker(已使用):[docker](github.com/moby/moby)

kubernetes:[/home/keke/goTest/src/gocore/README.md](github.com/kubernetes/kubernetes)

持续交付平台Drone:[Drone](github.com/drone/drone)

frp内网穿透支持http,tcp,udp(已使用):[frp](github.com/fatedier/frp)

ngrok内网穿透:[github.com/inconshreveable/ngrok ](github.com/inconshreveable/ngrok)

stun打洞服务器go实现:[stun打洞](github.com/ccding/go-stun)

基于KCP协议UDP TO TCP 网络加速通道(已使用):[基于KCP协议UDP](github.com/xtaci/kcptun )

持续文件同步:[持续文件同步](github.com/syncthing/syncthing)

文件同步(支持各种云):[文件同步](github.com/ncw/rclone)

请求流量复制:[goreplay](github.com/buger/goreplay)

redis集群解决方案:[redis集群](github.com/CodisLabs/codis)

consul服务发现:[consul](www.consul.io)

etcd K/V数据库:[github.com/coreos/etcd](github.com/coreos/etcd)

实时分布式消息传递平台:[nsq](nsq.io)

消息推送集群服务:[github.com/Terry-Mao/gopush-cluster ](github.com/Terry-Mao/gopush-cluster)

以太坊整套协议钱包Go实现:[go-ethereum](github.com/ethereum/go-ethereum)


## 开发套件

go-kit微服务套件:[go-kit](github.com/go-kit/kit)

桌面UI套件(基于CGO):[CGO](github.com/andlabs/ui)

桌面UI库(基于HTML):[app](github.com/murlokswarm/app)

Log库(已使用):[logrus](github.com/Sirupsen/logrus)

zap日志库(集成gocore):[zap](github.com/uber-go/zap)

图像处理库:[bild](github.com/anthonynsimon/bild)

图像处理库:[imaging](github.com/disintegration/imaging)

日期处理库:[now](github.com/jinzhu/now)

viper配置文件读取库(已使用):[viper](github.com/spf13/viper)

类型转换库(已使用):[cast](github.com/spf13/cast)

唯一识别码库(已使用):[uuid](github.com/satori/go.uuid)

压缩文件处理库:[archiver](github.com/mholt/archiver)

连接池库(已使用):[go-commons-pool ](github.com/jolestar/go-commons-pool)

程序内部系统资源,可以对不同的资源做出不同的规则调整:[gopsutil](github.com/shirou/gopsutil)


## 数据文件处理

文件嵌入到编译文件:[statik](github.com/rakyll/statik)

文件嵌入到编译文件(html,css,js):[rice](github.com/GeertJohan/go.rice)

内存敏感数据处理:[memguard](github.com/awnumar/memguard)


## 第三方软件使用

邮件发送(已使用):[gomail](github.com/go-gomail/gomail)

数据库操作(已使用):[github.com/jinzhu/gorm](github.com/jinzhu/gorm)

数据库操作:[github.com/go-xorm/xorm ](github.com/go-xorm/xorm )

redis操作库(已使用):[redis](gopkg.in/redis.v5)

rabbitmq使用框架(已使用):[amqp](github.com/streadway/amqp)

levelDB处理:[goleveldb ](https://github.com/syndtr/goleveldb )

bolt嵌入式数据库:[bolt](github.com/boltdb/bolt)



### 解析库
JSON解析库(已使用):[gjson](github.com/tidwall/gjson)

CSV处理库:[csvutil](github.com/jszwec/csvutil)

msgpack binc  cbor json解密库:[ugorji](github.com/ugorji/go)

golang解密php序列化库:[gphp_session_decoder](github.com/yvasiyarov/php_session_decoder )

高性能json库:[json-iterator](github.com/json-iterator/go)

google-protobuf库:[protobuf ](github.com/golang/protobuf )


## 网络框架
http网路框架(已使用):[echo](github.com/labstack/echo)

http网路框架(已使用):[gin](github.com/gin-gonic/gin)

http网络框架:[martini](https://github.com/go-martini/martini)

超级快的 http 网路框架(已使用):[fasthttp](github.com/valyala/fasthttp)

KCP协议golang实现(已使用):[kcp](github.com/xtaci/kcp-go)

IOT库 支持各种协议:[gobot](github.com/hybridgroup/gobot)

socket.io协议Go实现(已使用):[go-socket.io](github.com/googollee/go-socket.io)


## 深度学习
Go语言实现机器学习框架:[golearn](github.com/sjwhitworth/golearn )

GO机器学习图书馆,包含各种各样的算法:[gorgonia](github.com/gorgonia/gorgonia)

go语言对Tensorflow的封装:[tfgo](github.com/galeone/tfgo )


## 依赖管理
官方包管理:[dep](github.com/golang/dep )

包管理工具(本地打包到项目):[godep](github.com/tools/godep )

包管理工具(类似Composer和pip):[glide](github.com/Masterminds/glide )

## 容器沙箱
docker 运行时沙箱:[gvisor](https://github.com/google/gvisor)
Kubernetes的虚拟机管理附件:[KubeVirt](https://github.com/kubevirt/kubevirt)