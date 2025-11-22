type Context = any
type Next = () => Promise<void>
type Middleware = (ctx: Context, next: Next) => Promise<void>

export class Pipeline {
  private middlewares: Middleware[] = []

  use(fn: Middleware) {
    this.middlewares.push(fn)
    return this
  }

  async execute(ctx: Context) {
    const dispatch = async (index: number) => {
      if (index === this.middlewares.length) return

      const fn = this.middlewares[index]
      // 递归调用：传入 ctx 和下一个中间件的执行函数
      await fn(ctx, () => dispatch(index + 1))
    }

    await dispatch(0)
  }
}

// --- 业务实战：订单提交检查 ---
{
  const orderPipeline = new Pipeline()

  // 1. 检查库存
  orderPipeline.use(async (ctx, next) => {
    console.log('Checking stock...')
    if (ctx.count > 100) throw new Error('Out of stock')
    await next()
  })

  // 2. 检查优惠券
  orderPipeline.use(async (ctx, next) => {
    console.log('Validating coupon...')
    if (ctx.coupon === 'EXPIRED') ctx.discount = 0
    await next()
  })

  // 3. 上报日志 (洋葱模型：可以在结束后执行)
  orderPipeline.use(async (ctx, next) => {
    const start = Date.now()
    await next() // 等待后续逻辑执行完
    console.log(`Order processed in ${Date.now() - start}ms`)
  })

  // 执行
  orderPipeline.execute({ count: 5, coupon: 'VIP' }).catch(console.error)
}
