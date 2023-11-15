/**
 * 将时间轴上单点的 add 和 remove 转化为区间上的 add, 删除时元素必须要存在。
 * 如果 add 和 remove 按照时间顺序严格单增，那么可以使用 monotone = true 来加速。
 * !不能加入相同的元素，删除时元素必须要存在。
 */
class AddRemoveQuery<V extends string | number = number> {
  private readonly _mp: Map<V, number> = new Map()
  private readonly _events: { start: number; end: number; value: V }[] = []
  private readonly _adds: Map<V, number[]> = new Map()
  private readonly _removes: Map<V, number[]> = new Map()
  private readonly _monotone: boolean

  /**
   * @param monotone add 和 remove 是否按照时间顺序严格单增
   */
  constructor(monotone: boolean) {
    this._monotone = monotone
  }

  add(time: number, value: V): void {
    if (this._monotone) {
      this._addMonotone(time, value)
    } else {
      const adds = this._adds.get(value)
      if (adds == undefined) {
        this._adds.set(value, [time])
      } else {
        adds.push(time)
      }
    }
  }

  remove(time: number, value: V): void {
    if (this._monotone) {
      this._removeMonotone(time, value)
    } else {
      const removes = this._removes.get(value)
      if (removes == undefined) {
        this._removes.set(value, [time])
      } else {
        removes.push(time)
      }
    }
  }

  /**
   * @param lastTime 所有变更都结束的时间.例如INF.
   */
  work(lastTime: number): { start: number; end: number; value: V }[] {
    if (this._monotone) {
      return this._workMonotone(lastTime)
    }
    const res: { start: number; end: number; value: V }[] = []
    this._adds.forEach((addTimes, value) => {
      let removeTimes: number[] = []
      if (this._removes.has(value)) {
        removeTimes = this._removes.get(value)!
        this._removes.delete(value)
      }
      if (removeTimes.length < addTimes.length) {
        removeTimes.push(lastTime)
      }
      addTimes.sort((a, b) => a - b)
      removeTimes.sort((a, b) => a - b)
      for (let i = 0; i < addTimes.length; i++) {
        const t = addTimes[i]
        if (t < removeTimes[i]) {
          res.push({ start: t, end: removeTimes[i], value })
        }
      }
    })
    return res
  }

  private _addMonotone(time: number, value: V): void {
    if (this._mp.has(value)) {
      throw new Error('can not add a value that already exists')
    }
    this._mp.set(value, time)
  }

  private _removeMonotone(time: number, value: V): void {
    const startTime = this._mp.get(value)
    if (startTime == undefined) {
      throw new Error("can't remove a value that doesn't exist")
    }
    this._mp.delete(value)
    if (startTime !== time) {
      this._events.push({ start: startTime, end: time, value })
    }
  }

  private _workMonotone(lastTime: number): { start: number; end: number; value: V }[] {
    this._mp.forEach((startTime, value) => {
      if (startTime !== lastTime) {
        this._events.push({ start: startTime, end: lastTime, value })
      }
    })
    return this._events
  }
}

export { AddRemoveQuery }

if (require.main === module) {
  const Q = new AddRemoveQuery(false)
  Q.add(1, 1)
  Q.add(2, 2)
  Q.add(3, 3)
  Q.remove(4, 1)
  console.log(Q.work(5))
}
