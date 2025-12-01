mobx 的设计思想在计算机科学中被称为 **“基于依赖图的增量计算 (Incremental Computing on Dependency Graphs)”**，或者更具体地说是 **“惰性脏标记传播 (Lazy Dirty-Flag Propagation)”**。

我们可以将其抽象为一个 **“有向无环图 (DAG) 中的信号传播模型”**。

---

### 一、 核心抽象模型

在这个模型中，系统由无数个 **节点 (Node)** 组成，节点之间通过 **边 (Edge)** 连接，形成一个网络。

#### 1. 三种角色

- **源节点 (Source Node)**: 系统的输入端。它是“事实”的来源，不可派生，只能被外部修改。（对应 Atom）
- **派生节点 (Derived Node)**: 中间处理单元。它是一个纯函数 `f(x, y, ...)`，输入来自上游节点，输出供下游使用。它拥有**缓存 (Cache)** 能力。（对应 Computed）
- **终端节点 (Sink Node)**: 系统的输出端。它负责执行副作用（如写入磁盘、发送网络请求、驱动机械臂）。（对应 Reaction）

#### 2. 三种状态 (元数据)

每个派生节点都维护一个元数据状态，用于描述**“我的缓存是否可信”**：

1.  **`VALID` (可信)**: 我的缓存与上游完全一致，直接用。
2.  **`CHECK` (存疑)**: 上游有人变了，但我不知道那个变化是否会影响我。我需要“核查”。
3.  **`INVALID` (失效)**: 我直接依赖的源头变了，我的缓存肯定是错的，必须重算。

---

### 二、 运作机制：推拉结合 (Push-Pull)

这个模型的高明之处在于它将更新过程拆分为两个阶段：**廉价的通知** 和 **昂贵的计算**。

#### 阶段 1：失效传播 (Invalidation Phase) —— "推 (Push)"

- **触发**：源节点发生变化。
- **行为**：源节点沿着“边”向下游广播消息：“我变了”。
- **响应**：
  - 直接下游标记为 `INVALID`。
  - 再下游标记为 `CHECK`。
  - **关键点**：只改状态标记，**不进行任何数据计算**。这是一个 O(N) 的极速操作，瞬间完成。

#### 阶段 2：按需重算 (Revalidation Phase) —— "拉 (Pull)"

- **触发**：终端节点需要执行任务（或外部查询某个派生节点）。
- **行为**：终端节点向上游发起询问：“你的数据是最新的吗？”
- **响应 (递归逻辑)**：
  - 如果是 `VALID`：直接返回缓存。**（剪枝，性能极大提升）**
  - 如果是 `INVALID`：立即执行函数重算，更新缓存，设为 `VALID`。
  - 如果是 `CHECK`：**这是核心抽象**。
    1.  节点暂停回答，先去问它的上游。
    2.  如果上游重算后发现**值没变**（逻辑等价），那么本节点虽然曾被标记为 `CHECK`，但实际上不需要重算。直接切回 `VALID`。
    3.  如果上游重算后**值变了**，本节点才执行重算。

---

### 三、 泛化应用场景

这种思想广泛存在于计算机系统的各个领域：

#### 1. 编译系统 (Build Systems) —— 如 Make, Bazel, Ninja

- **源节点**: `.c` / `.cpp` 源代码文件。
- **派生节点**: `.o` 目标文件。
- **终端节点**: 可执行文件 (Binary)。
- **场景**:
  你修改了 `header.h` 中的一行注释。
  - **Push**: 所有 include 了这个头文件的 `.c` 文件对应的 `.o` 任务被标记为 `CHECK`。
  - **Pull**: 编译器检查发现 `header.h` 的**内容哈希**变了，但是预处理后的代码流（AST）没变（因为只是注释）。
  - **Result**: 编译器决定**不重新编译** `.o` 文件，链接器直接使用旧的 `.o`。这就是增量编译的精髓。

#### 2. 电子表格 (Spreadsheets) —— 如 Excel 引擎

- **源节点**: 单元格 A1 (输入 10)。
- **派生节点**: 单元格 B1 (=A1\*2), C1 (=if(B1>100, "High", "Low"))。
- **场景**:
  你把 A1 从 10 改为 11。
  - **Push**: B1 标记为 `INVALID`，C1 标记为 `CHECK`。
  - **Pull**: 屏幕刷新（终端节点）。计算 B1 得到 22。
  - **Check**: C1 发现 B1 虽然变了 (20 -> 22)，但公式逻辑 `22 > 100` 依然是 False。
  - **Result**: C1 的值保持 "Low" 不变，依赖 C1 的图表**不需要重绘**。

#### 3. 数据库物化视图 (Materialized Views)

- **源节点**: 基础表 (Base Tables)。
- **派生节点**: 物化视图 (预计算的 Join/Group By 结果)。
- **场景**:
  基础表插入了一条新记录。
  - 数据库不会立即全量重算视图（太慢）。
  - 它会标记视图为“脏”。
  - 当查询到来时，数据库检查增量日志 (Delta Log)，发现新插入的记录不满足视图的 `WHERE` 条件。
  - 视图直接返回旧缓存，零 IO 开销。

#### 4. 3D 渲染管线 (Game Engine / CGI)

