// 测试区间加法

import { SqrtDecomposition } from '../SqrtDecomposition'

describe('SqrtDecomposition', () => {
  it('should support range add range sum', () => {
    const n = Math.floor(Math.random() * 1000) + 1
    const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 1000))
    const sqrt = new SqrtDecomposition<number, number>(n, (_, start, end) => {
      const elements = nums.slice(start, end)
      let sum = 0
      let lazy = 0

      return {
        created() {
          this.updated()
        },
        updated() {
          sum = elements.reduce((a, b) => a + b, 0)
        },
        queryPart(start, end) {
          let res = 0
          for (let i = start; i < end; i++) {
            res += elements[i] + lazy
          }
          return res
        },
        updatePart(start, end, add) {
          for (let i = start; i < end; i++) {
            elements[i] += add
          }
        },
        queryAll() {
          return sum + lazy * (end - start)
        },
        updateAll(add) {
          lazy += add
        }
      }
    })

    for (let i = 0; i < 10; i++) {
      let l = Math.floor(Math.random() * n)
      let r = Math.floor(Math.random() * n)
      const add = Math.floor(Math.random() * 1000)
      sqrt.update(l, r, add)
      for (let j = l; j < r; j++) {
        nums[j] += add
      }

      let l2 = Math.floor(Math.random() * n)
      let r2 = Math.floor(Math.random() * n)
      let sum = 0
      sqrt.query(l2, r2, blockRes => {
        sum += blockRes
      })
      let res = 0
      for (let j = l2; j < r2; j++) {
        res += nums[j]
      }
      expect(sum).toBe(res)
    }
  })
})
