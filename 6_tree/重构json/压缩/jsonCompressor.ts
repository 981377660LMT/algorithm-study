type JsonValue = string | number | boolean | null | JsonObject | JsonArray
type JsonArray = Array<JsonValue>
interface JsonObject {
  [key: string]: JsonValue
}
type PathSegment = string | number

enum Action {
  /** 丢弃当前节点及其子树 */
  DROP = 'DROP',
  /** 保留当前节点，但继续递归处理子节点 */
  KEEP = 'KEEP',
  /** 保留当前节点及其所有后代（不再递归，性能优化关键） */
  KEEP_TREE = 'KEEP_TREE',
  /** 当前规则不处理，交给下一个规则 */
  SKIP = 'SKIP'
}

/**
 * 遍历上下文
 * @template U 用户自定义上下文类型 (User Context)
 */
interface WalkContext<U = unknown> {
  key: PathSegment | null // 根节点 key 为 null
  value: JsonValue
  parent: JsonValue | null // 根节点 parent 为 null
  path: PathSegment[] // 当前节点的路径
  userContext: U // 用户传入的全局上下文
}

type Rule<U = unknown> = (ctx: WalkContext<U>) => Action

export class JsonCompressor<U = unknown> {
  private rules: Rule<U>[]

  constructor(rules: Rule<U>[]) {
    this.rules = rules
  }

  /**
   * 执行压缩
   * @param root 根节点
   * @param userContext 用户上下文数据
   */
  compress(root: JsonValue, userContext: U): JsonValue {
    return this.walk(root, null, null, [], userContext)
  }

  private runRules(ctx: WalkContext<U>): Action {
    for (const rule of this.rules) {
      const action = rule(ctx)
      if (action !== Action.SKIP) {
        return action
      }
    }
    // 默认行为：如果是对象/数组，继续递归；如果是叶子，默认丢弃 (Strict Mode)
    // 或者：默认保留。这里采用“默认继续递归，叶子节点需显式保留”的策略更安全
    return Action.SKIP
  }

  private isValid(val: JsonValue): boolean {
    if (val === null || val === undefined) return false
    if (Array.isArray(val)) return val.length > 0
    if (typeof val === 'object') return Object.keys(val).length > 0
    return true
  }

  private walk(
    value: JsonValue,
    key: PathSegment | null,
    parent: JsonValue,
    path: PathSegment[],
    userContext: U
  ): JsonValue {
    if (value == null) return null

    const ctx: WalkContext<U> = { key, value, parent, path, userContext }

    // 1. 决策阶段
    let action = this.runRules(ctx)

    // 默认策略补充：如果规则链全部 SKIP
    if (action === Action.SKIP) {
      if (typeof value === 'object') {
        action = Action.KEEP // 对象默认进入递归
      } else {
        action = Action.DROP // 叶子节点默认丢弃（白名单模式）
        // 如果想要黑名单模式（默认保留），改为: action = Action.KEEP;
      }
    }

    if (action === Action.DROP) return null
    if (action === Action.KEEP_TREE) return value

    // 2. 递归阶段
    if (Array.isArray(value)) {
      const result = value
        .map((item, index) => this.walk(item, index, value, [...path, index], userContext))
        .filter(item => this.isValid(item))
      return result.length > 0 ? result : null
    }

    if (typeof value === 'object') {
      const result: JsonObject = {}
      for (const k of Object.keys(value)) {
        const v = this.walk(value[k], k, value, [...path, k], userContext)
        if (this.isValid(v)) {
          result[k] = v
        }
      }
      return Object.keys(result).length > 0 ? result : null
    }

    // 3. 叶子节点结果
    return action === Action.KEEP ? value : null
  }
}

