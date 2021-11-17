https://juejin.cn/post/6850037275579121671

通过 Chrome 开发者工具，我们可以知道当播放 「xgplayer-demo-720p.mp4」 视频文件时，发了 3 个 HTTP 请求：
此外，从图中可以清楚地看到，`头两个 HTTP 请求响应的状态码是 「206」`。这里我们来分析第一个 HTTP 请求的请求头和响应头：
在上面的请求头中，`有一个 range: bytes=0- 首部信息`，该信息用于检测服务端是否支持 Range 请求。`如果在响应中存在 Accept-Ranges 首部（并且它的值不为 “none”）`，那么表示该服务器支持范围请求。
![range请求](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/7/15/1734ff41d091fd45~tplv-t2oaga2asx-watermark.awebp)

对于使用 REST Client 发起的 「单一范围请求」，服务器端会返回状态码为 「206 Partial Content」 的响应。而响应头中的 「Content-Length」 首部现在用来表示先前请求范围的大小（而不是整个文件的大小）。「Content-Range」 响应首部则表示这一部分内容在整个资源中所处的位置。

1. 范围请求的响应
   与范围请求相关的有三种状态：

   1. 在请求成功的情况下，服务器会返回 「206 Partial Content」 状态码。
   2. 在请求的范围越界的情况下（范围值超过了资源的大小），服务器会返回 「416 Requested Range Not Satisfiable」 （请求的范围无法满足） 状态码。
   3. 在不支持范围请求的情况下，服务器会返回 「200 OK」 状态码。

若播放的视频文件太大或出现网络不稳定，`则会导致播放时`，需要等待较长的时间，这严重降低了用户体验。

那么如何解决这个问题呢？要解决该问题我们可以使用流媒体技术，接下来我们来介绍流媒体。

2. 流媒体
   当使用 HLS 流媒体网络传输协议时，<video> 元素 src 属性使用的是 blob:// 协议。
   HTTP Live Streaming（缩写是 HLS）是由苹果公司提出基于 HTTP 的流媒体网络传输协议。它的工作原理是把整个流分成一个个小的基于 HTTP 的文件来下载，每次只下载一些。
   HLS 是一种`自适应比特率流协议`。因此，HLS 流可以动态地使视频分辨率自适应每个人的网络状况。
3. DASH
   介绍完苹果公司推出的 HLS （HTTP Live Streaming）技术，接下来我们来介绍另一种基于 HTTP 的动态自适应流 —— DASH。
   「`基于 HTTP 的动态自适应流（英语：Dynamic Adaptive Streaming over HTTP，缩写 DASH，也称 MPEG-DASH）是一种自适应比特率流技术，使高质量流媒体可以通过传统的 HTTP 网络服务器以互联网传递。`」
   在国内 Bilibili 于 2018 年开始使用 DASH 技术，
4. FLV
   FLV 是 FLASH Video 的简称，FLV 流媒体格式是随着 Flash MX 的推出发展而来的视频格式。
   在浏览器中 HTML5 的 <video> 是不支持直接播放 FLV 视频格式，需要借助 flv.js 这个开源库来实现播放 FLV 视频格式的功能。
5. flv.js
   前面我们已经提到了 Bilibili，接下来不得不提其开源的一个著名的开源项目 —— flv.js
   flv.js 是用纯 JavaScript 编写的 HTML5 Flash Video（FLV）播放器，它底层依赖于 Media Source Extensions。在实际运行过程中，它会自动解析 FLV 格式文件并喂给原生 HTML5 Video 标签播放音视频数据，使浏览器在不借助 Flash 的情况下播放 FLV 成为可能。
6. 如何实现视频本地预览
   视频本地预览的功能主要利用 URL.createObjectURL() 方法来实现
7. 如何实现播放器截图
   播放器截图功能主要利用 `CanvasRenderingContext2D.drawImage()` +`canvas.toDataURL()` API 来实现。Canvas 2D API 中的 CanvasRenderingContext2D.drawImage() 方法提供了多种方式在 Canvas 上绘制图像。
   其中 image 参数表示绘制到上下文的元素。允许任何的 canvas 图像源（CanvasImageSource），例如：CSSImageValue，HTMLImageElement，SVGImageElement，HTMLVideoElement，HTMLCanvasElement，ImageBitmap 或者 OffscreenCanvas。
8. 如何实现 Canvas 播放视频
   使用 Canvas 播放视频主要是利用 `ctx.drawImage(video, x, y, width, height)` 来对视频当前帧的图像进行绘制，其中 video 参数就是页面中的 video 对象。所以如果我们按照`特定的频率不断获取 video 当前画面，并渲染到 Canvas 画布上`，就可以实现使用 Canvas 播放视频的功能。
9. 如何实现色度键控（绿屏效果）
   这是因为 Canvas 提供了 `getImageData` 和 `putImageData `方法使得开发者可以动态地更改每一帧图像的显示内容

10. B 站蒙版弹幕的实现原理
    https://juejin.cn/post/6844903766148284423
    基于用户数据和一些机器学习的相关应用，可以提炼出视频的关键主体
    服务端预先对视频进行处理，并生成相应的蒙版数据
    客户端播放视频时，实时地加载对应资源
    通过一些前端的技术手段，实现弹幕的蒙版处理

    客户端方面，由于 B 站弹幕是基于 div+css 的实现，因而采用了 svg 格式来传输矢量蒙版（至少目前是这样），通过 CSS 遮罩的方式实现渲染。
