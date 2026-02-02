# 更快的异步函数与 Promise：深入底层性能革命

# Faster async functions and promises: Deep technical dive

这篇 V8 博客是异步性能优化的行业标杆，详细解释了如何通过微调规范来提升巨量性能。

### 1. 核心矛盾：Promise 的“三倍开销”

在旧规范中，一个 `await` 表示三个 Promise 对象和三个任务队列循环（Microticks）。

1. 创建第一个中间 Promise。
2. 调用 `promiseResolve`。
3. 返回结果给外部等待者。

### 2. 优化方案：缩减 Microtick

V8 优化了等待逻辑：如果被 `await` 的对象本身就是一个原生 Promise 实例，V8 直接将后续任务挂载到该 Promise 的解析链上，跳过了不必要的封装。
这一改动将 `await` 的开销从 3 次 Microtick 减少到 1 次。

### 3. 避免堆分配

通过对异步状态机的精细控制，V8 减少了闭包和 Promise 对象的分配。这意味着更少的内存抖动和更轻的 GC 负担。

### 4. 一针见血的见解

`async`/`await` 的性能之争在 v7.3/v7.4 附近尘埃落定。V8 通过重构告诉大家：原生的 async 语法不仅比 Promise.then() 好看，而且它现在是 JS 异步执行的极致速度之选，任何手动模拟的异步逻辑都难以望其项背。
