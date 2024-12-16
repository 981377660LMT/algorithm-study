1.  为 REPL 添加支持，让用户同时输入语句和表达式。如果用户输入语句，则执行该语句。如果输入的是表达式，则对其进行评估并显示结果值。

    `表达式语法与语句语法的重叠问题：每个表达式都可以作为一个表达式语句（Expression Statement）的开始部分`
    It can be hard to do this in a clean way since the expression grammar
    overlaps the statement grammar so much (every expression is also the
    beginning of an expression statement containing that same expression).

    `常见的解析器解决方法：首先尝试解析为语句，如果解析失败，重新尝试解析为表达式。`
    One trick some parsers use is to simply _try_ to parse the syntax as a
    statement. If that fails, hide any parse errors and then try to parse it
    again as expression.

    `此处采用的不同方法：尝试解析为语句列表`
    I took a slightly different approach. Instead, the parser tries to parse a
    list of statements, but if it knows it's allowed to parse a single
    expression, and it reaches the end of the source right after parsing the
    expression part of an expression statement, then it stops early and returns
    that expression.
    `在完成解析后，解析器需要检查解析结果`
    All that's left is to see if the parsed value is an expression and, if so,
    evaluate it and print it.

    This isn't the cleanest implementation, but here goes. In Parser, add two
    new fields:

    ```java
    private boolean allowExpression;
    private boolean foundExpression = false;
    ```

    Then define this method:

    ```java
    Object parseRepl() {
      allowExpression = true;
      List<Stmt> statements = new ArrayList<>();
      while (!isAtEnd()) {
        statements.add(declaration());

        if (foundExpression) {
          Stmt last = statements.get(statements.size() - 1);
          return ((Stmt.Expression) last).expression;
        }

        allowExpression = false;
      }

      return statements;
    }
    ```

    And change expressionStatement() to:

    ```java
    private Stmt expressionStatement() {
      Expr expr = expression();

      if (allowExpression && isAtEnd()) {
        foundExpression = true;
      } else {
        consume(SEMICOLON, "Expect ';' after expression.");
      }
      return new Stmt.Expression(expr);
    }
    ```

    In Interpreter, add:

    ```java
    String interpret(Expr expression) {
      try {
        Object value = evaluate(expression);
        return stringify(value);
      } catch (RuntimeError error) {
        Lox.runtimeError(error);
        return null;
      }
    }
    ```

    Finally, in Lox, change runPrompt() to:

    ```java
    private static void runPrompt() throws IOException {
      InputStreamReader input = new InputStreamReader(System.in);
      BufferedReader reader = new BufferedReader(input);

      for (;;) {
        hadError = false;

        System.out.print("> ");
        Scanner scanner = new Scanner(reader.readLine());
        List<Token> tokens = scanner.scanTokens();

        Parser parser = new Parser(tokens);
        Object syntax = parser.parseRepl();

        // Ignore it if there was a syntax error.
        if (hadError) continue;

        if (syntax instanceof List) {
          interpreter.interpret((List<Stmt>)syntax);
        } else if (syntax instanceof Expr) {
          String result = interpreter.interpret((Expr)syntax);
          if (result != null) {
            System.out.println("= " + result);
          }
        }
      }
    }
    ```

    That should about do it.

2.  也许你希望 Lox 在变量初始化方面更明确一些。与其隐式地将变量初始化为 nil ，不如将访问未初始化或未赋值的变量定为运行时错误
    This is pretty simple. Instead of initializing variables with null if they
    have no initializer, we `use a special sentinel value to distinguish` it from
    Lox's nil. Then, we check for that when the variable is accessed.

    In Interpreter, add:

    ```java
    private static Object uninitialized = new Object();
    ```

    Change the first line of visitVarStmt() to:

    ```java
    Object value = uninitialized;
    ```

    Finally, change visitVariableExpr() to:

    ```java
    public Object visitVariableExpr(Expr.Variable expr) {
      Object value = environment.get(expr.name);
      if (value == uninitialized) {
        throw new RuntimeError(expr.name,
            "Variable must be initialized before use.");
      }
      return value;
    }
    ```

    The main downside is that checking for the uninitialized variable on every
    single access significantly slows execution for what is a very common
    operation. Not a big deal given that our Java interpreter isn't designed
    for speed anyway.
    主要缺点是，在每次访问变量时都要检查该变量是否已初始化，这会显著减慢执行速度，而变量访问是一项非常常见的操作。尽管如此，这并不是一个大问题，因为我们的 Java 解释器本身并不追求高速度。

3.  > What does the following program do?
    > 你希望它做什么？这是你认为它应该做的吗？您熟悉的其他语言中的类似代码是做什么的？您认为用户会期望它做什么？

    变量遮蔽(Naming Shadowing)

    ```lox
    var a = 1;
    {
      var a = a + 2;
      print a;
    }
    ```

    It prints 3. The shadowed variable doesn't come into scope until _after_ its
    initializer expression is evaluated, so `a + 2` is evaluated using the
    outer `a`, whose value is 1. Then the result is stored in the new `a`.

    > What did you expect it to do?

    Well, I wrote this book, so it's no surprise to me.

    > Is it what you think it should do?

    Code like this is rare in practice, so I don't care too much. But the
    current behavior is a little surprising. People read code left-to-right, so
    they probably expect the new variable to be in scope as soon as they scan
    over its name after `var`.

    Ideally, I'd make this kind of code a static error. Put the variable in
    scope as soon as its name is encountered but in a special "unusable" state.
    Then, once its initializer is done, make it available. If the initializer
    references it, make that a static error.

    > What does analogous code in other languages you are familiar with do?

    Java disallows shadowing local variables. C# allows shadowing, but doesn't
    allow multiple mentions of the same name in the same block to resolve to
    different variables.

    > What do you think users will expect this to do?

    I think they'd be surprised if the code was valid at all, and would
    probably consider it bad code even if it did do something.
