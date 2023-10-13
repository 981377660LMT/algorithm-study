/* eslint-disable max-len */
/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
// !由于lazy模板通用性 效率不如自己维护数组的线段树
// !注意如果是单点查询,可以去掉所有pushUp函数逻辑(js使用bigint会比较慢)
// !如果是单点修改,可以去掉所有懒标记逻辑

// !一些monoid (如果难以设计半群,就使用分块解决吧)
// https://maspypy.github.io/library/alg/acted_monoid/summax_assign.hpp

import { SegmentTreePointUpdateRangeQuery } from './SegmentTreePointUpdateRangeQuery'
import { SegmentTreeRangeUpdateRangeQuery } from './SegmentTreeRangeUpdateRangeQuery'

const INF = 2e9 // !超过int32使用2e15

/**
 * 区间加,查询区间最大值(幺元为0).
 */
function createRangeAddRangeMax(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => 0,
    id: () => 0,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => f + x,
    composition: (f, g) => f + g
  })
}

/**
 * 区间加,查询区间最小值(幺元为INF).
 */
function createRangeAddRangeMin(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => INF,
    id: () => 0,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => f + x,
    composition: (f, g) => f + g
  })
}

/**
 * 区间更新最大值,查询区间最大值(幺元为0).
 * RangeChmaxRangeMax
 */
function createRangeUpdateRangeMax(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => 0,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : Math.max(f, x)),
    composition: (f, g) => (f === -INF ? g : Math.max(f, g))
  })
}

/**
 * 区间更新最小值,查询区间最小值(幺元为INF).
 * RangeChminRangeMin
 */
function createRangeUpdateRangeMin(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => INF,
    id: () => INF,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => (f === INF ? x : Math.min(f, x)),
    composition: (f, g) => (f === INF ? g : Math.min(f, g))
  })
}

/**
 * 区间赋值,查询区间和(幺元为0).
 */
function createRangeAssignRangeSum(
  nOrNums: number | ArrayLike<{ sum: number; size: number }>
): SegmentTreeRangeUpdateRangeQuery<{ sum: number; size: number }, number> {
  return new SegmentTreeRangeUpdateRangeQuery<{ sum: number; size: number }, number>(nOrNums, {
    e: () => ({ sum: 0, size: 1 }),
    id: () => -1,
    op: (e1, e2) => ({ sum: e1.sum + e2.sum, size: e1.size + e2.size }),
    mapping: (f, e) => (f === -1 ? e : { sum: f * e.size, size: e.size }),
    composition: (f, g) => (f === -1 ? g : f)
  })
}

/**
 * 区间赋值,查询区间最大值(幺元为-INF).
 */
function createRangeAssignRangeMax(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => 0,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : f),
    composition: (f, g) => (f === -INF ? g : f)
  })
}

/**
 * 区间赋值,查询区间最小值(幺元为INF).
 */
function createRangeAssignRangeMin(
  nOrNums: number | ArrayLike<number>
): SegmentTreeRangeUpdateRangeQuery<number, number> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e: () => INF,
    id: () => INF,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => (f === INF ? x : f),
    composition: (f, g) => (f === INF ? g : f)
  })
}

/**
 * 01区间翻转,查询区间和.
 */
function createRangeFlipRangeSum(
  nOrNums: number | ArrayLike<{ sum: number; size: number }>
): SegmentTreeRangeUpdateRangeQuery<{ sum: number; size: number }, number> {
  return new SegmentTreeRangeUpdateRangeQuery<{ sum: number; size: number }, number>(nOrNums, {
    e: () => ({ sum: 0, size: 1 }),
    id: () => 0,
    op: (e1, e2) => ({ sum: e1.sum + e2.sum, size: e1.size + e2.size }),
    mapping: (f, e) => (f === 0 ? e : { sum: e.size - e.sum, size: e.size }),
    composition: (f, g) => f ^ g
  })
}

/**
 * 区间赋值区间加,区间和.
 */
