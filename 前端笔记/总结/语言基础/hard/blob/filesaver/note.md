1. FileSaver.js 是在客户端保存文件的解决方案，非常适合在客户端上生成文件的 Web 应用程序
   FileSaver.js 是 HTML5 的 saveAs() `FileSaver` 实现。它支持大多数主流的浏览器

```JS
FileSaver saveAs(Blob/File/Url, optional DOMString filename, optional Object { autoBom })

```

该方法支持字符串和 Blob 两种类型的参数

2. 在 FileSaver.js 内部提供了三种方案来实现文件保存
   当 FileSaver.js 在保存文件时，如果当前平台中 a 标签支持 `download` 属性且非 MacOS WebView 环境，则会优先使用 a[download] 来实现文件保存

   1. 字符串类型参数

   ```JS
   FileSaver.saveAs("https://httpbin.org/image", "image.jpg");
   ```

   下载资源的 URL 地址与当前站点是非同域的，则会先使用 同步的 HEAD 请求 来判断是否支持 CORS 机制，若支持的话，就会调用 download 方法进行文件下载

   ```JS
   function corsEnabled(url) {
     var xhr = new XMLHttpRequest();
     xhr.open("HEAD", url, false);
     try {
       xhr.send();
     } catch (e) {}
     return xhr.status >= 200 && xhr.status <= 299;
   }

   function download(url, name, opts) {
      var xhr = new XMLHttpRequest();
      xhr.open("GET", url);
      xhr.responseType = "blob";
      xhr.onload = function () {
        saveAs(xhr.response, name, opts);
      };
      xhr.onerror = function () {
        console.error("could not download file");
      };
      xhr.send();
    }


    我们需要设置 responseType 的类型为 blob。此外，因为返回的结果是 blob 类型的数据，所以在成功回调函数内部会继续调用 saveAs 方法来实现文件保存。

   ```

   而对于不支持 CORS 机制或同域的情形，它会调用内部的 click 方法来完成下载功能

   ```JS
       // `a.click()` doesn't work for all browsers (#465)
   function click(node) {
     try {
       node.dispatchEvent(new MouseEvent("click"));
     } catch (e) {
       var evt = document.createEvent("MouseEvents");
       evt.initMouseEvent(
         "click", true, true, window, 0, 0, 0, 80, 20,
         false, false, false, false, 0, null
       );
       node.dispatchEvent(evt);
     }
   }

   ```

   对于 blob 类型的参数，首先会通过 createObjectURL 方法来创建 Object URL，然后在通过 click 方法执行文件保存。为了能及时释放内存，在 else 处理分支中，会启动一个定时器来执行清理操作。
