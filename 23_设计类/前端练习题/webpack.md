要手写一个模块打包器（Bundler），我们需要深入理解 **AST（抽象语法树）**、**模块解析** 和 **代码生成**。

我们将实现一个迷你版的 Webpack，我给它起名叫 **`Minipack`**。

它将具备以下核心能力：

1.  **依赖分析**：读取入口文件，解析 `import` 语句，找到所有依赖。
2.  **构建图谱 (Dependency Graph)**：递归遍历，生成所有模块的依赖关系图。
3.  **打包 (Bundling)**：将所有模块合并成一个可以在浏览器运行的文件（IIFE 闭包）。

---

### 1. 准备工作：安装必要的编译器工具

我们需要三个工具来处理代码（这也是 Webpack 内部使用的）：

- `@babel/parser`: 把代码变成 AST（抽象语法树）。
- `@babel/traverse`: 遍历 AST，找到 `import` 语句。
- `@babel/core`: 把 ES6 代码转换成浏览器能跑的 ES5 代码。

_(注：为了演示原理，这里假设你已经安装了这些库。实际运行需要 `npm install @babel/parser @babel/traverse @babel/core`)。_

### 2. 第一步：处理单个文件 (Asset Creator)

这个函数的作用是：**给你一个文件路径，你把它的内容、依赖关系和转换后的代码吐出来。**

```javascript
const fs = require('fs')
const path = require('path')
const parser = require('@babel/parser')
const traverse = require('@babel/traverse').default
const { transformFromAst } = require('@babel/core')

let ID = 0 // 自增 ID，用于唯一标识每个模块

function createAsset(filename) {
  // 1. 读取文件内容
  const content = fs.readFileSync(filename, 'utf-8')

  // 2. 生成 AST (抽象语法树)
  // sourceType: 'module' 表示我们要解析的是 ES Module
  const ast = parser.parse(content, {
    sourceType: 'module'
  })

  // 3. 收集依赖
  // 我们需要找到所有的 'import ... from ...'
  const dependencies = []

  traverse(ast, {
    // 访问者模式：每当遇到 ImportDeclaration 节点时执行
    ImportDeclaration: ({ node }) => {
      // node.source.value 就是 import 后面引号里的路径 (e.g., './message.js')
      dependencies.push(node.source.value)
    }
  })

  // 4. 代码转换 (Transpiling)
  // 将 ES6 代码 (import/export) 转换为 CommonJS 或浏览器能跑的代码
  const { code } = transformFromAst(ast, null, {
    presets: ['@babel/preset-env']
  })

  // 5. 返回模块对象
  return {
    id: ID++,
    filename,
    dependencies,
    code
  }
}
```

### 3. 第二步：构建依赖图谱 (Graph Builder)

这个函数的作用是：**从入口文件开始，顺藤摸瓜，把所有文件都找出来，拍平变成一个数组。**

```javascript
function createGraph(entry) {
  // 1. 解析入口文件
  const mainAsset = createAsset(entry)

  // 2. 使用队列进行广度优先遍历 (BFS)
  const queue = [mainAsset]

  // 遍历队列
  for (const asset of queue) {
    asset.mapping = {} // 用于存储 "相对路径 -> 模块ID" 的映射

    // 获取当前模块所在的目录 (用于解析相对路径)
    const dirname = path.dirname(asset.filename)

    // 遍历当前模块的所有依赖
    asset.dependencies.forEach(relativePath => {
      // 拼接绝对路径 (e.g., './message.js' -> '/User/project/src/message.js')
      const absolutePath = path.join(dirname, relativePath)

      // 关键：递归解析依赖模块
      const childAsset = createAsset(absolutePath)

      // 记录映射关系：在当前模块里，'./message.js' 对应的 ID 是 1
      asset.mapping[relativePath] = childAsset.id

      // 将子模块加入队列，继续遍历它的依赖
      queue.push(childAsset)
    })
  }

  // 返回所有模块的数组 (这就是依赖图谱)
  return queue
}
```

### 4. 第三步：生成打包代码 (Bundle Generator)

这是最神奇的一步。我们需要生成一个 **立即执行函数 (IIFE)**，模拟一个浏览器端的 `require` 系统。

