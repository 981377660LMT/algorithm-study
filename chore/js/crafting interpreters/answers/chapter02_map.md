## 1. Find the various parts in an open source implementation.

选择一种你喜欢的语言的开源实现。下载源代码，并在其中探索。试着找到实现扫描器和解析器的代码。它们是手写的，还是使用 Lex 和 Yacc 等工具生成的？

---

选择一种开源语言的实现并探索其源代码是一个非常有趣且有益的学习过程。下面，我将以 **Go 语言**（也称为 Golang）为例，介绍其开源实现中的扫描器（Scanner）和解析器（Parser）的实现方式。

### 1. 获取源代码

Go 语言的官方实现是开源的，托管在 [GitHub 上的 Go 仓库](https://github.com/golang/go)。你可以通过以下命令克隆仓库：

```bash
git clone https://github.com/golang/go.git
```

### 2. 探索扫描器和解析器

在 Go 的源代码中，扫描器和解析器的实现主要位于 `src/go/parser` 和 `src/go/scanner` 目录下。

#### 扫描器（Scanner）

**位置：** `src/go/scanner`

Go 的扫描器负责将源代码转换为一系列的标记（tokens）。扫描器的实现是**手写的**，没有使用像 Lex 这样的词法分析生成工具。具体来说，`scanner/scanner.go` 文件中定义了扫描器的逻辑，包括识别关键字、标识符、字面量等。

**特点：**

- **手写实现：** 扫描器代码是纯手写的，采用了递归和状态机的方式来处理不同的字符和模式。
- **高效性：** 手写的扫描器在性能和内存使用上进行了优化，以适应 Go 语言对编译速度的高要求。

#### 解析器（Parser）

**位置：** `src/go/parser`

Go 的解析器负责将扫描器生成的标记序列转换为抽象语法树（AST）。解析器同样是**手写的**，没有使用 Yacc 或类似的语法分析工具。`parser/parser.go` 文件中详细描述了解析器的实现，包括语法规则和错误处理机制。

**特点：**

- **手写实现：** 解析器采用递归下降的方式实现，每个语法规则对应一个解析函数。
- **易于维护：** 手写的解析器使得代码结构清晰，便于开发者理解和维护。
- **错误恢复机制：** 解析器实现了健壮的错误恢复机制，能够在遇到语法错误时尽量继续解析，以提供更全面的错误报告。

### 3. 手写 vs. 生成工具

Go 语言的扫描器和解析器选择了**手写实现**，这在现代编程语言的编译器设计中并不少见。手写的优势包括：

- **更高的控制力：** 开发者可以针对特定的需求和优化进行细致的调整。
- **更好的错误处理：** 手写解析器可以实现更复杂和智能的错误恢复和报告机制。
- **性能优化：** 可以针对特定的语言特性进行专门的性能优化，而不受生成工具生成代码的限制。

尽管使用 Lex/Yacc 等工具可以加快开发速度，但对于像 Go 这样追求高性能和简洁性的语言，实现手写的扫描器和解析器更符合其设计理念。

### 4. 总结

通过对 Go 语言开源实现的探索，可以看到其扫描器和解析器都是**手写的**，这为语言的高效性和灵活性提供了保障。如果你对编译器设计感兴趣，深入研究 Go 的扫描器和解析器代码将是一个很好的学习机会。

希望这个示例对你有所帮助！如果你对其他语言的实现感兴趣，欢迎继续讨论。

## 2. Why not use a JIT?

即时编译往往是实现动态类型语言的最快方法，但并非所有语言都使用即时编译。不使用 JIT 的原因是什么？

1. 实现复杂性
   It's really complex to implement, debug, and maintain. Few people have the skill to do it.
2. 可移植性
   Like a native code compiler (which it is), it ties you to a specific CPU architecture.
   JIT 编译器通常需要针对特定的硬件架构和操作系统进行优化，这限制了语言在不同平台上的可移植性
3. 字节码通常比机器码更紧凑，占用更少的内存。在平台上
   像嵌入式设备一样，内存可能比速度更重要。
   Bytecode is generally more compact than machine code (since it's closer to the semantics of the language), so it takes up less memory. In platforms like embedded devices where memory may matter more than speed, that can be a worthwhile trade-off.
4. 安全性
   Some platforms, like iOS and most game consoles, expressly disallow executing code generated at runtime. The OS simply won't allow you to jump into memory that can be written to.
   一些平台，如 iOS 和大多数游戏控制台，明确禁止在运行时执行生成的代码
   代码注入风险：JIT 需要在运行时生成可执行代码，这可能增加代码注入攻击的风险。如果不加以严格控制，恶意代码可能利用 JIT 编译器的漏洞执行任意代码。

---

Go 语言：Go 选择使用 Ahead-Of-Time 编译，而不是 JIT。这一选择使得 Go 拥有极快的编译速度和高效的运行时性能，适合构建大型分布式系统和服务器应用。同时，避免了 JIT 带来的额外复杂性和资源消耗。
Python：Python 的主要实现（CPython）使用解释执行，而不是 JIT。虽然存在如 PyPy 这样的 JIT 实现，但由于 CPython 的广泛使用和生态系统的庞大，主要实现选择了稳定和简单的解释模式。然而，PyPy 的成功也表明，JIT 在某些情况下能够显著提升性能。

## 3. Why do Lisp compilers also contain an interpreter?

大多数编译为 C 语言的 Lisp 实现都包含一个解释器，可以让它们在运行中执行 Lisp 代码。为什么？

大多数 Lisp 实现支持`宏 —— 在编译时执行的代码`，因此**在编译过程中需要能够评估宏本身**。你可以通过“编译”宏然后运行它来实现这一点，但这会带来大量的开销。
