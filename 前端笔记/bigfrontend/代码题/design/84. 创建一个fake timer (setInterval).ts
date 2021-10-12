import { Func } from '../typings'

interface QueueItem {
  id: number
  cb: (...args: any[]) => void
  time: number
  args: any[]
}

/**
 * 模拟一个 priorityQueue 来 eventLoop
 */
class FakeTimer {
  private original: Record<'setTimeout' | 'clearTimeout' | 'now', Func>
  private timerId: number
  private currentTime: number
  private priorityQueue: QueueItem[] // 应使用小根堆

  constructor() {
    this.original = {
      setTimeout: window.setTimeout,
      clearTimeout: window.clearTimeout,
      now: Date.now,
    }
    this.timerId = 1
    this.currentTime = 0
    this.priorityQueue = []
  }

  install() {
    // replace window.setTimeout, window.clearTimeout, Date.now
    // with your implementation
    // @ts-ignore
    window.setTimeout = (cb: (args: void) => void, time: number, ...args: any[]): number => {
      const id = this.timerId++
      this.priorityQueue.push({
        id,
        cb,
        args,
        time: time + this.currentTime,
      })
      this.priorityQueue.sort((a, b) => a.time - b.time)
      return id
    }

    // @ts-ignore
    window.clearTimeout = (removeId: number) => {
      this.priorityQueue = this.priorityQueue.filter(({ id }) => id !== removeId)
    }

    Date.now = () => this.currentTime
  }

  uninstall() {
    // restore the original implementation of
    // window.setTimeout, window.clearTimeout, Date.now
    // @ts-ignore
    window.setTimeout = this.original.setTimeout
    window.clearTimeout = this.original.clearTimeout
    Date.now = this.original.now
  }

  // 执行时改变currentTime
  tick() {
    // run the scheduled functions without waiting
    while (this.priorityQueue.length) {
      const { cb, time, args } = this.priorityQueue.shift()!
      this.currentTime = time
      cb(...args) // 最好是保存this变量
    }
  }
}

const fakeTimer = new FakeTimer()
fakeTimer.install()

const logs: any[] = []
const log = (arg: string) => {
  logs.push([Date.now(), arg])
}

setTimeout(() => log('A'), 100)
// log 'A' at 100

const b = setTimeout(() => log('B'), 110)
clearTimeout(b)
// b is set but cleared

setTimeout(() => log('C'), 200)

expect(logs).toEqual([
  [100, 'A'],
  [200, 'C'],
])

fakeTimer.uninstall()

export {}
