HTTP 是一个在计算机世界里专门在「两点」之间「传输」文字、图片、音频、视频等「超文本」数据的「约定和规范」。
那「HTTP 是用于从互联网服务器传输超文本到本地浏览器的协议，这种说法正确吗?
这种说法是不正确的。因为也可以是「服务器<-->服务器」，所以采用两点之间的描述会更准确。

1. 读者问答
   读者问:“https 和 http 相比，就是传输的内容多了对称加密，可以这么理解吗?”

   1. 建立连接时候:https 比 http 多了 TLS 的握手过程;
   2. 传输内容的时候:https 会把数据进行加密，通常是对称加密数据;

   读者问:“我看文中 TLS 和 SSL 没有做区分，这两个需要区分吗?”
   **这两实际上是一个东西。**
   SSL 是洋文“Secure Sockets Layer 的缩写，中文叫做「安全套接层」。它是在上世纪 90 年代中期，由网景公司设计的。
   到了 1999 年，SSL 因为应用广泛，已经成为互联网上的事实标准。IETF 就在那年**把 SSL 标准化**。标准化之后的名称改为 TLS(是“Transport Layer Security”的缩写)，中文叫做「传输层安全协议」。
   很多相关的文章都把这两者并列称呼(SSLTLS)，**因为这两者可以视作同一个东西的不同阶段。**

2. http1.1 如何优化

   问你一句:「你知道 HTTP/1.1 该如何优化吗?J
   我想你第一时间想到的是，使用 KeepAlive 将 HITP/1.1 从短连接改成长链接。
   这个确实是一个优化的手段，它是从底层的传输层这一方向入手的，通过减少 TCP 连接建立和断开的次数，来减少了网络传输的延迟，从而提高 HIIP/1.1 协议的传输效率。
   但其实还可以从其他方向来优化 HTTP/1.1 协议，比如有如下 3 种优化思路

   1. 尽量避免发送 HTTP 请求;(缓存)
      缓存真的是性能优化的一把万能钥匙，小到 CPU Cache、Page Cache、Redis Cache，大到 HTTP 协议的缓存。

   2. 在需要发送 HTTP 请求时，考虑如何减少请求次数;(合并请求)
      通过将多个小图片合并成一个大图片来减少 HTTP 请求的次数，以减少 HTTP 请求的次数，从而减少网络的开销。

   3. 减少服务器的 HTTP 响应的数据大小;(压缩)
      无损压缩：gzip/brotli
      有损压缩:`Accept:audio/*;q=0.2,audio/basic `
      q 代表质量因子  
      例如 webP 格式

总结
这次主要从 3 个方面介绍了优化 HTTP/1.1 协议的思路。
第一个思路是，通过缓存技术来避免发送 HTTP 请求。客户端收到第一个请求的响应后，可以将其缓存在本地磁盘，下次请求的时候，如果缓存没过期，就直接读取本地缓存的响应数据。如果缓存过期，客户端发送请求的时候带上响应数据的摘要，服务器比对后发现资源没有变化，就发出不带包体的 304 响应，告诉客户端缓存的响应仍然有效。
第二个思路是，减少 HTIP 请求的次数，有以下的方法:
1．将原本由客户端处理的重定向请求，交给代理服务器处理，这样可以减少重定向请求的次数; 2.将多个小资源合并成一个大资源再传输，能够减少 HTTP 请求次数以及头部的重复传输，再来减少 TCP 连接数量，进而省去 TCP 握手和慢启动的网络消耗; 3.按需访问资源，只访问当前用户看得到/用得到的资源，当客户往下滑动，再访问接下来的资源，以此达到延迟请求，也就减少了同一时间的 HTIP 请求次数。
第三思路是，通过压缩响应资源，降低传输资源的大小，从而提高传输效率，所以应当选择更优秀的压缩算法。
不管怎么优化 HTTP/1.1 协议都是有限的，不然也不会出现 HTTP/2 和 HTTP/3 协议，后续我们再来介绍 HTTP/2 和 HITP/3 协议。

3. HTTP 报文结构是怎样的？
   起始行
   头部
   空行
   实体
4. 对 Accept 系列字段了解多少？
   对于 Accept 系列字段的介绍分为四个部分: 数据格式、压缩方式、支持语言和字符集。

   1. HTTP 从 **MIME type** 取了一部分来标记报文 body 部分的数据类型，这些类型体现在 Content-Type 这个字段，当然这是针对于发送端而言，接收端想要收到特定类型的数据，也可以用 Accept 字段。
      发送端:Content-Type
      接收端:Accept
      具体而言，这两个字段的取值可以分为下面几类:

      text： text/html, text/plain, text/css 等
      image: image/gif, image/jpeg, image/png 等
      audio/video: audio/mpeg, video/mp4 等
      application: application/json, application/javascript, application/pdf, application/octet-stream

   2. 采取什么样的压缩方式就体现在了发送方的 Content-Encoding 字段上， 同样的，接收什么样的压缩方式体现在了接受方的 Accept-Encoding 字段上
   3. 对于发送方而言，还有一个 Content-Language 字段，在需要实现国际化的方案当中，可以用来指定支持的语言，在接受方对应的字段为 Accept-Language。
   4. 在接收端对应为 `Accept-Charset`，指定可以接受的字符集，在发送端并没有对应的 Content-Charset, 而是直接放在了 `Content-Type` 中,以 charset 属性指定  
      // 发送端
      Content-Type: text/html; charset=utf-8
      // 接收端
      Accept-Charset: charset=utf-8

