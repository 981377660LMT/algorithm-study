/* eslint-disable no-useless-constructor */
import { Heap } from '../../../Heap'

// #region Event

interface IBaseEvent<T extends string> {
  type: T
}

interface ISleepEvent extends IBaseEvent<'sleep'> {
  duration: number
}

interface IAcquireEvent extends IBaseEvent<'acquire'> {
  resource: PropertyKey
  priority: IPriority
}

interface IPriority {
  less(other: this): boolean
}

interface IReleaseEvent extends IBaseEvent<'release'> {
  resource: PropertyKey
}

type Event = ISleepEvent | IAcquireEvent | IReleaseEvent

// #endregion

// #region Task

type Task<E extends Event = Event> = Generator<E, void, void>

// #endregion

// #region Scheduler

class Scheduler {
  currentTime = 0

  private readonly _timePQ = new Heap<{ wakeTime: number; task: Task }>({
    data: [],
    less: (a, b) => a.wakeTime < b.wakeTime
  })

  private readonly _waiting = new Map<PropertyKey, Heap<{ priority: IPriority; task: Task }>>()
  private _waitingCount = 0

  private readonly _locked = new Set<PropertyKey>()

  startTask(task: Task): void {
    const event = task.next().value
    this._dispatch(task, event)
  }

  run(): void {
    while (this._timePQ.size || this._waitingCount) {
      if (this._timePQ.size && this._timePQ.top().wakeTime <= this.currentTime) {
        const { wakeTime, task } = this._timePQ.pop()
        this.currentTime = wakeTime
        const event = task.next().value
        this._dispatch(task, event)
        continue
      }

      for (const [res, pq] of this._waiting.entries()) {
        if (pq.size && !this._locked.has(res)) {
          const { task } = pq.pop()
          this._waitingCount--
          this._locked.add(res)
          const event = task.next().value
          this._dispatch(task, event)
        }
      }

      if (this._timePQ.size) {
        this.currentTime = this._timePQ.top().wakeTime
      }
    }
  }

  private _dispatch(task: Task, event: Event | void) {
    if (!event) return

    if (event.type === 'sleep') {
      this._timePQ.push({ task, wakeTime: this.currentTime + event.duration })
    } else if (event.type === 'acquire') {
      const pq = this._waiting.get(event.resource)
      if (pq) {
        pq.push({ task, priority: event.priority })
      } else {
        const newPq = new Heap<{ priority: IPriority; task: Task }>({
          data: [{ task, priority: event.priority }],
          less: (a, b) => a.priority.less(b.priority)
        })
        this._waiting.set(event.resource, newPq)
      }
      this._waitingCount++
    } else {
      this._locked.delete(event.resource)
      const nextEvent = task.next().value
      this._dispatch(task, nextEvent)
    }
  }
}

// #endregion

// #region utils

const EMPTY_SLEEP: ISleepEvent = { type: 'sleep', duration: 0 }
const EMPTY_ACQUIRE: IAcquireEvent = {
  type: 'acquire',
  resource: '',
  priority: { less: () => false }
}
const EMPTY_RELEASE: IReleaseEvent = { type: 'release', resource: '' }

function sleep(duration: number): ISleepEvent {
  EMPTY_SLEEP.duration = duration
  return EMPTY_SLEEP
}

function acquire(resource: PropertyKey, priority: IPriority): IAcquireEvent {
  EMPTY_ACQUIRE.resource = resource
  EMPTY_ACQUIRE.priority = priority
  return EMPTY_ACQUIRE
}

function release(resource: PropertyKey): IReleaseEvent {
  EMPTY_RELEASE.resource = resource
  return EMPTY_RELEASE
}

// #endregion

function findCrossingTime(n: number, k: number, times: number[][]): number {
  const scheduler = new Scheduler()
  let remaining = n
  let lastTime = 0

  class BridgePriority implements IPriority {
    constructor(
      private readonly isLeft: boolean,
      private readonly crossTime: number,
      private readonly index: number
    ) {}

    less(other: BridgePriority): boolean {
      if (this.isLeft !== other.isLeft) return !this.isLeft
      if (this.crossTime !== other.crossTime) return this.crossTime > other.crossTime
      return this.index > other.index
    }
  }

  function* worker(id: number) {
    const cross = times[id][0] + times[id][2]
    while (true) {
      // 左岸申请
      yield acquire('bridge', new BridgePriority(true, cross, id))
      if (remaining <= 0) {
        yield release('bridge')
        return
      }
      // 过桥
      yield sleep(times[id][0])
      yield release('bridge')

      // 装货
      remaining--
      yield sleep(times[id][1])

      // 右岸申请
      yield acquire('bridge', new BridgePriority(false, cross, id))
      yield sleep(times[id][2])
      yield release('bridge')

      // 记录时间
      lastTime = scheduler.currentTime

      // 回左岸装下一件
      yield sleep(times[id][3])
    }
  }

  for (let i = 0; i < k; i++) {
    const task = worker(i)
    scheduler.startTask(task)
  }

  scheduler.run()
  return lastTime
}

// 测试
console.log(
  findCrossingTime(1, 3, [
    [1, 1, 2, 1],
    [1, 1, 3, 1],
    [1, 1, 4, 1]
  ])
) // => 6
