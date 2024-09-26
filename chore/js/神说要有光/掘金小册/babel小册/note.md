1. Program & Directive
   Program 是包裹具体执行语句的节点，而 Directive 则是代码中的指令部分。
2. File & Comment
   Comment：注释分为块注释和行内注释，对应 CommentBlock 和 CommentLine 节点。

---

如果想查看全部的 AST 可以在 babel parser 仓库里的 AST 文档里查，或者直接去看 @babel/types 的 typescript 类型定义
https://github.com/babel/babel/blob/main/packages/babel-types/src/ast-types/generated/index.ts
https://github.com/babel/babel/blob/main/packages/babel-parser/ast/spec.md

---

AST 的公共属性
https://astexplorer.net/#/gist/7267e806bfec60b48b9d39d039f29313/c343ad5a76a8dd78c22d39ce89f4d0733c2b17e4
