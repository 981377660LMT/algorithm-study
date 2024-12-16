1.  我们的解释器会仔细检查传递给函数的参数数是否与其期望的参数数一致。由于这种检查是在每次调用的运行时进行的，因此会产生性能代价。而Smalltalk实现则不存在这个问题。为什么没有？

    `在Smalltalk中，方法的选择器（selector）包含了参数的信息，这使得在编译时就能确定方法调用的正确性，从而避免了运行时对参数数量进行检查。`

    Smalltalk使用关键字选择器（keyword selectors），`方法名称本身就包含了参数的位置和数量。例如，insert:at:方法的名称已经明确表示它需要两个参数`。由于方法名称与参数数量紧密结合，调用时必须严格匹配相应的方法签名，否则编译器无法识别该方法，从而避免了错误的参数数量传递。

    Smalltalk has different call syntax for different arities. To define a
    method that takes multiple arguments, you use **keyword selectors**. Each
    argument has a piece of the method name preceding instead of using commas
    as a separator. For example, a method like:

    ```lox
    list.insert("element", 2)
    ```

    To insert "element" as index 2 would look like this in Smalltalk:

    ```smalltalk
    list insert: "element" at: 2
    ```

    Smalltalk doesn't use a dot to separate method name from receiver. More
    interestingly, the "insert:" and "at:" parts both form a single method
    call whose full name is "insert:at:". Since the selectors and the colons
    that separate them form part of the method's name, there's no way to call
    it with the wrong number of arguments. You can't pass too many or two few
    arguments to "insert:at:" because there would be no way to write that call
    while still actually naming that method.

