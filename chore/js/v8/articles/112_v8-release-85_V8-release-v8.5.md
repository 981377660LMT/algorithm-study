# V8 v8.5 发布：Promise 整合与 Wasm 多返回值

# V8 release v8.5: Promise.any and Wasm Multi-value

V8 v8.5 完善了 JS 异步模型并扩展了 WebAssembly 的函数表达力。

### 1. 异步基石：`Promise.any`

引入了 `Promise.any` 和 `AggregateError`。这补齐了 Promise 组合模式的最后一块拼图——即“只要有一个成功就返回”，特别适合处理多点并发的服务请求。

### 2. 标准增强：`String.prototype.replaceAll`

由于原生支持，开发者终于不再被强迫使用带 `/g` 的正则表达式来替换所有子串，代码的可读性和运行安全性得到了提升。

### 3. Wasm 重大升级：多返回值

WebAssembly 现在支持函数返回多个值。这一语义上的对齐，极大地简化了从 C++ 或 Go 等支持多返回值的语言编译到 Wasm 时的代码生成复杂度，也减少了由于模拟多返回值带来的栈操作损耗。

### 4. 一针见血的见解

v8.5 更多地是在提升“表达的效率”。通过提供更直观、更原生的 API，V8 让开发者能够以更少的代码行数表达更复杂的并行逻辑（JS）或计算流程（Wasm），从而减少了业务代码本身的脆弱性。
