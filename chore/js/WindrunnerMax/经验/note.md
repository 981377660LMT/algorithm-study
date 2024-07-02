## Git

Git 是目前世界上最先进的分布式版本控制系统

- 工作位置

Workspace：工作区。
Index/Stage：暂存区。
Repository：本地仓库。
Remote：远程仓库

- 操作

git blame [file]: 显示指定文件是什么人在什么时间修改过

## 初探 webpack 之编写 plugin

在 webpack 运行的生命周期中会广播出许多事件，plugin 可以 hook 这些事件，在合适的时机通过 webpack 提供的 API `改变其在处理过程中的输出结果`

## 竞态问题与 RxJs

虽然 Js 是单线程语言，但由于引入了异步编程，所以也会存在竞态的问题，而使用 RxJs 通常就可以解决这个问题，其使得编写异步或基于回调的代码更容易

## 富文本演进之路

L0、L1、L2

- L0 阶段：`contenteditable`+`document.execCommand`，最原始的富文本编辑器
- L1 阶段：多了自定义`数据模型`的抽离
  构建一个描述文档结构与内容的数据模型，并且使用自定义的 execCommand 对数据描述模型进行修改
- L2 阶段：多了自定义的排版引擎
  codeMirror、Google Docs、腾讯文档
  主流的 L2 富文本编辑器都是借助于 Canvas 来绘制所有的内容

## 协同算法

协同算法最主要的目的是在尽可能保持用户的意图的情况下提供`最终一致性`，重点在于提供最终一致性而不是保持用户的意图
OT 通常必须要有中央服务器进行协同调度
CRDT 更适合分布式系统，可以不需要中央服务器

## step 设计

type、payload
https://www.cnblogs.com/WindrunnerMax/p/18135375