2.  在 Lox 中添加匿名函数语法
    Lox 的函数声明语法执行两个独立的操作：它不仅创建了一个函数，还将其与名称绑定。
    `匿名函数(lambdas)：一种创建函数而不绑定名称的表达式语法。`

    **说明**

    - AST 修改：通过在 Expr 和 Stmt 中添加 Function 节点，使得匿名函数和命名函数都能被表示和处理。
      - `函数表达式：Expr.Function 节点表示匿名函数，包含参数列表和函数体。`
      - `函数语句：Stmt.Function 节点表示命名函数，包含函数名和 Expr.Function 节点。`
    - 解析器更新：`引入了新的方法 functionBody 来解析函数体`，无论是匿名函数还是命名函数。通过检查 fun 后是否跟有标识符来区分命名函数和匿名函数。
    - 解释器适配：在解释器中分别处理函数声明语句和函数表达式，确保匿名函数没有名称，而命名函数有名称，且在打印时能够正确显示。
    - 匿名函数的使用：通过 visitFunctionExpr 方法，匿名函数可以作为表达式被使用，如传递给其他函数或赋值给变量。

    **注意事项**

    - 参数数量限制：在解析函数参数时，限制最多不超过8个参数，以避免过于复杂的函数签名。
    - 语法冲突处理：通过 checkNext 方法，确保`只有在 fun 后跟有标识符时才解析为命名函数，否则解析为匿名函数`，避免语法冲突和解析错误。

    `function 表达式的ast`

    This requires juggling some code around. In GenerateAst, we need a node
    for function expressions. In the defineAst() call for Expr, add:

    ```java
    "Function : List<Token> parameters, List<Stmt> body",
    ```

    `function 语句的ast`

    While we're at it, we can reuse that for function statements. A function
    _statement_ is now just a name and a function expression:

    ```java
    "Function   : Token name, Expr.Function function",
    ```

    `LoxFunction 类需要存储 Expr.Function，以处理命名函数和匿名函数。`

    Over in LoxFunction, it will store an Expr.Function instead of a statement
    to handle both types. If the function does have a name, that is tracked
    separately, since lambdas won't have one:

    ```java
    class LoxFunction implements Callable {
      private final String name;
      private final Expr.Function declaration;
      private final Environment closure;

      LoxFunction(String name, Expr.Function declaration, Environment closure) {
        this.name = name;
        this.closure = closure;
        this.declaration = declaration;
      }
      @Override
      public String toString() {
        if (name == null) return "<fn>";
        return "<fn " + name + ">";
      }

      // ...
    }
    ```

    The parser changes are a little more complex. We move the logic to handle
    anonymous functions into a new method. Then the method to handle named
    functions becomes wrapper around that one:

    ```java
    private Stmt.Function function(String kind) {
      Token name = consume(IDENTIFIER, "Expect " + kind + " name.");
      return new Stmt.Function(name, functionBody(kind));
    }

    private Expr.Function functionBody(String kind) {
      consume(LEFT_PAREN, "Expect '(' after " + kind + " name.");
      List<Token> parameters = new ArrayList<>();
      if (!check(RIGHT_PAREN)) {
        do {
          if (parameters.size() >= 8) {
            error(peek(), "Can't have more than 8 parameters.");
          }

          parameters.add(consume(IDENTIFIER, "Expect parameter name."));
        } while (match(COMMA));
      }
      consume(RIGHT_PAREN, "Expect ')' after parameters.");

      consume(LEFT_BRACE, "Expect '{' before " + kind + " body.");
      List<Stmt> body = block();
      return new Expr.Function(parameters, body);
    }
    ```

    Now we can use `functionBody()` to parse lambdas. In `primary()`, add
    another clause:

    ```java
    if (match(FUN)) return functionBody("function");
    ```

    We've got one nasty little problem. We want lambdas to be a valid primary
    expression, and in theory any primary expression is allowed in a primary
    statement. But if you try to do:

    ```lox
    fun () {};
    ```

    Then the `declaration()` parser will match that `fun` and try to parse it
    as a named function declaration statement. It won't see a name and will
    report a parse error. Even though the above code is pointless, we want it
    to work to avoid a weird edge case in the grammar.

    To handle that, we only want to parse a function declaration if the current
    token is `fun` and the one past that is an identifier. That requires another
    token of lookahead, as we add:

    ```java
    private boolean checkNext(TokenType tokenType) {
      if (isAtEnd()) return false;
      if (tokens.get(current + 1).type == EOF) return false;
      return tokens.get(current + 1).type == tokenType;
    }
    ```

    Then, in `declaration()`, change the `match(FUN)) ...` line to:

    ```java
    if (check(FUN) && checkNext(IDENTIFIER)) {
      consume(FUN, null);
      return function("function");
    }
    ```

    `分别处理函数声明语句和函数表达式，确保匿名函数没有名称，而命名函数有名称，且在打印时能够正确显示。`

    ````java

    Now only a function with a name is parsed as such.

    Then our interpreter needs to handle both cases:

    ```java

    @Override
    public Void visitFunctionStmt(Stmt.Function stmt) {
        String fnName = stmt.name.lexeme;
        environment.define(fnName, new LoxFunction(fnName, stmt.function, environment));
        return null;
    }

    @Override
    public Object visitFunctionExpr(Expr.Function expr) {
        return new LoxFunction(null, expr, environment);
    }
    ````

    We could have re-used visitFunctionExpr but that would lose the function name if someone were to print it, this ensures we preserve it.

    ```lox
    fun whichFn(fn) {
      print fn;
    }

    whichFn(fun (b) {
     print b;
    });

    fun named(a) { print a; }
    whichFn(named);
    //
    // <fn>
    // <fn named>
    ```

3.  这个程序有效吗？

    ```lox
    fun scope(a) {
      var a = "local";
    }
    ```

    换句话说，`函数的参数是与其局部变量在同一作用域，还是在外层作用域？Lox 是怎么做的？你熟悉的其他语言呢？你认为语言应该怎么做？`

    No, it isn't. Lox uses the same scope for the parameters and local variables
    immediately inside the body. That's why Stmt.Function stores the body as a
    list of statements, not a single Stmt.Block that would create its own
    nested scope separate from the parameters.
    `Lox 对函数参数和函数体内部的局部变量使用相同的作用域`。这就是为什么 Stmt.Function 将函数体存储为语句列表，而不是单个 Stmt.Block，后者会创建一个与参数分开的嵌套作用域的原因。

    In Java, it's an error because you aren't allowed to shadow local variables
    inside a method or collide them.
    It's an error in C because parameters and locals share the same scope.
    It is allowed in Dart. There, parameters are in a separate scope surrounding
    the function body.

    I'm not a fan of Dart's choice. I think shadowing should be allowed in
    general because it helps ensure changes to code are encapsulated and don't
    affect parts of the program unrelated to the change. (See this design note
    for more: http://craftinginterpreters.com/statements-and-state.html#design-note).

    But shadowing still usually leads to more confusing code, so it should be
    avoided when possible. The only thing putting parameters in an outer scope
    allows is shadowing those parameters, but I think any code that did that
    would be _very_ hard to read. I would rather prohibit that outright.
    综合考虑，`我倾向于禁止在同一作用域内使用相同的变量名`，尤其是在函数参数和局部变量之间，以提高代码的可维护性和可读性。
