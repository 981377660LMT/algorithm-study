export {}

type Context = Record<string, any>
type Event = { type: string; payload?: any }
type Action<TContext> = (context: TContext, event: Event) => Partial<TContext> | void

interface StateNode<TContext> {
  // 核心：支持嵌套
  initial?: string // 如果有子状态，必须指定默认进入哪个子状态
  states?: Record<string, StateNode<TContext>> // 子状态表

  // 基础属性
  entry?: Action<TContext>[]
  exit?: Action<TContext>[]
  on?: {
    [eventType: string]: {
      target?: string // 目标状态路径，例如 'idle' 或 'moving.run'
      actions?: Action<TContext>[]
    }
  }
}

interface MachineConfig<TContext> {
  id: string
  initial: string
  context: TContext
  states: Record<string, StateNode<TContext>>
}

// HFSM
class HierarchicalStateMachine<TContext> {
  // 当前状态路径，例如 ['moving', 'run']
  private statePath: string[] = []
  private context: TContext
  private config: MachineConfig<TContext>

  constructor(config: MachineConfig<TContext>) {
    this.config = config
    this.context = { ...config.context }
    // 初始化：进入根状态的 initial
    this.enterState([config.initial])
  }

  // --- 核心逻辑 1: 递归查找节点 ---
  private getNode(path: string[]): StateNode<TContext> | undefined {
    let current: any = this.config
    for (const key of path) {
      if (!current.states || !current.states[key]) return undefined
      current = current.states[key]
    }
    return current
  }

  // --- 核心逻辑 2: 事件处理与冒泡 ---
  send(event: Event | string) {
    const evtObj = typeof event === 'string' ? { type: event } : event
    console.log(`\n[HFSM] Event: ${evtObj.type}`)

    // 从当前最深层的子状态开始，向上冒泡查找处理逻辑
    // 例如路径是 ['moving', 'run']
    // 1. 先看 'moving.run' 是否处理
    // 2. 再看 'moving' 是否处理
    // 3. 最后看 根节点 是否处理

    let handled = false
    // 这是一个从深到浅的循环
    for (let i = this.statePath.length; i >= 0; i--) {
      const currentPath = this.statePath.slice(0, i) // 当前层级的路径
      const node = i === 0 ? this.config : this.getNode(currentPath) // i=0 是根配置

      // 检查该节点是否定义了该事件的转换
      // @ts-ignore
      const transition = node?.on?.[evtObj.type]

      if (transition) {
        console.log(`  -> Handled by state: '${currentPath.join('.') || 'ROOT'}'`)

        // 1. 执行 Transition Actions
        transition.actions?.forEach((fn: any) => this.executeAction(fn, evtObj))

        // 2. 如果有 target，进行状态切换
        if (transition.target) {
          this.transitionTo(transition.target, currentPath, evtObj)
        }

        handled = true
        break // 停止冒泡
      }
    }

    if (!handled) {
      console.log(`  -> Ignored (No handler found in path: ${this.statePath.join('.')})`)
    }
  }

  // --- 核心逻辑 3: 状态切换 (LCA - 最近公共祖先算法) ---
  private transitionTo(targetStr: string, sourcePath: string[], event: Event) {
    // 解析目标路径，支持相对路径和绝对路径
    // 这里简化处理，假设 target 都是绝对路径，如 'idle' 或 'moving.run'
    const targetPath = targetStr.split('.')

    // 1. 找到最近公共祖先 (LCA)
    // 例如：从 ['moving', 'run'] 切换到 ['idle']，LCA 是 [] (根)
    // 例如：从 ['moving', 'run'] 切换到 ['moving', 'walk']，LCA 是 ['moving']
    let lcaIndex = 0
    while (
      lcaIndex < this.statePath.length &&
      lcaIndex < targetPath.length &&
      this.statePath[lcaIndex] === targetPath[lcaIndex]
    ) {
      lcaIndex++
    }

    // 2. Exit: 从当前状态向上退出，直到 LCA
    for (let i = this.statePath.length; i > lcaIndex; i--) {
      const path = this.statePath.slice(0, i)
      const node = this.getNode(path)
      console.log(`  << Exit: ${path.join('.')}`)
      node?.exit?.forEach(fn => this.executeAction(fn, event))
    }

    // 3. Entry: 从 LCA 向下进入，直到目标状态
    // 注意：这里需要处理 initial 状态的递归进入
    // 如果 target 是 'moving'，但 moving 有 initial: 'run'，则最终路径应该是 ['moving', 'run']

    const finalPath = [...targetPath]

    // 补全路径：如果目标节点还有子状态，必须进入其 initial
    let currentNode = this.getNode(finalPath)
    while (currentNode && currentNode.initial) {
      finalPath.push(currentNode.initial)
      currentNode = currentNode.states?.[currentNode.initial]
    }

    // 开始逐层进入
    for (let i = lcaIndex + 1; i <= finalPath.length; i++) {
      const path = finalPath.slice(0, i)
      const node = this.getNode(path)
      console.log(`  >> Entry: ${path.join('.')}`)
      node?.entry?.forEach(fn => this.executeAction(fn, event))
    }

    // 更新当前路径
    this.statePath = finalPath
    console.log(`  == Current State: ${this.statePath.join('.')}`)
  }

