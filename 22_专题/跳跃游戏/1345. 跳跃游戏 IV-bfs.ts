import { ArrayDeque } from '../../2_queue/Deque'

/**
 * @param {number[]} arr
 * @return {number}
 * 请你返回到达数组最后一个元素的下标处所需的 最少操作次数 。
 * 每一步，你可以从下标 i 跳到下标：
   i + 1 满足：i + 1 < arr.length
   i - 1 满足：i - 1 >= 0
   j 满足：arr[i] == arr[j] 且 i != j
   即：前跳 后跳 跳到另一个相等处
 */
const minJumps = function (arr: number[]): number {
  if (arr.length === 1) return 0
  const adjMap = new Map<number, number[]>()
  arr.forEach((val, index) => {
    !adjMap.has(val) && adjMap.set(val, [])
    adjMap.get(val)!.push(index)
  })

  const visited = new Set()
  const queue = new ArrayDeque(10000)
  queue.push(0)
  let steps = 0

  while (queue.length) {
    const len = queue.length

    for (let i = 0; i < len; i++) {
      const curIndex = queue.shift()!
      if (curIndex === arr.length - 1) return steps
      visited.add(curIndex)
      for (const next of [...adjMap.get(arr[curIndex])!, curIndex - 1, curIndex + 1]) {
        if (next >= 0 && next < arr.length && !visited.has(next)) {
          queue.push(next)
          // 这个高度的，已经visit了，防止后面同高度的点浪费时间，再去判断是否visit
          // 之前的已经不用去看了 类似于2_单词接龙 删除已经看过的词 之后才接近答案
          adjMap.set(arr[curIndex], [])
        }
      }
    }

    steps++
  }

  return -1
}

console.log(minJumps([100, -23, -23, 404, 100, 23, 23, 23, 3, 404]))
// 输出：3
// 解释：那你需要跳跃 3 次，下标依次为 0 --> 4 --> 3 --> 9 。下标 9 为数组的最后一个元素的下标。

// 超时
