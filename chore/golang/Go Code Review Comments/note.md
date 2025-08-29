https://go.dev/wiki/CodeReviewComments

- 在代码上运行 gofmt，可以自动修复大部分机械式的风格问题。
- 注释应以被描述对象的名称开头，并以句号结尾

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
- 错误字符串不应大写（除非以专有名词或首字母缩略词开头）或以标点结尾，因为它们通常在其他上下文之后被打印
- `不要使用 _ 变量丢弃错误`。如果函数返回错误，请检查它以确保函数成功。处理该错误、返回它，或在真正的例外情况下引发 panic。
- In-Band Errors
