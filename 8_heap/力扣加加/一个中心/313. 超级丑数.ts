import { PriorityQueue } from '../../../2_queue/todo优先级队列'

/**
 * @param {number} n
 * @param {number[]} primes
 * @return {number}
 * @description
 * 超级丑数 是一个正整数，并满足其所有质因数都出现在质数数组 primes 中。
   给你一个整数 n 和一个整数数组 primes ，返回第 n 个 超级丑数 。
   每次入队丑数与各个素因子的乘积，记得去重
   出队，去重，入队，有点像bfs
 */
const nthSuperUglyNumber = function (n: number, primes: number[]): number {
  let count = 0
  let res = 1
  const pq = new PriorityQueue<number>()
  pq.push(1)
  while (count < n) {
    const head = pq.shift()
    // 去重
    if (head === res && pq.length) continue
    res = head
    count++
    primes.forEach(prime => pq.push(prime * res))
  }

  return res
}

console.log(nthSuperUglyNumber(12, [2, 7, 13, 19]))
