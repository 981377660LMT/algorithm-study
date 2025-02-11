import { Mutex } from './Mutex'

export class Once {
  private _done = false
  private readonly _mutex = new Mutex()

  async do(fn: () => void | Promise<void>): Promise<void> {
    if (this._done) return

    await this._mutex.lock()
    try {
      if (!this._done) {
        await fn()
        this._done = true
      }
    } finally {
      this._mutex.unlock()
    }
  }

  done(): boolean {
    return this._done
  }
}

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  async function example() {
    const once = new Once()

    // 并发调用
    await Promise.all([
      once.do(() => console.log('只会执行一次')),
      once.do(() => console.log('这次不会执行')),
      once.do(() => console.log('这次也不会执行'))
    ])
  }
  example()
}
