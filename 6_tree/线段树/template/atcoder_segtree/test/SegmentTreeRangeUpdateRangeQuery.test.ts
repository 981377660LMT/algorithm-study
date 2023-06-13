import { SegmentTreeRangeUpdateRangeQuery } from '../SegmentTreeRangeUpdateRangeQuery'

const INF = 2e15

describe('SegmentTreeRangeUpdateRangeQuery', () => {
  // !叠加更新 区间最大值查询
  describe('MaxSegmentTree', () => {
    let tree: SegmentTreeRangeUpdateRangeQuery<number, number>
    beforeEach(() => {
      tree = new SegmentTreeRangeUpdateRangeQuery([1, 2, 3, 4, 5], {
        e: () => -INF,
        id: () => 0,
        op: (a, b) => Math.max(a, b),
        mapping: (data, lazy) => data + lazy,
        composition: (lazy1, lazy2) => lazy1 + lazy2
      })
    })

    it('should support query', () => {
      expect(tree.query(0, 5)).toBe(5)
      expect(tree.query(1, 3)).toBe(3)
      expect(tree.query(0, 0)).toBe(-INF)
    })

    it('should support update', () => {
      tree.update(2, 3, 10)
      expect(tree.query(0, 5)).toBe(13)
      tree.update(4, 5, 10)
      expect(tree.query(4, 5)).toBe(15)
      tree.update(3, 5, -10)
      expect(tree.query(0, 5)).toBe(13)
    })

    it('should support queryAll', () => {
      expect(tree.queryAll()).toBe(5)
    })

    it('should support maxRight', () => {
      expect(tree.maxRight(0, v => v < 3)).toBe(2)
      expect(tree.maxRight(0, v => v <= 4)).toBe(4)
    })

    it('should support minLeft', () => {
      expect(tree.minLeft(5, v => v > 3)).toBe(0)
    })

    it('should throw error when range is invalid', () => {
      try {
        tree.query(3, 2)
      } catch (error) {
        expect(error instanceof RangeError).toBeTruthy()
      }
    })
  })

  // !叠加更新 区间和查询
  describe('SumSegmentTree', () => {
    let tree: SegmentTreeRangeUpdateRangeQuery<[sum: number, length: number], number>
    beforeEach(() => {
      tree = new SegmentTreeRangeUpdateRangeQuery(
        Array.from({ length: 5 }, (_, i) => [i, 1]),
        {
          e: () => [0, 0],
          id: () => 0,
          op(data1, data2) {
            return [data1[0] + data2[0], data1[1] + data2[1]]
          },
          mapping(parentLazy, childData) {
            return [childData[0] + parentLazy * childData[1], childData[1]]
          },
          composition(parentLazy, childLazy) {
            return parentLazy + childLazy
          }
        }
      )
    })

    it('should support query', () => {
      expect(tree.query(0, 5)[0]).toBe(10)
      expect(tree.query(1, 3)[0]).toBe(3)
      expect(tree.query(0, 0)[0]).toBe(0)
    })

    it('should support update', () => {
      tree.update(2, 3, 10)
      expect(tree.query(0, 5)).toStrictEqual([20, 5])
      tree.update(3, 5, -10)
      expect(tree.query(0, 5)).toStrictEqual([0, 5])
    })

    it('should support queryAll', () => {
      expect(tree.queryAll()).toStrictEqual([10, 5])
    })
  })

  // 如题，已知一个数列，你需要进行下面三种操作：
  // 将某区间每一个数乘上 xx
  // 将某区间每一个数加上 xx
  // 求出某区间每一个数的和(奇妙序列)
  describe('https://www.luogu.com.cn/problem/P3373', () => {
    it('should pass test case', () => {
      const nums = [1, 5, 4, 2, 3]
      const MOD = 38n
      const queries = [
        [2, 1, 4, 1],
        [3, 2, 5],
        [1, 2, 4, 2],
        [2, 3, 5, 5],
        [3, 1, 4]
      ]

      type Data = [sum: bigint, length: bigint]
      type Lazy = [mul: bigint, add: bigint]

      const tree = new SegmentTreeRangeUpdateRangeQuery(
        nums.map<Data>(value => [BigInt(value), 1n]),
        {
          e: () => [0n, 1n],
          id: () => [1n, 0n],
          op(data1, data2) {
            return [(data1[0] + data2[0]) % MOD, data1[1] + data2[1]]
          },
          // 区间和等于原来的区间和乘以mul加上区间的长度乘以add
          mapping(parentLazy, childData) {
            return [
              (childData[0] * parentLazy[0] + BigInt(childData[1]) * parentLazy[1]) % MOD,
              childData[1]
            ]
          },
          composition(parentLazy, childLazy) {
            return [
              (parentLazy[0] * childLazy[0]) % MOD,
              (parentLazy[0] * childLazy[1] + parentLazy[1]) % MOD
            ]
          },
          equalsId(id1, id2) {
            return id1[0] === id2[0] && id1[1] === id2[1]
          }
        }
      )

      const expected = [17n, 2n]
      let ei = 0

      for (const [type, ...args] of queries) {
        if (type === 1) {
          const [left, right, mul] = args
          tree.update(left - 1, right, [BigInt(mul), 0n])
        } else if (type === 2) {
          const [left, right, add] = args
          tree.update(left - 1, right, [1n, BigInt(add)])
        } else {
          const [left, right] = args
          expect(tree.query(left - 1, right)[0]).toBe(expected[ei])
          ei++
        }
      }
    })
  })

  // 01串反转(flip)求区间1的个数
  // 若 op=0，则表示将01串的 [l, r] 区间内的 0 变成 1，1 变成 0。
  // 若 op=1，则表示询问01串的[l, r] 区间内有多少个字符 1。
  describe('https://www.luogu.com.cn/problem/P2574', () => {
    it('should pass test case', () => {
      const s = '1011101001'

      const queries = [
        [0, 2, 4],
        [1, 1, 5],
        [0, 3, 7],
        [1, 1, 10],
        [0, 1, 4],
        [1, 2, 6]
      ]

      type Data = [count0: number, count1: number]
      type Lazy = 0 | 1 // 0表示不反转，1表示反转

      const tree = new SegmentTreeRangeUpdateRangeQuery(
        s.split('').map<Data>(v => (v === '0' ? [1, 0] : [0, 1])),
        {
          e: () => [0, 0],
          id: () => 0,
          op(data1, data2) {
            return [data1[0] + data2[0], data1[1] + data2[1]]
          },
          mapping(parentLazy, childData) {
            if (parentLazy === 1) {
              // eslint-disable-next-line no-param-reassign
              ;[childData[0], childData[1]] = [childData[1], childData[0]]
            }
            return childData
          },
          composition(parentLazy, childLazy) {
            return (parentLazy ^ childLazy) as Lazy
          }
        }
      )

      const expected = [3, 6, 1]
      let ei = 0

      for (let [type, left, right] of queries) {
        left--
        if (type === 0) {
          tree.update(left, right, 1)
        } else {
          expect(tree.query(left, right)[1]).toBe(expected[ei])
          ei++
        }
      }
    })
  })

  // 01串反转(flip)求区间逆序对个数
  // 若 op=1，则表示将01串的 [l, r] 区间内的 0 变成 1，1 变成 0。
  // 若 op=2，则表示询问01串的[l, r] 区间内有多少个逆序对。
  describe('https://atcoder.jp/contests/practice2/tasks/practice2_l', () => {
    const nums = [0, 1, 0, 0, 1]
    const queries = [
      [2, 1, 5],
      [1, 3, 4],
      [2, 2, 5],
      [1, 1, 3],
      [2, 1, 2]
    ]

    type Data = [count0: number, count1: number, inv: number]
    type Lazy = 0 | 1 // 0表示不反转，1表示反转

    const tree = new SegmentTreeRangeUpdateRangeQuery<Data, Lazy>(
      nums.map<Data>(v => (v === 0 ? [1, 0, 0] : [0, 1, 0])),
      {
        e() {
          return [0, 0, 0]
        },
        id() {
          return 0
        },
        op(data1, data2) {
          return [
            data1[0] + data2[0],
            data1[1] + data2[1],
            data1[2] + data2[2] + data1[1] * data2[0]
          ]
        },
        mapping(parentLazy, childData) {
          if (parentLazy === 1) {
            // !4000ms => 2600ms 不创建新数组节省空间、时间
            // eslint-disable-next-line no-param-reassign
            ;[childData[0], childData[1], childData[2]] = [
              childData[1],
              childData[0],
              childData[0] * childData[1] - childData[2]
            ]
          }
          return childData
        },
        composition(parentLazy, childLazy) {
          return (parentLazy ^ childLazy) as Lazy
        }
      }
    )

    const expected = [2, 0, 1]
    let ei = 0
    for (let [type, left, right] of queries) {
      left--
      if (type === 1) {
        tree.update(left, right, 1)
      } else {
        expect(tree.query(left, right)[2]).toBe(expected[ei])
        ei++
      }
    }
  })
})
