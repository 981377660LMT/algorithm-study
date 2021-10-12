// Promise.prototype.finally() 可以在promise被settle的时候执行一个callback，
// 无论其是被fulfill或是被reject。
// 传入的callback不接受参数，也即该callback不影响promise链中的value（注意rejection)。
/**
 * @param {Promise<any>} promise
 * @param {() => void} onFinally
 * @returns {Promise<any>}
 * 异步需要套在Promise.resolve里
 */
function myFinally(promise: Promise<any>, onFinally: () => void): Promise<any> {
  // your code here
  return promise
    .then(data => Promise.resolve(onFinally()).then(() => data))
    .catch(e => Promise.resolve(onFinally()).then(() => Promise.reject(e)))
}

/**
 * @param {Promise<any>} promise
 * @param {() => void} onFinally
 * @returns {Promise<any>}
 */
async function myFinally2(promise: Promise<any>, onFinally: () => void): Promise<any> {
  try {
    const res = await promise
    await onFinally()
    return res
  } catch (error) {
    await onFinally()
    throw error
  }
}
