// 维护两个指针即可

class BrowserHistory {
  /**
   * @param {string} url
   * if url is set, it means new tab with url
   * otherwise, it is empty new tab
   */
  constructor(url) {
    this.history = url ? [url] : []
    this.cur = 0
    this.last = 0
  }

  /**
   * @param { string } url
   */
  visit(url) {
    this.cur++
    this.history[this.cur] = url
    this.last = this.cur
  }

  /**
   * @return {string} current url
   */
  get current() {
    return this.history[this.cur]
  }

  // go to previous entry
  goBack() {
    this.cur = Math.max(0, --this.cur)
    return this.current
  }

  // go to next visited url
  forward() {
    this.cur = Math.min(this.last, ++this.cur)
    return this.current
  }
}
