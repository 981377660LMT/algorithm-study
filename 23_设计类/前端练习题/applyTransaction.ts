// 核心类型定义 (MyState, MyAction, MyModule)。
// 核心 Reducer (coreReducer)，负责纯粹的状态转换。
// 状态管理器 (StateManager)，包含复现自 ProseMirror 的“事务反馈循环”的核心 dispatch 方法。
// 三个示例模块 (loggerModule, autoCorrectModule, readOnlyModule)，展示了插件如何工作。
// 使用示例，演示了如何初始化和运行整个系统。

// ===================================================================
// 1. 核心类型定义
// ===================================================================

/**
 * 你的项目的核心状态。
 * 必须是不可变的，每次更新都应创建新对象。
 */
export interface MyState {
  readonly data: any
  readonly version: number // 版本号用于追踪变化，便于调试
}

/**
 * 描述状态变化的 "动作" (Action)，对应 ProseMirror 的 Transaction。
 */
export interface MyAction {
  type: string
  payload?: any
}

/**
 * 模块/插件的规范，定义了模块可以实现的钩子。
 */
export interface MyModule {
  /**
   * 钩子1: 在动作应用前进行过滤 (可选)。
   * 返回 false 可以取消整个动作。
   */
  filterAction?: (action: MyAction, state: MyState) => boolean

  /**
   * 钩子2: 在状态变化后追加新的动作 (可选)。
   * `actions` 是刚刚被处理的动作。
   * `oldState` 是本轮 dispatch 开始前的状态。
   * `newState` 是应用 `actions` 后的最新状态。
   * 返回一个新的动作对象来触发连锁反应，或返回 null 不做任何事。
   */
  appendAction?: (
    actions: readonly MyAction[],
    oldState: MyState,
    newState: MyState
  ) => MyAction | null
}

// ===================================================================
// 2. 核心 Reducer
// ===================================================================

/**
 * 纯函数：根据动作和旧状态计算新状态。
 * 这是你项目的主要业务逻辑所在。
 * @param state 旧状态
 * @param action 要应用的动作
 * @returns 新状态
 */
function coreReducer(state: MyState, action: MyAction): MyState {
  // 重点：必须返回一个全新的对象以保证不可变性！
  switch (action.type) {
    case 'UPDATE_DATA':
      return {
        ...state,
        data: action.payload,
        version: state.version + 1
      }
    case 'RESET_DATA':
      return {
        ...state,
        data: null,
        version: state.version + 1
      }
    default:
      // 对于未知的动作类型，直接返回原状态
      return state
  }
}

// ===================================================================
// 3. 状态管理器 (StateManager)
// ===================================================================

export class StateManager {
  public currentState: MyState
  private readonly modules: readonly MyModule[]

  constructor(initialState: MyState, modules: readonly MyModule[] = []) {
    this.currentState = initialState
    this.modules = modules
  }

  /**
   * 分发一个动作，并驱动状态更新。
   * 这是复现 ProseMirror 事务反馈循环的核心。
   * @param rootAction 初始动作
   * @returns 返回最终的状态和所有被处理的动作
   */
  dispatch(rootAction: MyAction): { state: MyState; actions: readonly MyAction[] } {
    // --- 步骤 1: 初始过滤 ---
    // 任何一个模块的 filterAction 返回 false 都会阻止整个流程
    for (const module of this.modules) {
      if (module.filterAction && !module.filterAction(rootAction, this.currentState)) {
        console.log(`[StateManager] Action ${rootAction.type} was cancelled by a module.`)
        return { state: this.currentState, actions: [] }
      }
    }

    // --- 步骤 2: 初始应用 ---
    const appliedActions: MyAction[] = [rootAction]
    let newState = coreReducer(this.currentState, rootAction)
    const initialStateForAppend = this.currentState

    // --- 步骤 3: 进入反馈循环 ---
    // 这个循环会一直运行，直到没有任何模块再追加新的动作为止
    for (;;) {
      let haveNewActions = false

      // 遍历所有模块，看它们是否要对当前的变化做出反应
      for (const module of this.modules) {
        if (module.appendAction) {
          const appendedAction = module.appendAction(
            appliedActions, // 传递所有已应用的动作
            initialStateForAppend,
            newState
          )

          // 如果模块返回了一个新的动作...
          if (appendedAction) {
            let allowed = true
            // 这个新追加的动作也必须经过所有模块的过滤
            for (const filterModule of this.modules) {
              if (
                filterModule.filterAction &&
                !filterModule.filterAction(appendedAction, newState)
              ) {
                allowed = false
                break
              }
            }

            if (allowed) {
              // 如果允许，则应用它，并准备进行下一轮循环
              appliedActions.push(appendedAction)
              newState = coreReducer(newState, appendedAction)
              haveNewActions = true
            }
          }
        }
      }

      // --- 步骤 4: 退出循环 ---
      // 如果完整的一轮循环没有产生任何新动作，说明系统已稳定
      if (!haveNewActions) {
        this.currentState = newState
        return { state: this.currentState, actions: appliedActions }
      }
    }
  }
}

