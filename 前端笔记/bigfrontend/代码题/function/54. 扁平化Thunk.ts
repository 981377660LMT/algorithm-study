type Thunk = (callback: Callback) => void
type Callback = (error: Error | null, result: any | Thunk) => void

/**
 * @param {Thunk} thunk
 * @return {Thunk}
 * 当Error发生时，尚未被调用的函数需要被跳过。
 * @description
 * 把加强的callback传进thunk，连续执行
 */
function flattenThunk(thunk: Thunk): Thunk {
  return callback => {
    const inhancedCallback: Callback = (err, res) => {
      if (err) callback(err, undefined)
      else if (typeof res === 'function') res(inhancedCallback)
      else callback(err, res) // 遇到了最里层的thunk
    }

    thunk(inhancedCallback)
  }
}

if (require.main === module) {
  const func1: Thunk = (cb: Callback) => {
    setTimeout(() => cb(null, 'ok'), 10)
  }

  const func2: Thunk = (cb: Callback) => {
    setTimeout(() => cb(null, func1), 10)
  }

  const func3: Thunk = (cb: Callback) => {
    setTimeout(() => cb(null, func2), 10)
  }

  flattenThunk(func3)((error, data) => {
    if (error) console.error(error)
    else console.log(data) // 'ok'
  })
}

export { flattenThunk }