const Rules = {
  /**
   * 黑名单：命中 Key 则丢弃
   */
  omit(keys: string[]): Rule {
    const set = new Set(keys)
    return ({ key }) => (key && typeof key === 'string' && set.has(key) ? Action.DROP : Action.SKIP)
  },

  /**
   * 深层保留：命中 Key 则保留整棵子树 (性能优化)
   */
  deepKeep(keys: string[]): Rule {
    const set = new Set(keys)
    return ({ key }) =>
      key && typeof key === 'string' && set.has(key) ? Action.KEEP_TREE : Action.SKIP
  },

  /**
   * 白名单：命中 Key 则保留 (通常用于叶子节点)
   */
  pick(keys: string[]): Rule {
    const set = new Set(keys)
    return ({ key }) => (key && typeof key === 'string' && set.has(key) ? Action.KEEP : Action.SKIP)
  },

  /**
   * 正则匹配 Key
   */
  matches(pattern: RegExp, action: Action = Action.KEEP): Rule {
    return ({ key }) => (key && typeof key === 'string' && pattern.test(key) ? action : Action.SKIP)
  },

  /**
   * 动态值检测 (针对 {{value}} 场景)
   */
  dynamicValue<U>(predicate: (val: string, ctx: U) => boolean): Rule<U> {
    return ({ value, userContext }) => {
      if (typeof value === 'string' && predicate(value, userContext)) {
        return Action.KEEP
      }
      return Action.SKIP
    }
  },

  /**
   * 自定义断言
   */
  custom<U>(fn: (ctx: WalkContext<U>) => boolean, action: Action = Action.KEEP): Rule<U> {
    return ctx => (fn(ctx) ? action : Action.SKIP)
  },

  /**
   * 调试用：打印日志
   */
  debug(label: string = 'DEBUG'): Rule {
    return ({ key, path, value }) => {
      console.log(`[${label}]`, { path: path.join('.'), key, type: typeof value })
      return Action.SKIP
    }
  }
}

{
  type MyContext = Set<string>

  // 2. 准备业务规则
  const lowCodeRules = [
    // A. 黑名单 (优先级最高)
    Rules.omit(['style', 'className', 'meta', 'schema', 'dsl', '_internal']),

    // B. 深层保留 (性能优化，命中后不再递归)
    Rules.deepKeep(['event', 'events', 'eventList', 'Interaction', 'action', 'staticData']),

    // C. 浅层白名单 (保留特定属性)
    Rules.pick(['name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid', 'type']),

    // D. 正则匹配 (保留所有包含 vent 的 key)
    Rules.matches(/vent/i, Action.KEEP),

    // E. 动态值检测 (核心业务逻辑)
    Rules.dynamicValue<MyContext>((val, contextSet) => {
      // 1. 检查插值表达式 {{...}}
      if (val.includes('{{')) return true

      // 2. 检查是否在上下文变量中
      // 简单的全等检查，或者你可以复用之前的正则提取逻辑
      return contextSet.has(val)
    }),

    // F. 兜底策略：如果是对象或数组，继续递归；否则丢弃
    // (注意：JsonCompressor 内部默认逻辑已涵盖此点，但显式写出更清晰)
    Rules.custom(({ value }) => typeof value === 'object' && value !== null, Action.KEEP)
  ]

  // 3. 创建压缩器实例
  const compressor = new JsonCompressor<MyContext>(lowCodeRules)

  // ============================================================================
  // 测试数据
  // ============================================================================
  const componentTree = {
    id: 'root',
    type: 'Page',
    style: { width: '100px' }, // 应被删除 (omit)
    meta: { author: 'admin' }, // 应被删除 (omit)
    events: {
      // 应被深层保留 (deepKeep)
      click: {
        action: 'submit',
        payload: { a: 1 } // 即使不在白名单，因为父级是 events，也会保留
      }
    },
    children: [
      {
        type: 'Button',
        props: {
          text: 'Submit', // 普通字符串，不在白名单 -> 删除
          dynamicText: '{{title}}', // 动态值 -> 保留
          user: 'currentUser' // 上下文变量 -> 保留
        }
      }
    ]
  }

  const globalContextNames = new Set(['currentUser'])

  // 4. 执行
  const result = compressor.compress(componentTree, globalContextNames)

  console.log(JSON.stringify(result, null, 2))
}
