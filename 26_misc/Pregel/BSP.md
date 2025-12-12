**BSP (Bulk Synchronous Parallel，整体同步并行)** 模型

在 LangGraph/Pregel 的上下文中，它的核心在于**状态更新的“屏障”（Barrier）机制**。
即：所有的计算节点（Actors）在同一“步”中并行读取的是**旧状态**，计算产生的新数据不会立即生效，而是必须等到所有节点都运行完毕后，在“更新阶段”统一写入。这避免了并发读写导致的数据竞争和不确定性。

以下是一个简化的 TypeScript 抽象，展示了 **Plan (计划) -> Execution (执行) -> Update (更新)** 的核心循环：

```typescript
// 1. 定义通道 (Channel)：负责存储状态和处理状态更新（Reducer）
// 类似于 LangGraph 中的 Channels
class Channel<T> {
  private value: T | undefined
  private reducer: (current: T | undefined, update: T) => T

  constructor(reducer: (current: T | undefined, update: T) => T) {
    this.reducer = reducer
  }

  // 获取当前步骤的稳定值
  get(): T | undefined {
    return this.value
  }

  // 应用更新（在 Update 阶段调用）
  applyUpdate(update: T): void {
    this.value = this.reducer(this.value, update)
  }
}

// 2. 定义节点 (Actor/Node)：执行逻辑
interface Node {
  id: string
  subscribeTo: string[] // 订阅的通道 key
  writeTo: string[] // 写入的通道 key
  // 执行函数：接收输入，返回要写入的数据
  execute: (inputs: Record<string, any>) => Promise<Record<string, any>>
}

// 3. 定义运行时 (Runtime)：实现 BSP 循环
class BSPRuntime {
  private channels: Map<string, Channel<any>> = new Map()
  private nodes: Node[] = []
  private stepCount = 0

  constructor(channels: Record<string, Channel<any>>, nodes: Node[]) {
    Object.entries(channels).forEach(([k, v]) => this.channels.set(k, v))
    this.nodes = nodes
  }

  // 核心：执行单个 Superstep
  async step(): Promise<boolean> {
    this.stepCount++
    console.log(`--- Step ${this.stepCount} ---`)

    // --- Phase 1: Plan (计划) ---
    // 找出需要运行的节点（这里简化为：只要订阅的通道有值就运行）
    const activeNodes = this.nodes.filter(node => {
      return node.subscribeTo.some(chKey => this.channels.get(chKey)?.get() !== undefined)
    })

    if (activeNodes.length === 0) {
      return false // 没有节点需要运行，终止
    }

    // 暂存所有节点的输出，模拟“写入对其他 Actors 不可见”
    // 这就是 BSP 的“同步屏障”
    const pendingUpdates: Array<{ channelKey: string; value: any }> = []

    // --- Phase 2: Execution (执行) ---
    // 并行执行所有活跃节点
    const executions = activeNodes.map(async node => {
      // 1. 读取当前状态 (Snapshot)
      const inputs: Record<string, any> = {}
      node.subscribeTo.forEach(key => {
        inputs[key] = this.channels.get(key)?.get()
      })

      // 2. 运行逻辑
      const outputs = await node.execute(inputs)

      // 3. 收集输出（注意：此时不修改 Channel）
      Object.entries(outputs).forEach(([key, val]) => {
        if (node.writeTo.includes(key)) {
          pendingUpdates.push({ channelKey: key, value: val })
        }
      })
    })

    await Promise.all(executions) // 等待所有节点完成

    // --- Phase 3: Update (更新) ---
    // 统一应用所有更改
    if (pendingUpdates.length === 0) return false

    for (const update of pendingUpdates) {
      const channel = this.channels.get(update.channelKey)
      if (channel) {
        channel.applyUpdate(update.value)
        console.log(`Updated channel [${update.channelKey}] with value:`, channel.get())
      }
    }

    return true // 继续下一个循环
  }

  // 运行直到收敛
  async run() {
    let keepRunning = true
    while (keepRunning && this.stepCount < 10) {
      // 防止死循环
      keepRunning = await this.step()
    }
    console.log('Execution finished.')
  }
}

// --- 使用示例 ---

// 定义一个简单的 LastValue reducer
const lastValueReducer = (_: any, update: any) => update

// 1. 初始化通道
const channels = {
  input: new Channel<number>(lastValueReducer),
  middle: new Channel<number>(lastValueReducer),
  output: new Channel<number>(lastValueReducer)
}

// 2. 初始化节点
const nodes: Node[] = [
  {
    id: 'NodeA',
    subscribeTo: ['input'],
    writeTo: ['middle'],
    execute: async inputs => {
      console.log('NodeA running...')
      return { middle: inputs['input'] + 1 }
    }
  },
  {
    id: 'NodeB',
    subscribeTo: ['middle'],
    writeTo: ['output'],
    execute: async inputs => {
      console.log('NodeB running...')
      return { output: inputs['middle'] * 2 }
    }
  }
]

// 3. 运行
const app = new BSPRuntime(channels, nodes)

// 注入初始状态
channels['input'].applyUpdate(10)

app.run()
```

### 代码解析

