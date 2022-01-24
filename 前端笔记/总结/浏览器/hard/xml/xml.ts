const req = new XMLHttpRequest()

req.responseType
// type XMLHttpRequestResponseType = "" | "arraybuffer" | "blob" | "document" | "json" | "text";

req.responseText

req.timeout

req.withCredentials

req.abort()

req.open
// void open(
//   DOMString method,
//   DOMString url,
//   optional boolean async,
//   optional DOMString user,
//   optional DOMString password
// );
// async - 一个可选的布尔值参数，默认值为 true，表示执行异步操作。如果值为 false，则 send() 方法不会返回任何东西，直到接收到了服务器的返回数据

req.send
// void send();
// void send(ArrayBuffer data);
// void send(Blob data);
// void send(Document data);
// void send(DOMString? data);
// void send(FormData data);

req.setRequestHeader

req.onreadystatechange
