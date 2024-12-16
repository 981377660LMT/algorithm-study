1.  再过几章，当 Lox 支持函数和动态分派时，从技术上讲，我们就不需要在语言中内置分支语句了。请说明如何通过这些语句实现条件执行。请说出一种在控制流中使用这种技术的语言。
    这种方法借鉴了 Smalltalk 语言的控制流实现方式。Smalltalk 语言以其面向对象的设计著称，控制流操作通过消息传递和方法调用来实现.
    `条件逻辑将通过对象的方法调用和回调机制来完成。`
    `原理： 控制流操作（如条件判断）被转化为方法，这些方法接受回调函数（callbacks），分别对应条件为真或假的情况下要执行的代码块`

    ifThen 方法接受一个回调函数 thenBranch，在 True 类中调用该回调，而在 False 类中不执行任何操作。
    ifThenElse 方法接受两个回调函数 thenBranch 和 elseBranch，在 True 类中调用 thenBranch，在 False 类中调用 elseBranch。

    当 condition 为 t（True 的实例）时，ifThen 调用 ifThenFn，打印 "if then -> then"；ifThenElse 调用 ifThenElseThenFn，打印 "if then else -> then"。
    当 condition 为 f（False 的实例）时，ifThen 不执行任何操作；ifThenElse 调用 ifThenElseElseFn，打印 "if then else -> else"。

    The basic idea is that the control flow operations become methods that take
    callbacks for the blocks to execute when true or false. You define two
    classes with singleton instances, one for true and one for false. The
    implementations of the control flow methods on the true class invoke the
    then callbacks. The ones on the false class implement the else callbacks.

    Like so:

    ```lox
    class True {
      ifThen(thenBranch) {
        return thenBranch();
      }

      ifThenElse(thenBranch, elseBranch) {
        return thenBranch();
      }
    }

    class False {
      ifThen(thenBranch) {
        return nil;
      }

      ifThenElse(thenBranch, elseBranch) {
        return elseBranch();
      }
    }
    ```

    Then we make singleton instances of these classes:

    ```lox
    var t = True();
    var f = False();
    ```

    You can try them out like so:

    ```lox
    fun test(condition) {
      fun ifThenFn() {
        print "if then -> then";
      }

      condition.ifThen(ifThenFn);

      fun ifThenElseThenFn() {
        print "if then else -> then";
      }

      fun ifThenElseElseFn() {
        print "if then else -> else";
      }

      condition.ifThenElse(ifThenElseThenFn, ifThenElseElseFn);
    }

    test(t);
    test(f);
    ```

    This is famously how Smalltalk implements its control flow.

    It looks cumbersome because Lox doesn't have lambdas -- anonymous function
    expressions -- but those would be easy to add to the language if
    we wanted to go in this direction.

    Even more powerful would a nice terse syntax for defining and passing a
    closure to a method. The Grace language has a particularly nice notation
    for passing multiple blocks to a method. If we adapted that to Lox, we'd
    get something like:

    ```text
    fun test(condition) {
      condition.ifThen {
        print "if then -> then";
      };

      condition.ifThen {
        print "if then else -> then";
      } else {
        print "if then else -> else";
      };
    }

    test(t);
    test(f);
    ```

    It starts to look like this control flow is built into the language even
    though it's only method calls.

2.  只要我们的解释器支持一项重要的优化，就可以使用同样的工具实现循环。它是什么，为什么需要它？请说出一种使用这种技术进行迭代的语言。
    `通过递归和条件执行来实现循环结构的技术`，特别是依赖于一种关键的优化——`尾调用优化（Tail Call Optimization）。`
    **尾调用优化（Tail Call Optimization，TCO）**是一种编译器或解释器的优化技术，用于优化递归函数调用。具体来说，当一个函数的最后一个操作是调用另一个函数（或自身）时，优化器可以重用当前函数的栈帧，而无需为新的函数调用分配新的栈空间。这种优化可以显著减少递归调用带来的栈空间开销，避免因深度递归导致的栈溢出（Stack Overflow）问题。
    Scheme 是一种著名的支持尾调用优化的编程语言。`Scheme 强调函数式编程，广泛使用递归来实现迭代`，因此对尾调用优化有严格的要求。根据 Scheme 的标准，所有实现都必须支持尾调用优化，使得程序员可以安全地使用递归进行无限迭代，而无需担心栈溢出。

    Scheme is the language that famously shows that all iteration can be
    represented in terms of recursion and conditional execution. To execute a
    chunk of code more than once, hoist it out into a function that calls itself
    at the end of its body for the next iteration.

    For example, we could represent this `for` loop:

    ```lox
    for (var i = 0; i < 100; i = i + 1) {
      print i;
    }
    ```

    Like so:

    ```lox
    fun forStep(i) {
      print i;
      if (i < 99) forStep(i + 1);
    }
    ```

    When you see heavy use of recursion like here where there are almost a
    hundred recursive calls, the concern is overflowing the stack. However, in
    many cases, you don't need to preserve any information from the previous
    call when beginning a recursive call. If the recursive call is in _tail
    position_ -- it's the last thing in the body of the function -- then you
    can discard any stack space used by the previous call before beginning the
    next one.
    如果递归调用处于尾位置（Tail Position），即它是函数体中的最后一个操作，那么在开始下一个递归调用之前，可以释放之前调用所占用的栈空间。这意味着不需要为每次递归调用分配新的栈帧，从而避免栈空间的持续增长。

    This **tail call optimization** lets you use recursion for an unbounded
    number of iterations while consuming only a constant amount of stack space.
    Scheme and some other functional languages require an implementation to
    perform this optimization so that users can safely rely on recursion for
    iteration.

