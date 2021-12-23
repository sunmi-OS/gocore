## 安全最佳实践


### 网络链路安全
基础安全要求：
- 仅支持TLS > 1.2
- 不支持HTTP必须使用HTTPS方式进行访问
- 开启WAF防火墙防止XSS等攻击
- HTTPS启用严格传输安全（HSTS） 响应头
- 禁用CBC加密，开启CTR和GCM加密
- 响应报文中“X-Content-Type-Options”头配置为“nosniff“
- 响应报文设置X-XSS-Protection: 1; mode=block
- 对称加密使用相对安全的协议 AES256
- 签名算法使用相对安全的协议 Hmac256
- 生产环境模糊返回msg
- 请求端必须进行证书透明校验
- 隐藏返回涉及版本信息的header头


