这是一个手写的简易版 `interrupt` 机制实现。它模拟了 LangGraph 的核心原理：**异常抛出 (Exception Throwing)** + **状态重放 (Replay)**。

这个实现包含三个核心部分：
1.  **`GraphInterrupt` 异常**：用于携带数据并跳出当前执行栈。
2.  **`Runner` (运行器)**：负责捕获异常、保存状态、并在恢复时重新运行函数。
3.  **`interrupt` 函数**：用户调用的接口。

### 简易版 Interrupt 实现 (TypeScript)

```typescript
// 1. 定义特殊的异常类，用于携带中断数据
class GraphInterrupt extends Error {
  public payload: any;
  constructor(payload: any) {
    super("GraphInterrupt");
    this.payload = payload;
  }
}

// 2. 定义运行上下文，用于管理恢复值
class Context {
  // 存储恢复值的队列
  private resumeValues: any[] = [];
  // 当前消耗到第几个恢复值了 (索引)
  private pointer: number = 0;

  constructor(resumeValues: any[] = []) {
    this.resumeValues = resumeValues;
  }

  // 尝试获取恢复值
  // 如果有值，返回该值 (恢复模式)
  // 如果没值，返回 undefined (正常执行模式)
  nextResumeValue(): { hasValue: boolean; value?: any } {
    if (this.pointer < this.resumeValues.length) {
      return { hasValue: true, value: this.resumeValues[this.pointer++] };
    }
    return { hasValue: false };
  }
}

// 全局上下文 (简化演示用，实际库中会通过 AsyncLocalStorage 或参数传递)
let globalContext: Context | null = null;

// 3. 核心 interrupt 函数
export function interrupt(payload: any): any {
  if (!globalContext) {
    throw new Error("interrupt must be called within a Runner");
  }

  // 检查是否有恢复值 (Replay 逻辑)
  const { hasValue, value } = globalContext.nextResumeValue();

  if (hasValue) {
    console.log(`[System] 发现恢复值，跳过中断，直接返回: ${JSON.stringify(value)}`);
    return value; // 直接返回之前传入的 resume 值，不抛异常
  }

  // 没有恢复值，说明是第一次运行到这里 -> 抛出异常暂停
  console.log(`[System] 触发中断，抛出异常，Payload: ${JSON.stringify(payload)}`);
  throw new GraphInterrupt(payload);
}

// 4. 运行器 (Runner) - 模拟图的执行引擎
class Runner {
  // 模拟持久化存储 (Checkpointer)
  // 存储结构: thread_id -> resume_values[]
  private stateStore: Map<string, any[]> = new Map();

  async run(
    fn: () => Promise<any>, 
    threadId: string, 
    resumeValue?: any
  ): Promise<any> {
    console.log(`\n--- 开始运行 (Thread: ${threadId}) ---`);

    // 1. 加载历史状态
    let history = this.stateStore.get(threadId) || [];
    
    // 2. 如果有新的 resumeValue，追加到历史记录中
    if (resumeValue !== undefined) {
      history.push(resumeValue);
      this.stateStore.set(threadId, history); // 更新存储
    }

    // 3. 初始化上下文 (准备 Replay)
    globalContext = new Context(history);

    try {
      // 4. 执行用户函数
      // 如果历史中有值，interrupt 会直接返回而不抛异常
      const result = await fn();
      
      console.log("[System] 执行完成");
      globalContext = null; // 清理
      return { status: "completed", result };

    } catch (e) {
      globalContext = null; // 清理

      // 5. 捕获中断异常
      if (e instanceof GraphInterrupt) {
        console.log("[System] 捕获到中断请求，保存状态并暂停");
        return { 
          status: "interrupted", 
          interrupt_payload: e.payload 
        };
      }
      
      // 其他错误正常抛出
      throw e;
    }
  }
}

// ==========================================
// 测试代码
// ==========================================

// 用户定义的业务逻辑
async function myWorkflow() {
  console.log("Step 1: 开始工作...");
  
  // 第一个中断点
  const name = interrupt({ question: "What is your name?" });
  console.log(`Step 2: 你好, ${name}`);

  // 第二个中断点
  const age = interrupt({ question: "How old are you?" });
  console.log(`Step 3: 你今年 ${age} 岁了`);

  return `User ${name} is ${age} years old.`;
}

// 模拟执行过程
async function main() {
  const runner = new Runner();
  const threadId = "thread-123";

  // --- 第一次运行 ---
  // 遇到第一个 interrupt，抛出异常暂停
  const res1 = await runner.run(myWorkflow, threadId);
  console.log("Result 1:", res1); 
  // Output: status: "interrupted", payload: "What is your name?"

  // --- 第二次运行 (恢复) ---
  // 传入 "Alice"。
  // Replay: 
  // 1. 遇到第一个 interrupt -> 历史里有 "Alice" -> 直接返回 "Alice"
  // 2. 遇到第二个 interrupt -> 历史里没值 -> 抛出异常暂停
  const res2 = await runner.run(myWorkflow, threadId, "Alice");
  console.log("Result 2:", res2);
  // Output: status: "interrupted", payload: "How old are you?"

  // --- 第三次运行 (恢复) ---
  // 传入 "25"。
  // Replay:
  // 1. 遇到第一个 interrupt -> 历史里有 "Alice" -> 返回 "Alice"
  // 2. 遇到第二个 interrupt -> 历史里有 "25" -> 返回 "25"
  // 3. 函数执行完毕
  const res3 = await runner.run(myWorkflow, threadId, "25");
  console.log("Result 3:", res3);
  // Output: status: "completed", result: "User Alice is 25 years old."
}

main().catch(console.error);
```

### 关键点解析

1.  **`Context` 类中的 `pointer`**：
    这就是为什么 LangGraph 依赖**索引匹配**的原因。每次 Replay 时，`pointer` 重置为 0。
    *   第一次调用 `interrupt`，取 `resumeValues[0]`。
    *   第二次调用 `interrupt`，取 `resumeValues[1]`。
    *   如果 `pointer` 超出了数组长度，说明这是新的中断点，抛出异常。

2.  **`Runner` 中的 `history`**：
    这就是 Checkpointer 的作用。它保存了**所有**过去的中断恢复值。
    *   第一次运行：`history = []`
    *   第二次运行：`history = ["Alice"]`
    *   第三次运行：`history = ["Alice", "25"]`

3.  **无副作用原则**：
    你可以看到，在第三次运行时，`console.log("Step 1: 开始工作...")` 会被打印**三次**。这就是为什么在 `interrupt` 之前的代码必须是幂等的（或者无副作用的），否则会被重复执行。