更新日志
---

[comment]: <> (### 0.2.0)

[comment]: <> (> 发布时间：2021-10-17)

[comment]: <> (* Kubernetes最低兼容版本提高到`1.16`)

[comment]: <> (* 使用shadow pod代替shadow deployment)

[comment]: <> (* Windows的`socks`模式默认不再自动设置全局代理，新增开启该功能的`--setupGlobalProxy`参数)

[comment]: <> (* 新增`exchange`命令的`ephemeral`模式（for k8s 1.23+，感谢@[xyz-li]&#40;https://github.com/xyz-li&#41;）)

[comment]: <> (* 修复`exchange`命令连接时常卡顿的问题（issues #184，感谢@[xyz-li]&#40;https://github.com/xyz-li&#41;）)

[comment]: <> (* 当Port-forward的目标端口被占用时提供更优雅的报错信息（感谢@[xyz-li]&#40;https://github.com/xyz-li&#41;）)

[comment]: <> (* 自动根据用户权限控制生成路由的范围，去除Connect命令的`--global`参数)

[comment]: <> (* 优化Connect命令的`--cidr`参数，支持指定多个IP区段)

[comment]: <> (* 参数`--label`更名为`--withLabel`)

[comment]: <> (* 增加`--withAnnotation`参数为shadow pod增加额外标注)

[comment]: <> (* `connect`命令增加`--disablePodIp`参数支持禁用Pod IP路由)

[comment]: <> (* shadow pod增加`kt-user`标注用于记录本地用户名)

[comment]: <> (* 移除`check`命令)
