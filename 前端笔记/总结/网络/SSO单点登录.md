1. SSO 英文全称 Single Sign On，单点登录。SSO 是在多个应用系统中，用户只需要登录一次就可以访问所有相互信任的应用系统。
   例如访问在网易账号中心（https://reg.163.com/ ）登录之后 访问以下站点都是登录状态

   网易直播 https://v.163.com
   网易博客 https://blog.163.com
   网易花田 https://love.163.com
   网易考拉 https://www.kaola.com
   网易 Lofter http://www.lofter.com

2. 用户登录状态的存储与校验
   用户登录成功之后，生成 AuthToken 交给客户端保存。如果是浏览器，就保存在 Localstoarge 中(`Cookie 不能跨域`)。如果是手机 App 就保存在 App 本地缓存中。本篇主要探讨基于 Web 站点的 SSO。 用户在浏览需要登录的页面时，客户端将 AuthToken 提交给 SSO 服务校验登录状态/获取用户登录信息
   `对于登录信息的存储，建议采用 Redis`，使用 Redis 集群来存储登录信息，既可以保证高可用，又可以线性扩充。同时也可以让 SSO 服务满足负载均衡/可伸缩的需求。
   AuthToken 直接使用 UUID/GUID 即可，如果有验证 AuthToken 合法性需求，可以将 `UserName+时间戳加密生成`，服务端解密之后验证合法性
   登录信息 通常是将 UserId，UserName 缓存起来

3. 用户登出
   1. 服务端清除缓存（Redis）中的登录状态
   2. 客户端清除存储的 AuthToken
