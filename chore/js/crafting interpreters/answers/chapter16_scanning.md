## 1 考虑查看其他支持字符串模版（string interpolation）的语言实现，以了解它们是如何处理的。

I've implemented this in another language, Wren. You can see the code here:
https://github.com/munificent/wren/blob/8fae8e4f1e490888e2cc9b2ea6b8e0d0ff9dd60f/src/vm/wren_compiler.c#L118-L130
Poke around in that file for "interp" to see everything. The basic idea is you have two token types. TOKEN_STRING is for uninterpolated string literals, and the last segment of an interpolated string. Every piece of a string literal that precedes an interpolated expression uses a different TOKEN_INTERPOLATION type.
**基本思想是你有两种 token 类型。`TOKEN_STRING` 用于`未插值`的字符串字面量，以及插值字符串的最后一个部分。**
**每个`紧随插值表达式之前的字符串`字面量部分都使用不同的` TOKEN_INTERPOLATION` 类型。**

This:

```lox
"Tea will be ready in ${steep + cool} minutes."
```

Gets scanned like:

```text
TOKEN_INTERPOLATION "Tea will be ready in"
TOKEN_IDENTIFIER    "steep"
TOKEN_PLUS          "+"
TOKEN_IDENTIFIER    "cool"
TOKEN_STRING        "minutes."
```

(The interpolation delimiters themselves are discarded.)
`模版符号本身被丢弃。`

And this:

```lox
"Nested ${"interpolation?! Are you ${"mad?!"}"}"
```

Scans as:

```text
TOKEN_INTERPOLATION "Nested "
TOKEN_INTERPOLATION "interpolation?! Are you "
TOKEN_STRING        "mad?!"
TOKEN_STRING        ""
TOKEN_STRING        ""
```

The two empty TOKEN_STRING tokens are because the interpolation appears at the
very end of the string. They tell the parser that they've reached the end of
the interpolated expression.
最后两个空的 TOKEN_STRING token 是因为插值出现在字符串的末尾。`它们告诉解析器已经到达插值表达式的结束。`

## 2 如何处理泛型尖括号(angle brackets)和位移运算符的二义性？

几种语言使用尖括号表示泛型，并且也有一个 >> 右移运算符。这导致了 C++早期版本中的一个经典的编译错误：

```cpp
std::vector<std::vector<int>> nestedVector;
```

`用户被迫通过在闭合尖括号之间放置空格来避免这种情况。`
后来的 C++ 版本更智能，可以处理上述代码。Java 和 C# 从未遇到过这个问题。这些语言是如何指定和实现这一点的？

---

As far as I can tell, Java and C# don't actually specify it correctly. Unless
the verbiage is hidden away somewhere in the specs, I believe that this:
`据我所知，Java 和 C# 实际上并没有正确地在规范中指定它。`

```java
List<List<String>> nestedList;
```

Should technically by a syntax error in a fully spec-compliant implementation of Java or C#. However, all practical implementations don't follow the letter of the spec and instead do what users want.

C++, as of C++0x, does actually specify this:
`从 C++0x 开始，C++ 确实对这一点进行了规范：`

http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2005/n1757.html

It states that if a `<` has been scanned and no closing `>` has been scanned
yet, and there are no other intervening bracket characters, then a subsequent
`>>` is scanned as two `>` tokens instead of a single shift.
`文档指出，如果已经扫描到了一个 <，且尚未扫描到闭合的 >，并且中间没有其他括号字符，那么随后的 >> 会被当作两个 > token而不是一个单一的移位操作符来扫描。`

As far as implementation, I think javac handles this by scanning the `>>` as a
single shift token. When the parser is looking for a `>` to close the type
argument, if it sees a shift token, it splits it into two `>` tokens right then,
consumes the first, and then keeps parsing.
在实现方面，我认为 javac 通过将 >> 作为一个单一的移位操作符来扫描。
`当解析器在寻找关闭类型参数的 > 时，如果它遇到了一个移位操作符，它会立即将其分割成两个 > token，消耗第一个，然后继续解析。`

Microsoft's C# parser takes the opposite approach. It always scans `>>` as two
separate `>` tokens. Then, when parsing an expression, if it sees two `>` tokens
next to each other with no whitespace between them, it parses them as a shift
operator.
微软的 C# 解析器则采取了相反的方法。它始终将 >> 解析为两个独立的 > token。
`然后，在解析表达式时，如果它看到两个没有空格分隔的 > token，它会将其解析为一个移位操作符。`

## 3 列举一些其他语言中的上下文关键字（contextual keywords），以及它们有意义的上下文。拥有上下文关键字的优缺点是什么？如果需要，你将如何在你语言的前端实现它们？

### 上下文关键字是什么？

上下文关键字（Contextual Keywords）是在编程语言中，根据特定的上下文环境赋予特定含义的关键字。
与`保留关键字（Reserved Keywords）不同`，保留关键字在语言的任何地方都具有固定的含义，不能用作标识符（如变量名、函数名等）。
**而上下文关键字只有在特定的语法结构或上下文中才具有特殊意义，在其他地方则可以作为普通标识符使用。**
例如：

1. async 和 await 在方法声明和调用中具有特殊意义，但在其他上下文中（如变量命名）则可以作为普通标识符 => `await` 可以作为变量名使用 `const await = 5`;
2. Python 中的 match，只有在模式匹配语法中才有特殊含义，其他地方可以作为普通标识符使用。

### 答案

I don't generally like contextual keywords. It's fairly easy to write a real
parser that can handle them gracefully, but:
我通常不喜欢上下文关键字。虽然编写一个能够优雅处理它们的真正解析器相对容易，但：

- **用户常常感到困惑**
  Users are often confused by them. Many programmers don't even realize that
  contextual keywords exist. They assume all identifiers are either fully
  reserved by the language or fully available for use.
- **一旦一个标识符在某个上下文中成为关键字，它很快就会被读者理解为关键字的含义**
  Once an identifier becomes a keyword in some context, it quickly takes on
  that meaning to readers and becomes _very_ confusing if you use it for your
  own name outside of that context. Now that C# has async/await, you will
  just anger your fellow C# users if you name a variable `await` in some
  non-async method because they are so used to seeing `await` used for its
  keyword meaning.

  So even though it's _technically_ usable elsewhere, it's effectively fully
  reserved.

That being said, sometimes you have no other option. Once your language is in
wide use, reserving a new keyword is a breaking change to any code that was
previously using that name. If you can only reserve it inside a new context that
didn't previously exist (for example, async functions in C#), or in a context
where an identifier can't appear, then you can reserve it only in that context
and be confident that you didn't break any previous code.
尽管如此，有时候你别无选择。一旦你的语言被广泛使用，**保留一个新关键字会对之前使用该名称的任何代码造成破坏性更改。**
如果你只能在一个以前不存在的新上下文中保留它（例如，C# 中的异步函数），或者在标识符不能出现的上下文中保留它，那么你可以**仅在该上下文中保留它，并且可以确信不会破坏任何之前的代码。**

So they're sort of an inevitable compromise when evolving a language over time.
因此，在随着时间推移而演变语言时，上下文关键字是一种不可避免的折衷。

Implementing them is pretty easy. The scanner scans them like regular
identifiers, since it doesn't generally know the surrounding context. In the
parser, you recognize the keyword in that context by looking for an identifier
token and checking to see if its lexeme is the right string.
实现它们相当容易。扫描器会像处理普通标识符一样扫描它们，因为它通常不知道周围的上下文。
在解析器中，`你通过查找标识符令牌并检查其词素是否为正确的字符串，在特定的上下文中识别关键字。`
