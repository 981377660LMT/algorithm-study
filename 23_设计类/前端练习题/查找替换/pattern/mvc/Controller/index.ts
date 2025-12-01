import { ISearchModule, ISearchResult, ISearchStore } from '../types'
import { SearchTask } from './Task'

interface ISearchControllerDeps {
  store: ISearchStore
}

export class SearchController {
  private readonly _store: ISearchStore /** 存储. */

  private _modules: ISearchModule[] = []
  private _currentTask: SearchTask | undefined = undefined

  constructor(deps: ISearchControllerDeps) {
    this._store = deps.store
  }

  dispose(): void {
    this._currentTask?.cancel()
    this._currentTask = undefined
    this._modules.forEach(m => m.onDispose?.())
    this._modules = []
    this._store.dispose()
  }

  registerModule(module: ISearchModule): void {
    if (this._modules.some(m => m.id === module.id)) {
      console.error(`Module ${module.id} already registered.`)
      return
    }

    this._modules.push(module)
    module.onInit({
      research: () => {
        this._researchModule(module)
      }
    })

    this._currentTask?.searchModule(module)
  }

  unregisterModule(module: ISearchModule): void {
    this._modules = this._modules.filter(m => m.id !== module.id)
    this._currentTask?.cancelModule(module)
    module.onDispose?.()
    this._store.setModuleResults(module.id, [])
  }

  async search(keyword: string): Promise<void> {
    this._currentTask?.cancel()
    this._currentTask = this._createTask(keyword)
    const promises = this._modules.map(m => this._researchModule(m))
    await Promise.allSettled(promises)
  }

  /**
   * 设置当前选中的结果.
   *
   * 1. 不存在该结果，不做任何操作.
   * 2. 更新选中状态.
   * 3. 如果新旧结果来自不同模块，旧模块重新渲染.
   * 4. 新模块重新渲染.
   * 5. 新模块调用 reveal.
   */
  setActiveResult(result: ISearchResult): void {
    const module = this._modules.find(m => m.id === result.moduleId)
    if (!module) return
    const moduleResults = this._store.getModuleResults(module.id)
    if (!moduleResults.some(r => r.key === result.key)) return

    const oldActiveResult = this._store.getActiveResult()
    this._store.setActiveResult(result)

    if (oldActiveResult) {
      const oldModule = this._modules.find(m => m.id === oldActiveResult.moduleId)
      if (oldModule && oldModule.id !== module.id) {
        const oldModuleResults = this._store.getModuleResults(oldModule.id)
        this._renderModule(oldModule, oldModuleResults)
      }
    }

    this._renderModule(module, moduleResults)
    module.reveal?.(result)
  }

  /**
   * 替换当前选中的结果.
   */
  async replace(replaceText: string): Promise<void> {
    const activeResult = this._store.getActiveResult()
    if (!activeResult) return
    const module = this._modules.find(m => m.id === activeResult.moduleId)
    if (!module) return
    try {
      await module.replace?.(activeResult, replaceText)
    } catch (e) {
      console.error(`Module ${module.id} replace failed`, e)
    }
  }

  async replaceAll(replaceText: string): Promise<void> {
    const handleModule = async (module: ISearchModule) => {
      try {
        const results = this._store.getModuleResults(module.id)
        await module.replaceAll?.(results, replaceText)
      } catch (e) {
        console.error(`Module ${module.id} replaceAll failed`, e)
      }
    }
    await Promise.all(this._modules.map(handleModule))
  }

  /** Map 的 key 按照模块权重排序. */
  getSearchResults(): Map<string, ISearchResult[]> {
    const results = this._store.getAllModuleResults()
    const sortedModules = this._modules
      .filter(m => results.has(m.id))
      .sort((a, b) => a.priority() - b.priority())
    const sortedResults = new Map<string, ISearchResult[]>()
    sortedModules.forEach(m => {
      const moduleResults = results.get(m.id) || []
      if (moduleResults.length === 0) return
      sortedResults.set(m.id, moduleResults)
    })
    return sortedResults
  }

  private _createTask(keyword: string): SearchTask {
    const task = new SearchTask({
      keyword,
      didSearchFinished: this._didSearchFinished.bind(this)
    })
    return task
  }

  private _didSearchFinished(module: ISearchModule, results: ISearchResult[]): void {
    // 搜索完成后，校准 activeResult
    const activeResult = this._store.getActiveResult()
    if (activeResult?.moduleId === module.id) {
      const newResult = results.find(r => r.key === activeResult.key)
      if (newResult) {
        this._store.setActiveResult(newResult)
      } else {
        this._store.setActiveResult(undefined)
      }
    }

    this._store.setModuleResults(module.id, results)
    this._renderModule(module, results)
  }

  private async _researchModule(module: ISearchModule): Promise<void> {
    if (!this._currentTask) return
    await this._currentTask.searchModule(module)
  }

  private _renderModule(module: ISearchModule, results: ISearchResult[]): void {
    const activeResult = this._store.getActiveResult()
    const moduleActiveResult = () => {
      if (!activeResult) return undefined
      if (activeResult.moduleId !== module.id) return undefined
      return results.find(r => r.key === activeResult.key)
    }
    module.render?.(results, moduleActiveResult())
  }
}
