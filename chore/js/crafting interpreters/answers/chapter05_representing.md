1.  前面我说过，我们添加到语法元语法中的 | 、 \* 和 + 形式只是句法糖。就拿这个语法来说吧：

    ```text
    expr → expr ( "(" ( expr ( "," expr )* )? ")" | "." IDENTIFIER )+
         | IDENTIFIER
         | NUMBER
    ```

    生成与相同语言相匹配的语法，但不使用任何符号糖。这段语法编码了哪种表达方式？

    There are a few ways to do it. Here is one:

    ```text
    expr → expr calls
    expr → IDENTIFIER
    expr → NUMBER

    calls → calls call
    calls → call

    call → "(" ")"
    call → "(" arguments ")"
    call → "." IDENTIFIER

    arguments → expr
    arguments → arguments "," expr
    ```

    It's the syntax for a `function invocation`.

2.  访问者模式可以让你在面向对象的语言中模拟函数式风格。为函数式语言设计一种补充模式。它应能让你将一种类型的所有操作捆绑在一起，并能让你轻松定义新类型。
    One way is to create a record or tuple containing a function pointer for
    each `operation`. In order to allow defining new types and passing them to
    existing code, these functions need to encapsulate the type entirely -- the
    existing code isn't aware of it, so it can't type check. You can do that by
    having the functions be closures that all close over the same shared object,
    "this", basically.

3.  在逆波兰表达式（RPN）中，算术运算符的操作数都放在运算符之前，因此 1 + 2 变成 1 2 + 。运算从左到右进行。数字被推入隐式堆栈。算术运算符弹出最上面的两个数字，执行运算，然后推入结果。
    为我们的语法树类定义一个访问者类，该类接收表达式，将其转换为 RPN，并返回结果字符串。

    ```java
    class RpnPrinter implements Expr.Visitor<String> {
      String print(Expr expr) {
        return expr.accept(this);
      }

      @Override
      public String visitBinaryExpr(Expr.Binary expr) {
        return expr.left.accept(this) + " " +
               expr.right.accept(this) + " " +
               expr.operator.lexeme;
      }

      @Override
      public String visitGroupingExpr(Expr.Grouping expr) {
        return expr.expression.accept(this);
      }

      @Override
      public String visitLiteralExpr(Expr.Literal expr) {
        return expr.value.toString();
      }

      @Override
      public String visitUnaryExpr(Expr.Unary expr) {
        String operator = expr.operator.lexeme;
        if (expr.operator.type == TokenType.MINUS) {
          // Can't use same symbol for unary and binary.
          operator = "~";
        }

        return expr.right.accept(this) + " " + operator;
      }

      public static void main(String[] args) {
        Expr expression = new Expr.Binary(
            new Expr.Unary(
                new Token(TokenType.MINUS, "-", null, 1),
                new Expr.Literal(123)),
            new Token(TokenType.STAR, "*", null, 1),
            new Expr.Grouping(
                new Expr.Literal("str")));

        System.out.println(new RpnPrinter().print(expression));
      }
    }
    ```

    **在RPN中，同一个减号符号 - 可以表示两种不同的操作：一元减号和二元减号。当扫描器遇到减号 - 时，无法直接判断它是一元运算符还是二元运算符，因此不知道应该弹出栈中的一个操作数还是两个操作数。为了消除歧义，采用不同的符号来表示一元和二元减号是一种有效的方法。**

    Note that we have to handle unary "-" specially. In RPN, we can't use the
    same symbol for both unary and binary forms. When we encounter it, we
    wouldn't know whether to pop one or two numbers off the stack. So, to
    disambiguate, we pick a different symbol for negation.
