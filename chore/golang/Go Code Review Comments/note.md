https://go.dev/wiki/CodeReviewComments

- Interfaces (接口哲学): 这是 Go 设计哲学中非常独特的一点。
  **接口应由消费者定义，而非生产者**。**一个包应该返回具体的类型（struct 或 pointer），而不是接口**。这样，`生产者可以自由地添加新方法而不会破坏下游的实现`。消费者则根据自己的需求，定义自己需要的最小化接口。这被称为“接口分离原则”的 Go 化实践，促进了真正的解耦。
  **“如果只有一个实现，就不要使用接口”**。在没有实际需求前，不要过早抽象。
- 在代码上运行 gofmt，可以自动修复大部分机械式的风格问题。
- 注释应以被描述对象的名称开头，并以句号结尾
  这不仅仅是为了好看，更是因为 godoc 工具会提取这些注释生成 API 文档。规范的注释直接转化为高质量的文档。

```go
// Request represents a request to run a command.
type Request struct { ...

// Encode writes the JSON encoding of req to w.
func Encode(w io.Writer, req *Request) { ...
```

- 一个与具体请求无关的函数可以使用 context.Background()，但`即便你认为不需要，也应倾向于传递一个 Context`。默认情况下应传递 Context；只有在你有充分理由认为替代做法是错误时，才直接使用 context.Background()。
- **不要在结构体类型中添加 Context 成员**；相反，应在该类型上每个需要传递 Context 的方法中添加一个 ctx 参数。唯一的例外是方法签名必须与标准库或第三方库中的接口匹配的情况。
- 一个 nil 切片会被编码为 null ，而 []string{} 会被编码为 JSON 数组 []
- 不要在常规错误处理中使用 panic。使用 error 和多返回值。
- 错误字符串`不应大写（除非以专有名词或首字母缩略词开头）或以标点结尾`，因为它们通常在其他上下文之后被打印
- `不要使用 _ 变量丢弃错误`。如果函数返回错误，请检查它以确保函数成功。处理该错误、返回它，或在真正的例外情况下引发 panic。
- 当你生成 goroutine 时，要明确它们何时——或是否——会退出。
- In-Band Errors
  一般来说，Go 代码应通过返回额外的值来表示错误。
- 尽量让正常的代码路径`保持最小缩进`，并将错误处理缩进、优先处理。
- 测试：actual != expected

---

好的，我们来深入分析一下 Go 官方 Wiki 上的《Go Code Review Comments》这份文档。

这份文档是 Go 团队多年代码审查实践的结晶，可以看作是《Effective Go》的重要补充。它并非一套强制性的编码标准，而是一系列旨在提升代码**可读性、健壮性和可维护性**的最佳实践建议。其核心思想与 Go 语言本身的设计哲学一脉相承：**简洁、明确、高效**。

我将为你提炼和解读其中的关键部分，并将其整理成一份更详尽的学习笔记。

### 分析与解读

#### 1. 格式与命名 (Code's Appearance and Clarity)

这部分规范旨在通过统一的视觉风格和命名约定，降低阅读和理解代码的认知成本。

- **`gofmt` / `goimports`**: 这是 Go 语言的“铁律”。`gofmt` 统一代码格式，杜绝了关于空格、缩进的无谓争论。`goimports` 在此基础上还能自动管理 `import` 语句，是现代 Go 开发的必备工具。**核心价值在于，所有人的代码看起来都一样，开发者可以专注于逻辑而非风格。**

- **`Comment Sentences` (注释即文档)**: 要求注释是完整的句子，以被描述的实体名开头，以句号结尾。这不仅仅是为了好看，更是因为 `godoc` 工具会提取这些注释生成 API 文档。规范的注释直接转化为高质量的文档。

- **`Initialisms` (缩略词的大小写)**: `URL`, `ID`, `HTTP` 等缩略词应保持大小写一致（全大写 `URL` 或全小写 `url`），而不是 `Url` 或 `Id`。例如，`ServeHTTP` 而非 `ServeHttp`，`appID` 而非 `appId`。这保证了命名的一致性和专业性。

- **`Package Names` (包名)**: 包名是调用者使用你的代码时的前缀，因此应**简短、清晰、有意义**。避免使用 `util`, `common`, `helper` 等模糊的名称。好的包名能自解释其功能，例如 `http`, `json`, `time`。

- **`Variable Names` (变量名)**: Go 推崇短变量名，特别是对于作用域小的局部变量（如 `i` for index, `r` for reader）。规则是：**变量的作用域越大，其名称应越具描述性**。方法接收器（receiver）通常用一到两个字母（如 `c` for `Client`），而全局变量则需要一个完整的描述性名称。

