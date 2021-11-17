// JavaScript的内部编码方式:UTF-16编码,每个字符占用2个字节
function ab2str(buf: ArrayBuffer) {
  return String.fromCharCode.apply(null, new Uint16Array(buf) as any)
}

// 字符串转 ArrayBuffer对象
function str2ab(str: string) {
  let buf = new ArrayBuffer(str.length * 2) // 每个字符占用2个字节
  let bufView = new Uint16Array(buf)
  for (let i = 0; i < str.length; i++) {
    bufView[i] = str.charCodeAt(i)
  }
  return buf
}

// XHR 下载图片
const xhr = new XMLHttpRequest()
xhr.open('GET', 'https://avatars2.githubusercontent.com/u/4220799?v=3')
xhr.responseType = 'blob'

xhr.onload = function () {
  if (this.status == 200) {
    const img = document.createElement('img')
    img.src = window.URL.createObjectURL(this.response)
    img.onload = function () {
      window.URL.revokeObjectURL(img.src)
    }
    document.body.appendChild(img)
  }
}
xhr.send()

// XHR 上传数据
// var xhr = new XMLHttpRequest();
// xhr.open('POST','/upload');
// xhr.onload = function() { ... };
// xhr.send("text string");

// 发送FormData
// var formData = new FormData();
// formData.append('id', 123456);
// formData.append('topic', 'performance');

// var xhr = new XMLHttpRequest();
// xhr.open('POST', '/upload');
// xhr.onload = function() { ... };
// xhr.send(formData);

// 发送 Buffer
// var xhr = new XMLHttpRequest();
// xhr.open('POST', '/upload');
// xhr.onload = function() { ... };
// var uInt8Array = new Uint8Array([1, 2, 3]);
// xhr.send(uInt8Array.buffer);

// XHR 上传进度条
xhr.upload.onprogress = function (event) {
  if (event.lengthComputable) {
    const complete = ((event.loaded / event.total) * 100) | 0
    const progress = document.getElementById('uploadprogress') as HTMLProgressElement
    progress.value = complete
    progress.innerHTML = complete.toString()
  }
}
// 注意，progress事件不是定义在xhr，而是定义在xhr.upload，因为这里需要区分下载和上传，下载也有一个progress事件。

// sendAsBinary() polyfill
if (!XMLHttpRequest.prototype.sendAsBinary) {
  XMLHttpRequest.prototype.sendAsBinary = function (sData: string) {
    var nBytes = sData.length,
      ui8Data = new Uint8Array(nBytes)
    for (var nIdx = 0; nIdx < nBytes; nIdx++) {
      ui8Data[nIdx] = sData.charCodeAt(nIdx) & 0xff
    }
    this.send(ui8Data)
  }
}

// XHR 还支持大文件分块传输：每次通过 XHR 上传 1MB 数据块
const blob =new Blob() // 1

const BYTES_PER_CHUNK = 1024 * 1024; // 2
const SIZE = blob.size;

let start = 0;
let end_ = BYTES_PER_CHUNK;

while(start < SIZE) { // 3
  const xhr = new XMLHttpRequest();
  xhr.open('POST', '/upload');
  xhr.onload = function() { ... };

  xhr.setRequestHeader('Content-Range', start+'-'+end+'/'+SIZE); // 4
  xhr.send(blob.slice(start, end)); // 5

  start = end;
  end_ = start + BYTES_PER_CHUNK;
}


// 获取 XMLHttpRequest 响应体
function readBody(xhr:XMLHttpRequest) {
  let data:unknown;
  if (!xhr.responseType || xhr.responseType === "text") {
      data = xhr.responseText;
  } else if (xhr.responseType === "document") {
      data = xhr.responseXML;
  } else {
      data = xhr.response;
  }
  return data;
}
// 验证请求是否成功
export const isSuccess = (status: number): boolean => (status >= 200 && status < 300);

