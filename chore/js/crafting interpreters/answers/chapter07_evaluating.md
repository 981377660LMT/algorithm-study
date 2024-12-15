1.  您会扩展 Lox 以支持比较其他类型吗？如果是，您允许哪些类型对，如何定义它们的排序？请说明您的选择并与其他语言进行比较。
    `Python 3` allows comparing all of the various number types with each other,
    except for complex numbers. Booleans (True and False) are a subclass of
    int and work like 1 and 0 for comparison.

    Strings can be compared with each other and are ordered lexicographically.
    Likewise other sequences.

    Comparing sets is defined in terms of subsets and supersets, so that, for
    example `{1, 2} < {1, 2, 3}`. `This isn't a total order` since many pairs of
    sets are neither subsets nor supersets of each other.

    I think it would be reasonable to extend Lox to support comparing strings
    with each other. I wouldn't support comparing other built in types, nor
    mixing them. Allowing `"1" < 2` is a recipe for confusion.

2.  许多语言定义 + 时，如果其中一个操作数是字符串，另一个操作数会被转换为字符串，然后将结果连接起来。例如， "scone" + 4 将产生 scone4 。扩展 visitBinaryExpr() 中的代码，以支持该功能。
    Replace the Token.PLUS case with:

    ```java
    case PLUS:
      if (left instanceof String || right instanceof String) {
        return stringify(left) + stringify(right);
      }

      if (left instanceof Double && right instanceof Double) {
        return (double)left + (double)right;
      }

      throw new RuntimeError(expr.operator,
          "Operands must be two numbers or two strings.");
    ```

3.  如果你用一个数字除以零，现在会发生什么？你认为应该发生什么？请说明您的选择。您所知道的其他语言是如何处理除数为零的除法的？更改 visitBinaryExpr() 中的实现，以检测并报告这种情况下的运行时错误。
    It returns Infinity, -Infinity, or NaN based on sign of the dividend. Giventhat Lox is a high level scripting language, I think it would be better to
    raise a runtime error to let the user know something got weird. That's what Python and Ruby do.

    On the other hand, given that Lox gives the user no way to catch and handle runtime errors, not throwing one might be more flexible.

    - 当前 Lox 中除以零的行为
      当一个数字被零除时，结果会根据被除数的符号返回 Infinity、-Infinity 或 NaN（非数值）。具体行为如下：

      - 正数除以零：返回 Infinity
      - 负数除以零：返回 -Infinity
      - 零除以零：返回 NaN

    - 建议的改进：抛出运行时错误
      与 Python 和 Ruby 等语言的行为一致，增强语言的一致性和可靠性。
      Python：

      - 整数除以零：抛出 ZeroDivisionError。
      - 浮点数除以零：返回 inf/-inf/nan。

      ```java
      public class Interpreter implements Expr.Visitor<Object>, Stmt.Visitor<Void> {
          // ...

          @Override
          public Object visitBinaryExpr(Binary expr) {
              try {
                  Object left = evaluate(expr.left);
                  Object right = evaluate(expr.right);

                  switch (expr.operator.type) {
                      // 其他运算符处理
                      case SLASH:
                          checkNumberOperands(expr.operator, left, right);
                          double divisor = (double)right;
                          if (divisor == 0) {
                              throw new RuntimeError(expr.operator, "Division by zero.");
                          }
                          return (double)left / divisor;
                      // 其他运算符处理
                  }
              } catch (RuntimeError error) {
                  Lox.runtimeError(error);
              }

              return null;
          }

          // ...
      }

      ```