  // 初始化进入状态
  private enterState(path: string[]) {
    // 递归补全 initial
    const finalPath = [...path]
    let currentNode = this.getNode(finalPath)
    while (currentNode && currentNode.initial) {
      finalPath.push(currentNode.initial)
      currentNode = currentNode.states?.[currentNode.initial]
    }

    this.statePath = finalPath

    // 触发 Entry (简化版，仅触发最终路径上的)
    for (let i = 1; i <= finalPath.length; i++) {
      const p = finalPath.slice(0, i)
      const node = this.getNode(p)
      node?.entry?.forEach(fn => this.executeAction(fn, { type: 'INIT' }))
    }
    console.log(`[HFSM] Initialized at: ${this.statePath.join('.')}`)
  }

  private executeAction(fn: Action<TContext>, event: Event) {
    const partial = fn(this.context, event)
    // @ts-ignore
    if (partial) Object.assign(this.context, partial)
  }
}

{
  const playerMachine = new HierarchicalStateMachine({
    id: 'player',
    initial: 'idle',
    context: { hp: 100, speed: 0 },
    states: {
      // === 1. 空闲 ===
      idle: {
        entry: [() => console.log('  (Anim) Play Idle Animation')],
        on: {
          MOVE: { target: 'moving' } // 默认进入 moving.walk
        }
      },

      // === 2. 移动 (父状态) ===
      moving: {
        initial: 'walk', // 默认子状态

        // 父状态的 Entry/Exit 会包裹子状态
        entry: [() => console.log('  (Logic) Enable Physics')],
        exit: [() => console.log('  (Logic) Disable Physics')],

        // **关键：父状态处理公共事件**
        on: {
          STOP: { target: 'idle' }, // 无论在走还是跑，STOP 都回 idle
          HIT: {
            actions: [
              ctx => {
                ctx.hp -= 10
                console.log(`  (Logic) Ouch! HP: ${ctx.hp}`)
              }
            ]
          }
        },

        states: {
          // --- 2.1 走 ---
          walk: {
            entry: [ctx => ({ speed: 5 })], // 设置速度
            on: {
              SHIFT_DOWN: { target: 'moving.run' } // 切换兄弟状态
            }
          },

          // --- 2.2 跑 ---
          run: {
            entry: [ctx => ({ speed: 10 })],
            on: {
              SHIFT_UP: { target: 'moving.walk' }
            }
          }
        }
      }
    }
  })

  // --- 运行演示 ---

  // 1. 开始移动 (Idle -> Moving -> Walk)
  // 触发顺序：Idle Exit -> Moving Entry -> Walk Entry
  playerMachine.send('MOVE')

  // 2. 加速 (Walk -> Run)
  // 触发顺序：Walk Exit -> Run Entry (注意：Moving 没有退出！)
  playerMachine.send('SHIFT_DOWN')

  // 3. 受到攻击 (冒泡机制)
  // Run 状态没有定义 HIT，事件冒泡给 Moving，Moving 处理了 HIT
  playerMachine.send('HIT')

  // 4. 停止 (Run -> Idle)
  // Run 状态没有定义 STOP，冒泡给 Moving，Moving 定义了去 Idle
  // 触发顺序：Run Exit -> Moving Exit -> Idle Entry
  playerMachine.send('STOP')
}