#### 2. 错误处理与控制流 (Robustness and Readability)

Go 将错误处理视为一等公民，相关规范体现了对代码健壮性的高度重视。

- **`Handle Errors` (处理错误)**: **绝不使用 `_` 丢弃错误**。这是 Go 最核心的信条之一。每个返回 `error` 的地方都必须检查，然后进行处理、返回，或者在真正无法恢复的情况下 `panic`。忽略错误是万恶之源。

- **`Don't Panic` (不要滥用 Panic)**: `panic` 用于处理程序无法继续运行的、真正的异常情况（如空指针解引用、数组越界）。对于可预期的错误（如文件未找到、网络超时），应使用 `error` 返回值来处理。`panic/recover` 机制更像是其他语言中的 `try/catch`，但 Go 文化中极少用它来控制业务逻辑。

- **`Error Strings` (错误字符串格式)**: 错误信息应为小写字母开头，且不带标点符号。这是因为错误信息通常会被嵌入到其他日志或错误上下文中，例如 `log.Printf("failed to open file %s: %v", name, err)`。小写的 `err` 信息可以自然地融入句子中。

- **`Indent Error Flow` (错误处理优先)**: 优先处理错误情况，并提前返回（`return`, `continue`），让“快乐路径”（happy path）的代码保持最小的缩进。这使得主干逻辑非常清晰，一目了然。

  ```go
  // 推荐
  func process(file string) error {
      f, err := os.Open(file)
      if err != nil {
          return fmt.Errorf("opening file: %w", err)
      }
      defer f.Close()

      // ... 正常逻辑 ...
      return nil
  }

  // 不推荐
  func process(file string) error {
      f, err := os.Open(file)
      if err == nil {
          defer f.Close()
          // ... 正常逻辑 ...
      } else {
          return fmt.Errorf("opening file: %w", err)
      }
      return nil
  }
  ```

#### 3. API 设计与接口 (API Design and Interfaces)

这部分关乎如何设计出优雅、解耦且易于使用的 API。

- **`Contexts`**: `context.Context` 是现代 Go 服务中传递请求范围数据（如 deadline、cancellation、trace ID）的标准方式。

  - 它应该是函数的**第一个参数**，命名为 `ctx`。
  - **不要**将其作为结构体的字段。应将其作为方法参数显式传递。这使得依赖关系更清晰，并避免了 `Context` 的生命周期与结构体混淆。

- **`Interfaces` (接口哲学)**: 这是 Go 设计哲学中非常独特的一点。

  - **接口应由消费者定义，而非生产者**。一个包应该返回具体的类型（struct 或 pointer），而不是接口。这样，生产者可以自由地添加新方法而不会破坏下游的实现。消费者则根据自己的需求，定义自己需要的最小化接口。这被称为“接口分离原则”的 Go 化实践，促进了真正的解耦。
  - **“如果只有一个实现，就不要使用接口”**。在没有实际需求前，不要过早抽象。

- **`Receiver Type` (接收器类型)**:

  - **何时用指针接收器 (`*T`)**:
    1.  方法需要修改接收器（状态变更）。
    2.  接收器是包含 `sync.Mutex` 等同步原语的结构体（避免锁被复制）。
    3.  接收器是大型结构体或数组（避免大对象拷贝，提高性能）。
  - **何时用值接收器 (`T`)**:
    1.  接收器是 `map`, `func`, `chan`（它们本身就是引用类型）。
    2.  接收器是 `slice` 且方法不修改其长度或容量。
    3.  接收器是小型的、不变的结构体或基本类型。
  - **核心原则**: **保持一致性**。一个类型的所有方法，其接收器类型应该统一。**如果不确定，就用指针接收器**。

- **`Pass Values` (传递值)**: 不要仅仅为了节省几个字节而传递指针。如果函数内部总是对参数 `x` 使用 `*x`，那么参数类型就不应该是 `*T`。例如，`*string` 或 `*io.Reader` 几乎总是不必要的。值传递使得数据流更清晰，减少了意外的副作用。

### 总结与笔记生成

基于以上深入分析，我为你生成了一份更全面的 Markdown 笔记，你可以将其保存下来用于后续学习。

````markdown
# Go 官方代码审查指南深度解析

本文档是对 Go 官方 Wiki [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) 的深入解读和总结。这些是 Go 社区广泛遵循的最佳实践，旨在提升代码的可读性、健壮性和可维护性。

## 格式与命名 (Appearance and Clarity)

