class BrowserHistory {
  private cur: number
  private last: number
  private history: string[]

  constructor(homepage: string) {
    this.cur = 0
    this.last = 0
    this.history = [homepage]
  }

  // 从当前页跳转访问 url 对应的页面  。执行此操作会把浏览历史前进的记录全部删除(即将last定位至下一个元素)
  visit(url: string): void {
    this.cur++
    if (this.cur === this.history.length) {
      this.history.push('') // 多开一格
    }
    this.history[this.cur] = url
    this.last = this.cur
  }

  // 请返回后退 至多 steps 步以后的 url
  back(steps: number): string {
    this.cur = Math.max(0, this.cur - steps)
    return this.history[this.cur]
  }

  // 请返回前进 至多 steps步以后的 url
  forward(steps: number): string {
    this.cur = Math.min(this.last, this.cur + steps)
    return this.history[this.cur]
  }
}

export {}
