// localStorage是简单方便的client-side storage，但是因为其是synchronous，
// 所以localStorage尽量不要使用。
// 同时，Safari的 ITP 的影响下，如果safari7天以上被使用但是你的网页并没有被交互，
// 那么client-side script-writable storage 将会被删除，localStorage当然也是对象。
// 和Cookie不同，localStorage并没有过期的概念。
// 本题中，请实现使localStorage支持过期的wrapper。

// localStorage 与sessionStorage 的API 都是同步操作

/**
 * 思想是get时删除哪些过期的数据 不应该使用setTimeout(浏览器关闭时就没了，不能完全清除)
 */
class MyLocalStorage implements Omit<Storage, 'setItem'> {
  getItem(key: string) {
    const result = JSON.parse(window.localStorage.getItem(key) || '')

    if (result) {
      // 延迟删除
      if (result.maxAge <= Date.now()) {
        window.localStorage.removeItem(key)
        return null
      }

      return result.data
    }

    return null
  }

  /**
   *
   * @param key
   * @param value
   * @param maxAge 存活时间（max-age）
   * 负数：临时Cookie,有效期session(我们这里不管)；
   * 0：删除cookie；
   * 正数：有效期为创建时刻+ max-age
   */
  setItem(key: string, value: string, maxAge?: number) {
    const result = {
      data: value,
      maxAge: Infinity,
    }

    if (maxAge === 0) return
    if (maxAge) result.maxAge = Date.now() + maxAge
    window.localStorage.setItem(key, JSON.stringify(result))
  }

  removeItem(key: string) {
    window.localStorage.removeItem(key)
  }

  clear() {
    window.localStorage.clear()
  }
}

export {}

if (require.main === module) {
  const myLocalStorage = new MyLocalStorage()
  myLocalStorage.setItem('bfe', 'dev', 1000)
  console.log(myLocalStorage.getItem('bfe')) // 'dev'
  console.log(myLocalStorage.getItem('bfe')) // null
}
