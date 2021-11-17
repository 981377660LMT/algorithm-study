1. 为什么 GET 或 HEAD 请求，不能通过 send() 方法发送请求体？
   通过 XMLHttpRequest 规范，我们知道当请求方法是 GET 或 HEAD 时，send() 方法的 body 参数值将会被忽略。
   那么对于我们常用的 GET 请求，我们要怎么传递参数呢？解决参数传递可以使用以下两种方式：
   `URL 传参` - 常用方式，有大小限制大约为 2KB
   `请求头传参` - 一般用于传递 token 等认证信息
2. 什么情况下 XMLHttpRequest status 会为 0？
   如果状态是 UNSENT 或 OPENED，则返回 0
   如果错误标志被设置，则返回 0 (触发错误时)
   否则返回 HTTP 状态码
3. XMLHttpRequest 请求体支持哪些格式？
   浏览器可以为各种本地数据类型提供自动编码和解码，这样可以让应用程序将这些类型直接传递给 XHR
   `ABDDF:请求格式`
   `ABDJT:响应格式` JSON Text
   void send();

   void send(ArrayBuffer data);

   void send(Blob data);

   void send(Document data);

   void send(DOMString? data);

   void send(FormData data);

   POST 请求示例

   发送 POST 请求通常需要以下四个步骤：

   1. 使用 open() 方法打开连接时，设定 POST 请求方法和请求 URL 地址
   2. `设置正确的 Content-Type 请求头`
   3. 设置相关的事件监听
   4. 设置请求体，并使用 send() 方法，发送请求

```JS
var xhr = new XMLHttpRequest();
xhr.open("POST", '/server', true);

//Send the proper header information along with the request
xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

xhr.onreadystatechange = function() {
    if(xhr.readyState == 4 && xhr.status == 200) {
        // handle data
    }
}
xhr.send("foo=bar&lorem=ipsum");
// xhr.send('string');
// xhr.send(new Blob());
// xhr.send(new Int8Array());
// xhr.send({ form: 'data' });
// xhr.send(document);
```

4. 什么是简单请求和预请求 (preflight request) ？
   一些不会触发 CORS preflight 的请求被称为 "简单请求"
   只允许下列方法：
   `GET`
   `HEAD`
   `POST`

   除了用户代理自动设置的头部外（比如 Connection， User-Agent ，或者其他任意的 Fetch 规范定义的 禁止的头部名 ），唯一允许人工设置的头部是 Fetch 规范定义的 CORS-safelisted request-header，如下：
   Accept
   Accept-Language
   Content-Language
   Content-Type

   允许的 Content-Type 值有：
   application/x-www-form-urlencoded
   multipart/form-data
   text/plain

   **预请求**

   不同于上面讨论的简单请求，"预请求" 要求必须先发送一个 OPTIONS 方法请求给目的站点，来查明这个跨站请求对于目的站点是不是安全的可接受的。

5. 怎样防止重复发送 AJAX 请求？
   `setTimeout + clearTimeout` - 连续的点击会把上一次点击清除掉，也就是 ajax 请求会在最后一次点击后发出去
   `disable 按钮`
   `缓存已成功的请求`，若请求参数一致，则直接返回，不发送请求