```javascript
function bundle(graph) {
  let modules = ''

  // 1. 序列化模块
  // 我们要把 graph 数组转换成一个对象字符串：
  // {
  //   0: [ function(require, module, exports) { ...code... }, { './message.js': 1 } ],
  //   1: [ function(require, module, exports) { ...code... }, { } ]
  // }
  graph.forEach(mod => {
    // 这里的 mod.code 是 Babel 转译后的代码，它里面使用了 require() 和 exports
    // 我们需要把它们包裹在一个函数里，防止变量污染
    modules += `${mod.id}: [
      function (require, module, exports) {
        ${mod.code}
      },
      ${JSON.stringify(mod.mapping)},
    ],`
  })

  // 2. 构造最终的 IIFE 字符串
  const result = `
    (function(modules) {
      // 模拟 require 函数
      function require(id) {
        const [fn, mapping] = modules[id];

        // 定义一个局部的 require，用于处理相对路径
        // 因为模块内部代码写的是 require('./message.js')，而不是 require(1)
        function localRequire(name) {
          // 根据 mapping 找到模块 ID
          return require(mapping[name]);
        }

        const module = { exports: {} };

        // 执行模块代码
        // 这里的 fn 就是上面包裹的 function(require, module, exports) { ... }
        fn(localRequire, module, module.exports);

        return module.exports;
      }

      // 启动入口文件 (ID 为 0)
      require(0);
    })({${modules}})
  `

  return result
}
```

### 5. 实战演示

假设我们有以下三个文件：

**`src/name.js`**:

```javascript
export const name = 'World'
```

**`src/message.js`**:

```javascript
import { name } from './name.js'
export default `Hello ${name}!`
```

**`src/entry.js`** (入口):

```javascript
import message from './message.js'
console.log(message)
```

**运行打包器：**

```javascript
// 1. 构建图谱
const graph = createGraph('./src/entry.js')

// 2. 生成代码
const result = bundle(graph)

// 3. 输出结果
console.log(result)
```

### 6. 输出结果解析 (生成的 bundle.js)

运行上面的代码，你会得到类似下面的字符串。这就是浏览器能看懂的代码：

```javascript
;(function (modules) {
  function require(id) {
    // ... (上面定义的加载器逻辑) ...
  }
  require(0)
})({
  0: [
    function (require, module, exports) {
      'use strict'
      // Babel 转译后的 entry.js
      var _message = _interopRequireDefault(require('./message.js'))
      function _interopRequireDefault(obj) {
        return obj && obj.__esModule ? obj : { default: obj }
      }
      console.log(_message['default'])
    },
    { './message.js': 1 } // 映射表
  ],
  1: [
    function (require, module, exports) {
      'use strict'
      // Babel 转译后的 message.js
      Object.defineProperty(exports, '__esModule', { value: true })
      exports['default'] = void 0
      var _name = require('./name.js')
      var _default = 'Hello '.concat(_name.name, '!')
      exports['default'] = _default
    },
    { './name.js': 2 }
  ],
  2: [
    function (require, module, exports) {
      'use strict'
      // Babel 转译后的 name.js
      Object.defineProperty(exports, '__esModule', { value: true })
      exports.name = void 0
      var name = 'World'
      exports.name = name
    },
    {}
  ]
})
```

### 核心机制总结

1.  **AST 解析**：
    我们没有用正则去匹配 `import`，因为正则处理不了复杂的语法（比如注释里的 import）。AST 是唯一可靠的方法。

2.  **依赖图谱 (Graph)**：
    Webpack 的核心不是打包，而是**构建关系**。一旦有了 Graph，你想怎么打包都行（拆包、懒加载、Tree Shaking 都是基于 Graph 分析的）。

3.  **运行时模拟 (Runtime Shim)**：
    浏览器没有 `require` 和 `module.exports`。
    我们通过一个 IIFE 和一个自定义的 `require` 函数，在浏览器里**欺骗**了代码，让它们以为自己运行在 Node.js 环境里。
    - `localRequire` 是点睛之笔：它把代码里的相对路径 `'./name.js'` 转换成了打包器生成的数字 ID `2`。

这就是 Webpack、Vite (生产模式)、Rollup 等工具最底层的原理。掌握了这个，你就看透了前端工程化的本质。
