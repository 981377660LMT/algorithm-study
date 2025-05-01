import { Heap } from '../../../Heap'

type WorkerGenerator = Generator<void, void, WorkerGenerator>

interface ITimerItem {
  time: number
  index: number
  gen: WorkerGenerator
}

interface IBridgeLockItem {
  isLeft: boolean
  crossTime: number
  index: number
  gen: WorkerGenerator
}

function findCrossingTime(n: number, k: number, time: number[][]): number {
  const timerPq = new Heap<ITimerItem>({
    data: [],
    less(a, b) {
      if (a.time !== b.time) return a.time < b.time
      return a.index < b.index
    }
  })

  const bridgeLockPq = new Heap<IBridgeLockItem>({
    data: [],
    less(a, b) {
      if (a.isLeft !== b.isLeft) return !a.isLeft
      if (a.crossTime !== b.crossTime) return a.crossTime > b.crossTime
      return a.index > b.index
    }
  })

  let bridgeLocked = false
  let curTime = 0
  let res = 0
  let remaining = n
  for (let i = 0; i < k; i++) {
    const gen = worker(i)
    gen.next()
    gen.next(gen)
  }

  while (timerPq.size || bridgeLockPq.size) {
    if (timerPq.size && timerPq.top().time <= curTime) {
      const { gen } = timerPq.pop()
      gen.next()
    } else if (bridgeLockPq.size && !bridgeLocked) {
      const { gen } = bridgeLockPq.pop()
      gen.next()
    } else {
      curTime = timerPq.top().time
    }
  }

  return res

  function* worker(i: number): WorkerGenerator {
    const gen = yield
    while (remaining > 0) {
      yield* _acquireBridgeLock(i, true, gen)
      if (remaining > 0) {
        yield* _acquireTime(i, time[i][0], gen)
        bridgeLocked = false
      } else {
        bridgeLocked = false
        break
      }

      remaining--
      yield* _acquireTime(i, time[i][1], gen)

      yield* _acquireBridgeLock(i, false, gen)
      yield* _acquireTime(i, time[i][2], gen)
      bridgeLocked = false

      res = curTime

      yield* _acquireTime(i, time[i][3], gen)
    }
  }

  /**
   * 申请时间片.
   */
  function* _acquireTime(i: number, duration: number, gen: WorkerGenerator): WorkerGenerator {
    timerPq.push({
      gen,
      time: curTime + duration,
      index: i
    })
    yield
  }

  /**
   * 申请过桥的锁.
   */
  function* _acquireBridgeLock(i: number, isLeft: boolean, gen: WorkerGenerator): WorkerGenerator {
    bridgeLockPq.push({
      gen,
      isLeft,
      crossTime: time[i][0] + time[i][2],
      index: i
    })
    yield
    bridgeLocked = true
  }
}