- **源节点**: 骨骼位置、光照参数。
- **派生节点**: 世界坐标矩阵、包围盒 (Bounding Box)。
- **终端节点**: 像素绘制。
- **场景**:
  角色手臂动了一下。
  - 手臂的变换矩阵变了。
  - 但是角色的整体包围盒（用于视锥剔除）可能没有变（手臂还在盒子范围内）。
  - 渲染引擎检测到包围盒未变，**直接跳过**该角色所有子物体的可见性裁剪计算。

---

### 四、 总结：这是一种什么哲学？

这种抽象体现了系统设计中的 **“惰性 (Laziness)”** 和 **“短路 (Short-circuiting)”** 哲学。

**核心公理**：

1.  **变化是昂贵的**：不要因为“可能”变了就去重算，要确定“真的”变了才动。
2.  **传播是廉价的**：通知“我不确定”比证明“我是对的”要快得多。
3.  **稳定性优于精确性**：如果中间步骤剧烈震荡，但最终结果保持稳定，那么下游系统应当感知不到震荡（波动的吸收）。

这就是脱离了 MobX 和前端之后，这种机制在计算机科学中原本的样子。

---

这是一个脱离了具体业务场景（如前端 UI），纯粹的 **“基于依赖图的增量计算” (Incremental Computing on Dependency Graph)** 引擎实现。

这套代码实现了一个通用的 **DAG（有向无环图）计算调度器**，它可以用在构建系统、电子表格引擎、数据流处理等任何需要“按需计算”和“最小化更新”的场景中。

### 核心抽象设计

我们定义三种核心状态，对应之前的理论：

- `CLEAN` (Valid): 数据是新鲜的。
- `CHECK` (Possibly Stale): 上游变了，但我不知道我受不受影响。
- `DIRTY` (Invalid): 我直接依赖的源头变了，我肯定脏了。

### 代码实现 (TypeScript)

