import {
  CancellationToken,
  ISearchPlugin,
  ISearchPluginContext,
  ISearchResult,
  SearchController
} from './model'

type EditorRange = { start: number; end: number }

class EditorPlugin implements ISearchPlugin<EditorRange> {
  readonly id = 'main-editor'
  private context: ISearchPluginContext
  private reSearchTimer: any = null

  constructor(private controller: SearchController, private editorEngine: any) {
    this.context = controller.registerPlugin(this)

    this.editorEngine.on('content-change', () => {
      // 简单的防抖，避免长按按键时触发过多的 Task 创建
      if (this.reSearchTimer) clearTimeout(this.reSearchTimer)
      this.reSearchTimer = setTimeout(() => {
        this.context.reSearch()
      }, 50) // 50ms 的极短防抖即可大幅减轻 Controller 压力
    })
  }

  priority() {
    return 100
  }

  async search(keyword: string, token: CancellationToken): Promise<ISearchResult<EditorRange>[]> {
    // 这里可以分片
    const matches = await this.editorEngine.findAll(keyword)
    token.throwIfCancellationRequested()

    return matches.map((m: any, idx: number) => ({
      id: `editor-${idx}`,
      pluginId: this.id,
      matchContent: m.text,
      canReplace: true,
      payload: { start: m.start, end: m.end }
    }))
  }

  render(results: ISearchResult<EditorRange>[], activeResult: ISearchResult<EditorRange> | null) {
    const decorations = results.map(r => ({
      range: r.payload,
      color: r.id === activeResult?.id ? 'orange' : 'yellow'
    }))
    this.editorEngine.setDecorations(decorations)
  }

  scrollTo(result: ISearchResult<EditorRange>) {
    this.editorEngine.scrollTo(result.payload.start)
  }

  async replace(result: ISearchResult<EditorRange>, newText: string) {
    this.editorEngine.replace(result.payload, newText)
  }

  async replaceAll(results: ISearchResult<EditorRange>[], newText: string) {
    const ranges = results.map(r => r.payload)
    this.editorEngine.bulkReplace(ranges, newText)
  }

  clear() {
    this.editorEngine.setDecorations([])
  }
}
