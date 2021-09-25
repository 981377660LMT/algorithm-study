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
  const charToIndex = new Map<string, number[]>()
  for (let i = 0; i < ring.length; i++) {
    !charToIndex.has(ring[i]) && charToIndex.set(ring[i], [])
    charToIndex.get(ring[i])!.push(i)
  }

  const visited = new Set<string>('0#0')
  const queue = new MinHeap<[Distance, RingIndex, KeyIndex]>((a, b) => a[0] - b[0])
  queue.push([0, 0, 0])

  while (queue.size) {
    const [dis, ringIndex, keyIndex] = queue.shift()!

    if (keyIndex === key.length) return dis + key.length // 按下需要1步

    const k = `${ringIndex}#${keyIndex}`
    if (visited.has(k)) continue
    visited.add(k)

    for (const next of charToIndex.get(key[keyIndex]) || []) {
      queue.push([dis + getDistance(ringIndex, next), next, keyIndex + 1])
    }
  }

  return -1

  function getDistance(i: number, j: number, mod = ring.length) {
    return Math.min((i - j + mod) % mod, (j - i + mod) % mod)
  }
}

console.log(findRotateSteps('godding', 'gd'))

export {}