function createRangeAssignRangeAddRangeSum(
  nOrNums: number | ArrayLike<{ size: number; sum: number }>
): SegmentTreeRangeUpdateRangeQuery<{ size: number; sum: number }, { mul: number; add: number }> {
  return new SegmentTreeRangeUpdateRangeQuery<
    { size: number; sum: number },
    { mul: number; add: number }
  >(nOrNums, {
    e() {
      return { size: 1, sum: 0 }
    },
    id() {
      return { mul: 1, add: 0 }
    },
    op(e1, e2) {
      return { size: e1.size + e2.size, sum: e1.sum + e2.sum }
    },
    mapping(lazy, data) {
      return { size: data.size, sum: data.sum * lazy.mul + data.size * lazy.add }
    },
    composition(f, g) {
      return { mul: f.mul * g.mul, add: f.mul * g.add + f.add }
    },
    equalsId(id1, id2) {
      return id1.mul === id2.mul && id1.add === id2.add
    }
  })
}

/**
 * 区间仿射变换,区间和.
 */
function createRangeAffineRangeSum(
  nOrNums: number | ArrayLike<{ size: bigint; sum: bigint }>,
  bigMod: bigint
): SegmentTreeRangeUpdateRangeQuery<{ size: bigint; sum: bigint }, { mul: bigint; add: bigint }> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e() {
      return { size: 1n, sum: 0n }
    },
    id() {
      return { mul: 1n, add: 0n }
    },
    op(e1, e2) {
      return { size: e1.size + e2.size, sum: (e1.sum + e2.sum) % bigMod }
    },
    mapping(lazy, data) {
      return { size: data.size, sum: (data.sum * lazy.mul + data.size * lazy.add) % bigMod }
    },
    composition(f, g) {
      return { mul: (f.mul * g.mul) % bigMod, add: (f.mul * g.add + f.add) % bigMod }
    },
    equalsId(id1, id2) {
      return id1.mul === id2.mul && id1.add === id2.add
    }
  })
}

type Interval = {
  sum: number
  maxSum: number
  preMaxSum: number
  sufMaxSum: number
  minSum: number
  preMinSum: number
  sufMinSum: number
}

/**
 * 单点修改，区间最大子数组和最小子数组和.
 */
function createPointSetRangeMaxSumMinSum(arr: ArrayLike<number>): {
  fromElement: (value: number) => Interval
  tree: SegmentTreePointUpdateRangeQuery<Interval>
} {
  const leaves: Interval[] = Array(arr.length)
  const fromElement = (v: number) => ({
    sum: v,
    maxSum: v,
    preMaxSum: v,
    sufMaxSum: v,
    minSum: v,
    preMinSum: v,
    sufMinSum: v
  })
  for (let i = 0; i < arr.length; ++i) {
    leaves[i] = fromElement(arr[i])
  }

  return {
    fromElement,
    tree: new SegmentTreePointUpdateRangeQuery(
      leaves,
      () => ({
        sum: 0,
        maxSum: -INF,
        preMaxSum: -INF,
        sufMaxSum: -INF,
        minSum: INF,
        preMinSum: INF,
        sufMinSum: INF
      }),
      (e1, e2) => ({
        sum: Math.min(Math.max(e1.sum + e2.sum, -INF), INF),
        maxSum: Math.max(e1.maxSum, e2.maxSum, e1.sufMaxSum + e2.preMaxSum),
        preMaxSum: Math.max(e1.preMaxSum, e1.sum + e2.preMaxSum),
        sufMaxSum: Math.max(e2.sufMaxSum, e2.sum + e1.sufMaxSum),
        minSum: Math.min(e1.minSum, e2.minSum, e1.sufMinSum + e2.preMinSum),
        preMinSum: Math.min(e1.preMinSum, e1.sum + e2.preMinSum),
        sufMinSum: Math.min(e2.sufMinSum, e2.sum + e1.sufMinSum)
      })
    )
  }
}

type LongestRepeating<V> = {
  size: number
  max: number
  preMax: number
  sufMax: number
  leftValue?: V
  rightValue?: V
}

/**
 * 单点修改，区间最长相同元素连续个数.
 */
