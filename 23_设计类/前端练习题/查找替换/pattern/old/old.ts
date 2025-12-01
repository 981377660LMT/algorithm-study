// TODO:
// - research
// - Orchestor 太重了，抽出 Task/Modal/View 等概念
// replace 逻辑

// --- 1. 基础类型定义 ---

// 搜索配置
interface SearchOptions {
  matchCase: boolean
  useRegex: boolean
  wholeWord: boolean
}

// 搜索上下文
interface SearchContext {
  query: string
  options: SearchOptions
  cancellationToken: { isCancelled: boolean }
}

// 提供给子模块的上下文能力
interface ProviderContext {
  /**
   * 请求重新搜索（例如数据发生变化时）
   * Orchestrator 会自动进行防抖处理，子模块无需自己防抖
   */
  requestSearch: () => void
}

// 泛型位置接口
interface MatchLocation {
  id: string
  description: string
}

// 统一的匹配结果
// T: 具体的渲染上下文类型
interface SearchResult<T = any> {
  providerId: string // 来源模块 ID
  location: MatchLocation
  matchText: string
  range: [number, number]
  score: number
  renderData: T

  // 标记该结果是否已被替换/失效
  invalidated?: boolean
}

// --- 2. 子模块接口 (核心生命周期 + UI 能力 + 替换能力 + 数据监听) ---

interface SearchProvider<T = any> {
  readonly id: string
  readonly priority: number

  // --- 生命周期 ---

  /**
   * [修改] 注册时注入上下文能力
   * 子模块应保存 context.requestSearch 以便在数据变化时调用
   */
  onRegister?(context: ProviderContext): void | Promise<void>

  onSearchStart?(context: SearchContext): void
  search(context: SearchContext): Promise<SearchResult<T>[]>
  onSearchEnd?(): void
  onDispose?(): void

  // --- 替换能力 ---
  replace?(match: SearchResult<T>, newText: string): Promise<boolean>
  replaceAll?(matches: SearchResult<T>[], newText: string): Promise<void>

  // --- UI 能力 ---
  highlight(results: SearchResult<T>[]): void
  highlightActive?(result?: SearchResult<T>): void
  clearHighlights(): void
  focus(result: SearchResult<T>): Promise<void>
}

// --- 3. 搜索编排器 (Orchestrator) ---

class SearchOrchestrator<T = any> {
  private providers = new Map<string, SearchProvider<T>>()
  private currentCancellationToken = { isCancelled: false }

  // 状态管理
  private currentResults: SearchResult<T>[] = []
  private currentFocusIndex: number = -1

  // [新增] 数据变更处理相关
  private debounceTimer: any = null
  private lastQuery: string = ''
  private lastOptions: SearchOptions | null = null

  // 注册子模块
  async registerProvider(provider: SearchProvider<T>) {
    if (this.providers.has(provider.id)) {
      console.warn(`Provider ${provider.id} already registered.`)
      return
    }

    // [修改] 构造上下文并注入给 Provider
    const context: ProviderContext = {
      requestSearch: () => {
        console.log(`[System] Search requested by provider: ${provider.id}`)
        this.triggerDebouncedSearch()
      }
    }

    if (provider.onRegister) {
      await provider.onRegister(context)
    }
    this.providers.set(provider.id, provider)

    console.log(`[System] Provider registered: ${provider.id}`)
  }

  // 卸载子模块
  unregisterProvider(providerId: string) {
    const provider = this.providers.get(providerId)
    if (!provider) return

    provider.clearHighlights()
    if (provider.onDispose) {
      provider.onDispose()
    }
    this.providers.delete(providerId)
    console.log(`[System] Provider disposed: ${providerId}`)
  }

  // [新增] 防抖重搜
  private triggerDebouncedSearch() {
    if (this.debounceTimer) clearTimeout(this.debounceTimer)

    this.debounceTimer = setTimeout(() => {
      if (this.lastQuery && this.lastOptions) {
        this.refreshSearch()
      }
    }, 300) // 300ms 防抖
  }

