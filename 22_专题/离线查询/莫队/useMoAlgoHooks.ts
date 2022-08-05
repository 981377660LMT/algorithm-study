import assert from 'assert'
import { useMoAlgo } from './useMoAlgo'

/**
 * 在正整数数组中查询区间的mex
 */
const useQueryMex = useMoAlgo({
  mex: 1,
  counter: new Map<number, number>(),
  add(value: number) {
    this.counter.set(value, (this.counter.get(value) || 0) + 1)
    while ((this.counter.get(this.mex) || 0) > 0) this.mex++
  },
  remove(value: number) {
    this.counter.set(value, (this.counter.get(value) || 0) - 1)
    if ((this.counter.get(this.mex) || 0) === 0) this.mex = Math.min(this.mex, value)
  },
  query(): number {
    return this.mex
  }
})

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5]
  const Q = [
    [0, 0],
    [1, 1],
    [2, 2],
    [3, 3],
    [4, 4]
  ]

  const queryMex = useQueryMex(nums)
  Q.forEach(([left, right]) => queryMex.addQuery(left, right))
  assert.deepStrictEqual(queryMex.work(), [2, 1, 1, 1, 1])
}

export { useQueryMex }
