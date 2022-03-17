1. 模块的加载实质上就是，注入 exports、require、module 三个全局变量，然后执行模块的源码，然后将模块的 exports 变量的值输出。

```JS
console.log(555, arguments)
```

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

   总结:

   1. 相对路径转绝对路径
   2. 优先加载内置模块，即使有同名文件，也会优先使用内置模块。
   3. 根据路径判断是否有缓存
   4. 有缓存直接返回当前模块下缓存的 exports
   5. 无缓存则创建一个 Module 实例并缓存
   6. 取出模块后缀，根据后缀查找不同的方法并执行
   7. json 文件直接赋值给 module.exports
      js 文件:
      1. 包裹 js
      2. 使用 vm 模块执行被包裹的函数字符串，转化为真正的函数
      3. 利用 call 调用函数，从而修改 module.exports

3. require 的源码在 Node 的 lib/module.js 文件
4. require 函数返回的结果就是对应文件 module.exports 的值
5. 模块类型
   内置模块：就是 Node.js 原生提供的功能，比如 fs，http 等等，这些模块在 **Node.js 进程起来时就加载了**。
   文件模块：我们前面写的几个模块，还有第三方模块，即 node_modules 下面的模块都是文件模块。
6. 模块加载参数

```JS
[Arguments] {
  '0': {},
  '1': [Function: require] {
    resolve: [Function: resolve] { paths: [Function: paths] },
    main: Module {
      id: '.',
      path: 'e:\\test\\js\\react源码\\react-class-source-code',
      exports: {},
      parent: null,
      filename: 'e:\\test\\js\\react源码\\react-class-source-code\\函数长度.js',
      loaded: false,
      children: [],
      paths: [Array]
    },
    extensions: [Object: null prototype] {
      '.js': [Function (anonymous)],
      '.json': [Function (anonymous)],
      '.node': [Function (anonymous)]
    },
    cache: [Object: null prototype] {
      'e:\\test\\js\\react源码\\react-class-source-code\\函数长度.js': [Module]
    }
  },
  '2': Module {
    id: '.',
    path: 'e:\\test\\js\\react源码\\react-class-source-code',
    exports: {},
    parent: null,
    filename: 'e:\\test\\js\\react源码\\react-class-source-code\\函数长度.js',
    loaded: false,
    children: [],
    paths: [
      'e:\\test\\js\\react源码\\react-class-source-code\\node_modules',
      'e:\\test\\js\\react源码\\node_modules',
      'e:\\test\\js\\node_modules',
      'e:\\test\\node_modules',
      'e:\\node_modules'
    ]
  },
  '3': 'e:\\test\\js\\react源码\\react-class-source-code\\函数长度.js',
  '4': 'e:\\test\\js\\react源码\\react-class-source-code'
}
```

7. 循环引用

   1. main 加载 a，a 在真正加载前先去缓存中占一个位置
   2. a 在正式加载时加载了 b
   3. b 又去加载了 a，这时候缓存中已经有 a 了，所以直接返回 a.exports，即使这时候的 exports 是不完整的。

8. 总结

   1. require 不是黑魔法，整个 Node.js 的模块加载机制都是 JS 实现的。
   2. 每个模块里面的 exports, require, module, **filename, **dirname 五个参数都不是全局变量，而是模块加载的时候注入的。
   3. 为了注入这几个变量，我们需要将用户的代码用一个函数包裹起来，拼一个字符串然后调用沙盒模块 vm 来实现。
   4. 初始状态下，模块里面的 this, exports, module.exports 都指向同一个对象，如果你对他们重新赋值，这种连接就断了。
   5. 对 module.exports 的重新赋值会作为模块的导出内容，但是你`对 exports 的重新赋值并不能改变模块导出内容`，只是改变了 exports 这个变量而已，因为`模块始终是 module，导出内容是 module.exports。`
   6. 为了解决循环引用，模块在加载前就会被加入缓存，下次再加载会直接返回缓存，**如果这时候模块还没加载完，你可能拿到未完成的 exports**。
   7. Node.js 实现的这套加载机制叫 CommonJS。
