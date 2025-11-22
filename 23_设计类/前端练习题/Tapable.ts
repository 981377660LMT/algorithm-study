// Tapable 是 Webpack 插件系统的灵魂，它本质上是一个发布订阅模式的加强版，支持同步、异步、串行、并行、熔断等多种复杂的事件流控制。

// 要手写一套 Tapable 抽象，我们需要实现以下核心钩子（Hooks）：

// SyncHook: 同步串行执行（最基础）。
// SyncBailHook: 同步熔断执行（只要有一个插件返回非 undefined，就停止后续执行）。
// AsyncSeriesHook: 异步串行执行（前一个 Promise 完成，下一个才开始）。
// AsyncParallelHook: 异步并行执行（所有插件同时跑，都完成了才算结束）。

// !控制反转 (IoC)： Compiler 不知道具体要干什么，它只定义了“流程节点”（Hooks）。具体的逻辑全部由插件注册进来。
// 动态编译 (Dynamic Compilation)： 注意：真实的 Tapable 库比上面的实现要复杂得多。 真实的 Tapable 为了极致性能，使用了 new Function 动态生成代码。 当你调用 hook.call() 时，它不是遍历数组，而是直接执行一段编译好的字符串代码：
// 流控制的多样性： 通过 Series（串行）、Parallel（并行）、Bail（熔断）、Waterfall（数据传递）四种组合，几乎可以描述世界上所有的软件执行流程

export {}

// 插件接口
interface Tap {
  name: string
  type: 'sync' | 'async' | 'promise'
  fn: Function
}

class Hook {
  protected taps: Tap[] = []
  protected args: string[]

  constructor(args: string[] = []) {
    this.args = args // 记录参数名，主要用于调试或元编程，这里仅作占位
  }

  // 注册同步插件
  tap(name: string, fn: Function) {
    this.taps.push({ name, type: 'sync', fn })
  }

  // 注册异步 Promise 插件
  tapPromise(name: string, fn: (...args: any[]) => Promise<any>) {
    this.taps.push({ name, type: 'promise', fn })
  }

  // 清空插件
  clear() {
    this.taps = []
  }
}

class SyncHook<T extends any[]> extends Hook {
  call(...args: T): void {
    for (const tap of this.taps) {
      try {
        tap.fn(...args)
      } catch (e) {
        console.error(`[SyncHook] Error in plugin '${tap.name}':`, e)
        throw e // 同步钩子通常直接抛出错误
      }
    }
  }
}

class SyncBailHook<T extends any[], R> extends Hook {
  call(...args: T): R | undefined {
    for (const tap of this.taps) {
      const result = tap.fn(...args)
      // 熔断逻辑：只要有返回值，立即停止并返回
      if (result !== undefined) {
        return result
      }
    }
    return undefined
  }
}

class AsyncSeriesHook<T extends any[]> extends Hook {
  // 这里的 callAsync 模拟 tapable 的 promise 方法
  async promise(...args: T): Promise<void> {
    for (const tap of this.taps) {
      // 无论注册的是 tap 还是 tapPromise，都统一按 Promise 处理
      const result = tap.fn(...args)

      if (result instanceof Promise) {
        await result
      } else {
        // 如果是普通函数，直接执行完了，await undefined 会立即通过
        await Promise.resolve(result)
      }
    }
  }
}

class AsyncParallelHook<T extends any[]> extends Hook {
  async promise(...args: T): Promise<void> {
    const promises = this.taps.map(tap => {
      const result = tap.fn(...args)
      if (result instanceof Promise) {
        return result
      }
      return Promise.resolve(result)
    })

    await Promise.all(promises)
  }
}

class SyncWaterfallHook<T extends any[]> extends Hook {
  call(...args: T): any {
    if (this.taps.length === 0) return args[0]

    let [result, ...rest] = args

    for (const tap of this.taps) {
      // 上一次的结果 result 传给当前插件
      const nextResult = tap.fn(result, ...rest)

      // 如果插件返回了新值，更新 result；否则保持原样（视具体约定而定，通常 waterfall 必须返回值）
      if (nextResult !== undefined) {
        result = nextResult
      }
    }
    return result
  }
}

{
  class Compiler {
    hooks = {
      // 1. 生命周期开始（同步）
      initialize: new SyncHook<[]>(),

      // 2. 决定如何处理文件（熔断：只要有一个插件能处理，就不问别人了）
      shouldEmit: new SyncBailHook<[string], boolean>(),

      // 3. 编译过程（异步串行：先下载资源，再解析，顺序不能乱）
      compile: new AsyncSeriesHook<[string]>(),

      // 4. 压缩/上传（异步并行：图片压缩和 CSS 压缩可以同时做）
      optimize: new AsyncParallelHook<[]>(),

      // 5. 处理配置（瀑布流：插件可以修改配置对象）
      processConfig: new SyncWaterfallHook<[any]>()
    }

    async run() {
      console.log('--- Build Started ---')

      // 1. 初始化
      this.hooks.initialize.call()

      // 5. 处理配置 (演示 Waterfall)
      const initialConfig = { mode: 'development' }
      const finalConfig = this.hooks.processConfig.call(initialConfig)
      console.log('Final Config:', finalConfig)

      // 2. 检查是否需要发射
      const shouldEmit = this.hooks.shouldEmit.call('main.js')
      if (shouldEmit === false) {
        console.log('Skipping emit based on plugin decision.')
        return
      }

      // 3. 编译
      console.log('Compiling...')
      await this.hooks.compile.promise('main.js')

      // 4. 优化
      console.log('Optimizing...')
      await this.hooks.optimize.promise()

      console.log('--- Build Finished ---')
    }
  }

  // --- 插件编写 ---

  const compiler = new Compiler()

  // 插件 A: 初始化日志
  compiler.hooks.initialize.tap('LoggerPlugin', () => {
    console.log('[Plugin] System initialized.')
  })

  // 插件 B: 修改配置 (Waterfall)
  compiler.hooks.processConfig.tap('ProductionPlugin', config => {
    console.log('[Plugin] Switching to production mode')
    return { ...config, mode: 'production', minify: true }
  })

  // 插件 C: 决定是否发射 (Bail)
  compiler.hooks.shouldEmit.tap('FilterPlugin', filename => {
    if (filename.endsWith('.tmp')) return false // 拦截
    return undefined // 不拦截，交给下一个
  })

  // 插件 D: 异步编译任务 1
  compiler.hooks.compile.tapPromise('DownloadPlugin', async file => {
    console.log(`[Plugin] Downloading ${file}...`)
    await new Promise(r => setTimeout(r, 500))
  })

  // 插件 E: 异步编译任务 2
  compiler.hooks.compile.tapPromise('ParsePlugin', async file => {
    console.log(`[Plugin] Parsing ${file}...`)
    await new Promise(r => setTimeout(r, 300))
  })

  // 插件 F: 并行优化 1
  compiler.hooks.optimize.tapPromise('ImageMin', async () => {
    console.log('[Plugin] Minifying images...')
    await new Promise(r => setTimeout(r, 800)) // 慢
    console.log('[Plugin] Images done.')
  })

  // 插件 G: 并行优化 2
  compiler.hooks.optimize.tapPromise('CssMin', async () => {
    console.log('[Plugin] Minifying CSS...')
    await new Promise(r => setTimeout(r, 400)) // 快
    console.log('[Plugin] CSS done.')
  })

  // --- 运行 ---
  compiler.run()
}
