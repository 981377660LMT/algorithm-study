1. Python 和 Haskell 的词法语法并不规则(not regular)。这是什么意思？
   `Both of them have significant indentation.` To handle that, the scanner
   emits synthetic "{" and "}" tokens (or "indent" and "dedent" as Python
   calls them), as if there were explicit delimiters for each level of
   indentation.

   In order to know when a new line begins or ends one of more levels of
   indentation, the scanner has to track the _previous_ indentation value.
   That state has to be stored in the scanner, which means it has a little bit
   of _memory_. That makes it no longer a regular language, which is defined
   to only need to `store a single finite number identifying` which state it's
   in.

   You _could_ make a regular language for significant indentation, by having
   a hardcoded limit to the maximum amount of indentation, but that starts to
   split semantic hairs around the Chomsky hierarchy.

2. 空格确实会影响 CoffeeScript、Ruby 和 C 预处理器中代码的解析方式。在这些语言中，空格在哪些地方会产生影响？

   - 在 CoffeeScript、Ruby 中，空格的存在与否会影响函数调用的解析方式
   - 在 C 预处理器中，空格帮助区分简单宏和带参数的函数式宏，确保宏被正确展开。

   In CoffeeScript, parentheses are option in function calls. You can call a
   function like:

   ```coffeescript
   doStuff withThis
   ```

   Also, there is a nice syntax for lambda functions:

   ```coffeescript
   () -> someLambda
   -> anotherOne
   ```

   On the second line, you can see that you can omit the `()` if the lambda
   takes no parameters. So what does this do:

   ```coffeescript
   someFunction () -> someLambda
   ```

   Does it call `someFunction` with zero parameters and then call the result of
   _that_ with one parameter, a lambda? Or does it call `someFunction` with
   one parameter, the lambda? The answer depends on spaces:

   ```coffeescript
   someFunction() -> someLambda
   # Means the same as:
   someFunction()(() -> someLambda)

   someFunction () -> someLambda
   # Means the same as:
   someFunction(() -> someLambda)
   ```

   Ruby has similar corner cases because it also allow omitting the parentheses
   in method calls (which is where CoffeeScript gets it from).

   The C preprocessor relies on spaces to distinguish function macros from
   simple macros:

   ```c
   #define MACRO1 (p) (p)
   #define MACRO2(p) (p)
   ```

   Here, `MACRO1` is a simple text macro that expands to `(p) (p)` when used.
   `MACRO2(p)` is a function-like macro that takes a parameter and expands to
   `(p)` with `p` replaced by the parameter.

3. 和大多数扫描器一样，我们的扫描器也会丢弃注释和空白，因为解析器并不需要它们。为什么要编写一个不丢弃这些内容的扫描器？它对什么有用？
   Programmers often write `"doc comments"` above their functions and types. A
   `documentation generator` or an IDE that shows help text for declarations
   needs access to those comments, so a scanner for those should include them.
   程序员通常在他们的函数和类型上方写“文档注释”。文档生成器或显示声明帮助文本的集成开发环境需要访问这些注释，因此扫描器应该包括它们。

   An automated code formatter obviously needs to preserve comments and may
   want to be aware of the original whitespace if some of the author's
   formatting should be preserved.
   一个自动化代码格式化工具显然需要保留注释，并且可能希望了解原始空白，以便保留作者的一些格式。

   尽管在许多情况下，注释和空白字符对语法解析并不必要，因此扫描器通常会丢弃它们，但在以下几个重要的应用场景中，保留这些内容显得尤为重要：

   - 文档生成：自动提取和生成详细的API文档。
   - 集成开发环境：提供实时的代码提示和帮助信息。
   - 代码格式化：保持代码的可读性和一致性。
   - 代码重构和静态分析：确保代码逻辑和文档的一致性，提升代码质量。
   - 逆向工程和代码恢复：帮助理解和维护复杂的代码结构。

4. 为 Lox 扫描仪添加对 C 风格 /_ ... _/ 块注释的支持。确保能处理其中的换行符。考虑允许嵌套。添加嵌套支持是否比你预期的要费事？为什么？
   You can see where I've implemented them for a similar language here:

   https://github.com/munificent/wren/blob/c6eb0be99014d34085e2d24c696aed449e2fb171/src/vm/wren_compiler.c#L663

   The interesting part is the `nesting` variable. Like challenge #1, we
   `require some extra state to track the nesting`, which makes this not quite
   regular.

   Note also that we need to handle an unterminated block comment.

   为了支持嵌套注释，需要在扫描器中维护一个嵌套计数器。每次遇到 /_ 时，计数器加一；每次遇到 _/ 时，计数器减一。当计数器归零时，表示注释块结束。

   ```C
      void Scanner::skipBlockComment() {
        int nesting = 1; // 初始嵌套层级

        while (!isAtEnd() && nesting > 0) {
            if (peek() == '/' && peekNext() == '*') {
                advance(); // 跳过 '/'
                advance(); // 跳过 '*'
                nesting++;
            }
            else if (peek() == '*' && peekNext() == '/') {
                advance(); // 跳过 '*'
                advance(); // 跳过 '/'
                nesting--;
            }
            else {
                if (peek() == '\n') {
                    line++;
                }
                advance();
            }
        }

        if (nesting > 0) {
            Lox::error(line, "Unterminated block comment.");
        }
    }
   ```

   正则表达式（Regex）和有限状态自动机（Finite State Automaton, FSA）只能处理正规语言（Regular Languages），即那些不需要记忆前文状态的语言。然而，嵌套注释需要扫描器记住当前的嵌套层级，这超出了正规语言的处理能力。
