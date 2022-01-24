https://juejin.cn/post/6898098763772985352

1. FileReader API:选择本地图片 -> 图片预览
   `fileReader.readAsDataURL(event.target.files[0]);`
2. Fetch:网络下载图片 -> 图片预览
   `URL.createObjectURL(blob)`
3. DataView 与 ArrayBuffer
   DataView 视图是一个可以从二进制 ArrayBuffer 对象中读写多种数值类型的底层接口，
   使用它时，不用考虑不同平台的`字节序`问题。
   字节顺序，又称端序或尾序（英语：Endianness），在计算机科学领域中，指存储器中或在数字通信链路中，组成多字节的字的字节的排列顺序。

   字节的排列方式有两个通用规则。例如，一个多位的整数，按照存储地址从低到高排序的字节中，如果该整数的最低有效字节（类似于最低有效位）在最高有效字节的前面，则称小端序；反之则称大端序。在网络应用中，字节序是一个必须被考虑的因素，因为不同机器类型可能采用不同标准的字节序，所以均按照网络标准转化。

   例如假设上述变量 x 类型为 int，位于地址 0x100 处，它的值为 0x01234567，地址范围为 0x100~0x103 字节，其内部排列顺序依赖于机器的类型。大端法从首位开始将是：0x100: 01, 0x101: 23,..。而小端法将是：0x100: 67, 0x101: 45,..。

4. 图片灰度化:CanvasRenderingContext2D 提供的 `getImageData`/`putImageData`
5. 图片压缩:Canvas 对象提供的 `toDataURL()` 方法;该方法接收 type 和 encoderOptions 两个可选参数。
   它还提供了一个 toBlob() 方法

```JS
canvas.toBlob(callback, mimeType, qualityArgument)
```

和 toDataURL() 方法相比，toBlob() 方法是异步的，因此多了个 callback 参数，这个 callback 回调方法默认的第一个参数就是转换好的 blob 文件信息。

6. 图片上传：FormData
7. 如何查看图片的二进制数据
   Windows 平台下的 「WinHex」 或 macOS 平台下的 「Synalyze It! Pro」 十六进制编辑器
