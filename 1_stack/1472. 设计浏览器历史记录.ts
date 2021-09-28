class BrowserHistory {
  private first: number
  private last: number
  private history: string[]

  constructor(homepage: string) {
    this.first = 0
    this.last = 0
    this.history = [homepage]
  }

  // 从当前页跳转访问 url 对应的页面  。执行此操作会把浏览历史前进的记录全部删除(即将last定位至下一个元素)
  visit(url: string): void {
    this.first++
    if (this.first === this.history.length) {
      this.history.push('') // 多开一格
    }
    this.history[this.first] = url
    this.last = this.first
  }

  // 请返回后退 至多 steps 步以后的 url
  back(steps: number): string {
    this.first = Math.max(0, this.first - steps)
    return this.history[this.first]
  }

  // 请返回前进 至多 steps步以后的 url
  forward(steps: number): string {
    this.first = Math.min(this.last, this.first + steps)
    return this.history[this.first]
  }
}

export {}