  // [新增] 刷新搜索 (核心：保持焦点)
  private async refreshSearch() {
    console.log('[System] Refreshing search results due to data change...')

    // 1. 记住当前聚焦的 ID (Focus Preservation)
    let activeMatchId: string | null = null
    if (this.currentFocusIndex !== -1 && this.currentResults[this.currentFocusIndex]) {
      // 假设 location.id 是稳定的 (例如 "Row1:Col1")
      activeMatchId = this.currentResults[this.currentFocusIndex].location.id
    }

    // 2. 重新执行搜索
    // 注意：search 方法会重置 index 为 -1，所以我们需要在 await 之后恢复
    await this.search(this.lastQuery, this.lastOptions!)

    // 3. 尝试找回焦点
    if (activeMatchId) {
      const newIndex = this.currentResults.findIndex(r => r.location.id === activeMatchId)
      if (newIndex !== -1) {
        this.currentFocusIndex = newIndex
        // 重新绘制 Active 边框
        const match = this.currentResults[newIndex]
        const provider = this.providers.get(match.providerId)
        provider?.highlightActive?.(match)
        console.log(`[System] Focus preserved at index ${newIndex}`)
      } else {
        console.log(`[System] Previous focus lost (item might be deleted or not matching anymore).`)
        this.currentFocusIndex = -1
      }
    }
  }

  // 执行全局搜索
  async search(query: string, options: SearchOptions): Promise<SearchResult<T>[]> {
    // [新增] 记录查询参数
    this.lastQuery = query
    this.lastOptions = options

    // 1. 取消上一次
    this.currentCancellationToken.isCancelled = true
    const newToken = { isCancelled: false }
    this.currentCancellationToken = newToken

    // 2. 清理所有模块的 UI 状态
    this.providers.forEach(p => {
      p.clearHighlights()
      p.highlightActive?.(undefined)
    })
    this.currentResults = []
    this.currentFocusIndex = -1

    const context: SearchContext = {
      query,
      options,
      cancellationToken: newToken
    }

    console.log(`[System] Starting search for "${query}"...`)

    // 3. 收集并排序 Provider
    const sortedProviders = Array.from(this.providers.values()).sort(
      (a, b) => b.priority - a.priority
    )

    const allResults: SearchResult<T>[] = []

    // 4. 并行执行搜索
    const promises = sortedProviders.map(async provider => {
      try {
        provider.onSearchStart?.(context)
        if (newToken.isCancelled) return []

        const results = await provider.search(context)
        if (newToken.isCancelled) return []

        return results
      } catch (error) {
        console.error(`Error in provider ${provider.id}:`, error)
        return []
      } finally {
        provider.onSearchEnd?.()
      }
    })

    const resultsArrays = await Promise.all(promises)

    if (newToken.isCancelled) {
      console.log(`[System] Search cancelled.`)
      return []
    }

    // 拍平结果
    resultsArrays.forEach(arr => allResults.push(...arr))
    this.currentResults = allResults

    console.log(`[System] Search finished. Found ${allResults.length} items.`)

    // 5. 分发高亮指令
    const resultsByProvider = new Map<string, SearchResult<T>[]>()
    for (const res of this.currentResults) {
      if (!resultsByProvider.has(res.providerId)) {
        resultsByProvider.set(res.providerId, [])
      }
      resultsByProvider.get(res.providerId)!.push(res)
    }

    resultsByProvider.forEach((results, providerId) => {
      const provider = this.providers.get(providerId)
      if (provider) {
        provider.highlight(results)
      }
    })

    // 6. 自动聚焦第一个 (仅在非刷新模式下，或者刷新后没找到焦点时)
    // 这里简化处理：如果是全新的搜索(currentFocusIndex为-1)，则聚焦第一个
    if (this.currentResults.length > 0 && this.currentFocusIndex === -1) {
      this.focusNext()
    }

    return allResults
  }

  // --- 替换控制 ---

  async replaceCurrent(newText: string) {
    if (this.currentResults.length === 0 || this.currentFocusIndex === -1) return

    const match = this.currentResults[this.currentFocusIndex]
    if (match.invalidated) {
      this.focusNext()
      return
    }

    const provider = this.providers.get(match.providerId)
    if (!provider || !provider.replace) {
      console.warn(`Provider ${match.providerId} does not support replacement.`)
      return
    }

    const success = await provider.replace(match, newText)

    if (success) {
      match.invalidated = true
      console.log(`[Orchestrator] Replaced match ${this.currentFocusIndex + 1}`)
      this.focusNext()
    }
  }

