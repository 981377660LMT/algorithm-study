import { ISearchModule, ISearchResult } from '../../types'
import { CancellationTokenSource } from '../../utils'

interface ISearchTaskDeps {
  keyword: string
  didSearchFinished: (module: ISearchModule, results: ISearchResult[]) => void
}

export class SearchTask {
  readonly keyword: string
  private readonly _didSearchFinished: (module: ISearchModule, results: ISearchResult[]) => void

  private readonly _mainCancelSource = new CancellationTokenSource() // 全局任务取消源 (用户输入新词时触发)
  private readonly _moduleCancelSources = new Map<string, CancellationTokenSource>() // 插件级取消源 (reSearch 时触发，互斥同一插件的多次搜索)

  constructor(deps: ISearchTaskDeps) {
    this.keyword = deps.keyword
    this._didSearchFinished = deps.didSearchFinished
  }

  async searchModule(module: ISearchModule) {
    // 1. 如果全局任务已取消，直接终止
    if (this._mainCancelSource.token.isCancellationRequested) return

    // 2. 取消该插件上一次正在进行的搜索 (reSearch 互斥)
    const prevSource = this._moduleCancelSources.get(module.id)
    if (prevSource) {
      prevSource.cancel()
    }

    // 3. 创建本次搜索的专属 Token
    const localSource = new CancellationTokenSource()
    this._moduleCancelSources.set(module.id, localSource)

    try {
      // 4. 执行搜索
      // 注意：插件内部应该定期检查 token.isCancellationRequested
      const results =
        this.keyword === '' ? [] : await module.search(this.keyword, { token: localSource.token })

      // 5. 提交前检查：全局是否取消？局部是否被新 reSearch 取消？
      if (this._mainCancelSource.token.isCancellationRequested) return
      if (localSource.token.isCancellationRequested) return

      // 6. 通知更新结果
      this._didSearchFinished(module, results)
    } catch (e: any) {
      console.error(`Module ${module.id} search error:`, e)
    } finally {
      // 清理：只有当 map 中的 source 还是当前 source 时才删除
      // (防止误删了新发起的 reSearch source)
      if (this._moduleCancelSources.get(module.id) === localSource) {
        this._moduleCancelSources.delete(module.id)
      }
    }
  }

  cancelModule(module: ISearchModule) {
    if (this.isCancelled) return
    const source = this._moduleCancelSources.get(module.id)
    if (source) {
      source.cancel()
      this._moduleCancelSources.delete(module.id)
    }
  }

  cancel() {
    if (this.isCancelled) return
    this._mainCancelSource.cancel()
    this._moduleCancelSources.forEach(s => s.cancel())
    this._moduleCancelSources.clear()
  }

  get isCancelled() {
    return this._mainCancelSource.token.isCancellationRequested
  }
}
