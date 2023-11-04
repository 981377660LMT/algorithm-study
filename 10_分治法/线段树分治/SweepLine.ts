/* eslint-disable no-inner-declarations */
/* eslint-disable class-methods-use-this */

import { SortedListFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

/**
 * 给定一个时间轴（或者设想一个），
 * 有若干个操作(可交换,commutative)在时间 [start,end) 中起作用。
 * 询问某一个时间某个值是什么.
 * 如果修改操作可删除，那么可以使用'扫描线'来解决.
 */
class SweepLine {
  private readonly _mutate: (mutationId: number) => void
  private readonly _remove: (mutationId: number) => void
  private readonly _query: (queryId: number) => void
  private readonly _mutations: { start: number; end: number; id: number }[] = []
  private readonly _queries: { time: number; id: number }[] = []

  /**
   * 在每个位置处保存对应的变更和查询的编号.
   * 非负偶数:query, 非负奇数:mutate, 负数: remove.
   */
  private _nodes: number[][] = []

  private _mutationId = 0
  private _queryId = 0

  /**
   * 使用扫描线得到每个时间点的答案.
   * @param mutate 添加编号为`mutationId`的变更.
   * @param remove 删除编号为`mutationId`的变更.
   * @param query 响应编号为`queryId`的查询.
   * @complexity 一共调用 **O(n)** 次`mutate`、`remove` 和 **O(q)** 次`query`.
   */
  constructor(
    mutate: (mutationId: number) => void,
    remove: (mutationId: number) => void,
    query: (queryId: number) => void
  )
  constructor(
    options: {
      mutate: (mutationId: number) => void
      remove: (mutationId: number) => void
      query: (queryId: number) => void
    } & ThisType<void>
  )
  constructor(arg1: any, arg2?: any, arg3?: any) {
    if (typeof arg1 === 'object') {
      this._mutate = arg1.mutate
      this._remove = arg1.remove
      this._query = arg1.query
    } else {
      this._mutate = arg1
      this._remove = arg2
      this._query = arg3
    }
  }

  /**
   * 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
   */
  addMutation(startTime: number, endTime: number, id?: number): void {
    if (startTime >= endTime) return
    if (id == undefined) id = this._mutationId++
    this._mutations.push({ start: startTime, end: endTime, id })
  }

  /**
   * 在时间`time`时添加一个编号为`id`的查询.
   */
  addQuery(time: number, id?: number): void {
    if (id == undefined) id = this._queryId++
    this._queries.push({ time, id })
  }

  run(): void {
    if (!this._queries.length) return
    const times: number[] = Array(this._queries.length)
    for (let i = 0; i < this._queries.length; i++) times[i] = this._queries[i].time
    times.sort((a, b) => a - b)
    this._uniqueInplace(times)
    const usedTimes = new Uint8Array(times.length + 1)
    usedTimes[0] = 1
    for (let i = 0; i < this._mutations.length; i++) {
      const e = this._mutations[i]
      usedTimes[this._lowerBound(times, e.start)] = 1
      usedTimes[this._lowerBound(times, e.end)] = 1
    }
    for (let i = 1; i < times.length; i++) {
      if (!usedTimes[i]) times[i] = times[i - 1]
    }
    this._uniqueInplace(times)

    this._nodes = Array(times.length + 1)
    for (let i = 0; i < this._nodes.length; i++) this._nodes[i] = []
    for (let i = 0; i < this._mutations.length; i++) {
      const e = this._mutations[i]
      let left = this._lowerBound(times, e.start)
      let right = this._lowerBound(times, e.end)
      const eid = e.id * 2 + 1
      this._nodes[left].push(eid)
      this._nodes[right].push(-eid)
    }

    for (let i = 0; i < this._queries.length; i++) {
      const q = this._queries[i]
      this._nodes[this._upperBound(times, q.time) - 1].push(q.id * 2) // query
    }

    this._doSweep()
  }

  clear(): void {
    this._nodes = []
    this._queries.length = 0
    this._mutations.length = 0
    this._mutationId = 0
    this._queryId = 0
  }

  private _doSweep(): void {
    for (let time = 0; time < this._nodes.length; time++) {
      const events = this._nodes[time]
      for (let i = 0; i < events.length; i++) {
        const id = events[i]
        if (id >= 0) {
          if (id & 1) {
            this._mutate((id - 1) / 2) // mutate
          } else {
            this._query(id / 2) // query
          }
        } else {
          this._remove((-id - 1) / 2) // remove
        }
      }
    }
  }

  private _uniqueInplace(sorted: number[]): void {
    let slow = 0
    for (let fast = 0; fast < sorted.length; fast++) {
      if (sorted[fast] !== sorted[slow]) sorted[++slow] = sorted[fast]
    }
    sorted.length = slow + 1
  }

  private _lowerBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private _upperBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

export { SweepLine }

if (require.main === module) {
  demo()
  function demo(): void {
    const sweepLine = new SweepLine({
      mutate(id) {
        console.log(`mutate ${id}`)
      },
      remove(id) {
        console.log(`remove ${id}`)
      },
      query(id) {
        console.log(`query ${id}`)
      }
    })

    sweepLine.addMutation(-1, 10, 0)
    sweepLine.addMutation(5, 15, 1)
    sweepLine.addMutation(10, 20, 2)
    sweepLine.addQuery(-100, 7)
    sweepLine.addMutation(15, 25, 3)
    sweepLine.addMutation(20, 30, 4)
    sweepLine.addQuery(-2, 0)
    sweepLine.addQuery(5, 1)
    sweepLine.addQuery(10, 2)
    sweepLine.addQuery(15, 3)
    sweepLine.addQuery(20, 4)
    sweepLine.addQuery(25, 5)
    sweepLine.addQuery(25, 6)

    sweepLine.run()
  }

  // 2251. 花期内花的数目
  // https://leetcode.cn/problems/number-of-flowers-in-full-bloom/submissions/
  function fullBloomFlowers(flowers: number[][], people: number[]): number[] {
    const res = Array(people.length).fill(0)
    let count = 0
    const sweepLine = new SweepLine({
      mutate(id) {
        count++
      },
      remove(id) {
        count--
      },
      query(id) {
        res[id] = count
      }
    })

    for (let i = 0; i < flowers.length; i++) {
      const { 0: left, 1: right } = flowers[i]
      sweepLine.addMutation(left, right + 1, i)
    }
    for (let i = 0; i < people.length; i++) {
      sweepLine.addQuery(people[i], i)
    }

    sweepLine.run()
    return res
  }

  // 2747. 统计没有收到请求的服务器数目
  // https://leetcode.cn/problems/count-zero-request-servers/
  function countServers(n: number, logs: number[][], x: number, queries: number[]): number[] {
    const res = Array(queries.length).fill(0)
    let count = n
    const serverCount = new Map<number, number>()
    const sweepLine = new SweepLine({
      mutate(id) {
        const serverId = logs[id][0]
        const c = (serverCount.get(serverId) || 0) + 1
        serverCount.set(serverId, c)
        if (c === 1) count--
      },
      remove(id) {
        const serverId = logs[id][0]
        const c = (serverCount.get(serverId) || 0) - 1
        serverCount.set(serverId, c)
        if (c === 0) count++
      },
      query(id) {
        res[id] = count
      }
    })

    for (let i = 0; i < logs.length; i++) {
      const start = logs[i][1]
      sweepLine.addMutation(start, start + x + 1, i) // 可以覆盖[start, start + x] 的查询区间
    }
    for (let i = 0; i < queries.length; i++) {
      sweepLine.addQuery(queries[i], i)
    }

    sweepLine.run()
    return res
  }

  // 1851. 包含每个查询的最小区间
  // https://leetcode.cn/problems/minimum-interval-to-include-each-query/
  function minInterval(intervals: number[][], queries: number[]): number[] {
    const res = Array(queries.length).fill(0)
    const minLen = new SortedListFast<number>()
    const sweepLine = new SweepLine({
      mutate(id) {
        const interval = intervals[id]
        minLen.add(interval[1] - interval[0] + 1)
      },
      remove(id) {
        const interval = intervals[id]
        minLen.discard(interval[1] - interval[0] + 1)
      },
      query(id) {
        res[id] = minLen.length ? minLen.min : -1
      }
    })

    for (let i = 0; i < intervals.length; i++) {
      const { 0: left, 1: right } = intervals[i]
      sweepLine.addMutation(left, right + 1, i)
    }
    for (let i = 0; i < queries.length; i++) {
      sweepLine.addQuery(queries[i], i)
    }

    sweepLine.run()
    return res
  }

  // 1847. 最近的房间
  // https://leetcode.cn/problems/closest-room/
  function closestRoom(rooms: number[][], queries: number[][]): number[] {
    const res = Array(queries.length).fill(-1)
    const sl = new SortedListFast<number>()
    const sweepLine = new SweepLine({
      mutate(id) {
        const roomId = rooms[id][0]
        sl.add(roomId)
      },
      remove(id) {
        const roomId = rooms[id][0]
        sl.discard(roomId)
      },
      query(id) {
        const preferred = queries[id][0]
        let minDiff = Infinity
        const floor = sl.floor(preferred)
        if (floor != undefined) {
          const diff = Math.abs(floor - preferred)
          if (diff < minDiff || (diff === minDiff && floor < res[id])) {
            res[id] = floor
            minDiff = diff
          }
        }
        const ceiling = sl.ceiling(preferred)
        if (ceiling != undefined) {
          const diff = Math.abs(ceiling - preferred)
          if (diff < minDiff || (diff === minDiff && ceiling < res[id])) {
            res[id] = ceiling
            minDiff = diff
          }
        }
      }
    })

    for (let i = 0; i < rooms.length; i++) {
      const size = rooms[i][1]
      sweepLine.addMutation(0, size + 1, i)
    }
    for (let i = 0; i < queries.length; i++) {
      const minSize = queries[i][1]
      sweepLine.addQuery(minSize, i)
    }

    sweepLine.run()
    return res
  }
}
