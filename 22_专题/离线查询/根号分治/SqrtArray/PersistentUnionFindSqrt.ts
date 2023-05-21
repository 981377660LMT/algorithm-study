/* eslint-disable eqeqeq */
/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-constant-condition */

//
// Persistence Union Find
//
// Description:
//   Use persistent array instead of standard array in union find data structure
//
// Complexity:
//   O(a* T(n)), where T(n) is a complexity of persistent array

import { PersistentArraySqrtInt32 } from './PersistentArraySqrt'

/**
 * 完全可持久化并查集.
 */
class PersistentUnionFindSqrt {
  private readonly _parent: PersistentArraySqrtInt32
  constructor(nOrArray: number | PersistentArraySqrtInt32) {
    if (typeof nOrArray === 'number') {
      this._parent = new PersistentArraySqrtInt32(new Int32Array(nOrArray).fill(-1))
    } else {
      this._parent = nOrArray
    }
  }

  union(u: number, v: number): PersistentUnionFindSqrt {
    let root1 = this.find(u)
    let root2 = this.find(v)
    if (root1 === root2) return this
    let p1 = this._parent.get(root1)!
    let p2 = this._parent.get(root2)!
    if (p1 > p2) {
      root1 ^= root2
      root2 ^= root1
      root1 ^= root2
    }
    let newP = this._parent.set(root1, p1 + p2)
    newP = newP.set(root2, root1)
    return new PersistentUnionFindSqrt(newP)
  }

  find(u: number): number {
    while (true) {
      const p = this._parent.get(u)
      if (p == undefined) {
        throw new Error('invalid index')
      }
      if (p < 0) break
      u = p
    }
    return u
  }

  isConnected(u: number, v: number): boolean {
    return this.find(u) === this.find(v)
  }

  getSize(u: number): number {
    return -this._parent.get(this.find(u))!
  }
}

export { PersistentUnionFindSqrt }

if (require.main === module) {
  let n = 10
  let uf0 = new PersistentUnionFindSqrt(n)
  let uf1 = uf0.union(0, 1)
  let uf2 = uf1.union(1, 2)
  console.log(getGroups(uf0))
  console.log(getGroups(uf1))
  console.log(getGroups(uf2))

  console.time('test 1e5')
  n = 1e5
  let uf = new PersistentUnionFindSqrt(n)
  for (let i = 0; i < n - 1; i++) {
    uf = uf.union(i, i + 1)
  }
  console.timeEnd('test 1e5')

  function getGroups(uf: PersistentUnionFindSqrt): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < n; i++) {
      const root = uf.find(i)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(i)
    }
    return groups
  }

  // uf0 := NewPersistentUnionFindSqrt(3)
  // uf1 := uf0.Union(0, 1)
  // uf2 := uf0.Union(1, 2)
  // fmt.Println(uf1.IsConnected(0, 1), uf1.Find(0), uf1.Find(1))
  // fmt.Println(uf1.IsConnected(0, 2), uf1.Find(0), uf1.Find(2))
  // fmt.Println(uf1.IsConnected(1, 2), uf1.Find(1), uf1.Find(2))
  // fmt.Println(uf2.IsConnected(0, 1), uf2.Find(0), uf2.Find(1))
  // fmt.Println(uf2.IsConnected(0, 2), uf2.Find(0), uf2.Find(2))
  // fmt.Println(uf2.IsConnected(1, 2), uf2.Find(1), uf2.Find(2))
  uf0 = new PersistentUnionFindSqrt(3)
  uf1 = uf0.union(0, 1)
  uf2 = uf0.union(1, 2)
  console.log(uf1.isConnected(0, 1), uf1.find(0), uf1.find(1))
  console.log(uf1.isConnected(0, 2), uf1.find(0), uf1.find(2))
  console.log(uf1.isConnected(1, 2), uf1.find(1), uf1.find(2))
  console.log(uf2.isConnected(0, 1), uf2.find(0), uf2.find(1))
  console.log(uf2.isConnected(0, 2), uf2.find(0), uf2.find(2))
  console.log(uf2.isConnected(1, 2), uf2.find(1), uf2.find(2))
}
