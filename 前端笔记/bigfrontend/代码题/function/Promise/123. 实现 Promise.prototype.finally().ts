// Promise.prototype.finally() 可以在promise被settle的时候执行一个callback，
// 无论其是被fulfill或是被reject。
// 传入的callback不接受参数，也即该callback不影响promise链中的value（注意rejection)。
/**
 * @param {Promise<any>} promise
 * @param {() => void} onFinally
 * @returns {Promise<any>}
 * 异步需要套在Promise.resolve里
 */
function myFinally(promise: Promise<any>, onFinally: () => void): Promise<void> {
  // your code here
  return promise
    .then(data => Promise.resolve(onFinally()).then(() => data))
    .catch(e => Promise.resolve(onFinally()).then(() => Promise.reject(e)))
}
Promise.resolve().finally()
/**
 * @param {Promise<any>} promise
 * @param {() => void} onFinally
 * @returns {Promise<any>}
 */
async function myFinally2(promise: Promise<any>, onFinally: () => void): Promise<void> {
  try {
    const res = await promise
    await onFinally()
  } catch (error) {
    await onFinally()
    throw error
  }
}

// 如果promise的finally方法里面出错了怎么办
// 会返回一个 被reject 的promise对象。继续绑定catch函数可以捕获到