```typescript
/**
 * 节点状态枚举
 * 对应理论中的 VALID, CHECK, INVALID
 */
enum NodeState {
  CLEAN = 0, // 缓存完全可信
  CHECK = 1, // 上游有变动，需要核查
  DIRTY = 2 // 必须重算
}

/**
 * 依赖图中的通用节点基类
 * 维护图的拓扑结构和状态流转
 */
abstract class GraphNode<T> {
  public id: string
  public state: NodeState = NodeState.DIRTY

  // 缓存的值
  protected value: T | undefined

  // 图结构：谁依赖我 (Downstream)
  protected observers = new Set<GraphNode<any>>()

  // 图结构：我依赖谁 (Upstream)
  protected dependencies = new Set<GraphNode<any>>()

  constructor(id: string) {
    this.id = id
  }

  /**
   * [Push Phase] 信号传播：向下游广播“我变了”
   * 这是一个极其廉价的操作，只改状态，不计算
   */
  protected notifyObservers(newState: NodeState) {
    if (this.observers.size === 0) return

    for (const observer of this.observers) {
      // 如果下游已经是 DIRTY，那它肯定知道自己要重算，不用再通知
      // 如果下游是 CHECK，且我们这次传播的也是 CHECK，也不用重复通知
      if (observer.state === NodeState.DIRTY) continue
      if (observer.state === NodeState.CHECK && newState === NodeState.CHECK) continue

      observer.state = newState
      // 递归向下传播：一旦源头脏了，下游全部标记为 CHECK (或者 DIRTY)
      // 注意：这里通常传播 CHECK，因为下游不知道上游的具体变化是否会影响自己
      observer.notifyObservers(NodeState.CHECK)
    }
  }

  /**
   * 建立依赖关系 (连线)
   */
  public addDependency(node: GraphNode<any>) {
    if (this.dependencies.has(node)) return
    this.dependencies.add(node)
    node.observers.add(this)
  }
}

/**
 * 源节点 (Source Node)
 * 系统的输入端，外界直接修改它
 */
class InputNode<T> extends GraphNode<T> {
  constructor(id: string, initialValue: T) {
    super(id)
    this.value = initialValue
    this.state = NodeState.CLEAN
  }

  /**
   * 外界写入新值
   * 触发 Push Phase
   */
  public set(newValue: T) {
    if (newValue === this.value) return

    console.log(`[SET] ${this.id} = ${newValue}`)
    this.value = newValue
    // 标记自己为 CLEAN (因为值已经是最新的了)
    this.state = NodeState.CLEAN
    // 通知下游：你们脏了 (DIRTY)
    // 这里直接传 DIRTY 是因为直接依赖 Input 的节点肯定需要重算
    this.notifyObservers(NodeState.DIRTY)
  }

  public get(): T {
    return this.value!
  }
}

/**
 * 派生节点 (Derived Node)
 * 中间处理单元，拥有计算逻辑和缓存
 */
class ComputedNode<T> extends GraphNode<T> {
  private computeFn: () => T

  constructor(id: string, computeFn: () => T) {
    super(id)
    this.computeFn = computeFn
  }

  /**
   * [Pull Phase] 按需计算
   * 核心算法：惰性求值 + 脏标记检查
   */
  public get(): T {
    // 1. 快速路径：如果是 CLEAN，直接返回缓存
    if (this.state === NodeState.CLEAN) {
      console.log(`[HIT] ${this.id} cache hit`)
      return this.value!
    }

    // 2. 检查路径：如果是 CHECK，需要去问上游
    if (this.state === NodeState.CHECK) {
      console.log(`[CHECK] ${this.id} checking dependencies...`)
      let upstreamChanged = false

      for (const dep of this.dependencies) {
        // 递归调用上游的 get()
        // 如果上游是 ComputedNode，它会尝试重算
        if (dep instanceof ComputedNode) {
          // 这是一个比较 trick 的地方：
          // 我们需要知道上游是否 *刚刚* 发生了变化。
          // 在这个简化模型中，我们通过比较值来判断。
          // 在更复杂的系统中，可能会用版本号 (Version Vector)。
          const oldValue = dep.peek()
          const newValue = dep.get() // 触发上游重算
          if (oldValue !== newValue) {
            upstreamChanged = true
            break // 只要有一个上游变了，我就得重算
          }
        } else {
          // 如果是 InputNode，它如果是 CLEAN 的，说明没变（或者已经处理过了）
          // 但如果 InputNode 处于某种中间状态（在这个简化模型里 Input 总是 CLEAN），则需判断
          // 实际上 InputNode 变化会直接把下游设为 DIRTY，不会进入 CHECK 逻辑
          // 所以进入 CHECK 逻辑通常意味着依赖的是 ComputedNode
        }
      }

      if (!upstreamChanged) {
        console.log(`[KEEP] ${this.id} dependencies stable, keeping cache`)
        this.state = NodeState.CLEAN
        return this.value!
      }
    }

    // 3. 重算路径：DIRTY 或者 CHECK 失败
    console.log(`[CALC] ${this.id} re-computing...`)

    // 在执行 computeFn 之前，通常需要收集依赖（如果是动态图）
    // 这里为了简化，假设依赖关系是静态构建的 (addDependency)
    // 或者在 computeFn 内部手动调用依赖的 get()

    const newValue = this.computeFn()

    // 4. 变更吸收 (Change Absorption)
    // 如果算出来的新值和旧值一样，我们虽然重算了，但对下游来说，我没变！
    if (newValue !== this.value) {
      this.value = newValue
      // 注意：这里不需要再 notifyObservers 了
      // 因为下游已经是 CHECK 或 DIRTY 状态了，它们会在访问时自己来问我
    } else {
      console.log(`[ABSORB] ${this.id} value unchanged after calc`)
    }

    this.state = NodeState.CLEAN
    return this.value!
  }

  // 仅用于内部查看当前值，不触发计算
  public peek(): T | undefined {
    return this.value
  }
}

// ==========================================
// 场景演示：构建系统 / 电子表格
// ==========================================

// 场景：Z = if (A > 10) then B else C
// 初始：A=1, B=100, C=200 -> Z=200

console.log('--- Init Graph ---')

// 1. 定义节点
const A = new InputNode<number>('A', 1)
const B = new InputNode<number>('B', 100)
const C = new InputNode<number>('C', 200)

// 中间节点：判断条件
const Condition = new ComputedNode<boolean>('Cond', () => {
  return A.get() > 10
})
Condition.addDependency(A)

// 最终节点：结果
const Z = new ComputedNode<number>('Z', () => {
  // 动态依赖：根据条件决定依赖 B 还是 C
  // 在这个静态依赖模型中，我们需要预先声明所有可能的依赖
  // 或者在 computeFn 中动态访问
  if (Condition.get()) {
    return B.get()
  } else {
    return C.get()
  }
})
// 静态声明依赖（简化版）
Z.addDependency(Condition)
Z.addDependency(B)
Z.addDependency(C)

// 第一次计算
console.log(`Result Z: ${Z.get()}`)
// [CALC] Cond re-computing...
// [CALC] Z re-computing...
// Result Z: 200

console.log('\n--- Scenario 1: 无效更新 (Change Absorption) ---')
// 修改 A: 1 -> 5
// 逻辑上 5 > 10 依然是 False，所以 Cond 的结果不变，Z 也不应该重算
A.set(5)
// [SET] A = 5
// A 通知 Cond -> DIRTY
// Cond 通知 Z -> CHECK

console.log('>>> Reading Z...')
Z.get()
// 1. Z 是 CHECK。Z 检查依赖 Cond。
// 2. Cond 是 DIRTY。Cond 重算。 5 > 10 is False.
// 3. Cond 发现新值 False 等于旧值 False。[ABSORB] 发生。
// 4. Cond 变为 CLEAN。
// 5. Z 发现 Cond (上游) 没变。
// 6. Z 检查 B, C (InputNode 且没变)。
// 7. Z 决定不重算。[KEEP] 发生。
// 8. Z 变为 CLEAN。

console.log('\n--- Scenario 2: 有效更新 ---')
// 修改 A: 5 -> 15
// 逻辑上 15 > 10 is True，Cond 变了，Z 需要重算
A.set(15)
// A 通知 Cond -> DIRTY
// Cond 通知 Z -> CHECK

console.log('>>> Reading Z...')
console.log(`Result Z: ${Z.get()}`)
// 1. Z 是 CHECK。Z 检查依赖 Cond。
// 2. Cond 是 DIRTY。Cond 重算。 15 > 10 is True.
// 3. Cond 值变了 (False -> True)。
// 4. Z 发现上游变了。
// 5. Z 重算。读取 B (100)。
// 6. Z 更新为 100。
```

