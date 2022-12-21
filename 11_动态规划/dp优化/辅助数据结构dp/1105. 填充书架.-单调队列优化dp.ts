/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */

import { MinDeque } from '../../../2_queue/单调队列Monoqueue/MaxDequeue'

// 1 <= books.length <= 1000
// 1 <= thicknessi <= shelfWidth <= 1000
// 1 <= heighti <= 1000

type Item = {
  value: number
  left: number
  right: number
}

function minHeightShelves(books: number[][], shelfWidth: number): number {
  const n = books.length
  const heights = Array<number>(n + 1).fill(0)
  for (let i = 1; i <= n; i++) heights[i] = books[i - 1][1]
  const preSum = Array<number>(n + 1).fill(0)
  const dp = Array<number>(n + 1).fill(0)
  const queue = new MinDeque<Item>()

  for (let i = 1; i <= n; i++) {
    preSum[i] = preSum[i - 1] + books[i - 1][0]
    let left = i
    while (queue.size && heights[queue.at(-1)!.right] <= heights[i]) {
      left = queue.pop()!.left
    }

    queue.append({ left, value: dp[left - 1] + heights[i], right: i })
    while (queue.size && preSum[i] - preSum[queue.at(0)!.left - 1] > shelfWidth) {
      const item = queue.popLeft()!
      if (item.left + 1 <= item.right) {
        queue.appendLeft({
          left: item.left + 1,
          value: dp[item.left] + heights[item.right],
          right: item.right
        })
      }
    }

    dp[i] = queue.min
  }

  return dp[n]
}

if (require.main === module) {
  //  test 1e5
  const books = Array.from({ length: 1e5 }, () => [
    Math.floor(Math.random() * 1e5),
    Math.floor(Math.random() * 1e5)
  ])
  const shelfWidth = Math.floor(Math.random() * 1e5 + 1e5)
  console.time('minHeightShelves')
  console.log(minHeightShelves(books, shelfWidth))
  console.timeEnd('minHeightShelves')
}

export {}