  async replaceAll(newText: string) {
    const resultsByProvider = new Map<string, SearchResult<T>[]>()
    for (const res of this.currentResults) {
      if (!res.invalidated) {
        if (!resultsByProvider.has(res.providerId)) {
          resultsByProvider.set(res.providerId, [])
        }
        resultsByProvider.get(res.providerId)!.push(res)
      }
    }

    for (const [providerId, matches] of resultsByProvider) {
      const provider = this.providers.get(providerId)
      if (provider && provider.replaceAll) {
        await provider.replaceAll(matches, newText)
      }
    }

    this.currentResults = []
    this.currentFocusIndex = -1
    this.providers.forEach(p => {
      p.clearHighlights()
      p.highlightActive?.(undefined)
    })
    console.log('[Orchestrator] Replace All finished.')
  }

  // --- 导航控制 ---

  async focusNext() {
    if (this.currentResults.length === 0) return

    let nextIndex = (this.currentFocusIndex + 1) % this.currentResults.length
    let loopCount = 0

    while (this.currentResults[nextIndex].invalidated && loopCount < this.currentResults.length) {
      nextIndex = (nextIndex + 1) % this.currentResults.length
      loopCount++
    }

    if (loopCount === this.currentResults.length) {
      console.log('[Orchestrator] No valid matches left.')
      return
    }

    this.currentFocusIndex = nextIndex
    await this.performFocus()
  }

  async focusPrev() {
    if (this.currentResults.length === 0) return

    let prevIndex =
      (this.currentFocusIndex - 1 + this.currentResults.length) % this.currentResults.length
    let loopCount = 0

    while (this.currentResults[prevIndex].invalidated && loopCount < this.currentResults.length) {
      prevIndex = (prevIndex - 1 + this.currentResults.length) % this.currentResults.length
      loopCount++
    }

    if (loopCount === this.currentResults.length) {
      console.log('[Orchestrator] No valid matches left.')
      return
    }

    this.currentFocusIndex = prevIndex
    await this.performFocus()
  }

  private async performFocus() {
    const match = this.currentResults[this.currentFocusIndex]
    console.log(
      `[Orchestrator] Focusing match ${this.currentFocusIndex + 1}/${this.currentResults.length}`
    )

    const provider = this.providers.get(match.providerId)
    if (provider) {
      this.providers.forEach(p => p.highlightActive?.(undefined))
      await provider.focus(match)
      provider.highlightActive?.(match)
    } else {
      console.warn(`Provider ${match.providerId} not found for focusing.`)
    }
  }
}

// --- 4. 具体实现示例 ---

interface GridRenderContext {
  type: 'GRID'
  row: number
  col: number
}

interface CommentRenderContext {
  type: 'COMMENT'
  commentId: string
}

type AppRenderContext = GridRenderContext | CommentRenderContext

// 模拟：表格单元格内容查找模块
class CellContentProvider implements SearchProvider<GridRenderContext> {
  readonly id = 'CellData'
  readonly priority = 100

  // 模拟数据存储
  private dataStore = new Map<string, string>([['R1:C1', 'This is a test cell']])

  // [修改] 直接持有请求函数，不再维护 listener 数组
  private searchRequestor?: () => void

  // [修改] 接收上下文
  onRegister(context: ProviderContext) {
    this.searchRequestor = context.requestSearch
  }

  onDispose() {
    this.searchRequestor = undefined
  }

  // 模拟：外部调用此方法更新数据 (例如用户打字)
  public updateData(row: number, col: number, newVal: string) {
    const key = `R${row}:C${col}`
    this.dataStore.set(key, newVal)
    console.log(`[CellData] User typed "${newVal}" at ${key}`)
    // [修改] 直接调用注入的能力
    this.searchRequestor?.()
  }

  async search(ctx: SearchContext): Promise<SearchResult<GridRenderContext>[]> {
    const results: SearchResult<GridRenderContext>[] = []
    // 简单的遍历模拟
    for (const [key, value] of this.dataStore.entries()) {
      if (value.includes(ctx.query)) {
        const [rStr, cStr] = key.split(':')
        const row = parseInt(rStr.substring(1))
        const col = parseInt(cStr.substring(1))

        // 显式构建 GridRenderContext
        const renderData: GridRenderContext = { type: 'GRID', row, col }

        results.push({
          providerId: this.id,
          location: { id: key, description: `Row ${row}, Col ${col}` },
          matchText: ctx.query,
          range: [value.indexOf(ctx.query), value.indexOf(ctx.query) + ctx.query.length],
          score: 1,
          renderData
        })
      }
    }
    return results
  }