### 代码深度解析

这段代码展示了增量计算引擎的三个核心特征：

1.  **Push Phase (通知)**:
    在 `InputNode.set` 中，我们只做了一件事：`notifyObservers(NodeState.DIRTY)`。
    这个过程是**瞬时**的。无论图有多大，它只是沿着边把下游节点的状态位改成 `CHECK` 或 `DIRTY`。它不执行任何用户逻辑（`computeFn`）。

2.  **Pull Phase (拉取)**:
    在 `ComputedNode.get` 中，计算是**惰性**的。只有当你调用 `get()` 时，计算才会发生。

3.  **Change Absorption (变更吸收/波动的阻断)**:
    这是最精彩的部分。在 Scenario 1 中：

    - `A` 变了 (1 -> 5)。
    - `Condition` 被迫重算。
    - 但是 `Condition` 的结果 (False) 没变。
    - 于是 `Z` **完全跳过了重算**。

    想象一下，如果 `Z` 是一个耗时 1 秒的复杂计算，或者是触发一次昂贵的数据库写入，这个机制就节省了巨大的资源。

### 总结

这套代码是所有现代响应式系统（React Fiber, Vue Reactivity, MobX, Bazel Build System, Excel Engine）的**最小公约数**。它剥离了所有语法糖，只保留了**依赖图**和**状态机**这两个核心概念。

---

基于刚才那个“最小公约数”版本的增量计算引擎，如果要在生产环境（如大型前端框架、构建工具、金融计算引擎）中使用，还有几个**致命的优化点**和**高级特性**需要补充。

目前的实现虽然逻辑正确，但在**易用性**、**内存管理**和**极端性能**上还有很大提升空间。

以下是四个核心优化方向：

---

### 1. 动态依赖收集 (Dynamic Dependency Collection)

**痛点**：目前的实现需要手动调用 `Z.addDependency(B)`。这在复杂的业务逻辑中是不可能的（比如 `if (cond) A else B`，依赖是动态变化的）。

**优化方案**：引入**全局上下文栈**。

- **原理**：
  1.  在 `ComputedNode` 开始计算前，把自己推入全局栈 `Stack.push(this)`。
  2.  任何节点被读取 (`get()`) 时，检查栈顶是谁。如果是 `Node X`，说明 `Node X` 正在依赖我。
  3.  自动建立 `X -> Me` 的连接。
  4.  计算结束后，`Stack.pop()`。

**代码抽象**：

```typescript
// 全局栈
const contextStack: ComputedNode[] = []

class GraphNode {
  get() {
    // 自动收集依赖
    if (contextStack.length > 0) {
      const runningNode = contextStack[contextStack.length - 1]
      runningNode.addDependency(this) // 自动连线
    }
    // ... 原有逻辑
  }
}

class ComputedNode {
  get() {
    // ...
    contextStack.push(this) // 入栈
    try {
      this.computeFn()
    } finally {
      contextStack.pop() // 出栈
    }
    // ...
  }
}
```

---

### 2. 依赖剪枝与内存回收 (Dependency Pruning / Diffing)

**痛点**：对于 `Z = cond ? A : B`。

- T1 时刻 `cond=true`，Z 依赖 A。
- T2 时刻 `cond=false`，Z 依赖 B。
- **问题**：如果不清理，Z 依然订阅着 A。当 A 变化时，Z 会被标记为 `CHECK`，甚至触发重算，但实际上 Z 根本不关心 A 了。这叫 **“僵尸依赖” (Zombie Dependencies)**，会导致内存泄漏和无效计算。

**优化方案**：**双缓冲依赖对比 (Double-Buffer Diffing)**。

- **原理**：
  1.  每次重算前，创建一个新的空集合 `newDependencies`。
  2.  计算过程中，所有访问到的节点加入 `newDependencies`。
  3.  计算结束后，对比 `oldDependencies` 和 `newDependencies`。
      - 在 Old 但不在 New 的：**解绑** (removeObserver)。
      - 在 New 但不在 Old 的：**绑定** (addObserver)。
  4.  `oldDependencies = newDependencies`。

---

### 3. 批处理与防抖 (Batching & Glitch Freedom)

**痛点**：

```typescript
A.set(1)
B.set(2)
// 假设 C = A + B
```

如果我们连续修改 A 和 B，C 可能会收到两次通知，甚至在中间状态下（A 新 B 旧）被计算一次。这种中间的不一致状态被称为 **"Glitch"**。

**优化方案**：**事务 (Transaction)**。

- **原理**：
  引入一个全局锁或计数器。
  1.  `startBatch()`: 停止所有 `notifyObservers` 的实际触发，只把脏节点放入一个 `pendingQueue`。
  2.  修改 A，修改 B...
  3.  `endBatch()`: 遍历 `pendingQueue`，统一进行一次传播。

这不仅解决了 Glitch，还极大地减少了图遍历的次数。

---

### 4. 位掩码状态压缩 (Bitmask Optimization)

