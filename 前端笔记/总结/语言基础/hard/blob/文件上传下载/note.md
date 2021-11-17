1. JavaScript 中如何实现并发控制？
   async-pool
2. JavaScript 中如何实现大文件并行下载？
   HTTP 范围请求
   如果在响应中存在 Accept-Ranges 首部（并且它的值不为 “none”），那么表示该服务器支持范围请求。
   在一个 Range 首部中，可以一次性请求多个部分，服务器会以 multipart 文件的形式将其返回。如果服务器返回的是范围响应，需要使用 206 Partial Content 状态码。假如所请求的范围不合法，那么服务器会返回 416 Range Not Satisfiable 状态码，表示客户端错误。服务器允许忽略 Range 首部，从而返回整个文件，状态码用 200 。

   Range 语法
   Range: <unit>=<range-start>-
   Range: <unit>=<range-start>-<range-end>
   Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>
   Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>, <range-start>-<range-end>
