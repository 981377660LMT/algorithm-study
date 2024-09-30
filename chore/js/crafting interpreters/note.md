crafting interpreters
一本制作程序设计语言的指南

英文版：https://craftinginterpreters.com/
中文版：craftinginterpreters-zh-jet.vercel.app；https://confucianzuoyuan.github.io/craftinginterpreters/scanning.html
笔记：
https://misakatang.cn/2024/01/18/crafting-interpreters-notes/
https://xffxff.github.io/posts/crafting_interpreters_1
https://timothya.com/learning/crafting-interpreters/

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

设计笔记：隐藏的分号

- 几乎每一种新语言都会放弃一个小的语法点（一些古老的语言，比如 BASIC 从来没有过），那就是**将;作为显式的语句结束符**
  许多风格指南要求在每条语句后都显式地使用分号，尽管理论上该语言允许您省略它们，但 JavaScript 是我所知道的唯一一种（省略分号的）语言。

- 在不同的语言中，有各种不同的规则来决定哪些换行符是分隔符。

- 为什么 Python 中的 lambda 只允许单行的表达式体：
  如果要求进入一个嵌套在括号内的语句中，并且要求其中的换行是有意义的，**那么 Python 将需要一套不同的隐式连接行的规则**

## 5 Representing Code 代码的表示形式

