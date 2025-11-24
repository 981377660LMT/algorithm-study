export {}

// ==========================================
// 1. 基础类型定义
// ==========================================

/**
 * 简单的事件订阅发布器
 */
type Listener<T> = (event: T) => void
type Disposable = () => void
class EventEmitter<T> {
  private listeners: Set<Listener<T>> = new Set()

  subscribe(listener: Listener<T>): Disposable {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  emit(event: T) {
    this.listeners.forEach(l => l(event))
  }

  clear() {
    this.listeners.clear()
  }
}

/**
 * 标准中断令牌 (CancellationToken)
 */
export interface CancellationToken {
  readonly isCancellationRequested: boolean
  throwIfCancellationRequested(): void
}

class CancellationTokenSource {
  private _isCancelled = false

  get token(): CancellationToken {
    return {
      isCancellationRequested: this._isCancelled,
      throwIfCancellationRequested: () => {
        if (this._isCancelled) throw new Error('OperationCancelled')
      }
    }
  }

  cancel() {
    this._isCancelled = true
  }
}

// 1. 搜索结果 (泛型 T: 插件自定义的定位数据)
export interface ISearchResult<T = unknown> {
  readonly id: string // 全局唯一 ID
  readonly pluginId: string // 来源插件 ID
  readonly matchContent: string // 匹配文本
  readonly payload: T // 定位数据
  readonly canReplace: boolean
}

// 2. 插件上下文 (Controller -> Plugin)
export interface ISearchPluginContext {
  // 数据变动时，插件主动请求局部重搜
  reSearch: () => void
}

type SyncOrAsync<T> = T | Promise<T>
// 3. 插件接口 (Plugin)
export interface ISearchPlugin<T = unknown> {
  readonly id: string
  /**
   * 动态优先级：每次搜索开始时调用，决定模块在结果列表中的顺序。
   * 数值越小越靠前。
   */
  readonly priority: () => number

  search(keyword: string, token: CancellationToken): SyncOrAsync<ISearchResult<T>[]>
  render(results: ISearchResult<T>[], activeResult: ISearchResult<T> | null): void
  scrollTo(result: ISearchResult<T>): void

  replace?(result: ISearchResult<T>, newText: string): SyncOrAsync<void>
  replaceAll?(results: ISearchResult<T>[], newText: string): SyncOrAsync<void>

  clear?(): void
}

// ==========================================
// 2. Model 层: SearchStore (纯状态容器)
// ==========================================

interface SearchState {
  readonly keyword: string
  readonly isSearching: boolean
  readonly results: ReadonlyArray<ISearchResult<any>>
  readonly currentIndex: number
  readonly currentResult: ISearchResult<any> | null
}

const INITIAL_STATE: SearchState = {
  keyword: '',
  isSearching: false,
  results: [],
  currentIndex: -1,
  currentResult: null
}

export class SearchStore extends EventEmitter<SearchState> {
  private _state: SearchState = INITIAL_STATE

  get state() {
    return this._state
  }

  /**
   * 唯一修改状态的入口
   */
  setState(partial: Partial<SearchState>) {
    const newState = { ...this._state, ...partial }
    // 简单的浅比较，避免无意义的 emit (可选优化)
    // if (JSON.stringify(newState) === JSON.stringify(this._state)) return;

    this._state = newState
    this.emit(this._state)
  }

  reset() {
    this.setState(INITIAL_STATE)
  }
}

// ==========================================
// 3. Worker 层: SearchTask
// ==========================================

class SearchTask {
  // 全局任务取消源 (用户输入新词时触发)
  private mainCancelSource = new CancellationTokenSource()

  // 插件级取消源 (reSearch 时触发，互斥同一插件的多次搜索)
  private pluginCancelSources = new Map<string, CancellationTokenSource>()

  // 结果缓存
  private resultsMap = new Map<string, ISearchResult<any>[]>()

  private sortedPlugins: ISearchPlugin<any>[]

  private timeoutTimer: any = null

  constructor(
    public readonly keyword: string,
    // 传入已按 priority 排序的插件列表
    plugins: ISearchPlugin<any>[],
    private onUpdate: (flatResults: ISearchResult<any>[]) => void,
    private timeoutMs: number = 10000
  ) {
    this.sortedPlugins = plugins.slice()
  }

  /**
   * 启动全量搜索
   */
  async run() {
    // 启动超时保护
    this.timeoutTimer = setTimeout(() => {
      if (!this.mainCancelSource.token.isCancellationRequested) {
        console.warn(`SearchTask timeout after ${this.timeoutMs}ms`)
        this.cancel()
      }
    }, this.timeoutMs)

    const promises = this.sortedPlugins.map(p => this.searchSinglePlugin(p))
    await Promise.allSettled(promises)

    if (this.timeoutTimer) clearTimeout(this.timeoutTimer)
  }