  async replace(match: SearchResult<GridRenderContext>, newText: string): Promise<boolean> {
    const { row, col } = match.renderData
    const key = `R${row}:C${col}`
    const currentVal = this.dataStore.get(key) || ''
    // 简单替换
    const newVal = currentVal.replace(match.matchText, newText)
    this.updateData(row, col, newVal) // 这会触发 onDidChange -> refreshSearch
    return true
  }

  highlight(results: SearchResult<GridRenderContext>[]): void {
    console.log(`[CellData UI] Rendering ${results.length} yellow backgrounds.`)
  }

  highlightActive(result?: SearchResult<GridRenderContext>): void {
    if (result) {
      console.log(
        `[CellData UI] Drawing ORANGE border on R${result.renderData.row}:C${result.renderData.col}`
      )
    } else if (!result) {
      console.log(`[CellData UI] Removing active border.`)
    }
  }

  clearHighlights(): void {
    console.log(`[CellData UI] Clearing all highlights.`)
  }

  async focus(result: SearchResult<GridRenderContext>): Promise<void> {
    const data = result.renderData
    if (data.type === 'GRID') {
      console.log(`[CellData UI] Scrolling grid to Row ${data.row}, Col ${data.col}...`)
      await new Promise(r => setTimeout(r, 50))
      console.log(`[CellData UI] Scroll done.`)
    }
  }
}

// 模拟：表格批注/评论查找模块
class CommentProvider implements SearchProvider<CommentRenderContext> {
  readonly id = 'Comments'
  readonly priority = 50

  async search(ctx: SearchContext): Promise<SearchResult<CommentRenderContext>[]> {
    if (ctx.query === 'test') {
      const renderData: CommentRenderContext = { type: 'COMMENT', commentId: 'c-123' }
      return [
        {
          providerId: this.id,
          location: { id: 'R5:C2:Comment', description: 'Comment on R5:C2' },
          matchText: 'Please test this value',
          range: [7, 11],
          score: 0.8,
          renderData
        }
      ]
    }
    return []
  }

  highlight(results: SearchResult<CommentRenderContext>[]): void {
    console.log(`[Comments UI] Highlighting ${results.length} comments.`)
  }

  highlightActive(result?: SearchResult<CommentRenderContext>): void {
    if (result && result.renderData.type === 'COMMENT') {
      console.log(`[Comments UI] Focusing comment box ${result.renderData.commentId}`)
    }
  }

  clearHighlights(): void {
    console.log(`[Comments UI] Clearing highlights.`)
  }

  async focus(result: SearchResult<CommentRenderContext>): Promise<void> {
    const data = result.renderData
    if (data.type === 'COMMENT') {
      console.log(`[Comments UI] Expanding sidebar...`)
      await new Promise(r => setTimeout(r, 30))
      console.log(`[Comments UI] Scrolling to comment ${data.commentId}`)
    }
  }
}

// --- 5. 运行演示 ---

async function runDemo() {
  const orchestrator = new SearchOrchestrator<AppRenderContext>()
  const cellProvider = new CellContentProvider()

  await orchestrator.registerProvider(cellProvider)
  await orchestrator.registerProvider(new CommentProvider())

  console.log('--- Step 1: Search "test" ---')
  await orchestrator.search('test', {
    matchCase: false,
    useRegex: false,
    wholeWord: false
  })

  console.log('\n--- Step 2: User updates data (Change "test" to "demo") ---')
  // 模拟用户输入，这应该导致搜索结果变为空
  cellProvider.updateData(1, 1, 'This is a demo cell')

  // 等待防抖和搜索完成
  await new Promise(r => setTimeout(r, 500))

  console.log('\n--- Step 3: User updates data again (Add "test" back) ---')
  // 模拟用户再次输入，应该重新找到结果
  cellProvider.updateData(1, 1, 'New test value')

  await new Promise(r => setTimeout(r, 500))
}

runDemo()
