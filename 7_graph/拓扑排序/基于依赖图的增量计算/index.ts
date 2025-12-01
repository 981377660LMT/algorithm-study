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
