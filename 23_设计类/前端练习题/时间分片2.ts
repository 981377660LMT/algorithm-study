export interface ISchedulerOptions {
  /** 每帧的时间预算 (ms)，默认 5ms (React Concurrent Mode 标准) */
  frameBudget?: number
  /** 取消信号 */
  signal?: AbortSignal
}

/**
 * 核心调度器：用于将长任务切片，避免阻塞主线程
 */
export class TimeSlicer {
  /**
   * 1. Generator 模式：适合复杂的、有状态的流程控制（如状态机、递归算法）
   * @returns Generator 的最终返回值
   */
  static async runGenerator<T, R = any>(
    gen: Generator<T, R, unknown>,
    options: ISchedulerOptions = {}
  ): Promise<R> {
    const { frameBudget = 5, signal } = options
    let lastYieldTime = performance.now()

    let res = gen.next()

    while (!res.done) {
      // 检查取消
      if (signal?.aborted) throw new DOMException('Aborted', 'AbortError')

      // 检查时间预算
      const now = performance.now()
      if (now - lastYieldTime > frameBudget) {
        // 超时，让出主线程
        await this.yieldToMain()
        lastYieldTime = performance.now() // 重置计时
      }

      // 继续执行
      res = gen.next()
    }

    return res.value as R
  }

  /**
   * 2. 数组遍历模式 (forEach)：适合表格行渲染、大数据处理
   * 替代 Array.prototype.forEach
   */
  static async forEach<T>(
    items: T[],
    iterator: (item: T, index: number) => void,
    options: ISchedulerOptions = {}
  ): Promise<void> {
    const { frameBudget = 5, signal } = options
    const len = items.length
    let lastYieldTime = performance.now()

    for (let i = 0; i < len; i++) {
      // 批量执行检查：每处理一个元素检查一次可能太频繁，
      // 可以在这里加一个计数器优化 (例如 i % 10 === 0 才检查时间)，但为了通用性先每次检查

      const now = performance.now()
      if (now - lastYieldTime > frameBudget) {
        if (signal?.aborted) throw new DOMException('Aborted', 'AbortError')
        await this.yieldToMain()
        lastYieldTime = performance.now()
      }

      iterator(items[i], i)
    }
  }

  /**
   * 3. 数组映射模式 (map)：适合编辑器文本解析、数据转换
   * 替代 Array.prototype.map
   */
  static async map<T, R>(
    items: T[],
    mapper: (item: T, index: number) => R,
    options: ISchedulerOptions = {}
  ): Promise<R[]> {
    const results: R[] = new Array(items.length)

    await this.forEach(
      items,
      (item, index) => {
        results[index] = mapper(item, index)
      },
      options
    )

    return results
  }

  /**
   * 4. 手动打点模式：适合无法拆分为数组的 while 循环
   * 返回一个检查函数，在循环内部调用
   */
  static createChecker(options: ISchedulerOptions = {}) {
    const { frameBudget = 5, signal } = options
    let lastYieldTime = performance.now()

    return async () => {
      if (signal?.aborted) throw new DOMException('Aborted', 'AbortError')

      const now = performance.now()
      if (now - lastYieldTime > frameBudget) {
        await this.yieldToMain()
        lastYieldTime = performance.now()
      }
    }
  }

  // 私有：让出主线程
  private static yieldToMain() {
    // 优先使用 MessageChannel (宏任务)，比 setTimeout(0) 更快，延迟更低
    if (typeof MessageChannel !== 'undefined') {
      return new Promise<void>(resolve => {
        const channel = new MessageChannel()
        channel.port1.onmessage = () => resolve()
        channel.port2.postMessage(null)
      })
    }
    // 降级方案
    return new Promise<void>(resolve => setTimeout(resolve, 0))
  }
}

{
  // 假设有 10万条数据
  const largeTableData = Array.from({ length: 100000 }, (_, i) => ({ id: i, val: Math.random() }))

  async function renderTable() {
    const container = document.getElementById('table-body')
    const controller = new AbortController()

    try {
      // 使用 TimeSlicer 替代 data.forEach
      await TimeSlicer.forEach(
        largeTableData,
        row => {
          const tr = document.createElement('tr')
          tr.innerHTML = `<td>${row.id}</td><td>${row.val}</td>`
          container?.appendChild(tr)
        },
        {
          frameBudget: 8, // 稍微放宽一点预算
          signal: controller.signal
        }
      )

      console.log('表格渲染完成，全程无卡顿')
    } catch (e) {
      console.log('渲染被取消')
    }
  }
}

{
  declare function expensiveTokenizer(line: string): string[]
  declare function renderTokens(tokens: string[][]): void
  const sourceCode = `...` // 假设是一个很大的源代码字符串
  const codeLines = sourceCode.split('\n') // 假设有 5000 行代码

  async function highlightCode() {
    // 使用 TimeSlicer 替代 lines.map
    // 这样即使文件很大，编辑器也能保持响应（打字不卡）
    const tokens = await TimeSlicer.map(
      codeLines,
      line => {
        return expensiveTokenizer(line) // 耗时的正则匹配
      },
      { frameBudget: 5 }
    )

    renderTokens(tokens)
  }
}

{
  async function parseStream(stream: string) {
    const checkTime = TimeSlicer.createChecker({ frameBudget: 5 })
    let i = 0

    while (i < stream.length) {
      // 在循环内部手动打点
      // 如果当前帧耗时过长，这里会自动 await 暂停
      await checkTime()

      processChar(stream[i])
      i++
    }
  }
}
