crafting interpreters
一本制作程序设计语言的指南

英文版：https://craftinginterpreters.com/
中文版：craftinginterpreters-zh-jet.vercel.app；https://confucianzuoyuan.github.io/craftinginterpreters/scanning.html
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

1. Lox 采用与 C、Java 一脉相承的类 C 语法
2. JavaScript 最初诞生的那会儿，布兰登·艾奇为了能让网页上的按钮动起来，仅花了十天时间就设计实现了第一支 JS 解释器，放进网景浏览器里
3. 动态类型。将类型检查下推到运行时，可以让我们快速构建起一支可以执行代码的语言解释器。
4. 自动内存管理。
5. 数据类型：Boolean、Number(只有 f64)、String、NIL
   nil：如果我们使用的是静态类型语言，那么禁止它是值得的。然而，在动态类型中，消除它往往比保留它更加麻烦。
6. 表达式
   算术运算
   比较与相等:我通常是反对隐式转换的
   逻辑运算
   优先级与分组
7. 语句

   表达式的主要作用是产生一个值，语句的主要作用是产生一个效果
   表达式后跟分号（;）可以将表达式提升为语句状态。这被称为(很有想象力)表达式语句。

8. 变量
   你可以使用 var 语句声明变量。如果你省略了初始化操作，变量的值默认为 nil
9. 控制流
   if 语句
   while 语句
   for 语句
10. 函数
    fun 定义
    一等公民
    如果执行到达代码块的末尾而没有 return 语句，则会隐式返回 nil。
    支持闭包。由于闭包的存在，我们不能再假定变量作用域严格地像堆栈一样工作，在函数返回时局部变量就消失了
11. 类

    - 为什么任何语言都想要面向对象？
      对于动态类型语言来说，对象是`非常方便`的。我们需要某种方式来定义复合数据类型，用来将一堆数据组合在一起。
      如果我们也能把方法挂在这些对象上，那么我们就不需要把函数操作的数据类型的名字作为函数名称的前缀，以`避免与不同类型的类似函数发生冲突`。方法的作用域是对象，所以这个问题就不存在了。

    - 为什么 Lox 是面向对象的？
      我们很多人整天都在使用 OOP 语言
    - 类还是原型？
      我们将省去用户的麻烦，直接把`类`包含进去。
    - Lox 中的类:first-class
    - 实例化和初始化(Instantiation and initialization)
      如果您的类中包含一个名为 init() 的方法，则在构造对象时会自动调用该方法。传递给类的任何参数都会转发给它的初始化器：
    - 单继承
      当你声明一个类时，你可以使用小于(<)操作符指定它继承的类

      > 为什么用<操作符？我不喜欢引入一个新的关键字，比如 extends。

      super:子类通常也想定义自己的 init()方法。但还需要调用原始的初始化方法，以便超类能够维护其状态

      Lox 不是一种纯粹的面向对象的语言。`在真正的 OOP 语言中，每个对象都是一个类的实例，即使是像数字和布尔值这样的基本类型。`从类实例的意义上说，基本类型的值并不是真正的对象。它们没有方法或属性。如果以后我想让 Lox 成为真正的用户使用的语言，我会解决这个问题。

12. 标准库

---

表达式和语句
"一切都是表达式" 的语言往往具有函数式的血统，包括大多数 Lisps、SML、Haskell、Ruby 和 CoffeeScript。
要做到这一点，对于语言中的每一个 "类似于语句" 的构造，你需要决定它所计算的值是什么。其中有些很简单：

- if 表达式的计算结果是所选分支的结果。同样，switch 或其他多路分支的计算结果取决于所选择的情况。
- 变量声明的计算结果是变量的值。
- 块的计算结果是序列中最后一个表达式的结果。

# II A TREE-WALK INTERPRETER jlox 介绍

## 4 Scanning 扫描

1. 错误处理
   **把产生错误的代码和报告错误的代码分开**是一个很好的工程实践
   it’s good engineering practice to separate the code that generates the errors from the code that reports them
   功能齐全的语言实现中，您可能会通过多种方式显示错误：在 stderr 上、在 IDE 的错误窗口中、记录到文件中等。您不希望该代码遍布扫描仪和解析器。
2. Token
   - 类型
   - 字面量：数字、字符串等
   - 位置信息
     两个数字：**偏移量、token 长度。**
     知道偏移量之后，可以二分求出行号和列号.
3. 正则语言和表达式
   扫描器的核心是一个循环
   决定一门语言如何将字符分组为词素的规则被称为它的词法语法`(lexical grammar)`
   调库：像 Lex 或 Flex 这样的工具就是专门为实现这一功能而设计的——`向其中传入一些正则表达式，它可以为您提供完整的扫描器。`
   由于我们的目标是了解扫描器是如何工作的，`所以我们不会把这个任务交给正则表达式`。我们要亲自动手实现。
4. Scanner 类
   关键属性：
   `start、current、line`
5. 识别 token
   `_scanToken`
   `_advance()、_addToken()`

   - 词法错误
     错误的字符仍然会被前面调用的 advance()方法消费。这一点很重要，这样我们就不会陷入无限循环了
   - 单字符 token
   - 多字符 token
     操作符：需要 peek 下一个
     更长：lookahead 多个。大多数广泛使用的语言只需要提前一到两个字符。
   - 字面量

     - 字符串字面量： lox 中的字符串以"开头结尾
     - 数字字面量：lox 中的数字是 IEEE 754 双精度浮点数
       本可以让 peek()方法接受一个参数来表示要前瞻的字符数，而不需要定义两个函数。但这样做就会允许前瞻任意长度的字符。**提供两个函数可以让读者更清楚地知道，我们的扫描器最多只能向前看两个字符**

   - 保留字和标识符
     剩下的词素是 Boolean 和 nil，但我们把它们作为关键字来处理
     - maximal munch(最长匹配原则)
       每个符号序列总是以合法符号序列中最长的那个解释。当两个语法规则都能匹配扫描器正在处理的一大块代码时，哪个规则相匹配的字符最多，就使用哪个规则，尽管这样做会在语法分析器中导致后面的语法错误。
       `a+++p 会被解释为a++ +p`
       `<= 会被解释为 <= 而不是 < 和 =`
     - 如果匹配的话，就使用关键字的标记类型。否则，就是一个普通的`用户定义的标识符。`

   至此，我们就有了一个完整的扫描器，可以扫描整个 Lox 词法语法

---

为什么 Python 中的 lambda 只允许单行的表达式体

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
