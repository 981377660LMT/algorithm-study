![babel插件通关秘籍](image.png)

实战部分

1. 自动国际化：自动转换代码为国际化之后的，自动引入资源包并替换代码中的文本
2. 自动生成文档：自动生成 api 文档，不再需要手动去维护
3. 自动埋点：自动进行函数插桩，埋点是一种常见的函数插桩
4. linter： 探索 eslint、stylelint 等 lint 工具的实现原理，能够实现各种 lint 插件
5. type checker：实现简单的 ts 类型检查，会对 typescript 的类型检查的原理会有更深的理解
6. 压缩混淆：前端必用工具之一，探索它的实现原理，压缩怎么做，混淆怎么做，怎么用 babel 实现，开阔下思路
7. js 解释器： AST 除了转译、静态分析外，还可以直接解释执行，学完这个案例可以知道解释器是怎么解释代码的
8. 模块遍历器：基于 babel 做模块的遍历，理解打包工具的依赖图构建原理
9. 手写 babel： 手写 babel 是为了加深对 babel 的理解，真正掌握 babel

# babel 插件基础

## 介绍

一般编译器 Compiler 是指高级语言到低级语言的转换工具。
而从高级语言到高级语言的转换工具，被叫做转换编译器，简称转译器 (Transpiler)。
babel 就是一个 Javascript Transpiler。`(也可以说，transpiler 是一种特殊的 compiler)`

babel 最开始叫 6to5，顾名思义是 es6 转 es5，但是后来随着 es 标准的演进，有了 es7、es8 等， 6to5 的名字已经不合适了，所以改名为了 babel。
babel 是巴别塔的意思，来自圣经中的典故。
我们平时主要用 babel 来做 3 种事情：

1. `转译` esnext、typescript、flow 等到目标环境支持的 js
2. 一些特定用途的代码`转换`
   例如：小程序转译工具 taro
3. 代码的`静态分析`
   例如:linter、api 文档生成工具、type checker、代码混淆工具、js 解释器

## Babel 的编译流程

![babel](image-1.png)
熟悉的三步：parse、transform、generate

1. parse：通过 parser 把源码转成抽象语法树（AST）
   ![词法分析出token+语法分析出ast](image-2.png)
2. transform：遍历 AST，调用各种 transform 插件对 AST 进行增删改
3. generate：把转换后的 AST 打印成目标代码，并生成 sourcemap

## Babel 的 AST

https://astexplorer.net/ 可以查看各种语言
点击这里的 save 就可以保存下来，然后把 url 分享出去
![Alt text](image-11.png)

1. Literal(字面量)
   ![Literal](image-3.png)
2. Identifier(标识符)
   变量名、属性名、参数名等各种声明和引用的名字，都是 Identifer。
   `JS 中的标识符只能包含字母或数字或下划线或美元符号，且不能以数字开头`。这是 Identifier 的词法特点。
   ![Alt text](image-4.png)
3. Statement(语句)
   我们写的每一条可以独立执行的代码都是语句。`语句末尾一般会加一个分号分隔，或者用换行分隔。`
   语句是代码执行的最小单位，可以说，代码是由语句（Statement）构成的。
   **特点是能够单独执行**
   ![Alt text](image-5.png)
4. Declaration
   声明语句用于定义变量
   ![Alt text](image-6.png)
   例如：

   - ImportDeclaration
   - ExportDefaultDeclaration、ExportNamedDeclaration、ExportAllDeclaration

5. Expression(表达式)
   执行完以后有返回值
   ![Alt text](image-7.png)
   ![a=1](image-8.png)

   ```jsonc
   // a=1
   {
     "type": "Program",
     "start": 0,
     "end": 3,
     "body": [
       {
         "type": "ExpressionStatement",
         "start": 0,
         "end": 3,
         "expression": {
           "type": "AssignmentExpression",
           "start": 0,
           "end": 3,
           "operator": "=",
           "left": {
             "type": "Identifier",
             "start": 0,
             "end": 1,
             "name": "a"
           },
           "right": {
             "type": "Literal",
             "start": 2,
             "end": 3,
             "value": 1,
             "raw": "1"
           }
         }
       }
     ],
     "sourceType": "module"
   }
   ```

6. Class
   ![Alt text](image-9.png)
7. Modules

- import: ImportSpicifier、ImportDefaultSpecifier、ImportNamespaceSpcifier
- export: ExportNamedDeclaration、ExportDefaultDeclaration、ExportAllDeclaration

8. Program & Directive
   program 是代表整个程序的节点，它有 body 属性代表程序体，存放 statement 数组；还有 directives 属性，存放 Directive 节点
   ![Alt text](image-10.png)
