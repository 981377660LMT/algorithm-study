/* eslint-disable no-param-reassign */

// 重链剖分专用树状数组

// !Range Add Range Sum, 0-based.
class BITArray2 {
  private readonly _n: number
  private readonly _tree1: number[]
  private readonly _tree2: number[]

  constructor(n: number) {
    this._n = n
    this._tree1 = Array(n + 1).fill(0)
    this._tree2 = Array(n + 1).fill(0)
  }

  // 切片内[start, end)的每个元素加上delta.
  //  0<=start<=end<=n
  add(start: number, end: number, delta: number) {
    end--
    this._add(start, delta)
    this._add(end + 1, -delta)
  }

  // 求切片内[start, end)的和.
  //  0<=start<=end<=n
  query(start: number, end: number) {
    end--
    return this._query(end) - this._query(start - 1)
  }

  private _add(index: number, delta: number) {
    index++
    const rawIndex = index
    for (let i = index; i <= this._n; i += i & -i) {
      this._tree1[i] += delta
      this._tree2[i] += (rawIndex - 1) * delta
    }
  }

  private _query(index: number) {
    index++
    if (index > this._n) {
      index = this._n
    }
    const rawIndex = index
    let res = 0
    for (let i = index; i > 0; i -= i & -i) {
      res += rawIndex * this._tree1[i] - this._tree2[i]
    }
    return res
  }
}

if (require.main === module) {
  // check with bruteforce
  const arr = Array(100).fill(0)
  const bit = new BITArray2(100)
  for (let i = 0; i < 100; i++) {
    for (let j = i; j < 100; j++) {
      const rand = Math.floor(Math.random() * 100)
      arr[j] += rand
      bit.add(j, j + 1, rand)
    }

    for (let j = 0; j < 100; j++) {
      for (let k = j; k < 100; k++) {
        const sum = arr.slice(j, k + 1).reduce((a, b) => a + b, 0)
        if (sum !== bit.query(j, k + 1)) {
          throw new Error('wrong')
        }
      }
    }
  }

  console.log('ok')
}

export {}