1.  **隔离性 (Isolation)**: 在 `Phase 2` 中，节点读取的是 `inputs`，产生的是 `pendingUpdates`。即使 Node A 写入了 `middle` 通道，Node B 在**同一个 Step** 中也读不到这个新值。
2.  **同步屏障 (Barrier)**: `await Promise.all(executions)` 确保了所有并行任务都结束后，才进入 `Phase 3`。
3.  **原子更新 (Atomic Update)**: 在 `Phase 3` 中，系统统一遍历 `pendingUpdates` 并调用 Channel 的 Reducer 更新状态。这就是文档中提到的“统一将 Actors 写入的数据应用到 Channels 中”。

---

这种算法的底层本质是 **BSP (Bulk Synchronous Parallel，整体同步并行)** 模型的一个变体，或者更抽象地说，是一个 **基于同步屏障（Barrier）的离散时间步进系统**。

剥离掉 LangGraph 的业务概念（节点、通道），其最底层的数学抽象是：**$S_{t+1} = Update(S_t, Execute(Plan(S_t)))$**。

这是一个**“读写分离”**的状态机：在时间步 $T$ 内，所有计算都基于 $S_T$（只读快照），产生的变更 $\Delta$ 只有在时间步结束时才合并，形成 $S_{T+1}$。

以下是完全独立于业务逻辑的 TypeScript 抽象：

```typescript
/**
 * 泛型定义：
 * S = State (系统状态)
 * D = Delta (状态变更/增量)
 */
abstract class BarrierStepLoop<S, D> {
  constructor(protected state: S) {}

  /**
   * 阶段 1: Plan (计划)
   * 观察当前状态，决定这一步谁需要“动”。
   * 返回一组待执行的“纯计算任务”。
   */
  protected abstract plan(currentState: S): Array<() => Promise<D>>

  /**
   * 阶段 3: Update (更新)
   * 将收集到的所有变更（Deltas）合并到状态中，推进时间步。
   * 这是一个“归约”过程。
   */
  protected abstract update(currentState: S, deltas: D[]): S

  /**
   * 核心循环：运行直到收敛 (Run until convergence)
   */
  public async run(maxSteps: number = 100): Promise<S> {
    let step = 0

    while (step < maxSteps) {
      console.log(`--- Step ${step} ---`)

      // 1. Plan: 基于当前状态快照，生成任务列表
      const tasks = this.plan(this.state)

      // 收敛检测：如果没有任务需要执行，说明系统稳定，退出循环
      if (tasks.length === 0) {
        console.log('Converged (No more tasks).')
        break
      }

      // 2. Execution: 并行执行 (同步屏障)
      // 关键点：并行执行期间，任务无法修改 this.state，只能返回 delta
      const deltas = await Promise.all(tasks.map(task => task()))

      // 3. Update: 原子性状态推进
      // 只有在这里，时间才从 T 变为 T+1
      this.state = this.update(this.state, deltas)

      step++
    }

    return this.state
  }
}

// --- 示例实现：康威生命游戏 (Game of Life) 的一维简化版 ---
// 规则：如果一个细胞的值 < 10，它会给自己和邻居 +1。直到所有细胞 >= 10。

type CellState = number[] // 状态是数组
type CellDelta = { index: number; add: number }[] // 变更是索引和增量

class SimpleSimulation extends BarrierStepLoop<CellState, CellDelta> {
  // 1. Plan: 找出所有值 < 10 的细胞
  protected plan(currentState: CellState): Array<() => Promise<CellDelta>> {
    const tasks: Array<() => Promise<CellDelta>> = []

    currentState.forEach((val, idx) => {
      if (val < 10) {
        // 创建一个闭包任务，捕获当前的 idx 和 val
        tasks.push(async () => {
          // 模拟计算延迟
          return [
            { index: idx, add: 1 }, // 给自己加 1
            { index: idx + 1, add: 1 } // 给右边邻居加 1
          ]
        })
      }
    })

    return tasks
  }

  // 3. Update: 汇总所有加法操作
  protected update(currentState: CellState, deltasList: CellDelta[]): CellState {
    // 复制一份新状态（保持不可变性原则，虽然这里可以直接改，但为了演示）
    const newState = [...currentState]

    // 扁平化所有变更并应用
    deltasList.flat().forEach(change => {
      if (change.index < newState.length) {
        newState[change.index] += change.add
      }
    })

    console.log('State updated to:', newState)
    return newState
  }
}

// 运行
const sim = new SimpleSimulation([8, 8, 12])
// 初始状态: [8, 8, 12]
// Step 0 Plan: 索引0(8)和索引1(8) 激活。
// Step 0 Exec: 索引0产生(+1 idx0, +1 idx1)，索引1产生(+1 idx1, +1 idx2)
// Step 0 Update: idx0+=1, idx1+=2, idx2+=1 -> [9, 10, 13]
// ... 循环直到所有 >= 10
sim.run()
```

### 这种抽象的三个关键特征：

1.  **快照隔离 (Snapshot Isolation)**: `plan` 和 `tasks` 看到的都是 `currentState`。在 `Execution` 阶段，没有任何任务能看到其他任务产生的结果。
2.  **同步屏障 (Barrier Synchronization)**: `await Promise.all(...)` 是一堵墙。只有所有并行的计算都撞到这堵墙并停下来，系统才会进入下一个阶段。
3.  **原子状态转换 (Atomic State Transition)**: `update` 函数是唯一允许修改状态的地方，它将离散的计算结果一次性“坍缩”为新的系统状态。