3.  增加对 break 语句的支持。
    这包括对抽象语法树（AST）的修改、词法分析器（Lexer）、语法解析器（Parser）以及解释器（Interpreter）的相应调整。

    As usual, we start with the AST:

    ```java
    defineAst(outputDir, "Stmt", Arrays.asList(
      "Block      : List<Stmt> statements",
      "Break      : ",  // <--
      "Expression : Expr expression",
      "If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
      "Print      : Expr expression",
      "Var        : Token name, Expr initializer",
      "While      : Expr condition, Stmt body"
    ));
    ```

    Break doesn't have any fields, which actually breaks the little generator
    script, so you also need to change defineType() to:
    Break 语句不需要任何字段，因为它只是一个简单的跳出循环的指令

    ```java
    // Store parameters in fields.
    String[] fields;
    if (fieldList.isEmpty()) {
      fields = new String[0];
    } else {
      fields = fieldList.split(", ");
    }
    ```

    Run that to get the new AST class. Now we need to push the syntax through the
    front end, starting with the new keyword. In TokenType, add `BREAK`:
    接下来，我们需要在词法分析器中识别 break 关键字。

    ```java
    // Keywords.
    AND, BREAK, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, OR,
    ```

    And then define it in the lexer:

    ```java
    keywords.put("break",  BREAK);
    ```

    In the parser, we match the keyword in `statement()`:

    ```java
    if (match(BREAK)) return breakStatement();
    ```

    Which calls:

    ```java
    private Stmt breakStatement() {
      consume(SEMICOLON, "Expect ';' after 'break'.");
      return new Stmt.Break();
    }
    ```

    We need some additional parser support. It should be a syntax error to use
    `break` outside of a loop. We do that by adding a field in Parser to track
    how many enclosing loops there currently are:
    `为了确保 break 语句只能在循环内部使用，我们需要在解析器中跟踪当前嵌套的循环数量。(try-finally)`

    ```java
    private int loopDepth = 0;
    ```

    In `forStatement()`, we update that when parsing the loop body:

    ```java
    try {
      loopDepth++;
      Stmt body = statement();

      if (increment != null) {
        body = new Stmt.Block(Arrays.asList(
            body,
            new Stmt.Expression(increment)));
      }

      if (condition == null) condition = new Expr.Literal(true);
      body = new Stmt.While(condition, body);

      if (initializer != null) {
        body = new Stmt.Block(Arrays.asList(initializer, body));
      }

      return body;
    } finally {
      loopDepth--;
    }
    ```

    Likewise `whileStatement()`:

    ```java
    try {
      loopDepth++;
      Stmt body = statement();

      return new Stmt.While(condition, body);
    } finally {
      loopDepth--;
    }
    ```

    Now we can check that when parsing the `break` statement:

    ```java
    private Stmt breakStatement() {
      if (loopDepth == 0) {
        error(previous(), "Must be inside a loop to use 'break'.");
      }
      consume(SEMICOLON, "Expect ';' after 'break'.");
      return new Stmt.Break();
    }
    ```

    To interpret this, we'll use exceptions to jump from the break out of the
    loop. In Interpreter, define a class:
    `最后，我们需要在解释器中实现对 break 语句的处理。由于 break 语句需要跳出循环，我们可以通过抛出异常来实现这一点。`

    ```java
    private static class BreakException extends RuntimeException {}
    ```

    Executing a `break` simply throws that:
    当解释器访问到一个 Break 语句时，抛出 BreakException 异常，以中断当前循环的执行

    ```java
    @Override
    public Void visitBreakStmt(Stmt.Break stmt) {
      throw new BreakException();
    }
    ```

    That gets caught by the `while` loop code and then proceeds from there.
    在执行 while 循环时，使用 try-catch 块捕获 BreakException 异常。
    当捕获到 BreakException 时，停止当前循环的执行，继续执行循环之后的代码。

    ```java
    @Override
    public Void visitWhileStmt(Stmt.While stmt) {
      try {
        while (isTruthy(evaluate(stmt.condition))) {
          execute(stmt.body);
        }
      } catch (BreakException ex) {
        // Do nothing.
      }
      return null;
    }
    ```