  /**
   * 搜索单个插件 (支持 reSearch)
   */
  async searchSinglePlugin(plugin: ISearchPlugin<any>) {
    // 1. 如果全局任务已取消，直接终止
    if (this.mainCancelSource.token.isCancellationRequested) return

    // 2. 取消该插件上一次正在进行的搜索 (reSearch 互斥)
    const prevSource = this.pluginCancelSources.get(plugin.id)
    if (prevSource) {
      prevSource.cancel()
    }

    // 3. 创建本次搜索的专属 Token
    const localSource = new CancellationTokenSource()
    this.pluginCancelSources.set(plugin.id, localSource)

    try {
      // 4. 执行搜索
      // 注意：插件内部应该定期检查 token.isCancellationRequested
      const results = await plugin.search(this.keyword, localSource.token)

      // 5. 提交前检查：全局是否取消？局部是否被新 reSearch 取消？
      if (this.mainCancelSource.token.isCancellationRequested) return
      if (localSource.token.isCancellationRequested) return

      // 6. 更新结果并通知
      this.resultsMap.set(plugin.id, results)
      this.aggregateAndNotify()
    } catch (e: any) {
      if (e.message !== 'OperationCancelled') {
        console.error(`Plugin ${plugin.id} search error:`, e)
      }
    } finally {
      // 清理：只有当 map 中的 source 还是当前 source 时才删除
      // (防止误删了新发起的 reSearch source)
      if (this.pluginCancelSources.get(plugin.id) === localSource) {
        this.pluginCancelSources.delete(plugin.id)
      }
    }
  }

  /**
   * 新增：动态添加插件
   * 场景：搜索进行中，某个懒加载模块被注册
   */
  addPlugin(plugin: ISearchPlugin<any>) {
    // 直接 push，依赖 aggregateAndNotify 中的 sort 保证顺序
    this.sortedPlugins.push(plugin)
    this.searchSinglePlugin(plugin)
  }

  /**
   * 新增：动态移除插件
   * 场景：搜索进行中，某个模块被卸载
   */
  removePlugin(pluginId: string) {
    // 1. 从列表中移除 (防止内存泄漏和无效遍历)
    this.sortedPlugins = this.sortedPlugins.filter(p => p.id !== pluginId)

    // 2. 清理结果
    const source = this.pluginCancelSources.get(pluginId)
    source?.cancel()
    this.pluginCancelSources.delete(pluginId)

    const cached = this.resultsMap.get(pluginId)
    if (cached && cached.length > 0) {
      this.resultsMap.delete(pluginId)
      this.aggregateAndNotify()
    } else {
      this.resultsMap.delete(pluginId)
    }
  }

  private aggregateAndNotify() {
    this.sortedPlugins.sort((a, b) => a.priority() - b.priority())

    const flatResults: ISearchResult<any>[] = []

    // 按照 sortedPlugins 的顺序聚合，确保 priority 生效
    for (const plugin of this.sortedPlugins) {
      const pluginResults = this.resultsMap.get(plugin.id)
      if (pluginResults) {
        flatResults.push(...pluginResults)
      }
    }

    this.onUpdate(flatResults)
  }

  cancel() {
    if (this.timeoutTimer) clearTimeout(this.timeoutTimer)
    this.mainCancelSource.cancel()
    this.pluginCancelSources.forEach(s => s.cancel())
    this.pluginCancelSources.clear()
  }

  get isCancelled() {
    return this.mainCancelSource.token.isCancellationRequested
  }
}

// ==========================================
// 4. Controller 层: SearchController (业务逻辑)
// ==========================================
export class SearchController {
  public readonly store = new SearchStore()

  private plugins = new Map<string, ISearchPlugin<any>>()
  private currentTask: SearchTask | null = null

  registerPlugin<T>(plugin: ISearchPlugin<T>): ISearchPluginContext {
    this.plugins.set(plugin.id, plugin)

    if (this.currentTask && !this.currentTask.isCancelled) {
      this.currentTask.addPlugin(plugin)
    }

    return {
      reSearch: () => {
        // 仅当任务未被全局取消时，允许局部重搜
        if (this.currentTask && !this.currentTask.isCancelled) {
          this.currentTask.searchSinglePlugin(plugin)
        }
      }
    }
  }

  unregisterPlugin(pluginId: string) {
    this.plugins.get(pluginId)?.clear?.()
    this.plugins.delete(pluginId)
    if (this.currentTask && !this.currentTask.isCancelled) {
      this.currentTask.removePlugin(pluginId)
    }
  }

