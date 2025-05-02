import { Heap } from '../../../Heap'

interface IEvent<R> {}
class Sleep implements IEvent<unknown> {
  constructor(public duration: number) {}
}
class Acquire<R> implements IEvent<R> {
  constructor(
    public resource: R,
    public priority: number
  ) {}
}
class Release<R> implements IEvent<R> {
  constructor(public resource: R) {}
}

type TaskGenerator<R> = Generator<IEvent<R>, void, TaskGenerator<R>>

class Task<R> {
  started = false
  constructor(public gen: TaskGenerator<R>) {}
}

class Scheduler<R> {
  currentTime = 0

  private readonly _timePQ = new Heap<{ wakeTime: number; task: Task<R> }>({
    data: [],
    less: (a, b) => a.wakeTime < b.wakeTime
  })

  private readonly _waiting = new Map<R, Heap<{ priority: number; task: Task<R> }>>()

  private _waitingCount = 0
  private readonly _locked = new Map<R, boolean>()

  startTask(task: Task<R>): void {
    if (task.started) return
    task.started = true
    // 启动到第一个 yield（Acquire / Sleep）
    const event = task.gen.next(task.gen).value
    this._dispatch(task, event)
  }

  run(): void {
    while (this._timePQ.size() || this._waitingCount) {
      // 1) 时间优先
      if (this._timePQ.size() && this._timePQ.peek().wakeTime <= this.currentTime) {
        const { wakeTime, task } = this._timePQ.pop()!
        this.currentTime = wakeTime
        const ev = task.gen.next().value
        this._dispatch(task, ev)
        continue
      }

      // 2) 资源唤醒
      for (const [res, heap] of this._waiting.entries()) {
        if (!this._locked.get(res) && heap.size()) {
          const { task } = heap.pop()!
          this._waitingCount--
          this._locked.set(res, true)
          const ev = task.gen.next().value
          this._dispatch(task, ev)
        }
      }

      // 3) 快进时间
      if (this._timePQ.size()) {
        this.currentTime = this._timePQ.peek().wakeTime
      }
    }
  }

  private _dispatch(task: Task<R>, event?: IEvent<R>) {
    if (!event) return

    if (event instanceof Sleep) {
      this._timePQ.push({ task, wakeTime: this.currentTime + event.duration })
    } else if (event instanceof Acquire) {
      const heap =
        this._waiting.get(event.resource) ??
        new Heap<{ priority: number; task: Task<R> }>({
          data: [],
          less: (a, b) => a.priority < b.priority
        })
      heap.push({ task, priority: -event.priority })
      this._waiting.set(event.resource, heap)
      this._waitingCount++
    } else if (event instanceof Release) {
      this._locked.set(event.resource, false)
      // 立即让当前协程继续
      const nextEvent = task.gen.next().value as IEvent<R>
      this._dispatch(task, nextEvent)
    } else {
      throw new Error('Unknown event')
    }
  }
}

function findCrossingTime(n: number, k: number, times: number[][]): number {
  type Res = 'bridge'
  const scheduler = new Scheduler<Res>()
  const nRef = { value: n }
  const resRef = { value: 0 }

  function* worker(id: number): TaskGenerator<Res> {
    while (nRef.value > 0) {
      // 左岸申请
      yield new Acquire<Res>('bridge', times[id][0] + times[id][2])
      // 过桥
      yield new Sleep(times[id][0])
      yield new Release<Res>('bridge')

      // 装货
      nRef.value--
      yield new Sleep(times[id][1])

      // 右岸申请
      yield new Acquire<Res>('bridge', times[id][0] + times[id][2])
      yield new Sleep(times[id][2])
      yield new Release<Res>('bridge')

      // 记录时间
      resRef.value = scheduler.currentTime

      // 回左岸装下一件
      yield new Sleep(times[id][3])
    }
  }

  for (let i = 0; i < k; i++) {
    const task = new Task<Res>(worker(i))
    scheduler.startTask(task)
  }

  scheduler.run()
  return resRef.value
}

// 测试
console.log(
  findCrossingTime(1, 3, [
    [1, 1, 2, 1],
    [1, 1, 3, 1],
    [1, 1, 4, 1]
  ])
) // => 6
