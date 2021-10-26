1. 模块的加载实质上就是，注入 exports、require、module 三个全局变量，然后执行模块的源码，然后将模块的 exports 变量的值输出。

```JS
(function (exports, require, module, __filename, __dirname) {
  // 模块源码
});
```

2. require 的 用法
   当 Node 遇到 require(X) 时，按下面的顺序处理

   1. 如果是内置模块 返回
   2. 如果是相对路径(以 "./" 或者 "/" 或者 "../" 开头)
      - 根据 X 所在的父模块，确定 X 的绝对路径。
      - 将 X 当成文件，依次查找下面文件，只要其中有一个存在，就返回该文件，不再继续执行。
        `X` `X.js` `X.json` `X.node`
      - 将 X 当成目录，依次查找下面文件，只要其中有一个存在，就返回该文件，不再继续执行。
        `X/package.json(main 字段)` `X/index.js` `X/index.json` `X/index.node`
   3. 如果是绝对路径
      - 根据 X 所在的父模块，确定 X 可能的安装目录。
      - 依次在每个目录中，将 X 当成文件名或目录名加载。
   4. 抛出 "not found"

   例子：`/home/ry/projects/foo.js 执行了 require('bar') ` 上面的第三种情况
   /home/ry/projects/node_modules/bar
   /home/ry/node_modules/bar
   /home/node_modules/bar
   /node_modules/bar
   搜索时，Node 先将 bar 当成文件名，依次尝试加载下面这些文件，只要有一个成功就返回。
   bar
   bar.js
   bar.json
   bar.node
   如果都不成功，说明 bar 可能是目录名，于是依次尝试加载下面这些文件。
   bar/package.json（main 字段）
   bar/index.js
   bar/index.json
   bar/index.node

3. require 的源码在 Node 的 lib/module.js 文件