**痛点**：在超大规模图（如 Excel 的 100 万个单元格）中，对象属性的访问速度和内存占用至关重要。使用 `enum` 或字符串存储状态太慢。

**优化方案**：**位运算**。

- **原理**：
  使用一个 `int` 类型的 `flags` 字段存储所有状态。

  ```typescript
  const STATE_CLEAN = 0;
  const STATE_CHECK = 1 << 0; // 1
  const STATE_DIRTY = 1 << 1; // 2
  const IS_COMPUTING = 1 << 2; // 4 (防止循环依赖)

  // 检查是否脏
  if (node.flags & (STATE_CHECK | STATE_DIRTY)) { ... }

  // 标记为 Clean
  node.flags &= ~STATE_DIRTY & ~STATE_CHECK;
  ```

  位运算比布尔逻辑快得多，且节省内存。

---

### 5. 环检测 (Cycle Detection)

**痛点**：用户写出了 `A = B + 1; B = A + 1;`。
在目前的实现中，这会导致无限递归，栈溢出崩溃。

**优化方案**：**染色法**。

- **原理**：
  在 `get()` 开始时，标记当前节点为 `COMPUTING` (灰色)。
  在 `get()` 结束时，标记为 `COMPUTED` (黑色)。
  如果在 `get()` 过程中遇到了一个标记为 `COMPUTING` 的节点，说明**检测到了环**。
  此时应抛出异常 `CycleDetectedException`，而不是让程序崩溃。

---

### 终极形态的代码结构预览

如果加上上述所有优化，一个工业级的 `ComputedNode` 大概长这样：

```typescript
class OptimizedNode {
  // 优化4: 位掩码
  flags: number = STATE_DIRTY

  // 优化2: 双缓冲依赖
  dependencies: Node[] = []
  newDependencies: Node[] = []

  get() {
    // 优化5: 环检测
    if (this.flags & IS_COMPUTING) throw new Error('Cycle Detected')

    // 优化1: 动态收集 - 建立反向连接
    if (globalStack.length > 0) {
      globalStack.peek().newDependencies.push(this)
    }

    if (this.flags & (STATE_CHECK | STATE_DIRTY)) {
      // ... 执行 Pull Phase ...
      // 优化2: 在这里进行依赖 Diff 和解绑操作
    }

    return this.value
  }
}
```

这些优化点正是 MobX、Vue 3 Reactivity、React Fiber 内部最复杂、最核心的部分。

---

除了上述提到的动态依赖收集、依赖剪枝、批处理和位掩码优化外，如果我们要将这个引擎推向**极致性能**或**特定领域的极端场景**（如游戏引擎、高频交易、超大规模数据可视化），还有以下几个深水区的优化点：

### 6. 拓扑排序执行 (Topological Sort Execution)

**痛点**：
在目前的“递归拉取（Pull）”模式中，如果图非常深，递归调用栈会很深。更重要的是，可能会出现**重复访问**。
例如：`D` 依赖 `B` 和 `C`，`B` 和 `C` 都依赖 `A`。
当读取 `D` 时，可能会先递归计算 `B`（进而访问 `A`），回来后再递归计算 `C`（再次访问 `A`）。虽然有缓存，但函数调用的开销依然存在。

**优化方案**：**基于层级的优先级队列**。

- **原理**：
  1.  给每个节点维护一个 `depth`（深度）或 `rank`。源节点深度为 0，下游节点深度 = `max(上游深度) + 1`。
  2.  当源节点变化时，不立即递归通知。而是将所有受影响的节点放入一个**优先级队列**（按深度从小到大排序）。
  3.  调度器按顺序取出节点执行。
- **效果**：确保在计算 `D` 之前，`B` 和 `C` 肯定都已经算好了。`A` 只会被访问一次。这被称为 **"Glitch-free by construction"**（天然无中间态）。

### 7. 脏区裁剪 (Dirty Region Clipping) —— 针对 UI/渲染

**痛点**：
假设你有一个 10000 行的列表。数据源变了，导致 10000 个 Row 组件都标记为 `CHECK`。
但实际上，屏幕上只显示了前 20 行。计算剩下 9980 行的状态是浪费的。

**优化方案**：**视口感知 (Viewport Awareness) / 惰性订阅**。

- **原理**：
  引入“激活”状态。只有当一个节点被“最终消费者”（如屏幕上的像素、活动的 DOM 节点）订阅时，它才处于 `Active` 状态。
  - **Inactive Node**：即使上游变了，我也不标记为 Dirty，我直接忽略。
  - **On Activate**：当节点从屏幕外滑入屏幕内（变为 Active）时，立即检查一次上游，如果过期了再重算。
- **效果**：性能与数据总量无关，只与**可见数据量**有关。

### 8. 结构共享与持久化数据结构 (Structural Sharing)

**痛点**：
如果节点的值是一个巨大的对象或数组（比如 1MB 的 JSON）。
`newValue !== oldValue` 的比较代价极高（深度比较太慢，引用比较容易误判）。
且每次生成新对象都会导致大量的 GC（垃圾回收）压力。

**优化方案**：**不可变数据 (Immutable Data) + 结构共享**。

