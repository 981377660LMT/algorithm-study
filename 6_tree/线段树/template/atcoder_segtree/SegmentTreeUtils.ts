/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
// !由于lazy模板通用性 效率不如自己维护数组的线段树
// !注意如果是单点查询,可以去掉所有pushUp函数逻辑(js使用bigint会比较慢)
// !如果是单点修改,可以去掉所有懒标记逻辑

// !一些monoid (如果难以设计半群,就使用分块解决吧)
// https://maspypy.github.io/library/alg/acted_monoid/summax_assign.hpp

import { SegmentTreeRangeUpdateRangeQuery } from './SegmentTreeRangeUpdateRangeQuery'

const INF = 2e15

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
  nOrNums: number | ArrayLike<[sum: number, size: number]>
): SegmentTreeRangeUpdateRangeQuery<[sum: number, size: number], number> {
  return new SegmentTreeRangeUpdateRangeQuery<[sum: number, size: number], number>(nOrNums, {
    e: () => [0, 1],
    id: () => -1,
    op: ([sum1, size1], [sum2, size2]) => [sum1 + sum2, size1 + size2],
    mapping: (f, [sum, size]) => (f === -1 ? [sum, size] : [f * size, size]),
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
  nOrNums: number | ArrayLike<[sum: number, size: number]>
): SegmentTreeRangeUpdateRangeQuery<[sum: number, size: number], number> {
  return new SegmentTreeRangeUpdateRangeQuery<[sum: number, size: number], number>(nOrNums, {
    e: () => [0, 1],
    id: () => 0,
    op: ([sum1, size1], [sum2, size2]) => [sum1 + sum2, size1 + size2],
    mapping: (f, [sum, size]) => (f === 0 ? [sum, size] : [size - sum, size]),
    composition: (f, g) => f ^ g
  })
}

/**
 * 区间赋值区间加,区间和.
 */
function createRangeAssignRangeAddRangeSum(
  nOrNums: number | ArrayLike<[size: number, sum: number]>
): SegmentTreeRangeUpdateRangeQuery<[size: number, sum: number], [mul: number, add: number]> {
  return new SegmentTreeRangeUpdateRangeQuery<
    [size: number, sum: number],
    [mul: number, add: number]
  >(nOrNums, {
    e() {
      return [1, 0]
    },
    id() {
      return [1, 0]
    },
    op(e1, e2) {
      return [e1[0] + e2[0], e1[1] + e2[1]]
    },
    mapping(lazy, data) {
      return [data[0], data[1] * lazy[0] + data[0] * lazy[1]]
    },
    composition(f, g) {
      return [f[0] * g[0], f[0] * g[1] + f[1]]
    },
    equalsId(id1, id2) {
      return id1[0] === id2[0] && id1[1] === id2[1]
    }
  })
}

/**
 * 区间仿射变换,区间和.
 */
function createRangeAffineRangeSum(
  nOrNums: number | ArrayLike<[size: bigint, sum: bigint]>,
  bigMod: bigint
): SegmentTreeRangeUpdateRangeQuery<[size: bigint, sum: bigint], [mul: bigint, add: bigint]> {
  return new SegmentTreeRangeUpdateRangeQuery(nOrNums, {
    e() {
      return [1n, 0n]
    },
    id() {
      return [1n, 0n]
    },
    op(e1, e2) {
      return [e1[0] + e2[0], (e1[1] + e2[1]) % bigMod]
    },
    mapping(lazy, data) {
      return [data[0], (data[1] * lazy[0] + data[0] * lazy[1]) % bigMod]
    },
    composition(f, g) {
      return [(f[0] * g[0]) % bigMod, (f[0] * g[1] + f[1]) % bigMod]
    },
    equalsId(id1, id2) {
      return id1[0] === id2[0] && id1[1] === id2[1]
    }
  })
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
  createRangeAffineRangeSum
}
