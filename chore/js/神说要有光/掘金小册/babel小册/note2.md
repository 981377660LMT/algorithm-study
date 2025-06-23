# babel 插件实战

https://github.com/QuarkGluonPlasma/babel-plugin-exercize

## 自动埋点

函数插桩是在函数中插入一段逻辑但不影响函数原本逻辑，埋点就是一种常见的函数插桩，我们完全可以用 babel 来自动做。

实现思路分为`引入 tracker 模块`和`函数插桩`两部分：

- 引入 tracker 模块需要判断 ImportDeclaration 是否包含了 tracker 模块，没有的话就用 @babel/helper-module-import 来引入。
  如果已经引入过就不引入，没有的话就引入，并且生成个唯一 id 作为标识符。
- 函数插桩就是在函数体开始插入一段代码，如果没有函数体，需要包装一层，并且处理下返回值。、

实际上可能有的函数不需要埋点，这种可以自己做一下过滤，或者在函数上写上注释，然后根据注释来过滤，就像 eslint 支持 /_ eslint-disable / 来配置 rule 的开启关闭，teser 支持 `/ @**PURE**_/` 来配置纯函数一样。

## 自动国际化

需求：
如果没有引入 intl 模块，就自动引入，并且生成唯一的标识符，不和作用域的其他声明冲突
把字符串和模版字符串替换为 intl.t 的函数调用的形式
带有 `/*i18n-disable*/` 注释的字符串就忽略
把收集到的值收集起来，输出到一个资源文件中

---

![滴滴的 di18n](image-25.png)

替换字符串和模版字符串(StringLiteral 和 TemplateLiteral 节点)为对应的函数调用语句，要做模块的自动引入。
引入的 id 要生成全局唯一的，注意 jsx 中如果是属性的替换要用 {} 包裹。

## 自动生成 API 文档

对外提供 sdk 的话，那么自动文档生成是个刚需，不然每次都要人工同步改。
自动文档生成主要是信息的提取和渲染两部分，提取源码信息我们只需要分别处理 ClassDeclaration、FunctionDeclaration 或其他节点，然后从 ast 取出名字、注释等信息，之后通过 renderer 拼接成不同的字符串。

## Linter

## 类型检查

## 压缩混淆

## JS 解释器

## 模块遍历

## Babel Macros

## 如何调试 Babel 源码？