### 1. `gofmt` / `goimports`

- **规则**: 提交代码前必须使用 `gofmt` (或 `goimports`) 格式化。
- **原因**: 消除所有关于代码风格的争论，让所有 Go 代码看起来都一样，使开发者能专注于业务逻辑。

### 2. `Comment Sentences` (注释即文档)

- **规则**: 公共 API 的注释应是完整句子，以被描述的实体名开头，以句号结尾。
- **原因**: `godoc` 会提取这些注释生成 API 文档。规范的注释等于高质量的文档。
- **示例**:

  ```go
  // Package http provides HTTP client and server implementations.
  package http

  // Get issues a GET to the specified URL.
  func Get(url string) (*Response, error) { ... }
  ```

### 3. `Initialisms` (缩略词大小写)

- **规则**: 像 `URL`, `ID`, `HTTP` 这样的缩略词，其大小写应保持一致（全大写或全小写），而不是驼峰式。
- **原因**: 保持命名一致性和专业性。
- **示例**: `ServeHTTP`, `appID`, `urlPony` (而非 `ServeHttp`, `appId`, `UrlPony`)。

### 4. `Package Names` (包名)

- **规则**: 包名应简短、清晰、有意义。避免使用 `util`, `common`, `helpers` 等模糊名称。
- **原因**: 包名是调用者使用代码时的前缀，好的包名能自解释其功能。

### 5. `Variable Names` (变量名)

- **规则**: 变量作用域越小，名字可以越短。作用域越大，名字应越具描述性。
- **原因**: 减少局部范围内的认知负担，同时让全局或大范围的变量易于理解。
- **示例**:
  - 循环变量: `for i, v := range s`
  - 方法接收器: `func (c *Client) Send(...)`
  - 全局变量: `var defaultTransport http.RoundTripper`

## 错误处理与控制流 (Robustness and Readability)

### 6. `Handle Errors` (必须处理错误)

- **规则**: **绝不**使用 `_` 丢弃 `error` 返回值。必须检查并处理。
- **原因**: 这是 Go 保证程序健壮性的核心机制。忽略错误是未定义行为和 bug 的主要来源。

### 7. `Don't Panic` (不要滥用 Panic)

- **规则**: `panic` 只用于程序无法恢复的异常状态。常规业务错误应使用 `error` 返回值。
- **原因**: `panic` 会中断正常的控制流，滥用会使代码难以理解和维护。

### 8. `Error Strings` (错误字符串格式)

- **规则**: 错误信息字符串应以小写字母开头，且不带结尾标点。
- **原因**: 错误信息通常被嵌入到其他日志上下文中，小写开头能自然地融入句子。
- **示例**: `fmt.Errorf("something bad happened")`，使用时 `log.Printf("request failed: %v", err)`。

### 9. `Indent Error Flow` (错误处理优先，保持主路清晰)

- **规则**: 优先处理错误情况并提前返回，让正常逻辑路径保持最小的缩进。
- **原因**: 极大地提高了代码的可读性，使主干逻辑一目了然。

## API 设计与接口 (API Design and Interfaces)

### 10. `Contexts`

- **规则**:
  1.  作为函数/方法的**第一个参数**，命名为 `ctx`。
  2.  **不要**将 `context.Context` 作为结构体的字段。
- **原因**: 显式传递 `Context` 使请求的生命周期和依赖关系清晰可见。将其嵌入结构体会导致生命周期混乱。

### 11. `Interfaces` (接口哲学)

- **规则**:
  1.  **接口由消费者定义，而非生产者**。生产者应返回具体类型。
  2.  在没有实际需求（例如，多个实现或需要 mock）之前，不要创建接口。
- **原因**: 促进真正的解耦。消费者只关心它需要的方法，生产者可以自由演进。

### 12. `Receiver Type` (接收器类型)

- **规则**:
  - **指针接收器 `*T`**: 当方法需要修改接收器、接收器含锁、或接收器是大型结构体时使用。
  - **值接收器 `T`**: 当接收器是 map/chan/func/slice、或小型不可变结构体时使用。
  - **关键**: **保持一致性**。如果不确定，**优先使用指针接收器**。
- **原因**: 明确方法的副作用（是否修改对象），并兼顾性能和正确性。

### 13. `Pass Values` (传递值而非指针)

- **规则**: 不要仅仅为了节省几个字节而传递指针。如果函数总是解引用来使用参数，那么它应该接收值。
- **原因**: 值传递使数据流更清晰，减少了因意外修改共享数据而导致的 bug。`*string`, `*int` 等通常是坏味道。
````