- **原理**：
  类似 Immutable.js 或 Immer。
  当修改大对象的某个属性时，不复制整个对象，而是复用未修改部分的引用，只创建修改路径上的新节点。
- **效果**：
  1.  **极速比较**：`oldValue === newValue` 引用相等即代表内容完全一致，O(1) 复杂度。
  2.  **历史回溯**：天然支持 Undo/Redo（时间旅行）。

### 9. 弱引用 (Weak References)

**痛点**：
在动态依赖图中，`A` 依赖 `B`。如果不手动解绑，`B` 的 `observers` 集合会一直持有 `A` 的引用。
导致 `A` 即使不再使用了，也无法被垃圾回收（Memory Leak）。

**优化方案**：**`WeakMap` / `WeakSet` / `WeakRef`**。

- **原理**：
  在存储反向依赖（`observers`）时，使用弱引用。
  如果 `A` 没有其他强引用（比如被销毁的组件），GC 会自动回收 `A`，同时 `B` 的 `observers` 列表中对应的条目也会自动消失（或被清理）。
- **效果**：彻底杜绝“忘记取消订阅”导致的内存泄漏。

### 10. 细粒度并发 (Fine-grained Concurrency)

**痛点**：
一次计算任务太重（比如 Excel 中修改一个单元格触发 10 万个公式重算），导致主线程卡死（UI 掉帧）。

**优化方案**：**时间切片 (Time Slicing) / 可中断计算**。

- **原理**：
  借鉴 React Fiber 的思想。
  1.  将巨大的依赖图遍历任务拆分成一个个小单元（比如每个节点的 `get` 是一个单元）。
  2.  调度器执行几个单元后，检查时间（`performance.now()`）。
  3.  如果超过 16ms（一帧的时间），**挂起**当前任务，把控制权交还给浏览器绘制 UI。
  4.  下一帧空闲时，**恢复**计算。
- **效果**：即使计算量巨大，界面依然流畅响应。

### 11. 静态分析与预编译 (AOT Compilation)

**痛点**：
动态收集依赖（运行时的 `stack.push/pop`）虽然灵活，但有运行时开销。
对于像 Svelte 或 SolidJS 这样的框架，模板结构是固定的。

**优化方案**：**编译时依赖推导**。

- **原理**：
  在代码编译阶段（Build time），分析代码结构：
  `<div>{count() * 2}</div>`
  编译器直接生成代码：
  `createEffect(() => updateText(node, count() * 2))`
  直接建立 `count -> effect` 的连接，跳过运行时的依赖收集过程。
- **效果**：运行时体积极小，速度极快（No Virtual DOM）。

---

### 总结：优化的尽头

如果你把以上所有优化都做到了，你就得到了一个：

- **MobX** 的易用性（动态收集）
- **React Fiber** 的流畅度（时间切片）
- **Svelte** 的轻量级（预编译）
- **Immutable.js** 的高效比对（结构共享）
- **Rust/C++** 级别的内存管理（弱引用）

这就是现代前端工程化和高性能计算引擎的**技术深水区**。

---

这是一个非常好的要求。为了不“硬凑”，我们需要剔除那些边角料的微优化，直击**架构设计层面**的痛点。

在超大规模（百万级节点）或超高频（60FPS+）的场景下，目前的架构还有三个**真正的瓶颈**：

1.  **深层图的遍历开销**（即使不计算，光是问“你变了吗”就很慢）。
2.  **内存碎片与 GC 压力**（对象太多了）。
3.  **缓存的内存占用**（所有计算结果都存着，内存会爆）。

针对这三个痛点，还有三个**核弹级**的优化方案：

---

### 1. 全局版本控制 (Global Versioning / Clock)

**痛点：O(Depth) 的检查成本**
在目前的 `CHECK` 状态逻辑中，如果依赖链很深：`A -> B -> C -> ... -> Z`。
当读取 `Z` 时，`Z` 问 `Y`，`Y` 问 `X` ... 一直问到 `A`。
即使 `A` 根本没变，光是这层层递归的函数调用和状态检查，在深层图中也是巨大的开销。

**优化方案：全局时钟 + 版本号比对**

- **原理**：
  1.  维护一个全局递增计数器 `GlobalVersion`。
  2.  **写操作**：每当任何 `InputNode` 发生改变，`GlobalVersion++`，并将该 Node 的 `lastChangedVersion` 设为当前全局版本。
  3.  **读操作**：`ComputedNode` 记录自己上次计算时的 `lastVerifiedVersion`。
  4.  **极速剪枝**：
      在 `get()` 的第一行：
      ```typescript
      if (this.lastVerifiedVersion === GlobalVersion) {
        return this.value // O(1) 返回，完全跳过图遍历
      }
      ```
- **效果**：
  如果整个系统在这一帧没有任何变化（或者变化已经处理过了），无论图有多深，所有节点的访问瞬间变成 **O(1)**。这对于高频轮询的场景（如游戏循环）是质的飞跃。

---

### 2. 扁平化内存布局 (Structure of Arrays - SoA)

**痛点：指针跳跃与 GC 压力**
目前的实现中，每个 Node 都是一个独立的 JS 对象 (`new Node()`)。

