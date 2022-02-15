import { MinHeap } from '../../../2_queue/minheap'

type Distance = number
type RingIndex = number
type KeyIndex = number

/**
 * @param {string} ring  ring，表示刻在外环上的编码 最初，ring 的第一个字符与12:00方向对齐
 * @param {string} key  key，表示需要拼写的关键词
 * @return {number}  能够拼写关键词中所有字符的最少步数
 * ring 和 key 的字符串长度取值范围均为 1 至 100；
 *
 * 将 ring 顺时针或逆时针旋转一个位置，计为1步
 * 按下中心按钮进行拼写，这也将算作 1 步
 * 旋转的最终目的是将字符串 ring 的一个字符与 12:00 方向对齐，并且这个字符必须等于字符 key[i] 。
 */
const findRotateSteps = function (ring: string, key: string): number {
  const charPos = new Map<string, number[]>()
  for (let i = 0; i < ring.length; i++) {
    !charPos.has(ring[i]) && charPos.set(ring[i], [])
    charPos.get(ring[i])!.push(i)
  }

  const visited = new Set<string>()
  const queue = new MinHeap<[Distance, RingIndex, KeyIndex]>((a, b) => a[0] - b[0])
  queue.push([0, 0, 0])

  while (queue.size) {
    const [dis, cur, steps] = queue.shift()!

    if (steps === key.length) return dis + key.length // 按下需要1步

    const curKey = `${cur}#${steps}`
    if (visited.has(curKey)) continue
    visited.add(curKey)

    // 不能用dist数组是因为权值不是确定的，需要用visited；直接全部入堆
    for (const next of charPos.get(key[steps]) || []) {
      queue.push([dis + getDistance(cur, next), next, steps + 1])
    }
  }

  return -1

  function getDistance(i: number, j: number, mod = ring.length) {
    return Math.min((i - j + mod) % mod, (j - i + mod) % mod)
  }
}

console.log(findRotateSteps('godding', 'gd'))
console.log(findRotateSteps('caotmcaataijjxi', 'oatjiioicitatajtijciocjcaaxaaatmctxamacaamjjx'))

export {}