// ===================================================================
// 4. 示例模块
// ===================================================================

/**
 * 示例模块1: 一个日志模块，用于观察流程。
 */
const loggerModule: MyModule = {
  appendAction: (actions, oldState, newState) => {
    // 只对最新发生的动作打印日志
    const lastAction = actions[actions.length - 1]
    console.log(
      `[LoggerModule] Action processed: ${lastAction.type}, State version: ${oldState.version} -> ${newState.version}`
    )
    return null // 日志模块不追加新动作
  }
}

/**
 * 示例模块2: 一个自动纠错模块，演示追加新动作。
 */
const autoCorrectModule: MyModule = {
  appendAction: (actions, oldState, newState) => {
    // 检查最新的状态数据是否为 "teh"
    if (newState.data === 'teh') {
      // 检查是否已经纠正过，防止无限循环
      const alreadyCorrected = actions.some(a => a.type === 'UPDATE_DATA' && a.payload === 'the')
      if (!alreadyCorrected) {
        console.log('[AutoCorrectModule] Detected "teh", appending a correction action.')
        // 追加一个新动作来纠错
        return { type: 'UPDATE_DATA', payload: 'the' }
      }
    }
    return null
  }
}

/**
 * 示例模块3: 一个只读模块，演示过滤动作。
 */
const readOnlyModule: MyModule = {
  filterAction: (action, state) => {
    // 拒绝所有会更新数据的动作
    if (action.type === 'UPDATE_DATA') {
      console.log(`[ReadOnlyModule] Blocked action of type: ${action.type}.`)
      return false
    }
    return true
  }
}

// ===================================================================
// 5. 使用示例
// ===================================================================

function runDemo() {
  console.log('--- DEMO 1: Auto-correction feedback loop ---')
  const initialState: MyState = { data: 'initial', version: 0 }
  const modules = [loggerModule, autoCorrectModule]
  const stateManager = new StateManager(initialState, modules)

  // 用户输入了 "teh"
  const result1 = stateManager.dispatch({ type: 'UPDATE_DATA', payload: 'teh' })

  console.log('\n--- DEMO 1 FINAL RESULT ---')
  console.log('Final state:', result1.state)
  console.log(
    'All processed actions:',
    result1.actions.map(a => a.type)
  )
  console.log('----------------------------------------\n')

  console.log('--- DEMO 2: Action filtering ---')
  const modulesWithReadOnly = [loggerModule, autoCorrectModule, readOnlyModule]
  const stateManager2 = new StateManager(initialState, modulesWithReadOnly)

  // 尝试更新数据，但会被 readOnlyModule 阻止
  const result2 = stateManager2.dispatch({ type: 'UPDATE_DATA', payload: 'some new data' })

  console.log('\n--- DEMO 2 FINAL RESULT ---')
  console.log('Final state (should be unchanged):', result2.state)
  console.log('Processed actions (should be empty):', result2.actions)
  console.log('----------------------------------------\n')
}

runDemo()
