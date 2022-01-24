消息系统中握手的作用及如何实现握手；
消息模型的设计及如何实现消息验证来保证通信安全；
postMessage 的使用及如何利用它实现父子页面的消息通信；
消息通信 API 的设计与实现。

Postmate 是一个强大，简单，基于 Promise 的 postMessage 库。它允许父页面以最小的成本与跨域的子 iframe 进行通信。

1. 需要实现父页面与 iframe 加载的子页面之间的消息通信

```JS
otherWindow.postMessage(message, targetOrigin, [transfer]);
otherWindow：其他窗口的一个引用，比如 iframe 的 contentWindow 属性、执行 window.open 返回的窗口对象等。
message：将要发送到其他 window 的数据，它将会被结构化克隆算法序列化。
targetOrigin：通过窗口的 origin 属性来指定哪些窗口能接收到消息事件，其值可以是字符串 "*"（表示无限制）或者一个 URI。
transfer（可选）：是一串和 message 同时传递的 Transferable 对象。这些对象的所有权将被转移给消息的接收方，而发送一方将不再保有所有权。

发送方通过 postMessage API 来发送消息，而接收方可以通过监听 message 事件，来添加消息处理回调函数，具体使用方式如下：
window.addEventListener("message", receiveMessage, false);

function receiveMessage(event) {
  let origin = event.origin || event.originalEvent.origin;
  if (origin !== "http://semlinker.com") return;
}

```

2. Postmate 握手的实现:在主应用和子应用双方完成握手之后，就可以进行双向消息通信了
