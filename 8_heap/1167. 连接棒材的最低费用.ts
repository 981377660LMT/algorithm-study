import { MinHeap } from './MinHeap'

// 哈夫曼编码的思想
function connectSticks(sticks: number[]): number {
  if (sticks.length === 1) return 0
  if (sticks.length === 2) return sticks[0] + sticks[1]
  const queue = new MinHeap((a, b) => a - b, Infinity, sticks)
  queue.heapify()

  let res = 0
  while (queue.size >= 2) {
    const [a, b] = [queue.heappop()!, queue.heappop()!]
    const sum = a + b
    res += sum
    queue.heappush(sum)
  }
  return res
}

console.log(connectSticks([1, 8, 3, 5]))
