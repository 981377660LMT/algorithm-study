type Function<T = any> = (...args: any[]) => T
type Request = Record<any, any>
type Response = any
type NextFunc = (error?: any) => void
type MiddlewareFunc = (req: Request, res: Response, next: NextFunc) => void
type ErrorHandler = (error: Error, req: Request, res: Response, next: NextFunc) => void

class Middleware {
  private callbacks: MiddlewareFunc[]
  private errorhandlers: ErrorHandler[]

  constructor() {
    this.callbacks = []
    this.errorhandlers = []
  }

  use(func: MiddlewareFunc | ErrorHandler) {
    if (func.length === 3) this.callbacks.push(func as MiddlewareFunc)
    else if (func.length === 4) this.errorhandlers.push(func as ErrorHandler)
  }

  listen(req = {}, res = {}) {
    const callback = this.compose(this.callbacks, 'callback')
    const errHandler = this.compose(this.errorhandlers, 'handleError')
    try {
      callback(req, res, () => {})
    } catch (error: any) {
      errHandler(error, req, res, () => {})
    }
  }

  private compose(func: MiddlewareFunc[], type: 'callback'): MiddlewareFunc
  private compose(func: ErrorHandler[], type: 'handleError'): ErrorHandler
  private compose(funcs: Function[], type: 'callback' | 'handleError'): Function {
    if (funcs.length === 0) return (args: any) => args
    if (funcs.length === 1) return funcs[0]
    if (type === 'callback') {
      return funcs.reduce((a, b) => (req, res, next) => a(req, res, () => b(req, res, next)))
    } else {
      return funcs.reduce(
        (a, b) => (err, req, res, next) => a(err, req, res, () => b(err, req, res, next))
      )
    }
  }
}

export {}

if (require.main === module) {
  const app = new Middleware()

  // middleware.use((req: any, next: NextFunc) => {
  //   req.a = 1
  //   next()
  // })

  // middleware.use((req: any, next: NextFunc) => {
  //   req.b = 2
  //   next()
  // })

  // middleware.use((req: any, next: NextFunc) => {
  //   console.log(req)
  // })

  // middleware.start({})
  // {a: 1, b: 2}
  ///////////////////////////////////////////////////////////////////
  // throw an error at first function
  app.use((req: any, res: any, next: NextFunc) => {
    req.a = 1
    next()
    // throw new Error('sth wrong')
    // or `next(new Error('sth wrong'))`
  })

  // since error occurs, this is skipped
  app.use((req: any, res: any, next: NextFunc) => {
    req.b = 2
    next()
  })

  // since error occurs, this is skipped
  app.use((req: any, res: any, next: NextFunc) => {
    console.log(req)
    next()
  })

  // since error occurs, this is called
  app.use((error, req, res, next) => {
    console.log(error)
    console.log(req)
    next()
  })

  app.listen()
  // Error: sth wrong
  // {a: 1}
}

export {}
