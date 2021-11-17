1. 位图(bitmap)
   「位图图像（bitmap），亦称为点阵图像或栅格图像，是由称作**像素**（图片元素）的单个点组成的。」
   常用的位图处理软件有 Photoshop、Painter 和 Windows 系统自带的画图工具等。
   分辨率是位图不可逾越的壁垒，在对位图进行缩放、旋转等操作时，无法生产新的像素，因此会放大原有的像素填补空白，这样会让图片显得不清晰。
2. 矢量图
   它们都是通过`数学公式计算获得的`，具有编辑后不失真的特点。
   「矢量图以几何图形居多，图形可以无限放大，不变色、不模糊。」 常用于图案、标志、VI、文字等设计。常用软件有：CorelDraw、Illustrator、Freehand、XARA、CAD 等。
   可缩放矢量图形（英语：Scalable Vector Graphics，SVG）是一种基于可扩展标记语言（XML），用于描述二维矢量图形的图形格式。SVG 由 W3C 制定，是一个开放标准。
3. 位图的数学表示
   根据位深度，可将位图分为 1、4、8、16、24 及 32 位图像等。每个像素使用的信息位数越多，可用的颜色就越多，颜色表现就越逼真，相应的数据量越大。

   > 位深度为 1 的像素位图只有两个可能的值（黑色和白色），所以又称为二值图像。
   > 彩色图像，就是常说的 24 位真彩，约为 1670 万色。
   > 32 位：基于 24 位而生，增加 8 个位（256 种）的透明通道，共 4,294,967,296 种颜色。

   「图像处理的本质实际上就是对这些像素矩阵进行计算。」

4. 如何区分图片的类型
   「计算机并不是通过图片的后缀名来区分不同的图片类型，而是通过 “**魔数**”（Magic Number）来区分。」
   对于某一些类型的文件，`起始的几个字节内容都是固定的`，跟据这几个字节的内容就可以判断文件的类型。
   文件类型 文件后缀 魔数
   JPEG jpg/jpeg 0xFFD8FF
   PNG png 0x89504E47
   GIF gif 0x47494638（GIF8）
   BMP bmp 0x424D
   你想要判断一张图片是否为 PNG 类型，这时你可以使用 is-png 这个库。它同时支持浏览器和 Node.js，使用示例如下：
5. 如何获取图片的尺寸
   图片的尺寸、位深度、色彩类型和压缩算法都会存储在文件的二进制数据中
6. 如何预览本地图片
   利用 HTML `FileReader` API，我们也可以方便的实现图片本地预览功能

```HTML
<input type="file" accept="image/*" onchange="loadFile(event)">
<img id="output"/>
<script>
  const loadFile = function(event) {
    const reader = new FileReader();
    reader.onload = function(){
      const output = document.querySelector('output');
      output.src = reader.result;
    };
    reader.readAsDataURL(event.target.files[0]);
  };
</script>

```

服务端需要做一些相关处理，才能正常保存上传的图片，这里以 Express 为例(`通过 buffer`)

```JS
const app = require('express')();

app.post('/upload', function(req, res){
    let imgData = req.body.imgData; // 获取POST请求中的base64图片数据
    let base64Data = imgData.replace(/^data:image\/\w+;base64,/, "");
    let dataBuffer = Buffer.from(base64Data, 'base64');
    fs.writeFile("image.png", dataBuffer, function(err) {
        if(err){
          res.send(err);
        }else{
          res.send("图片上传成功！");
        }
    });
});

```

7. 如何实现图片压缩(canvas toDataUrl=>dataUrlToBlob=>通过 FormData 发送)
   我们希望在上传本地图片时，先对图片进行一定的压缩，然后再提交到服务器，从而减少传输的数据量。在前端要实现图片压缩，`我们可以利用 Canvas 对象提供的 toDataURL()` 方法，该方法接收 type 和 encoderOptions 两个可选参数。
   其中 type 表示图片格式，默认为 image/png。而 encoderOptions 用于表示图片的质量，在指定图片格式为 image/jpeg 或 image/webp 的情况下，`可以从 0 到 1 的区间内选择图片的质量`。如果超出取值范围，将会使用默认值 0.92，其他参数会被忽略。

```JS
function compress(base64, quality, mimeType) {
  let canvas = document.createElement("canvas");
  let img = document.createElement("img");
  img.crossOrigin = "anonymous";
  return new Promise((resolve, reject) => {
    img.src = base64;
    img.onload = () => {
      let targetWidth, targetHeight;
      if (img.width > MAX_WIDTH) {
        targetWidth = MAX_WIDTH;
        targetHeight = (img.height * MAX_WIDTH) / img.width;
      } else {
        targetWidth = img.width;
        targetHeight = img.height;
      }
      canvas.width = targetWidth;
      canvas.height = targetHeight;
      let ctx = canvas.getContext("2d");
      ctx.clearRect(0, 0, targetWidth, targetHeight); // 清除画布
      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
      let imageData = canvas.toDataURL(mimeType, quality / 100);
      resolve(imageData);
    };
  });
}

  对于返回的 Data URL 格式的图片数据，为了进一步减少传输的数据量，我们可以把它转换为 Blob 对象：
  function dataUrlToBlob(base64, mimeType) {
      let bytes = window.atob(base64.split(",")[1]);
      let ab = new ArrayBuffer(bytes.length);
      let ia = new Uint8Array(ab);
      for (let i = 0; i < bytes.length; i++) {
        ia[i] = bytes.charCodeAt(i);
      }
      return new Blob([ab], { type: mimeType });
  }


  在转换完成后，我们就可以压缩后的图片对应的 Blob 对象封装在 FormData 对象中，然后再通过 AJAX 提交到服务器上：

  function uploadFile(url, blob) {
      let formData = new FormData();
      let request = new XMLHttpRequest();
      formData.append("image", blob);
      request.open("POST", url, true);
      request.send(formData);
  }


```

8.  如何操作位图像素数据
    我们可以利用 CanvasRenderingContext2D 提供的 `getImageData` 来获取图片像素数据

```JS
const imageData=ctx.getImageData(sx, sy, sw, sh);
console.log(imageData.data) // Uint8ClampedArray(1920000)
Uint8Array 与 Uint8ClampedArray 的区别
Uint8ClampedArray:如果你指定一个在 [0,255] 区间外的值，它将`被替换为0或255`；如果你指定一个非整数，那么它将被设置为最接近它的整数。
Uint8Array:它将输入数`与256取模`，将8个比特位转化为正整数，它也不会进行四舍五入。
```

当完成处理后，若要在页面上显示处理效果，则我们需要利用 CanvasRenderingContext2D 提供的另一个 API —— putImageData。

```JS
void ctx.putImageData(imagedata, dx, dy);
void ctx.putImageData(imagedata, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight);

```

9. 如何实现图片隐写
   「隐写术是一门关于信息隐藏的技巧与科学，所谓信息隐藏指的是不让除预期的接收者之外的任何人知晓信息的传递事件或者信息的内容。」 隐写术的英文叫做 Steganography，来源于特里特米乌斯的一本讲述密码学与隐写术的著作 Steganographia，该书书名源于希腊语，意为 “隐秘书写”。
