let handleError: (err: Error) => void | undefined

const utils = {
  foo(fn: () => void) {
    callWithErrorHandling(fn)
  },
  // 用户可以调用该函数注册统一的错误处理函数
  registerErrorHandle(fn: (err: Error) => void) {
    handleError = fn
  }
}

function callWithErrorHandling(fn: () => void): void {
  try {
    fn && fn()
  } catch (e) {
    // 将捕获到的错误传递给用户的错误处理程序
    handleError && handleError(e as Error)
  }
}

export {}

if (require.main === module) {
  // 这时错误处理的能力完全由用户控制，用户既可以选择忽略错误，也可以调用上报程序将错误上报给监控系统
  utils.registerErrorHandle(err => {
    console.error('统一错误处理：', err.message)
  })

  utils.foo(() => {
    throw new Error('foo error')
  })
}
