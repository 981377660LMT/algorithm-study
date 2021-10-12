type SomeConfig = Record<string, any>
// getAPI is bundled with your code, config will only be some plain objects.
declare const getAPI: <T>(path: string, config: SomeConfig) => Promise<T>
// const list1 = await getAPI('/list', { keyword: 'bfe'})
// const list2 = await getAPI('/list', { keyword: 'dev'})

// 这本来运行地很好。直到页面UI变得如此复杂，以至于相同的API在极短的时间内被调用了多次。
// 你想要基于如下的假设避免没必要的API请求。
// 在1000ms以内的GET API的response几乎肯定是一样的
////////////////////////////////////////////////////////////////////////////////////////////////
const cache = new Map<string, Promise<any>>()

/**
 * @param {string} path
 * @param {Record<string,any>} config  假定config只有一层 无嵌套的object
 * @returns {Promise<any>}
 */
function getAPIWithMerging(path: string, config: Record<string, any>): Promise<any> {
  const cacheKey = getHashKey(path, config)
  if (cache.has(cacheKey)) return cache.get(cacheKey)!
  return getAPIAndStoreInCache(path, config)
}

function getHashKey(path: string, config: Record<string, any>): string {
  const keys = Object.keys(config).sort()
  return `${path}:${keys.map(key => `${key}-${config[key]}`).join('&')}`
}

function getAPIAndStoreInCache(path: string, config: Record<string, any>): Promise<any> {
  if (cache.size >= 5) cache.delete(cache.keys().next().value) // 限制容量

  const key = getHashKey(path, config)
  const promise = getAPI(path, config)
  cache.set(key, promise)

  // 过期时间
  setTimeout(() => {
    cache.delete(key)
  }, 1000)

  return promise
}

getAPIWithMerging.clearCache = () => {
  cache.clear()
}

if (require.main === module) {
  getAPIWithMerging('/list', { keyword: 'bfe' }).then()
  // 第1次调用。getAPI被触发。

  getAPIWithMerging('/list', { keyword: 'bfe' }).then()
  // 第2次调用和第1次是相同的
  // 所以 getAPI 被没有被触发
  // 被合并到了第1次调用中

  getAPIWithMerging('/list', { keyword: 'dev' }).then()
  // 第3次调用是不同的参数，所以getAPI被触发

  // 1000ms 过后
  getAPIWithMerging('/list', { keyword: 'bfe' }).then()
  // 第4次调用和第1次相同
  // 不过已经过了1000ms，getAPI还是被触发
  // 请注意内存泄漏！
  // 你的缓存机制需要设置上限。该问题中，请设置为最多5条缓存记录。意味着：

  getAPIWithMerging('/list1', { keyword: 'bfe' }).then()
  // 第1次调用，触发callAPI()，添加1条缓存记录
  getAPIWithMerging('/list2', { keyword: 'bfe' }).then()
  // 第2次调用，触发callAPI()，添加1条缓存记录
  getAPIWithMerging('/list3', { keyword: 'bfe' }).then()
  // 第3次调用，触发callAPI()，添加1条缓存记录
  getAPIWithMerging('/list4', { keyword: 'bfe' }).then()
  // 第4次调用，触发callAPI()，添加1条缓存记录
  getAPIWithMerging('/list5', { keyword: 'bfe' }).then()
  // 第5次调用，触发callAPI()，添加1条缓存记录

  getAPIWithMerging('/list6', { keyword: 'bfe' }).then()
  // 第6次调用，触发callAPI()，添加1条缓存记录
  // 第1次调用的缓存记录被删除

  getAPIWithMerging('/list1', { keyword: 'bfe' }).then()
  // 和第1次调用相同，但是没有缓存的数据
  // 所以新的缓存记录被添加

  // clear()
  // 为了测试，请提供上述方法来清楚所有缓存。调用方式如下
  getAPIWithMerging.clearCache()

  const p1 = new Promise(resolve => {
    console.log('共用的promise只请求一次')
    resolve(666)
  })

  p1.then(console.log)
  p1.then(console.log)
  p1.then(console.log)
  p1.then(console.log)
}
