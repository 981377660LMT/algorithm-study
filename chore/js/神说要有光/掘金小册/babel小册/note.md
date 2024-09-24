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

# babel 插件进阶

# babel 插件实战

https://github.com/QuarkGluonPlasma/babel-plugin-exercize

# 手写建议的 babel

https://github.com/QuarkGluonPlasma/babel-plugin-exercize
