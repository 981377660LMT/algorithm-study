// 类型定义
type UnaryHandler<Req, Res> = (req: Req, ctx: any) => Promise<Res>

type UnaryInterceptor<Req, Res> = (req: Req, ctx: any, next: UnaryHandler<Req, Res>) => Promise<Res>

/**
 * 链式组合多个拦截器
 */
function chainInterceptors<Req, Res>(
  interceptors: UnaryInterceptor<Req, Res>[]
): UnaryInterceptor<Req, Res> {
  return async (request, context, next) => {
    // 创建递归调用链
    let index = -1

    const dispatch = async (i: number): Promise<Res> => {
      if (i <= index) {
        throw new Error('next() called multiple times')
      }

      index = i

      // 到达链末尾，调用原始处理函数
      if (i === interceptors.length) {
        return next(request, context)
      }

      // 执行当前拦截器，传入下一个dispatch作为next
      const interceptor = interceptors[i]
      return interceptor(request, context, (req, ctx) => dispatch(i + 1))
    }

    // 从第一个拦截器开始调用
    return dispatch(0)
  }
}

/**
 * 应用拦截器链到原始处理器
 */
function applyInterceptor<Req, Res>(
  handler: UnaryHandler<Req, Res>,
  interceptor: UnaryInterceptor<Req, Res>
): UnaryHandler<Req, Res> {
  return (request, context) => interceptor(request, context, handler)
}

// 类型定义
type StreamHandler<Req, Res> = (reqStream: AsyncIterable<Req>, ctx: any) => AsyncIterable<Res>

type StreamInterceptor<Req, Res> = (
  reqStream: AsyncIterable<Req>,
  ctx: any,
  next: StreamHandler<Req, Res>
) => AsyncIterable<Res>

/**
 * 链式组合多个流拦截器
 */
function chainStreamInterceptors<Req, Res>(
  interceptors: StreamInterceptor<Req, Res>[]
): StreamInterceptor<Req, Res> {
  return async function* (requestStream, context, next) {
    let index = -1

    async function* dispatch(i: number): AsyncIterable<Res> {
      if (i <= index) {
        throw new Error('next() called multiple times')
      }

      index = i

      if (i === interceptors.length) {
        yield* next(requestStream, context)
        return
      }

      const interceptor = interceptors[i]
      yield* interceptor(requestStream, context, (req, ctx) => dispatch(i + 1))
    }

    yield* dispatch(0)
  }
}

// 一元调用示例
async function testUnaryInterceptors() {
  // 原始处理函数
  const handler = async (req: any, ctx: any) => ({ message: `Hello, ${req.name}!` })

  // 示例拦截器
  const interceptors = [
    // 日志拦截器
    async (req, ctx, next) => {
      console.log('Interceptor 1: Before call')
      try {
        const result = await next(req, ctx)
        console.log('Interceptor 1: After call')
        return result
      } catch (error) {
        console.log('Interceptor 1: Error')
        throw error
      }
    },
    // 认证拦截器
    async (req, ctx, next) => {
      console.log('Interceptor 2: Authentication')
      if (!ctx.token) {
        throw new Error('Unauthorized')
      }
      ctx.user = { id: '123' }
      return next(req, ctx)
    }
  ]

  // 链式组合并应用
  const chain = chainInterceptors(interceptors)
  const handlerWithInterceptors = applyInterceptor(handler, chain)

  // 调用
  const result = await handlerWithInterceptors({ name: 'World' }, { token: 'valid-token' })
  console.log('Result:', result)
}

testUnaryInterceptors().catch(console.error)

export {}
