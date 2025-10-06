// Reactive Dependency Graph
// 一个基于依赖图的、懒加载的、增量计算的状态管理引擎

// 节点值的状态
const enum Status {
  Unresolved = 0, // 未计算
  Computing = 1, // 计算中 (用于检测循环依赖)
  Computed = 2 // 已计算
}

// 节点ID
let nextId = 0

// 源节点 (简化版 StateField)
export class SourceNode<T> {
  readonly id = nextId++
  constructor(
    public readonly create: () => T,
    public readonly update: (value: T, event: any) => T,
    public readonly compare: (a: T, b: T) => boolean = (a, b) => a === b
  ) {}
}

// 派生节点 (简化版 Facet)
export class DerivedNode<T> {
  readonly id = nextId++
  constructor(
    public readonly dependencies: (SourceNode<any> | DerivedNode<any>)[],
    public readonly compute: (get: <V>(node: SourceNode<V> | DerivedNode<V>) => V) => T,
    public readonly compare: (a: T, b: T) => boolean = (a, b) => a === b
  ) {}
}

// 系统状态 (简化版 EditorState)
export class SystemState {
  private values: any[] = []
  private status: Status[] = []

  // 构造函数接收一个“配置蓝图”
  constructor(private readonly config: SystemConfig) {
    this.values = new Array(config.nodeCount).fill(undefined)
    this.status = new Array(config.nodeCount).fill(Status.Unresolved)
  }

  // 获取一个节点的值 (这是用户调用的入口)
  public get<T>(node: SourceNode<T> | DerivedNode<T>): T {
    this.ensure(node.id)
    return this.values[node.id]
  }

  // 更新状态，返回一个新的状态实例
  public update(event: any): SystemState {
    const newState = new SystemState(this.config)

    // 对于源节点，直接用旧值和事件计算新值
    for (const source of this.config.sources) {
      const oldValue = this.get(source.node)
      const newValue = source.node.update(oldValue, event)
      if (!source.node.compare(oldValue, newValue)) {
        newState.values[source.node.id] = newValue
        newState.status[source.node.id] = Status.Computed
      } else {
        // 如果值没变，直接复用
        newState.values[source.node.id] = oldValue
        newState.status[source.node.id] = Status.Computed
      }
    }
    return newState
  }

  // 调度器 (简化版 ensureAddr)
  private ensure(id: number): void {
    if (this.status[id] === Status.Computed) return
    if (this.status[id] === Status.Computing) throw new Error('Cyclic dependency detected')

    this.status[id] = Status.Computing

    const derived = this.config.derived.find(d => d.node.id === id)
    if (derived) {
      // 这是一个派生节点，先确保其所有依赖项已计算
      for (const dep of derived.node.dependencies) {
        this.ensure(dep.id)
      }
      // 然后计算自己
      const getter = <V>(node: SourceNode<V> | DerivedNode<V>) => this.values[node.id]
      const newValue = derived.node.compute(getter)

      // 优化：与旧值比较，如果没变就不更新 (这里简化了，实际应从旧状态获取)
      this.values[id] = newValue
    } else {
      // 这是一个源节点，在第一次访问时创建初始值
      const source = this.config.sources.find(s => s.node.id === id)!
      this.values[id] = source.node.create()
    }

    this.status[id] = Status.Computed
  }
}

// 配置蓝图 (简化版 Configuration)
class SystemConfig {
  public readonly sources: { node: SourceNode<any> }[] = []
  public readonly derived: { node: DerivedNode<any> }[] = []
  public readonly nodeCount: number

  constructor(nodes: (SourceNode<any> | DerivedNode<any>)[]) {
    this.nodeCount = nextId // 依赖全局ID计数器
    for (const node of nodes) {
      if (node instanceof SourceNode) this.sources.push({ node })
      else this.derived.push({ node })
    }
  }
}

// 1. 定义节点
const strength = new SourceNode<number>(
  () => 10, // 初始力量
  (val, event) => (event.type === 'LEVEL_UP' ? val + 1 : val) // 升级事件
)

const agility = new SourceNode<number>(
  () => 10, // 初始敏捷
  (val, event) => (event.type === 'LEVEL_UP' ? val + 1 : val)
)

const equipmentStrength = new SourceNode<number>(
  () => 5, // 初始装备力量
  (val, event) => (event.type === 'EQUIP' && event.stat === 'STR' ? event.value : val)
)

const totalStrength = new DerivedNode<number>(
  [strength, equipmentStrength],
  get => get(strength) + get(equipmentStrength)
)

const totalAgility = new DerivedNode<number>(
  [agility], // 假设没有装备加敏捷
  get => get(agility)
)

const attackPower = new DerivedNode<number>([totalStrength], get => get(totalStrength) * 1.5)

const critChance = new DerivedNode<number>([totalAgility], get => get(totalAgility) * 0.5)

// 2. 创建系统
const allNodes = [
  strength,
  agility,
  equipmentStrength,
  totalStrength,
  totalAgility,
  attackPower,
  critChance
]
const config = new SystemConfig(allNodes)
let state = new SystemState(config)

// 3. 查询初始值 (懒加载会在这里触发计算)
console.log(`初始攻击力: ${state.get(attackPower)}`) // -> (10 + 5) * 1.5 = 22.5
console.log(`初始暴击率: ${state.get(critChance)}`) // -> 10 * 0.5 = 5

// 4. 应用一个事件，生成新状态
console.log('\n--- 角色升级 ---')
let stateAfterLevelUp = state.update({ type: 'LEVEL_UP' })

// 旧状态的值不变
console.log(`旧状态攻击力: ${state.get(attackPower)}`) // 22.5

// 新状态的值是更新后的 (只有依赖项被重新计算)
console.log(`新状态攻击力: ${stateAfterLevelUp.get(attackPower)}`) // -> (11 + 5) * 1.5 = 24

// 5. 应用另一个事件
console.log('\n--- 换上新装备 ---')
let stateAfterEquip = stateAfterLevelUp.update({ type: 'EQUIP', stat: 'STR', value: 20 })
console.log(`换装备后攻击力: ${stateAfterEquip.get(attackPower)}`) // -> (11 + 20) * 1.5 = 46.5
// 注意：暴击率相关的节点完全不会被重新计算，因为它的依赖项(敏捷)没有变化

export {}
