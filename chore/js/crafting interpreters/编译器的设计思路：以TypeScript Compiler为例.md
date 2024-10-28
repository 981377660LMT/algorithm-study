编译器的设计思路：以 TypeScript Compiler 为例
https://zhuanlan.zhihu.com/p/636256231

- 学习 ts compiler api 的设计
- 难点在于模块划分、概念理解、架构设计

作者以 TypeScript Compiler（tsc）的架构为核心进行讲解，但会在本篇和接下来的一些文章中探讨下一些其它语言的某些设计思路，例如 Python 的 PEG parser、rustc 的一些设计

看完感觉，还是 crafting interpreters 的编译器设计思路更系统一些

# Part1：如何造一个编译器？

## 编译器是做什么的

- 两个任务：CodeGen、TypeCheck(Diagnostic)

## 表征语言的结构

## Syntactic analysis

## 定义语法并自动生成 parser

- 一个比较常用的工具是 Antlr4
- 从 Python3.9 起，其采用了一类叫做 PEG 的语法来自动生成 cpython 的 parser。Python 的创始人 Guido 写了一个文章系列来讲解这个叫做 pegen 的库的设计。
  https://docs.python.org/3/reference/grammar.html

## 还要检查什么？

## 树形结构的限制

## 嵌套 Scope

## 基于 Scope 的查表

## Symbol, Node, Symbol Table

## Binder

## resolveName()

## Program

AST -> Binder -> Checker -> Diagnostics

在类型检查前，会有几遍的 Resolve

## Checker/Semantic Analysis

# Part2：相关源码及索引

## 核心类型 src/compiler/types.ts

## Constructors

## 主干代码

## Binder src/compiler/binder.ts

## Checker src/compiler/checker.ts

# 编译器相关的参考资料