![https://www.nosuchfield.com/2017/07/30/Talk-about-compilation-principles-2/](image-4.png)

- 代码的表示形式。它应该易于语法分析器生成，也易于解释器使用
  It should be simple for the parser to produce and easy for the interpreter to consume.
- 有一种方法可以将优先级进行可视化，那就是使用树。还有一种方式是字节码，这是另一种对人类不友好但更接近机器的表示方式。
  > 在美国，运算符优先级常缩写为 PEMDAS，分别表示 Parentheses(括号), Exponents(指数), Multiplication/Division(乘除), Addition/Subtraction(加减)。为了便于记忆，将缩写词扩充为“Please Excuse My Dear Aunt Sally”。

1. 上下文无关文法(CFG)
   https://www.nosuchfield.com/2017/07/30/Talk-about-compilation-principles-2/
   正则语言 (regular language) 是可以用正则表达式或自动机描述的语言(连接、选择、重复)。
   正则语言还不够强大，**无法处理可以任意深度嵌套的表达式(不能处理递归)**。例如，正则语法可以表达重复，但它们无法统计有多少重复。
   我们还需要一个更强大的工具，就是`上下文无关文法`（context-free grammar，CFG）。它是`形式化文法`的工具箱中又一个重量级的工具。
   | Terminology<br/>术语 | | Lexical grammar 词法 | Syntactic grammar 语法 |
   | ---------------------------------------- | --- | ----------------------------------- | ---------------------- |
   | The “alphabet” is . . . <br />字母表 | → | Characters<br />字符 | Tokens<br />词法标记 |
   | A “string” is . . . <br />字符串 | → | Lexeme or token<br />词素或词法标记 | Expression<br />表达式 |
   | It's implemented by the . . . <br />实现 | → | Scanner<br />扫描器 | Parser<br />解析器 |

   CFG 是一个四元组（N, T, P, S），组成为：

   - N 是非终止符集合 (Non-terminal)
   - T 是终止符集合 (Terminal)
   - P 是产生式集合 (Production rules)
   - S 是(唯一的)开始符号 (Start symbol)

   比较有名的是**巴科斯范式(BNF)**，它们的本质其实是一样的，都是对 **CFG 的四元组**的描述

   ```
   S –> AB
   A –> aA | ε
   B –> b | bB
   ```

   其中 S A B 就是非终结符，代表可以继续扩展或产生的符号；a b ε 是终结符，表示其无法再产生新的符号了，其中 ε 表示一个空句子；上面的每一行就是一个产生式规则，代表了一种非终结符的转移方式；而 S 就是开始符号。

   `形式化文法的工作是指定哪些字符串有效，哪些无效`。如果我们要为英语句子定义一个文法，“eggs are tasty for breakfast”会包含在文法中，但“tasty breakfast for are eggs”可能不会。

   - 语法规则(Rules for grammars)
     CFG 产生语言的基本方法 —— 推导
     如果你从规则入手，你可以用它们生成语法中的字符串。以这种方式`创建的字符串被称为推导式（derivations）`，因为每个字符串都是从语法规则中推导出来的。`规则被称为产生式(productions)`，因为它们生成了语法中的字符串。

     我试图提出一个简单的形式。 每个规则都是一个名称，后跟一个箭头（→），后跟一系列符号，最后以分号（;）结尾。 终止符是带引号的字符串，非终止符是小写的单词。

     早餐表达式语法：

     ```js
     breakfast  → protein "with" breakfast "on the side" ;
     breakfast  → protein ;
     breakfast  → bread ;

     protein    → crispiness "crispy" "bacon" ;
     protein    → "sausage" ;
     protein    → cooked "eggs" ;

     crispiness → "really" ;
     crispiness → "really" crispiness ;

     cooked     → "scrambled" ;
     cooked     → "poached" ;
     cooked     → "fried" ;

     bread      → "toast" ;
     bread      → "biscuits" ;
     bread      → "English muffin" ;
     ```

   - 增强符号(enhancing our notation，syntactic sugar)

     1. 我们将允许一系列由管道符(|)分隔的生成式，避免在每次在添加另一个生成式时重复规则名称。
        bread → "toast" | "biscuits" | "English muffin" ;
     2. 此外，我们允许用括号进行分组，然后在分组中可以用|表示从一系列生成式中选择一个。
        protein → ( "scrambled" | "poached" | "fried" ) "eggs" ;
     3. 我们也使用后缀`*`来允许前一个符号或组重复零次或多次。
        crispiness → "really" "really"`*` ;
     4. 后缀+与此类似，但要求前面的生成式至少出现一次。
        crispiness → "really"+ ;
     5. 后缀？表示可选生成式，它之前的生成式可以出现零次或一次，但不能出现多次。

     有了所有这些语法上的技巧，我们的早餐语法浓缩为：

     ```js
     breakfast → protein ( "with" breakfast "on the side" )?
               | bread ;

     protein   → "really"+ "crispy" "bacon"
               | "sausage"
               | ( "scrambled" | "poached" | "fried" ) "eggs" ;

     bread     → "toast" | "biscuits" | "English muffin" ;
     ```

     在本书的其余部分中，我们将使用这种表示法来精确地描述 Lox 的语法。当您使用编程语言时，您会发现上下文无关的语法(使用此语法或 EBNF 或其他一些符号)可以帮助您将非正式的语法设计思想具体化。它们也是与其他语言黑客交流语法的方便媒介。

   - Lox 表达式语法 (A Grammar for Lox expressions)
     现在，我们只关心几个表达式：

     - 字面量(Literals)：数字、字符串、布尔值以及 nil。
     - 一元表达式(Unary expressions)：!、-。
     - 二元表达式(Binary expressions)：+、-、`*`、/、>、>=、<、<=、==、!=。
     - 括号(Parentheses)

     使用我们的新符号，下面是语法的表示：

     ```js
     expression → literal
                | unary
                | binary
                | grouping ;

     literal    → NUMBER | STRING | "true" | "false" | "nil" ;
     grouping   → "(" expression ")" ;
     unary      → ( "-" | "!" ) expression ;
     binary     → expression operator expression ;
     operator   → "==" | "!=" | "<" | "<=" | ">" | ">="
                 | "+"  | "-"  | "*" | "/" ;
     ```

     除了与精确词素相匹配的终止符会加引号外，我们还对表示单一词素的终止符进行大写化，这些词素的文本表示方式可能会有所不同。NUMBER 是任何数字字面量，STRING 是任何字符串字面量。稍后，我们将对 IDENTIFIER 进行同样的处理
     这个语法实际上是有歧义的，我们在解析它时就会看到这一点。但现在这已经足够了。

2. 实现语法树

- 树节点数据结构
  - 非面向对象
    树节点是只有数据，没有方法的 dataClass。
    为什么？
    因为树节点不属于任何单个的领域。
    树是在解析的时候创建的，难道类中应该有解析对应的方法？或者因为树结构在解释的时候被消费，其中是不是要提供解释相关的方法？`树跨越了这些领域之间的边界，这意味着它们实际上不属于任何一方。`
    `这些类型的存在是为了让parser和interpreter能够进行交流`。
    这就适合于那些只是简单的数据而没有相关行为的类型。这种风格在 Lisp 和 ML 这样的函数式语言中是非常自然的，因为在这些语言中，所有的数据和行为都是分开的，但是在 Java 中感觉很奇怪。
- 节点树元编程
  generateAst 自动化生成节点类的代码
  `一个启示：dataClass 可以用脚本生成`

3. visitor 模式
   作者先提及解释器模式(其实就是模板方法)：可以在 Expr 上添加一个抽象的 interpret()方法，然后每个子类将实现该方法来解释自身。`这对于小型项目来说没问题，但扩展性很差。(This works alright for tiny projects, but it scales poorly. )`为什么？因为树节点跨越了几个领域。至少，解析器和解释器都会弄乱它们。
   如果我们为每个操作的表达式类添加实例方法，那么就会`将一堆不同的域混在一起`。这违反了关注点分离(separation of concerns)并导致代码难以维护。

   - 表达式问题
     ![行是类型，列是操作，单元格是实现代码](image-7.png)
     `添加新操作非常简单` —— 只需定义另一个与所有类型模式匹配的的函数即可
     ![添加新操作](image-5.png)
     但是，反过来说，`添加新类型是困难的`。您必须回头向已有函数中的所有模式匹配添加一个新的 case。
     ![添加新类型](image-6.png)

     **面向对象的语言希望你按照类型的行来组织你的代码。而函数式语言则鼓励你把每一列的代码都归纳为一个函数。**
     一群聪明的语言迷注意到，这两种风格都不容易向表格中添加行和列。`他们称这个困难为“表达式问题”`
     人们已经抛出了各种各样的语言特性、设计模式和编程技巧，试图解决这个问题，但还没有一种完美的语言能够解决它。与此同时，`我们所能做的就是尽量选择一种与我们正在编写的程序的自然架构相匹配的语言。`

   - visitor 模式
     **Visitor 模式让你可以在面向对象的语言中模仿函数式**
     访问者模式是所有设计模式中最容易被误解的模式。
     问题出在术语上。这个模式不是关于“visiting（访问）”，它的 “accept”方法也没有让人产生任何有用的想象。
     访问者模式实际上**近似于 OOP 语言中的函数式。它让我们可以很容易地向表中添加新的列。**
     我们可以在一个地方定义针对一组类型的新操作的所有行为，而不必触及类型本身。`这与我们解决计算机科学中几乎所有问题的方式相同：添加中间层。(adding a layer of indirection)`

4. 一个（不是很）漂亮的打印器(pretty printer)
   我们希望字符串非常明确地显示树的嵌套结构。

## 6 Parsing Expressions 解析表达式

有一整套解析技术，其名称大多是“L”和“R”的组合——LL(k)、LR(1)、LALR——以及更奇特的野兽，如解析器组合器、Earley 解析器、调车场算法和 Packrat 解析。对于我们的第一个解释器，一种技术就足够了：递归下降。

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

---

问题

- 即时编译(JIT)往往是实现动态类型语言的最快方法，但并非所有语言都使用它。有哪些理由不采用 JIT？
  Slower startup 启动速度较慢
  Memory overhead 内存开销
  Implementation complexity 实施复杂度

  https://stackoverflow.com/q/3221861

- The lexical grammars of Python and Haskell are not regular. What does that mean, and why aren’t they? Python 和 Haskell 的词法语法并不规则。这意味着什么？为什么不呢？
  正则语言是可以用正则表达式或确定性或非确定性有限自动机或状态机来表达的语言。
  **Python 基于缩进的作用域无法用正则表达式来表达。**
