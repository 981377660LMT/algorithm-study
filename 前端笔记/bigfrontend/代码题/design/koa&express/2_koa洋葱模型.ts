// 本质上是dfs前序 后后序回溯

type Function = (...args: any[]) => any
type Middleware = (context: any, next: Function) => any

class Koa {
  private middlewares: Middleware[] = []

  use(fn: Middleware) {
    if (typeof fn !== 'function') throw new Error('入参必须是函数')
    this.middlewares.push(fn)
  }

  listen(...args: any[]) {
    const fn = this.compose(this.middlewares)('context')
    return fn.then(() => console.log('over')).catch(err => console.log(err))
  }

  /**
   *
   * @param middlwwares 核心是递归  middleware传参next表示下一个中间件
   * @returns
   */
  private compose(middlwwares: Middleware[]): (context: any) => Promise<void> {
    return (context: any) => {
      return dfs(0)

      async function dfs(index: number): Promise<void> {
        if (index == middlwwares.length) return
        const middleware = middlwwares[index]
        try {
          // next 相当于 ()=>dfs(index+1)
          await middleware(context, () => dfs(index + 1))
        } catch (error) {
          throw error
        }
      }
    }
  }
}

// next的含义为下一个中间件
if (require.main === module) {
  const app = new Koa()
  app.use(async (ctx, next) => {
    console.log(1)
    await next()
    console.log(2)
  })
  app.use(async (ctx, next) => {
    console.log(3)
    await next()
    console.log(4)
  })
  app.use(async (ctx, next) => {
    console.log(5)
    await next()
    console.log(6)
  })

  app.listen(3000)
  // 1 // 3 // 5 // 6 // 4 // 2
}

export {}