function createPointSetRangeLongestRepeating<V>(arr: ArrayLike<V>): {
  fromElement: (value: V) => LongestRepeating<V>
  tree: SegmentTreePointUpdateRangeQuery<LongestRepeating<V>>
} {
  const leaves: LongestRepeating<V>[] = Array(arr.length)
  const fromElement = (v: V) => ({
    size: 1,
    max: 1,
    preMax: 1,
    sufMax: 1,
    leftValue: v,
    rightValue: v
  })
  for (let i = 0; i < arr.length; ++i) {
    leaves[i] = fromElement(arr[i])
  }

  return {
    fromElement,
    tree: new SegmentTreePointUpdateRangeQuery(
      leaves,
      () => ({ size: 0, max: 0, preMax: 0, sufMax: 0 } as LongestRepeating<V>),
      (a, b) => {
        const res: LongestRepeating<V> = {
          size: a.size + b.size,
          max: 0,
          preMax: 0,
          sufMax: 0,
          leftValue: a.leftValue,
          rightValue: b.rightValue
        }
        if (a.rightValue === b.leftValue) {
          res.preMax = a.preMax
          if (a.preMax === a.size) res.preMax += b.preMax
          res.sufMax = b.sufMax
          if (b.sufMax === b.size) res.sufMax += a.sufMax
          res.max = Math.max(a.max, b.max, a.sufMax + b.preMax)
        } else {
          res.preMax = a.preMax
          res.sufMax = b.sufMax
          res.max = Math.max(a.max, b.max)
        }
        return res
      }
    )
  }
}

type LongestOne = {
  size: number
  preOne: number
  sufOne: number
  longestOne: number
  leftValue: 0 | 1
  rightValue: 0 | 1

  /**
   * 区间内所有极长连续1段的贡献和 sum(len_i*(len_i+1)/2)
   */
  pairCount: number
}

/**
 * 单点修改，区间最长连续1.
 */
function createPointSetRangeLongestOne(nOrArr: number | ArrayLike<0 | 1>): {
  fromElement: (bit: 0 | 1) => LongestOne
  tree: SegmentTreePointUpdateRangeQuery<LongestOne>
} {
  if (typeof nOrArr === 'number') nOrArr = new Uint8Array(nOrArr) as ArrayLike<0 | 1>
  const leaves: LongestOne[] = Array(nOrArr.length)
  const fromElement = (v: 0 | 1) => {
    if (v === 1) {
      return {
        size: 1,
        preOne: 1,
        sufOne: 1,
        longestOne: 1,
        leftValue: 1,
        rightValue: 1,
        pairCount: 1
      } as const
    }
    return {
      size: 1,
      preOne: 0,
      sufOne: 0,
      longestOne: 0,
      leftValue: 0,
      rightValue: 0,
      pairCount: 0
    } as const
  }
  for (let i = 0; i < nOrArr.length; ++i) {
    leaves[i] = fromElement(nOrArr[i])
  }

  return {
    fromElement,
    tree: new SegmentTreePointUpdateRangeQuery(
      leaves,
      () => ({
        size: 0,
        preOne: 0,
        sufOne: 0,
        longestOne: 0,
        leftValue: 0,
        rightValue: 0,
        pairCount: 0
      }),
      (a, b) => {
        const res: LongestOne = {
          size: a.size + b.size,
          preOne: 0,
          sufOne: 0,
          longestOne: 0,
          leftValue: a.leftValue,
          rightValue: b.rightValue,
          pairCount: 0
        }

        if (a.rightValue === b.leftValue) {
          res.preOne = a.preOne
          if (a.preOne === a.size) res.preOne += b.preOne
          res.sufOne = b.sufOne
          if (b.sufOne === b.size) res.sufOne += a.sufOne
          res.longestOne = Math.max(a.longestOne, b.longestOne, a.sufOne + b.preOne)

          // TODO
          const n1 = a.sufOne
          const n2 = b.preOne
          const n3 = n1 + n2
          res.pairCount =
            a.pairCount +
            b.pairCount +
            (n3 * (n3 + 1)) / 2 -
            (n1 * (n1 + 1)) / 2 -
            (n2 * (n2 + 1)) / 2
        } else {
          res.preOne = a.preOne
          res.sufOne = b.sufOne
          res.longestOne = Math.max(a.longestOne, b.longestOne)

          // TODO
          res.pairCount = a.pairCount + b.pairCount
        }
        return res
      }
    )
  }
}

export {
  //
  createRangeAddRangeMax,
  createRangeAddRangeMin,
  //
  createRangeUpdateRangeMax,
  createRangeUpdateRangeMin,
  //
  createRangeAssignRangeMax,
  createRangeAssignRangeMin,
  createRangeAssignRangeSum,
  //
  createRangeFlipRangeSum,
  //
  createRangeAssignRangeAddRangeSum,
  //
  createRangeAffineRangeSum,
  //
  createPointSetRangeMaxSumMinSum,
  //
  createPointSetRangeLongestRepeating,
  //
  createPointSetRangeLongestOne
}