9. File & Comment
   babel 的 AST 最外层节点是 File，它有 program、comments、tokens 等属性，分别存放 Program 程序体、注释、token 等，是最外层节点
   块注释：CommentBlock
   行内注释：CommentLine

想查看全部的 AST 可以在 babel parser 仓库里的 AST 文档里查，或者直接去看 @babel/types 的 typescript 类型定义
https://github.com/babel/babel/blob/main/packages/babel-parser/ast/spec.md
https://github.com/babel/babel/blob/main/packages/babel-types/src/ast-types/generated/index.ts

---

AST 的公共属性

- type： AST 节点的类型
- start、end、loc：start 和 end 代表该节点在源码中的开始和结束下标。而 loc 属性是一个对象，有 line 和 column 属性分别记录开始和结束的行列号
- leadingComments、innerComments、trailingComments：注释
- extra：记录一些额外的信息，用于处理一些特殊情况。记录一些额外的信息，用于处理一些特殊情况。比如 StringLiteral 的 value 只是值的修改，而`修改 extra.raw 则可以连同单双引号一起修改`。
  ![Alt text](image-12.png)

## babel 的 Api

文档：https://www.babeljs.cn/docs/babel-parser

1. `@babel/parser` 对源码进行 parse，可以通过 plugins、sourceType 等来指定 parse 语法
   babel parser 叫 babylon，是基于 **acorn** 实现的，扩展了很多语法，可以支持 es next（现在支持到 es2020）、jsx、flow、typescript 等语法的解析。默认只能 parse js 代码，jsx、flow、typescript 这些非标准的语法的解析需要指定语法插件。
   提供了有两个 api：parse 和 parseExpression
   ```ts
   function parse(input: string, options?: ParserOptions): File
   function parseExpression(input: string, options?: ParserOptions): Expression
   ```
   ![options](image-13.png)
2. `@babel/traverse` 通过 visitor 函数对遍历到的 ast 进行处理，分为 enter 和 exit 两个阶段，具体操作 AST 使用 path 的 api，还可以通过 state 来在遍历过程中传递一些数据
   ```ts
   traverse(ast, {
     'FunctionDeclaration|VariableDeclaration'(path, state) {}
   })
   ```
   - path 有很多属性和方法，比如记录父子、兄弟等关系的、增删改 AST 的、判断 AST 类型的
   - state 是 Context，传递一些数据
3. `@babel/types` 用于创建、判断 AST 节点，提供了 xxx、isXxx、assertXxx 的 api
4. `@babel/template` 用于批量创建节点
   通过 @babel/types 创建 AST 还是比较麻烦的，要一个个的创建然后组装
   简化了创建 AST 的过程
   支持变量
5. `@babel/generator` 打印 AST 成目标代码字符串，支持 comments、minified、sourceMaps 等选项。
   ```ts
   function (ast: Object, opts: Object): {code, map}
   ```
6. @babel/code-frame 可以创建友好的报错信息
   控制台打印代码格式的功能就叫做 code frame(例如：高亮显示错误的代码行)

7. @babel/core 基于上面的包来完成 babel 的编译流程，可以从源码字符串、源码文件、AST 开始。
   ```ts
   transformSync(code, options) // => { code, map, ast }
   transformFileSync(filename, options) // => { code, map, ast }
   transformFromAstSync(parsedAst, sourceCode, options) // => { code, map, ast }
   ```
   options 主要配置 plugins 和 presets，指定具体要做什么转换。
   这些 api 也同样提供了异步的版本，异步地进行编译，返回一个 promise；明确是同步还是异步

## 实战案例：插入函数调用参数

通过 babel 能够自动在 console.log 等 api 中插入文件名和行列号的参数，方便定位到代码。

1. 分析代码 ast，确定思路
   https://astexplorer.net/#/gist/09113e146fa04044e99f8a98434a01af/80bef2b9068991f7a8e4f113ff824f56e3292253
   函数调用表达式的 AST 是 CallExpression。
   CallExrpession 节点有两个属性，**callee 和 arguments**，分别对应调用的函数名和参数， 所以我们要判断当 callee 是 console.xx 时，在 arguments 的数组中中插入一个 AST 节点。
   ![CallExpression](image-14.png)

   **babel parser 是自动生成的类型文件，类似条件编译**
   ![条件编译](image-15.png)

# babel 插件进阶

# babel 插件实战

https://github.com/QuarkGluonPlasma/babel-plugin-exercize

# 手写建议的 babel

https://github.com/QuarkGluonPlasma/babel-plugin-exercize