1.  **内存碎片**：这些对象散落在堆内存的各个角落，CPU 缓存命中率（Cache Locality）极低。
2.  **GC 压力**：创建 100 万个节点就是 100 万个对象，GC 标记清除时会卡顿。

**优化方案：使用 TypedArray 模拟对象池**

- **原理**：
  不再 `new Class`，而是分配几个巨大的数组（如 `Int32Array`, `Float64Array`）。用 **索引 (Index)** 代表节点 ID。

  ```typescript
  // 假设最大 10万个节点
  const state = new Uint8Array(100000); // 存 NodeState
  const values = new Float64Array(100000); // 存 Value (假设是数字)
  const dependencyHeads = new Int32Array(100000); // 存依赖链表的头指针
  // ... 其他属性

  function getNodeValue(index) {
     if (state[index] === DIRTY) { ... }
     return values[index];
  }
  ```

- **效果**：
  1.  **零 GC**：整个生命周期只有几个大数组，没有对象创建和销毁。
  2.  **极致缓存亲和性**：数据在内存中是连续的，CPU 预取极其高效。
      _这是高性能物理引擎（如 Rapier.js）和 ECS 游戏架构的标准做法。_

---

### 3. 易失性缓存 (Volatile / Weak Caching)

**痛点：内存爆炸 (OOM)**
在电子表格或数据分析工具中，可能有 1000 万个中间计算步骤。如果每个 `ComputedNode` 都永久缓存其结果 `this.value`，内存很快就会耗尽。
但如果完全不缓存，计算又太慢。

**优化方案：LRU 策略 / 引用计数清理**

- **原理**：
  区分 **"Hot Node"** (有人订阅) 和 **"Cold Node"** (无人订阅，只是临时读了一下)。
  1.  **Hot Node**：保持缓存，行为照旧。
  2.  **Cold Node**：
      - 策略 A (立即释放)：读取完后，不保存 `value`，下次读重新算。
      - 策略 B (LRU 缓存)：保留 `value`，但放入一个 LRU 队列。当内存压力大或队列满时，丢弃最久未使用的 `value`（将状态重置为 `DIRTY` 或特殊状态）。
- **效果**：
  在有限的内存中运行无限的计算图。MobX 的 `keepAlive: false` (默认行为) 其实就是一种激进的易失性缓存（无人订阅就销毁）。

---

### 总结

这三个优化点分别解决了：

1.  **全局版本** -> 解决 **CPU** 在深层图上的空转。
2.  **SoA 布局** -> 解决 **GC** 和 **CPU Cache** 的瓶颈。
3.  **易失缓存** -> 解决 **RAM** 的物理限制。

加上之前的优化，这基本上涵盖了从 Web 前端到 3A 游戏引擎在“状态管理”这一层面的所有核心技术壁垒。

---

从算法层面深入理解 MobX，本质上是在探讨**有向无环图（DAG）上的动态增量计算问题**。MobX 的核心算法目标是：在图结构动态变化的前提下，以最小的代价维护数据的一致性。

我们可以将 MobX 的算法核心拆解为三个关键部分：**图的构建（依赖收集）**、**图的遍历（变更传播）** 和 **图的剪枝（惰性求值与去重）**。

---

### 一、 图的构建：动态依赖收集算法

MobX 的依赖图不是静态定义的，而是在运行时动态生成的。这对应于算法中的 **“运行时图构建”**。

#### 1. 算法描述

- **输入**：一个正在执行的函数 $F$（Reaction 或 Computed）。
- **过程**：
  1.  设置全局指针 `CurrentNode = F`。
  2.  执行 $F$。
  3.  在 $F$ 执行过程中，每当读取一个节点 $A$（Observable）的值时：
      - 建立有向边 $A \to F$。
      - 将 $A$ 加入 $F$ 的依赖列表 $D_{new}$。
  4.  $F$ 执行结束。
  5.  **Diff 算法**：对比 $F$ 上一次执行时的依赖列表 $D_{old}$ 和本次的 $D_{new}$。
      - $\forall n \in D_{old} \setminus D_{new}$：移除边 $n \to F$（解绑）。
      - $\forall n \in D_{new} \setminus D_{old}$：保持边 $n \to F$（新绑定的在读取时已建立）。
- **输出**：更新后的局部依赖图结构。

#### 2. 算法复杂度

假设 $F$ 读取了 $N$ 个 Observable。

- **时间复杂度**：$O(N)$。MobX 使用了优化的 Diff 算法（基于标记位或双缓冲），避免了 $O(N^2)$ 的集合对比。
- **空间复杂度**：$O(N)$，用于存储依赖关系。

#### 3. 算法意义

支持了 `if (cond) A else B` 这种条件依赖。图结构随数据变化而重构，保证了依赖关系的**完备性**和**最小性**（不依赖不需要的数据）。

---

### 二、 图的遍历：两阶段传播算法 (Two-Phase Propagation)

当源节点发生变化时，MobX 采用了一种**推拉结合（Push-Pull）**的策略来更新图。这是一种**拓扑排序**的变体，但针对动态图进行了优化。

#### 阶段 1：标记传播 (Push Phase / Invalidation)

