## 1 写下这些函数被调用的跟踪。显示它们被调用的顺序，哪个调用了哪个，以及传递给它们的参数。

```
(-1 + 2) * 3 - -4
```

It's:

```
expression
| parsePrecedence(PREC_ASSIGNMENT)
| | grouping
| | | expression
| | | | parsePrecedence(PREC_ASSIGNMENT)
| | | | | unary // for "-"
| | | | | | parsePrecedence(PREC_UNARY)
| | | | | | | number
| | | | | binary // for "+"
| | | | | | parsePrecedence(PREC_FACTOR) // PREC_TERM + 1
| | | | | | | number
| | binary // for "*"
| | | parsePrecedence(PREC_UNARY) // PREC_FACTOR + 1
| | | | number
| | binary // for "-"
| | | parsePrecedence(PREC_FACTOR) // PREC_TERM + 1
| | | | unary // for "-"
| | | | | parsePrecedence(PREC_UNARY)
| | | | | | number
```

## 2 TOKEN_MINUS 的 ParseRule 规则同时具有前缀和中缀函数回调。这是因为 - 既是前缀运算符（单目取反），也是中缀运算符（减法）。在完整的 Lox 语言中，还有哪些其他标记可以在前缀和中缀位置使用？在 C 语言或您选择的其他语言中呢？

- Lox only has one other: `left parenthesis` is used as a prefix expression for grouping, and as an infix expression for invoking a function.
  lox 中还有另外一个兼具前缀和中缀位置的标记：`左括号`，作为前缀表达式时表示分组，作为中缀表达式时表示调用。

- Several languages allow `+` as a prefix unary operator as a parallel to `-` and then also of course use infix `+` for addition.
  有的语言允许 `+` 作为前缀一元运算符，与 - 类似，然后当然也使用中缀 + 进行加法运算。

- A number of languages use square brackets for list or array literals, which makes `[` a prefix expression and then also use square brackets as a subscript operator to access elements from a list.
  许多语言使用方括号 [ 来表示`列表或数组字面量`，这使得 [ 成为前缀表达式，然后也使用方括号作为下标运算符从列表中`访问元素`。
- C uses `*` as a prefix operator to dereference a pointer and as infix for multiplication. Likewise, `&` is a prefix address-of operator and infix bitwise and.
  C 语言使用 \* 作为前缀运算符来解引用指针，以及作为中缀运算符进行乘法。同样，& 是前缀的取地址运算符和中缀的按位与运算符。
- `*` and `&` aren't prefix _expressions_ in Ruby, but they can appear in prefix position before an argument in an argument list.
  在 Ruby 中，\* 和 & 不是前缀 表达式，但它们可以出现在参数列表中的前缀位置。

## 3 支持三元表达式（ternary expression）

为编译器添加对该运算符的支持。您不必生成任何字节码，只需展示如何将其连接到解析器并处理操作数。

The `?:` operator has lower precedence than almost anything, so we add a new `PREC_CONDITIONAL` level between `PREC_ASSIGN` and `PREC_OR`.
`?: 运算符的优先级比几乎所有其他运算符都低，因此我们在 PREC_ASSIGN 和 PREC_OR 之间添加一个新的 PREC_CONDITIONAL 级别。`

I'll skip adding the new TokenType enums for `?` and `:`. That part is pretty obvious.
In the new row in the table for the `?` token type, we call:

```c
static void conditional()
{
  // Compile the then branch.
  parsePrecedence(compiler, PREC_CONDITIONAL);

  consume(compiler, TOKEN_COLON,
          "Expect ':' after then branch of conditional operator.");

  // Compile the else branch.
  parsePrecedence(compiler, PREC_ASSIGNMENT);
}
```

Of course a full implementation needs more code to actually do the conditional evaluation, but that should compile the operands with the right precedence. Note that the precedence of the operands is a little unusual. The precedence of the last operand is _lower_ than the conditional expression itself.
`在解析三元运算符时，最后一个操作数需要以较低的优先级进行解析`

That might be surprising, but it's how C rolls.
这可能令人惊讶，但这就是 C 语言的做法。