5. 对于定长和不定长的数据，HTTP 是怎么传输的？
   定长包体
   对于定长包体而言，发送端在传输的时候一般会带上 Content-Length, 来指明包体的长度。

   ```JS
      const http = require('http');

      const server = http.createServer();

      server.on('request', (req, res) => {
         if(req.url === '/') {
            res.setHeader('Content-Type', 'text/plain');
            res.setHeader('Content-Length', 10);
            res.write("helloworld");
         }
      })

      server.listen(8081, () => {
         console.log("成功启动");
      })

   ```

   浏览器中显示如下:
   helloworld

   `res.setHeader('Content-Length', 8);`之后
   **hellowor**
   后面的 ld 哪里去了呢？实际上在 http 的响应体中直接被截去了。

   然后我们试着将这个长度设置得大一些:
   `res.setHeader('Content-Length', 12);`
   **直接无法显示了。**
   可以看到 Content-Length 对于 http 传输过程起到了十分关键的作用，如果设置不当可以直接导致传输失败。

   不定长包体
   另外一个 http 头部字段了:
   `Transfer-Encoding: chunked`
   表示分块传输数据，设置这个字段后会`自动产生两个效果`:
   Content-Length 字段会被忽略
   基于长连接持续推送动态内容

6. HTTP 如何处理大文件的传输？
   对于几百 M 甚至上 G 的大文件来说，如果要一口气全部传输过来显然是不现实的，会有大量的等待时间，严重影响用户体验。因此，HTTP 针对这一场景，采取了`范围请求`的解决方案，允许客户端仅仅请求一个资源的一部分。

   当然，前提是服务器要支持范围请求，要支持这个功能，就必须加上这样一个响应头:
   `Accept-Ranges: none`用来告知这边是支持范围请求的。

   而对于客户端而言，它需要指定请求哪一部分，通过 Range 这个请求头字段确定，格式为`bytes=x-y`
   服务器收到请求之后，首先验证范围是否合法，如果越界了那么返回 416 错误码，否则读取相应片段，返回 206 状态码。
   同时，服务器`需要添加 Content-Range 字段`，这个字段的格式根据请求头中 Range 字段的不同而有所差异。

7. 对 websocket 的了解
   了解 WebSocket 的诞生背景：实现推送技术、WebSocket 是什么及它的优点；
   了解 WebSocket 含有哪些 API 及如何使用 WebSocket API 发送普通文本和二进制数据；
   了解 WebSocket 的握手协议和数据帧格式、掩码算法等相关知识；
   了解如何实现一个支持发送普通文本的 WebSocket 服务器。

8. WebSocket 优点
   - 在 WebSocket API 中，浏览器和服务器`只需要完成一次握手`，两者之间就可以创建持久性的连接，并进行双向数据传输。
   - 更强的实时性。由于协议是全双工的，所以服务器可以随时主动给客户端下发数据。相对于 HTTP 请求需要等待客户端发起请求服务端才能响应，延迟明显更少。
9. websocket 服务端
   1. 握手协议
      `WebSocket 通过 HTTP/1.1 协议的 101` 状态码进行握手。
      GET ws://echo.websocket.org/ HTTP/1.1
      Host: echo.websocket.org
      Origin: file:// 浏览器中发起此 WebSocket 连接所在的页面
      Connection: Upgrade
      Upgrade: websocket  
      Sec-WebSocket-Version: 13
      Sec-WebSocket-Key: Zx8rNEkBE4xnwifpuh8DHQ==
      Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
   2. 消息通信基础
      数据是通过一系列数据帧来进行传输的
10. WebSocket 与 HTTP 有什么关系
    WebSocket 是一种与 HTTP 不同的协议。两者都位于 OSI 模型的应用层，并且都依赖于传输层的 TCP 协议。 虽然它们不同，但是 RFC 6455 中规定：WebSocket 被设计为在 HTTP 80 和 443 端口上工作，并支持 HTTP 代理和中介，从而使其与 HTTP 协议兼容。 为了实现兼容性，WebSocket 握手使用 HTTP Upgrade 头，从 HTTP 协议更改为 WebSocket 协议。
11. 什么是 WebSocket 心跳
    心跳包就是客户端定时发送简单的信息给服务器端告诉它我还在而已
    在 WebSocket 协议中定义了 心跳 Ping 和 心跳 Pong 的控制帧
12. Socket 是什么
    ip+port
    在 Internet 上的主机一般运行了多个服务软件，同时提供几种服务。每种服务都打开一个 Socket，并绑定到一个端口上，不同的端口对应于不同的服务。
    `Socket 正如其英文原义那样，像一个多孔插座`。一台主机犹如布满各种插座的房间，每个插座有一个编号，有的插座提供 220 伏交流电， 有的提供 110 伏交流电，有的则提供有线电视节目。 客户软件将插头插到不同编号的插座，就可以得到不同的服务。
