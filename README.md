# gocore
SUNMI go开发封装核心库 

## 终端开发

CLI库(ETH等使用)  :  [github.com/urfave/cli](github.com/urfave/cli)

CLI库(K8S,ETCD等使用)  :  [github.com/spf13/cobra](github.com/spf13/cobra)

终端仪表盘  :  [github.com/gizak/termui](github.com/gizak/termui)

终端文字美化输出各种色彩终端  :  [github.com/fatih/color](github.com/fatih/color)

在终端上输出进度条  [github.com/schollz/progressbar ](github.com/schollz/progressbar )


## 系统组件

DNS库  :  [github.com/miekg/dns](github.com/miekg/dns)

docker(已使用)  :  [github.com/moby/moby ](github.com/moby/moby )

k8s  :  [github.com/kubernetes/kubernetes ](github.com/kubernetes/kubernetes )

持续交付平台  :  [github.com/drone/drone](github.com/drone/drone)

内网穿透支持http,tcp,udp(已使用)  :  [github.com/fatedier/frp](github.com/fatedier/frp)

内网穿透  :  [github.com/inconshreveable/ngrok ](github.com/inconshreveable/ngrok )

stun打洞服务器go实现  :  [github.com/ccding/go-stun ](github.com/ccding/go-stun )

基于KCP协议UDP TO TCP 网络加速通道(已使用)  :  [github.com/xtaci/kcptun ](github.com/xtaci/kcptun )

持续文件同步  :  [github.com/syncthing/syncthing](github.com/syncthing/syncthing)

文件同步(支持各种云)  :  [github.com/ncw/rclone](github.com/ncw/rclone)

请求流量复制  :  [github.com/buger/goreplay](github.com/buger/goreplay)

redis集群解决方案  :  [github.com/CodisLabs/codis](github.com/CodisLabs/codis)

服务发现  :  [www.consul.io ](www.consul.io )

K/V数据库  :  [github.com/coreos/etcd ](github.com/coreos/etcd )

实时分布式消息传递平台  :  [nsq.io](nsq.io)

消息推送集群服务  :  [github.com/Terry-Mao/gopush-cluster ](github.com/Terry-Mao/gopush-cluster )

以太坊整套协议钱包Go实现  :  [github.com/ethereum/go-ethereum ](github.com/ethereum/go-ethereum )


## 开发套件
微服务套件  :  [github.com/go-kit/kit](github.com/go-kit/kit)

桌面UI套件(基于CGO)  :  [github.com/andlabs/ui](github.com/andlabs/ui)

桌面UI库(基于HTML)  :  [github.com/murlokswarm/app](github.com/murlokswarm/app)

LOG库(已使用)  :  [github.com/Sirupsen/logrus ](github.com/Sirupsen/logrus )

zap - LOG库(集成gocore)  :  [github.com/uber-go/zap](github.com/uber-go/zap)

图像处理库  :  [github.com/anthonynsimon/bild ](github.com/anthonynsimon/bild )

图像处理库  :  [github.com/disintegration/imaging](github.com/disintegration/imaging)

日期处理库  :  [github.com/jinzhu/now ](github.com/jinzhu/now )

配置文件读取库(已使用)  :  [github.com/spf13/viper ](github.com/spf13/viper )

类型转换库(已使用)  :  [github.com/spf13/cast ](github.com/spf13/cast )

UUID库(已使用)  :  [github.com/satori/go.uuid](github.com/satori/go.uuid)

压缩文件处理库  :  [github.com/mholt/archiver ](github.com/mholt/archiver )

连接池库(已使用)  :  [github.com/jolestar/go-commons-pool ](github.com/jolestar/go-commons-pool )

程序内部系统资源,可以对不同的资源做出不同的规则调整  :  [github.com/shirou/gopsutil ](github.com/shirou/gopsutil )


## 数据文件处理

文件嵌入到编译文件  :  [github.com/rakyll/statik ](github.com/rakyll/statik )

文件嵌入到编译文件(html,css,js)  :  [github.com/GeertJohan/go.rice](github.com/GeertJohan/go.rice)

内存敏感数据处理  :  [github.com/awnumar/memguard ](github.com/awnumar/memguard )


## 第三方软件使用

邮件发送(已使用)  :  [github.com/go-gomail/gomail](github.com/go-gomail/gomail)

数据库操作(已使用)  :  [github.com/jinzhu/gorm](github.com/jinzhu/gorm)

数据库操作  :  [github.com/go-xorm/xorm ](github.com/go-xorm/xorm )

redis操作库(已使用)  :  [gopkg.in/redis.v5](gopkg.in/redis.v5)

rabbitmq使用框架(已使用)  :  [github.com/streadway/amqp](github.com/streadway/amqp)

levelDB处理  :  [https://github.com/syndtr/goleveldb ](https://github.com/syndtr/goleveldb )

bolt嵌入式数据库  :  [github.com/boltdb/bolt](github.com/boltdb/bolt)



### 解析库
JSON解析库(已使用)  :  [github.com/tidwall/gjson](github.com/tidwall/gjson)

CSV处理库  :  [github.com/jszwec/csvutil](github.com/jszwec/csvutil)

msgpack binc  cbor json   解密库  :  [github.com/ugorji/go](github.com/ugorji/go)

golang解密php序列化库  :  [github.com/yvasiyarov/php_session_decoder ](github.com/yvasiyarov/php_session_decoder )

高性能json库  :  [github.com/json-iterator/go](github.com/json-iterator/go)

google-protobuf库  :  [github.com/golang/protobuf ](github.com/golang/protobuf )


## 网络框架
http网路框架(已使用)  :  [github.com/labstack/echo](github.com/labstack/echo)

http网路框架(已使用)  :  [github.com/gin-gonic/gin](github.com/gin-gonic/gin)

http网络框架  :  [https://github.com/go-martini/martini ](https://github.com/go-martini/martini )

超级快的 http 网路框架(已使用)  :  [github.com/valyala/fasthttp ](github.com/valyala/fasthttp )

KCP协议golang实现(已使用)  :  [github.com/xtaci/kcp-go ](github.com/xtaci/kcp-go )

IOT库 支持各种协议  :  [github.com/hybridgroup/gobot ](github.com/hybridgroup/gobot )

socket.io协议Go实现(已使用)  :  [github.com/googollee/go-socket.io](github.com/googollee/go-socket.io)


## 深度学习
Go语言实现机器学习框架  :  [github.com/sjwhitworth/golearn ](github.com/sjwhitworth/golearn )

GO机器学习图书馆,包含各种各样的算法  :  [github.com/gorgonia/gorgonia](github.com/gorgonia/gorgonia)

go语言对Tensorflow的封装  :  [github.com/galeone/tfgo ](github.com/galeone/tfgo )


## 依赖管理
官方包管理  :  [github.com/golang/dep ](github.com/golang/dep )

包管理工具(本地打包到项目)  :  [github.com/tools/godep ](github.com/tools/godep )

包管理工具(类似Composer 和 pip)  :  [github.com/Masterminds/glide ](github.com/Masterminds/glide )
