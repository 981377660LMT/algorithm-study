1.  添加对逗号表达式的支持。编写语法，然后实现必要的解析代码。
    The comma operator has the `lowest precedence`, so it goes between expression
    and equality:

    ```ebnf
    expression → comma ;
    comma      → equality ( "," equality )* ;
    equality   → comparison ( ( "!=" | "==" ) comparison )* ;
    comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
    term       → factor ( ( "-" | "+" ) factor )* ;
    factor     → unary ( ( "/" | "*" ) unary )* ;
    unary      → ( "!" | "-" | "--" | "++" ) unary
               | postfix ;
    postfix    → primary ( "--" | ++" )* ;
    primary    → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
    ```

    We could define a new syntax tree node by adding this to the `defineAst()`
    call:

    ```java
    "Comma    : Expr left, Expr right",
    ```

    But a simpler choice is to treat it like any other binary operator and
    reuse Expr.Binary.

    Parsing is similar to other infix operators (except that we don't bother to
    keep the operator token):

    ```java
    private Expr expression() {
      return comma();
    }

    private Expr comma() {
      Expr expr = equality();

      while (match(COMMA)) {
        Token operator = previous();
        Expr right = equality();
        expr = new Expr.Binary(expr, operator, right);
      }

      return expr;
    }
    ```

    `处理函数调用中的逗号歧义`
    `解决方案：调整参数解析的优先级`

    Keep in mind that commas are already used in the grammar to separate
    arguments in function calls. With the above change, this:

    ```lox
    foo(1, 2)
    ```

    Now gets parsed as:

    ```lox
    foo((1, 2))
    ```

    In other words, pass a single argument to `foo`, the result of evaluating
    `1, 2`. That's not what we want. To fix that, we simply change the way we
    parse function arguments to require a higher precedence expression than the
    comma operator:
    `在解析函数调用的参数时，应该使用比 comma 更高优先级的表达式解析函数，如 equality()。由于 equality() 的优先级高于 comma()，这意味着在函数参数列表中，逗号将被视为参数分隔符，而不是运算符的一部分。`

    ```java
    if (!check(RIGHT_PAREN)) {
      do {
        if (arguments.size() >= 8) {
          error(peek(), "Can't have more than 8 arguments.");
        }
        arguments.add(equality()); // <-- was expression().
      } while (match(COMMA));
    }
    ```

2.  添加对 C 风格条件运算符或 三元运算符 ?: 的支持。 ? 和 : 之间允许什么优先级？整个运算符是左关联还是右关联？

    - 三元运算符具有`较低的优先级，仅高于赋值运算符`。这意味着它会在大多数其他运算符（如算术运算符、比较运算符等）之后进行求值
    - 三元运算符是`右结合`的。这意味着在多个三元运算符连用时，会从右向左进行解析
      a ? b = c : d 应被解析为 a ? (b = c) : d 而不是 (a ? b : d) = c
      a ? b : c ? d : e 应被解析为 a ? b : (c ? d : e) 而不是 (a ? b : c) ? d : e

    ```ebnf
    expression  → conditional ;
    conditional → equality ( "?" expression ":" conditional )? ;
    // Other rules...
    ```

    The precedence of the operands is pretty interesting. The left operand has
    higher precedence than the others, and the middle operand has lower
    precedence than the condition expression itself. That allows:

          a ? b = c : d

    Again, I won't bother showing the scanner and token changes since they're
    pretty obvious.

    ```java
    private Expr expression() {
      return conditional();
    }

    private Expr conditional() {
      Expr expr = equality();  // 最低优先级

      if (match(QUESTION)) {
        Expr thenBranch = expression();
        consume(COLON,
            "Expect ':' after then branch of conditional expression.");
        Expr elseBranch = conditional(); // 右结合
        expr = new Expr.Conditional(expr, thenBranch, elseBranch);
      }

      return expr;
    }
    ```

3.  如果一个二元运算符出现在表达式的开头，如 + a - b + c - d，解析器会无法识别缺失的左操作数，这应被视为一个语法错误。为了提高解析器的容错能力，您希望在遇到这种错误时：

    - 报告错误：通知开发者缺少左操作数。
    - 继续解析：丢弃错误运算符，并解析其右操作数，以便继续解析后续表达式。

    为此，我们需要改进语法。`在 primary 规则中添加错误产生式，允许解析器在遇到缺失左操作数的运算符时进行错误处理。`

    Here's an updated grammar. The grammar itself doesn't "know" that some of
    these productions are errors. The parser handles that.

    ```ebnf
    expression → equality ;
    equality   → comparison ( ( "!=" | "==" ) comparison )* ;
    comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
    term       → factor ( ( "-" | "+" ) factor )* ;
    factor     → unary ( ( "/" | "*" ) unary )* ;
    unary      → ( "!" | "-" | "--" | "++" ) unary
               | postfix ;
    postfix    → primary ( "--" | ++" )* ;
    primary    → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
               // Error productions...
               | ( "!=" | "==" ) equality
               | ( ">" | ">=" | "<" | "<=" ) comparison
               | ( "+" ) term
               | ( "/" | "*" ) factor ;
    ```

    - 注意`-`不是错误产生式，因为它是一个有效的前缀表达式。
      Note that "-" isn't an error production because that _is_ a valid prefix
      expression.

    - 错误产生式（Error Productions）与正常产生式的区别:
      `错误产生式的右侧操作数规则与运算符本身的优先级相同`，因此解析器能够继续处理同级别的运算符，维护整体语法结构的完整性。
      而在标准的中缀运算符规则中，`操作数（operand）的语法规则优先级比运算符本身的优先级高一级`。为了处理多个相同优先级的运算符连续出现的情况，语法规则明确允许这些运算符和操作数重复出现。

      With the normal infix productions, the operand non-terminals are one
      precedence level higher than the operator's own precedence. In order to
      handle a series of operators of the same precedence, the rules explicitly
      allow repetition.

      With the error productions, though, the right-hand operand rule is the same
      precedence level. That will effectively strip off the erroneous leading
      operator and then consume a series of infix uses of operators at the same
      level by reusing the existing correct rule. For example:

      ```lox
      + a - b + c - d
      ```

      The error production for `+` will match the leading `+` and then use
      `term` to also match the rest of the expression.

      ```java
      private Expr primary() {
        if (match(FALSE)) return new Expr.Literal(false);
        if (match(TRUE)) return new Expr.Literal(true);
        if (match(NIL)) return new Expr.Literal(null);

        if (match(NUMBER, STRING)) {
          return new Expr.Literal(previous().literal);
        }

        if (match(LEFT_PAREN)) {
          Expr expr = expression();
          consume(RIGHT_PAREN, "Expect ')' after expression.");
          return new Expr.Grouping(expr);
        }

        // Error productions.
        if (match(BANG_EQUAL, EQUAL_EQUAL)) {
          error(previous(), "Missing left-hand operand.");
          equality();
          return null;
        }
        if (match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)) {
          error(previous(), "Missing left-hand operand.");
          comparison();
          return null;
        }
        if (match(PLUS)) {
          error(previous(), "Missing left-hand operand.");
          term();
          return null;
        }
        if (match(SLASH, STAR)) {
          error(previous(), "Missing left-hand operand.");
          factor();
          return null;
        }

        throw error(peek(), "Expect expression.");
      }
      ```

      这种方法确保了即使在表达式开头存在语法错误，解析器仍能继续处理后续的部分，减少了因单个错误导致的整体解析失败。
