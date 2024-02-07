## 安全最佳实践


### 网络链路安全

基础安全要求：
- 仅支持TLS > 1.2
- 不支持 HTTP 必须使用 HTTPS 方式进行访问
- 开启 WAF 防火墙防止 XSS 等攻击
- HTTPS 启用严格传输安全（HSTS） 响应头
- 禁用 CBC加密，开启 CTR 和 GCM 加密
- 响应报文中“X-Content-Type-Options”头配置为“nosniff“
- 响应报文设置 X-XSS-Protection: 1; mode=block
- 对称加密使用相对安全的协议 AES256
- 签名算法使用相对安全的协议 Hmac256
- 生产环境模糊返回 msg
- 请求端必须进行证书透明校验
- 隐藏返回涉及版本信息的 header 头


