# V8 与 Node.js 的故事：十年协同进化

## The story of V8 and Node.js

- **Original Link**: [https://v8.dev/blog/v8-nodejs](https://v8.dev/blog/v8-nodejs)
- **Publication Date**: 2019-11-20
- **Summary**: 本文回顾了 V8 与 Node.js 之间长达十年的共生关系。从 Ryan Dahl 最初选择 V8 到 Node.js 推动 V8 改进发布工程流程，这种合作重新定义了服务器端 JavaScript 的地位。

---

### 1. 为什么 Ryan Dahl 选择 V8？

2ten 年前，Node.js 诞生时，Ryan Dahl 比较了 SpiderMonkey 和 V8。

- **原因**：V8 虽年轻，但其代码非常整洁，且 API 设计允许清晰地处理多个并发上下文（Isolates）。V8 的 C++ 风格与 Node.js 追求的轻量 I/O 十分配套。

### 2. 长久以来的痛点：发布节奏

- **浏览器 VS 服务器**：Chrome 每 6 周更新一次，用户无感知。但 Node.js 是服务器组件，需要长期稳定支持（LTS）。
- **ABI/API 稳定性**：V8 曾经频繁修改 API 导致 Node.js 模块生态大面积崩溃。

### 3. V8 的回应：Release Engineering

为了支持 Node.js，V8 做出了以下调整：

- **Lts 分支支持**：V8 团队开始与 Node.js 核心团队协作，在 Node.js 即将发布的版本中提前测试 V8 的新特性。
- **Embedder 接口稳定性**：引入了更稳定的 C++ 导出接口，减少了 Node.js 升级 V8 时的工程量。
- **性能关注点转移**：Node.js 对垃圾回收的停顿极其敏感，这促使 V8 加大了对增量回收和并发回收的投入。

### 4. 关键里程碑

- **async/await 的落地**：V8 针对 Node.js 社区对异步编程的极度渴求，深入优化了 Promise 的执行速度。
- **64 位支持**：早期 Node.js 对大内存的支持需求，直接推动了 V8 在 64 位系统上的成熟。

### 5. 结论

Node.js 不是简单的把 V8 搬到了服务器，它是 V8 演进的重要拉动力。没有 Node.js，V8 可能只是一个优秀的浏览器组件；没有 V8，Node.js 可能无法实现其高性能异步 I/O 的愿景。
