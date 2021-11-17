1. ZIP 通常使用后缀名 “.zip”，它的 MIME 格式为 “application/zip”。
2. JSZip 这个库的浏览器解压方案。
   JsZip 又是怎么实现浏览器端解压的
   在浏览器端通过 FileReader API 就可以获取文件的二进制数据，剩下的就是按照文件的格式来解析。
3. JavaScript 如何检测文件的类型？
   我们希望能限制文件上传的类型，比如限制只能上传 PNG 格式的图片。针对这个问题，我们会想到通过 input 元素的 accept 属性来限制上传的文件类型：

```HTML
<input type="file" id="inputFile" accept="image/png" />
```

这种方案虽然可以满足大多数场景，但如果用户把 JPEG 格式的图片后缀名更改为 .png 的话，就可以成功突破这个限制。那么应该如何解决这个问题呢？
其实我们可以通过读取`文件的二进制数据来识别正确的文件类型`。
