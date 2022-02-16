import { MinHeap } from '../../../minheap'

/**
 * @param {number} n 1 <= n <= 106
 * @param {number[]} primes
 * @return {number}
 * 超级丑数 是一个正整数，并满足其所有质因数都出现在质数数组 primes 中
 */
var nthSuperUglyNumber = function (n: number, primes: number[]): number {
  const res = [1]
  // val row col
  const pq = new MinHeap<[number, number, number]>((a, b) => a[0] - b[0])
  for (let i = 0; i < primes.length; i++) {
    pq.heappush([primes[i], i, 0])
  }

  while (res.length < n) {
    const [val, row, col] = pq.heappop()!
    if (val !== res[res.length - 1]) res.push(val)
    // pq.push([primes[row] + val, row, col + 1])
    pq.heappush([primes[row] * res[col + 1], row, col + 1])
  }
  // console.log(res)
  return res[res.length - 1]
}

export {}
console.log(nthSuperUglyNumber(12, [2, 7, 13, 19]))
