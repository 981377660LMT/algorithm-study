/* eslint-disable no-inner-declarations */
// Van Tree.
// 梵峨眉大悲寺树.

const INF = 2e15

/**
 * !建立一个元素范围为(-INF,INF)的VanEmdeBoasTree.
 * @warning 允许插入的元素范围为[0,INF).
 */
class VanEmdeBoasTree {
  private readonly _root: _VNode
  private _size = 0

  /**
   * @param depth 树的深度.默认为16.一般取16或32.
   */
  constructor(depth = 16) {
    this._root = new _VNode(depth)
  }

  has(x: number): boolean {
    return this._root.has(x)
  }

  /**
   * !0<=x<INF.
   */
  insert(x: number): boolean {
    if (x < 0) throw new Error('inserted value must be non-negative!')
    if (this.has(x)) return false
    this._size++
    this._root.insert(x)
    return true
  }

  erase(x: number): boolean {
    if (!this.has(x)) return false
    this._size--
    this._root.erase(x)
    return true
  }

  /**
   * 返回小于等于i的最大元素.如果不存在,返回-INF.
   */
  prev(x: number): number {
    return this._root.prev(x)
  }

  /**
   * 返回大于等于i的最小元素.如果不存在,返回INF.
   */
  next(x: number): number {
    return this._root.next(x)
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x >= end) break
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(-INF, INF, v => sb.push(v.toString()))
    return `VanEmdeBoasTree(${this.size}){${sb.join(', ')}}`
  }

  /**
   * 如果没有元素,返回INF.
   */
  get min(): number {
    return this._root.min
  }

  /**
   * 如果没有元素,返回-INF.
   */
  get max(): number {
    return this._root.max
  }

  get size(): number {
    return this._size
  }
}

class _VNode {
  dep: number
  min = INF
  max = -INF
  aux: _VNode | undefined = undefined
  son: Map<number, _VNode> = new Map()

  constructor(dep: number) {
    this.dep = dep
  }

  has(x: number): boolean {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x === vMin || x === vMax) return true
    if (!vDep || x < vMin || x > vMax) return false
    const i = x >>> vDep
    const soni = this.son.get(i)
    if (!soni) return false
    return soni.has(x - (i << vDep))
  }

  insert(x: number): void {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (vMin > vMax) {
      this.min = x
      this.max = x
      return
    }
    if (vMin === vMax) {
      if (x < vMin) {
        this.min = x
        return
      }
      if (x > vMax) {
        this.max = x
        return
      }
    }
    if (x < vMin) {
      const tmp = x
      x = vMin
      this.min = tmp
    }
    if (x > vMax) {
      const tmp = x
      x = vMax
      this.max = tmp
    }
    const i = x >>> vDep
    let soni = this.son.get(i)
    if (!soni) {
      soni = new _VNode(vDep >>> 1)
      this.son.set(i, soni)
    }
    if (soni.empty()) {
      if (!this.aux) this.aux = new _VNode(vDep >>> 1)
      this.aux.insert(i)
    }
    soni.insert(x - (i << vDep))
  }

  erase(x: number): void {
    const { min: vMin, max: vMax, dep: vDep, aux: vAux } = this
    if (vMin === x && vMax === x) {
      this.min = INF
      this.max = -INF
      return
    }
    if (x === vMin) {
      if (!vAux || vAux.empty()) {
        this.min = vMax
        return
      }
      const auxMin = vAux.min
      x = (auxMin << vDep) + this.son.get(auxMin)!.min
      this.min = x
    }
    if (x === vMax) {
      if (!vAux || vAux.empty()) {
        this.max = vMin
        return
      }
      const auxMax = vAux.max
      x = (auxMax << vDep) + this.son.get(auxMax)!.max
      this.max = x
    }
    const i = x >>> vDep
    const soni = this.son.get(i)!
    soni.erase(x - (i << vDep))
    if (soni.empty()) vAux!.erase(i)
  }

  prev(x: number): number {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x < vMin) return -INF
    if (x >= vMax) return vMax
    const i = x >>> vDep
    const hi = i << vDep
    const lo = x - hi
    const soni = this.son.get(i)
    if (soni && lo >= soni.min) return hi + soni.prev(lo)
    let y = -INF
    if (this.aux && i > 0) y = this.aux.prev(i - 1)
    if (y === -INF) return vMin
    return (y << vDep) + this.son.get(y)!.max
  }

  next(x: number): number {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x <= vMin) return vMin
    if (x > vMax) return INF
    const i = x >>> vDep
    const hi = i << vDep
    const lo = x - hi
    const soni = this.son.get(i)
    if (soni && lo <= soni.max) return hi + soni.next(lo)
    let y = INF
    if (this.aux) y = this.aux.next(i + 1)
    if (y === INF) return vMax
    return (y << vDep) + this.son.get(y)!.min
  }

  empty(): boolean {
    return this.min > this.max
  }
}

export { VanEmdeBoasTree }

if (require.main === module) {
  const van = new VanEmdeBoasTree()
  console.log(van.min, van.max, van.size)
  van.insert(0)
  console.log(van.prev(0.5))
  console.log(van.prev(-1))
  console.log(van.prev(1))

  const n = 2e5
  const set2 = new VanEmdeBoasTree()
  console.time('VanEmdeBoasTree')
  for (let i = 0; i < n; i++) {
    set2.insert(i)
    set2.next(i)
    set2.prev(i)
    set2.has(i)
    set2.erase(i)
    set2.insert(i)
  }
  console.timeEnd('VanEmdeBoasTree') // 360ms

  testWithBruteForce()
  function testWithBruteForce(): void {
    class Mocker {
      private readonly _set = new Set<number>()
      insert(x: number): void {
        this._set.add(x)
      }
      erase(x: number): void {
        this._set.delete(x)
      }
      has(x: number): boolean {
        return this._set.has(x)
      }
      prev(x: number): number {
        let res = -INF
        for (const v of this._set) {
          if (v > x) continue
          if (v > res) res = v
        }
        return res
      }
      next(x: number): number {
        let res = INF
        for (const v of this._set) {
          if (v < x) continue
          if (v < res) res = v
        }
        return res
      }
    }

    const mocker = new Mocker()
    const set = new VanEmdeBoasTree()
    for (let i = 0; i < 4e4; i++) {
      const x = Math.floor(Math.random() * 1000)
      if (Math.random() < 0.5) {
        mocker.insert(x)
        set.insert(x)
      } else {
        mocker.erase(x - 500)
        set.erase(x - 500)
      }
      const x2 = Math.floor(Math.random() * 1000) - 500
      const prev1 = mocker.prev(x2)
      const prev2 = set.prev(x2)
      if (prev1 !== prev2) {
        console.log('prev1:', prev1, 'prev2:', prev2)
        throw new Error()
      }
      const next1 = mocker.next(x2)
      const next2 = set.next(x2)
      if (next1 !== next2) {
        console.log('next1:', next1, 'next2:', next2)
        throw new Error()
      }
    }
  }
  console.log('ok')
}
