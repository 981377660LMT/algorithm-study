/**
 * 标准中断令牌 (CancellationToken)
 */
export interface CancellationToken {
  readonly isCancellationRequested: boolean
  throwIfCancellationRequested(): void
}

export interface ISearchState {
  keyword: string
  results: Map<string, ISearchResult[]>
  activeResult: ISearchResult | undefined
}

/**
 * 实际项目可用 mobx，方法更加精细化.
 */
export interface ISearchStore {
  subscribe(listener: (state: ISearchState, prevState: ISearchState) => void): () => void
  dispose(): void

  setModuleResults(moduleId: string, results: ISearchResult[]): void
  getModuleResults(moduleId: string): ISearchResult[]
  getAllModuleResults(): Map<string, ISearchResult[]>
  setActiveResult(result: ISearchResult | undefined): void
  getActiveResult(): ISearchResult | undefined
}

export interface ISearchResult {
  /**
   * 结果的唯一标识 (用于 React key).
   * 需要返回稳定的指纹（例如 文件路径 + 行号 + 列号 的哈希).
   */
  readonly key: string
  /** 来源模块实例，用于调用 reveal/replace. */
  readonly moduleId: string
  readonly canReplace: boolean
}

export interface ISearchModuleContext {
  /**
   * 模块内容变更时，主动请求重搜.
   */
  research: () => void
}

type SyncOrAsync<T> = T | Promise<T>
export interface ISearchModule<R extends ISearchResult = ISearchResult> {
  readonly id: string
  /**
   * 数值越小越靠前.
   */
  readonly priority: () => number

  /**
   * 生命周期：插件注册时调用，注入上下文能力.
   */
  onInit(context: ISearchModuleContext): void

  /**
   * 生命周期：插件注销时调用，清理资源.
   */
  onDispose?(): void

  search(
    keyword: string,
    context: {
      token: CancellationToken
    }
  ): SyncOrAsync<R[]>

  /**
   * 告诉视图层“现在有哪些结果”。通常用于更新侧边栏列表、编辑器中的高亮装饰器 (Decorations)。
   * 触发时机：结果集发生变化时（搜索完成、清空搜索、替换导致文档变化）。
   * 注意是覆盖式更新，而不是增量式更新。
   */
  render?(results: R[], activeResult: R | undefined): SyncOrAsync<void>

  /**
   * 告诉视图层“请看这一个结果”。通常用于滚动视口、展开折叠代码块、移动光标。
   * 当前激活项 (activeResult) 发生变化时（用户点击列表、点击“下一个”、替换后自动跳到下一个）
   *
   * 语义比 scrollTo 更丰富，隐含了：
   * 1. 如果在折叠区域内，需要先展开
   * 2. 如果在可视区域外，需要滚动
   * 3. 可能伴随高亮或光标移动
   */
  reveal?(result: R): SyncOrAsync<void>

  replace?(result: R, newText: string): SyncOrAsync<void>
  replaceAll?(results: R[], newText: string): SyncOrAsync<void>
}
