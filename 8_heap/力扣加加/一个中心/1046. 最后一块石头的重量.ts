import { PriorityQueue } from '../../../2_queue/todo优先级队列'

/**
 * @param {number[]} stones
 * @return {number}
 * @description 每一回合，从中选出两块 最重的 石头，然后将它们一起粉碎。
 * 最重的石头在模拟过程中是动态变化的。
   这种动态取极值的场景使用堆就非常适合
 */
const lastStoneWeight = function (stones: number[]) {
  const pq = new PriorityQueue<number>((a, b) => b - a)
  stones.forEach(stone => pq.push(stone))

  while (pq.length >= 2) {
    const head1 = pq.shift()!
    const head2 = pq.shift()!
    if (head1 === head2) continue
    pq.push(head1 - head2)
  }

  return pq.length === 0 ? 0 : pq.shift()
}

console.log(lastStoneWeight([2, 7, 4, 1, 8, 1]))

export default 1