  /**
   * 搜索入口 (无防抖，由 View 层负责)
   */
  search(keyword: string) {
    // 1. 取消旧任务
    this.currentTask?.cancel()

    // 2. 更新状态
    this.store.setState({
      keyword,
      isSearching: !!keyword,
      results: [],
      currentIndex: -1,
      currentResult: null
    })
    this.clearAllHighlights()

    if (!keyword) {
      this.currentTask = null
      return
    }

    // 3. 准备插件列表 (按 priority 排序)
    const sortedPlugins = Array.from(this.plugins.values()).sort((a, b) => {
      const pA = a.priority() // 动态调用
      const pB = b.priority() // 动态调用
      return pA - pB
    })

    // 4. 创建并启动新任务
    this.currentTask = new SearchTask(keyword, sortedPlugins, results =>
      this.handleTaskUpdate(results)
    )

    this.currentTask.run().finally(() => {
      // 只有当任务还是当前任务时，才更新 loading 状态
      if (this.currentTask?.keyword === keyword && !this.currentTask.isCancelled) {
        this.store.setState({ isSearching: false })
      }
    })
  }

  next() {
    this.moveIndex(1)
  }
  prev() {
    this.moveIndex(-1)
  }

  async replace(newText: string) {
    const { currentResult } = this.store.state
    if (!currentResult || !currentResult.canReplace) return

    const plugin = this.plugins.get(currentResult.pluginId)
    if (plugin?.replace) {
      await plugin.replace(currentResult, newText)
      // 插件数据变化后，应通过 context.reSearch() 自动触发更新
    }
  }

  async replaceAll(newText: string) {
    const { results } = this.store.state
    if (results.length === 0) return

    // 按插件分组
    const resultsByPlugin = new Map<string, ISearchResult<any>[]>()
    for (const res of results) {
      if (!res.canReplace) continue
      const list = resultsByPlugin.get(res.pluginId) || []
      list.push(res)
      resultsByPlugin.set(res.pluginId, list)
    }

    // 并行执行
    const promises: SyncOrAsync<void>[] = []
    for (const [pluginId, pluginResults] of resultsByPlugin) {
      const plugin = this.plugins.get(pluginId)
      if (plugin?.replaceAll) {
        promises.push(plugin.replaceAll(pluginResults, newText))
      }
    }
    await Promise.all(promises)
  }

  /**
   * 资源清理
   */
  dispose() {
    this.currentTask?.cancel()
    this.plugins.forEach(p => p.clear?.())
    this.plugins.clear()
  }

  private handleTaskUpdate(newResults: ISearchResult<any>[]) {
    let newIndex = -1
    const { currentResult, currentIndex } = this.store.state

    if (newResults.length > 0) {
      if (currentResult) {
        newIndex = newResults.findIndex(
          r => r.id === currentResult.id && r.pluginId === currentResult.pluginId
        )
      }
      if (newIndex === -1) {
        newIndex = currentIndex >= 0 && currentIndex < newResults.length ? currentIndex : 0
      }
    }

    this.store.setState({
      results: newResults,
      currentIndex: newIndex,
      currentResult: newResults[newIndex] || null
    })

    this.refreshView(true)
  }

  private moveIndex(delta: number) {
    const { results, currentIndex } = this.store.state
    if (results.length === 0) return

    const nextIndex = (currentIndex + delta + results.length) % results.length

    this.store.setState({
      currentIndex: nextIndex,
      currentResult: results[nextIndex]
    })

    this.refreshView(true)
  }

  // 不做复杂的竞态保护
  // 绝大多数前端场景下，JavaScript 是单线程执行的，且 UI 渲染（尤其是 DOM 操作）通常很快。
  // 除非插件内部有非常耗时的异步操作（比如网络请求或复杂的 WebWorker 计算），否则引入版本号确实显得“过重”且增加了维护成本。
  private refreshView(shouldScroll: boolean) {
    const { results, currentResult } = this.store.state

    this.plugins.forEach(plugin => {
      try {
        const pluginResults = results.filter(r => r.pluginId === plugin.id)
        const active = currentResult && currentResult.pluginId === plugin.id ? currentResult : null

        plugin.render(pluginResults, active)

        if (shouldScroll && active) {
          plugin.scrollTo(active)
        }
      } catch (e) {
        console.error(`Plugin ${plugin.id} render/scrollTo failed:`, e)
      }
    })
  }

  private clearAllHighlights() {
    this.plugins.forEach(p => {
      try {
        p.render([], null)
      } catch (e) {
        console.error(e)
      }
    })
  }
}
