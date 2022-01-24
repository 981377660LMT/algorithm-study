1. Blob 表示的不一定是 JavaScript 原生格式的数据。比如 File 接口基于 Blob，继承了 blob 的功能并将其扩展使其支持用户系统上的文件。
2. Blob 由一个可选的字符串 type（通常是 MIME 类型）和 blobParts 组成
3. Blob 对象是不可改变的。我们不能直接在一个 Blob 中更改数据，但是我们可以对一个 Blob 进行分割，从其中创建新的 Blob 对象，将它们混合到一个新的 Blob 中。这种行为类似于 JavaScript 字符串：我们无法更改字符串中的字符，但可以创建新的更正后的字符串。
4. Blob URL/Object URL
   Blob URL/Object URL 是一种伪协议，允许 Blob 和 File 对象用作图像，下载二进制数据链接等的 URL 源。
   URL.createObjectURL 方法来创建 Blob URL，该方法接收一个 Blob 对象，并为其创建一个唯一的 URL，其形式为 blob:<origin>/<uuid>，对应的示例如下：
   blob:https://example.org/40a5fb5a-d56d-4a33-b4e2-0acf6a8e5f641
   浏览器内部为每个通过 URL.createObjectURL 生成的 URL `存储了一个 URL → Blob 映射`
   生成的 URL 仅在当前文档打开的状态下才有效
   因此，如果我们创建一个 Blob URL，即使不再需要该 Blob，它也会存在内存中。
   针对这个问题，我们可以调用 `URL.revokeObjectURL(url) 方法，从内部映射中删除引用`，从而允许删除 Blob（如果没有其他引用），并释放内存。
5. base64
   在传输编码方式中指定 base64。使用的字符包括大小写拉丁字母各 26 个、数字 10 个、加号 + 和斜杠 /，共 64 个字符，等号 = 用来作为后缀用途。
   绝大多数现代浏览器都支持一种名为 Data URLs 的特性，允许使用 base64 对图片或其他文件的二进制数据进行编码，将其作为文本字符串嵌入网页中。`逗号分隔`
   data:[<mediatype>][;base64],<data>

```HTML
<img alt="logo" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUg...">
```

如果数据是文本类型，你可以直接将文本嵌入（根据文档类型，使用合适的实体字符或转义字符）。如果是二进制数据，你可以将数据进行 base64 编码之后再进行嵌入。比如嵌入一张图片。
利用 FileReader API，我们也可以方便的实现图片本地预览功能(FileReader 对象的 readAsDataURL() 方法，把本地图片对应的 File 对象转换为 Data URL。)

6. Blob 与 ArrayBuffer 的区别
   ArrayBuffer 对象用于表示通用的，固定长度的原始二进制数据缓冲区。你不能直接操纵 ArrayBuffer 的内容，而是需要创建一个类型化数组对象或 DataView 对象，该对象以特定格式表示缓冲区，并使用该对象读取和写入缓冲区的内容。
   Blob 类型的对象表示不可变的类似文件对象的原始数据。Blob 表示的不一定是 JavaScript 原生格式的数据。File 接口基于 Blob，继承了 Blob 功能并将其扩展为支持用户系统上的文件。

Blob 与 ArrayBuffer 对象之间是可以相互转化的：

使用 `FileReader 的 readAsArrayBuffer()` 方法，可以把 Blob 对象转换为 ArrayBuffer 对象；
使用 Blob 构造函数，如 new Blob([new Uint8Array(data]);，可以把 ArrayBuffer 对象转换为 Blob 对象。

ArrayBuffer 是存在内存中的，可以直接操作。而 Blob 可以位于磁盘、高速缓存内存和其他不可用的位置。
比如 MySQL 或 Oracle 数据库中，blob 类型是可以直接存储二进制的，常见是图片，这些数据最终会存储在硬盘中。
