import { Func } from '../typings'

type Request = Record<PropertyKey, any>
type NextFunc = (error?: any) => void
type MiddlewareFunc = (req: Request, next: NextFunc) => void
type ErrorHandler = (error: Error, req: Request, next: NextFunc) => void

class Middleware {
  private callbacks: MiddlewareFunc[]
  private errorhandlers: ErrorHandler[]

  constructor() {
    this.callbacks = []
    this.errorhandlers = []
  }

  use(func: MiddlewareFunc | ErrorHandler) {
    if (func.length === 2) this.callbacks.push(func as MiddlewareFunc)
    else if (func.length === 3) this.errorhandlers.push(func as ErrorHandler)
  }

  start(req: Request) {
    const callback = this.compose(this.callbacks, 'callback')
    const errHandler = this.compose(this.errorhandlers, 'handleError')
    try {
      callback(req, () => {})
    } catch (error: any) {
      errHandler(error, req, () => {})
    }
  }

  private compose(func: MiddlewareFunc[], type: 'callback'): MiddlewareFunc
  private compose(func: ErrorHandler[], type: 'handleError'): ErrorHandler
  private compose(funcs: Func[], type: 'callback' | 'handleError'): Func {
    if (funcs.length === 0) return (args: any) => args
    if (funcs.length === 1) return funcs[0]
    if (type === 'callback') {
      return funcs.reduce((a, b) => (req, next) => a(req, () => b(req, next)))
    } else {
      return funcs.reduce((a, b) => (err, req, next) => a(err, req, () => b(err, req, next)))
    }
  }
}

export {}

if (require.main === module) {
  const middleware = new Middleware()

  // middleware.use((req: Request, next: NextFunc) => {
  //   req.a = 1
  //   next()
  // })

  // middleware.use((req: Request, next: NextFunc) => {
  //   req.b = 2
  //   next()
  // })

  // middleware.use((req: Request, next: NextFunc) => {
  //   console.log(req)
  // })

  // middleware.start({})
  // {a: 1, b: 2}
  ///////////////////////////////////////////////////////////////////
  // throw an error at first function
  middleware.use((req: Request, next: NextFunc) => {
    req.a = 1
    throw new Error('sth wrong')
    // or `next(new Error('sth wrong'))`
  })

  // since error occurs, this is skipped
  middleware.use((req: Request, next: NextFunc) => {
    req.b = 2
  })

  // since error occurs, this is skipped
  middleware.use((req: Request, next: NextFunc) => {
    console.log(req)
  })

  // since error occurs, this is called
  middleware.use((error, req, next) => {
    console.log(error)
    console.log(req)
  })

  middleware.start({})
  // Error: sth wrong
  // {a: 1}
}