- **触发**：源节点 $S$ 的值被修改。
- **算法**：深度优先搜索 (DFS) 或 广度优先搜索 (BFS)。
- **过程**：
  从 $S$ 出发，沿着有向边 $A \to B$ 向下遍历。
  - 对于直接下游 $N$：标记状态为 `STALE`（脏）。
  - 对于间接下游 $M$：标记状态为 `POSSIBLY_STALE`（可能脏）。
- **终止条件**：遇到已经标记过的节点，停止该路径的遍历。
- **复杂度**：$O(V + E)$，其中 $V, E$ 是受影响的子图的节点数和边数。由于只改状态位不计算，速度极快。

#### 阶段 2：按需求值 (Pull Phase / Revalidation)

- **触发**：读取某个节点 $T$ 的值（通常是 Reaction 自动触发，或手动读取 Computed）。
- **算法**：递归求值 + 状态机检查。
- **过程**：
  `evaluate(T)`:
  1.  若 $T.state == UP\_TO\_DATE$：返回 $T.value$（缓存命中）。
  2.  若 $T.state == STALE$：
      - 执行 $T$ 的计算函数。
      - 更新 $T.value$。
      - $T.state \leftarrow UP\_TO\_DATE$。
  3.  若 $T.state == POSSIBLY\_STALE$：
      - **核心剪枝逻辑**：遍历 $T$ 的所有直接上游 $D_i$。
      - 递归调用 `evaluate(D_i)`。
      - 若所有 $D_i$ 的值在重算后都未变（引用相等），则 $T$ 不需要重算，$T.state \leftarrow UP\_TO\_DATE$。
      - 若任一 $D_i$ 变了，则 $T$ 转为 `STALE`，执行步骤 2。

---

### 三、 图的剪枝：变更吸收算法 (Change Absorption)

这是 MobX 算法中最精妙的部分，用于解决**无效更新传播**问题。

#### 1. 问题模型

考虑菱形依赖：$A \to B$, $A \to C$, $(B, C) \to D$。
假设 $A: 1 \to 2$。

- $B = A \% 2$。$B: 1 \to 0$（变了）。
- $C = A > 5$。$C: false \to false$（没变）。
- $D = B + (C ? 1 : 0)$。

#### 2. 算法执行流

1.  $A$ 变，$B, C$ 标 `STALE`，$D$ 标 `POSSIBLY_STALE`。
2.  求值 $D$。
3.  $D$ 检查 $B$。$B$ 重算变为 0。$B$ 确实变了。
4.  $D$ 检查 $C$。$C$ 重算仍为 `false`。
    - **变更吸收**：$C$ 的计算结果拦截了 $A$ 的变化传播。对于 $D$ 而言，$C$ 就像没变一样。
5.  $D$ 根据新的 $B$ 和旧的 $C$ 重算。

#### 3. 算法价值

在数据流图中，很多中间节点的计算结果是离散的（如 boolean, enum）。变更吸收算法能**在图的中间层截断变化流**，防止下游进行无意义的重算和渲染。这在 UI 编程中至关重要（避免不必要的 DOM 操作）。

---

### 四、 事务算法：批处理 (Batching)

为了保证图的一致性（Consistency）和避免中间态（Glitch），MobX 引入了事务算法。

#### 1. 算法描述

- 引入全局计数器 `BatchDepth = 0`。
- 引入待处理队列 `PendingReactions = []`。
- **StartBatch**: `BatchDepth++`。
- **EndBatch**: `BatchDepth--`。若 `BatchDepth == 0`，执行 `RunReactions()`。
- **RunReactions**:
  1.  对 `PendingReactions` 中的节点进行去重。
  2.  (可选) 根据图的拓扑序对 Reaction 进行排序（MobX 实际上没有严格做全图拓扑排序，而是依赖 `POSSIBLY_STALE` 机制来保证顺序正确性）。
  3.  依次执行 Reaction。
  4.  若执行过程中触发了新的 Reaction，加入队列，循环直到队列为空（不动点算法）。

#### 2. 解决的问题 (Glitch Freedom)

考虑 $A \to B$, $A \to C$, $(B, C) \to D$。
若无事务，修改 $A$ 可能导致：

1.  $A$ 更新。
2.  $B$ 更新 -> 触发 $D$ 更新（此时 $D$ 读到新 $B$ 和旧 $C$ -> **Glitch! 数据不一致**）。
3.  $C$ 更新 -> 触发 $D$ 更新（此时 $D$ 读到新 $B$ 和新 $C$ -> 正确）。
    $D$ 被计算了两次，且第一次是错误的。

MobX 的算法保证：在最外层 Action 结束前，$D$ 不会运行。当 $D$ 运行时，它会通过 Pull 机制拉取到 $B$ 和 $C$ 的最新值。

---

### 总结：MobX 的算法画像

从算法角度看，MobX 是一个：
**基于动态构建的有向无环图 (DAG)，采用推拉结合 (Push-Pull) 的传播策略，具备变更吸收 (Change Absorption) 能力和事务一致性 (Transactional Consistency) 的增量计算引擎。**

- **空间换时间**：通过维护复杂的双向依赖链表（图结构），换取 $O(1)$ 或 $O(k)$ 的变更感知能力。
- **惰性求值**：通过状态机标记，将计算推迟到最后一刻，最大化利用缓存。
