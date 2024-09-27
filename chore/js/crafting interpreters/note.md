crafting interpreters
一本制作程序设计语言的指南

英文版：https://craftinginterpreters.com/
中文版：craftinginterpreters-zh-jet.vercel.app
笔记：
https://misakatang.cn/2024/01/18/crafting-interpreters-notes/
https://xffxff.github.io/posts/crafting_interpreters_1

# I WELCOME 导读

真正实现 lox 语言两次

## 1 Introduction 前言

不在编译理论部分着墨过多。当我们逐步搭建起编译系统的各个部分之时，我再向你介绍关于这部分的历史与其背后的概念。

在大部分情况下，我们都将绞尽脑汁先让我们的程序运行起来。这并不意味着理论不重要，能够正确且形式化地推导语法和语义是开发程序设计语言一项至关重要的技能。但是，我个人还是喜欢边做边学。通读大段大段充斥抽象概念的段落并充分理解其含义对我来说实在是太过困难，但如果让我写下几行代码，跑一跑，跟踪 debug 一下，我很快就可以掌握它。

动手实操，对一门真正程序设计语言的诞生过程留下一个坚实的印象。

### 为什么要学习它

- 小众语言(DSL)无处不在
  对于任何一支成功的通用程序设计语言，都会有上千种其他语言与其相辅相成。我们通常将它们称为“小众语言(Domain-specific Languages，DSL)”。
  用途：应用程序脚本、模版引擎、标记格式、配置文件等
  ![DSL](image.png)
  如果当现有的领域特定语言代码库无法满足你的需求时，你还是需要勉为其难地手写一枚解析器或者类似的工具以满足需求。即使你想要重用一些现有的代码，你也将不可避免地对其进行调试和维护，深入研究其原理。
- 提高编程能力
  长跑运动员负重、高原训练的例子

## 2 A Map of the Territory 解释器简介

## 3 The Lox Language Lox 语言介绍

# II A TREE-WALK INTERPRETER jlox 介绍

## 4 Scanning 扫描器相关

## 5 Representing Code 表示代码

## 6 Parsing Expressions 解析表达式

## 7 Evaluating Expressions 执行表达式

## 8 Statements and State 语句和状态

## 9 Control Flow 控制流

## 10 Functions 函数

## 11 Resolving and Binding 解析和绑定

## 12 Classes 类

## 13 Inheritance 继承

# III A BYTECODE VIRTUAL MACHINE clox 介绍

## 14 Chunks of Bytecode 字节码

## 15 A Virtual Machine 虚拟机

## 16 Scanning on Demand 扫描

## 17 Compiling Expressions 编译表达式

## 18 Types of Values 值类型

## 19 Strings 字符串

## 20 Hash Tables 哈希表

## 21 Global Variables 全局变量

## 22 Local Variables 局部变量

## 23 Jumping Back and Forth 来回跳转

## 24 Calls and Functions 调用和函数

## 25 Closures 闭包

## 26 Garbage Collection 垃圾回收

## 27 Classes and Instances 类和实例

## 28 Methods and Initializers 方法和初始化

## 29 Superclasses 超类

## 30 Optimization 优化

# BACKMATTER 后记

## A1 Appendix I: Lox Grammar Lox 语法

## A2 Appendix II: Generated Syntax Tree Classes 语法树类
