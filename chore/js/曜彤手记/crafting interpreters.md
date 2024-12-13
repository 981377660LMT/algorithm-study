https://www.hangyu.site/2023/07/18/%E3%80%8ACrafting-Interpreters%E3%80%8B%E8%AF%BB%E4%B9%A6%E7%AC%94%E8%AE%B0%EF%BC%88%E7%AC%AC%201-10%20%E7%AB%A0%EF%BC%89/

龙书是编译入门的理论圣经，这本书可以算是编译入门的实践宝典，好书。

1. self-hosting（自举）: implement a compiler in the same language it compiles.
   自举（self-hosting）：用与其编译相同的语言实现编译器。

2. Runtime: provide services like GC, reflection. `In a fully compiled language, the code implementing the runtime gets inserted directly into the resulting executable`.
   In, say, Go, each compiled application has its own copy of Go’s runtime directly embedded in it. If the language is run inside an interpreter or VM, then the runtime lives there, .e.g Java, Python, and JS.
   运行时：提供 GC、反射等服务。`在完全编译的语言中，实现运行时的代码会直接插入生成的可执行文件中`。例如，在 Go 中，每个编译后的应用程序都直接嵌入了自己的 Go 运行时副本。
   如果语言是在解释器或虚拟机中运行，那么运行时就存在于解释器或虚拟机中，例如 Java、Python 和 JS。

3. lox language

   - 表达式的主要作用是产生数值(produce a value)，而语句的作用是产生效果(produce an effect)。
     `表达式后跟一个分号（;），可将表达式提升为语句。`
   - first class members: it means they are real values that you can get a reference to, store in variables, pass around, etc.
     一等公民：这意味着它们是真实的值，可以获取引用、存储在变量中、四处传递等。
   - 实现object的两种方法
     1. class
     2. prototype

4. scanning
   - token = lexeme + metadata
   - “Maximal munch” principle
     当两个词法规则都能匹配scanner正在查看的代码块时，匹配字符数最多的规则获胜。
     ```c
     // For below C code, "-" and "--" are both valid lexical grammar rules, -
     // but only the last one is the correct scanning result.
     ---a => - --a;  // ✗
     ---a => -- -a;  // ✓
     ```
5. Representing Code
6. Parsing Expressions
   BNF 语法具有优先级（最低 -> 最高，每行规则分开）和关联性：
   ```bnf
   expression     → equality ;
   equality       → comparison ( ( "!=" | "==" ) comparison )* ;
   comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
   term           → factor ( ( "-" | "+" ) factor )* ;
   factor         → unary ( ( "/" | "*" ) unary )* ;
   unary          → ( "!" | "-" ) unary
                  | primary ;
   primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
   ```
   - `Each rule needs to match expressions at that precedence level or higher.`
     每条规则都需要匹配该优先级或更高的表达式。
   - Avoid having “left-recursive” in the BNF grammar for certain implementation algorithms `(recursion -> iteration)`.
     避免在某些实现算法（递归 -> 迭代）的 BNF 语法中出现 "左递归"。
     ```
     factor         → factor ( "/" | "*" ) unary | unary ;
     ```
7. Evaluating Expressions
   用户选择静态类型语言的一个重要原因是，静态类型语言能让他们确信，程序运行时绝不会出现某些类型的错误。
   如果将过多的类型检查推迟到运行时，就会削弱这种信心。
8. Statements and State
   - 表达式语句：可让您将表达式置于预期语句的位置。它们用于评估具有副作用的表达式。
     任何时候，只要你看到函数或方法调用后面跟着一个 ;，你就看到了一个表达式语句。
   - 关于变量和作用域的一条规则是：`"如果有疑问，就学 Scheme 的做法"`。
     Scheme 的主要目标之一就是将`词法作用域（静态作用域）`引入世界，因此如果你跟随他们的脚步，就很难出错。
     > 区别于动态作用域，词法作用域是指变量的作用域由程序的结构决定，与作用域所在的位置有关，而与程序的执行过程无关。
   - 赋值(assignment) 是优先度最低的`表达式`
9. Control Flow
   - Dangling else
     The else is bound to the nearest if that precedes it.
   - "logic_or "和 "logic_and "的优先级并不相同，**"logic_and "的优先级高于 "logic_or"**
     `在布尔代数中，"and"通常被看作是乘法，而 "or"则被看作是加法，这就是为什么它们通常继承了同名的优先级。`
     > logic_and, bitwise_and
   - 语法本身并不关心它们是否短路，这是语义问题，应在运行时处理。
   - "forStmt "可以作为 "whileStmt "的语法糖来实现。
10. Functions

- Native functions (primitives/external functions/foreign functions):
  functions that the interpreter exposes to user code but that are implemented in the host language, not the language being implemented.
  本地函数（primitivesd/external functions/foreign functions）：
  **解释器向用户代码公开的函数，但这些函数是用宿主语言而不是正在执行的语言实现的。**
- 命名函数声明并不是一个真正的原始操作。它是两个不同步骤（匿名函数 + 赋值）的语法糖：
  1. Creating a new function object.
     创建一个新的函数对象。
  2. Binding a new variable to it.
     绑定一个新变量。
- Return Statements:
  It could be implemented simply by Exception **(unwinding the visited call)**.
  可以简单地通过 Exception（解除访问调用）来实现。
