import { useArrayXORTrie } from './XORTrieArray'

/**
 * @param {number[]} nums  0 <= nums[i] <= 2**31 - 1
 * @return {number}
 * 你可以在 O(n) 的时间解决这个问题吗？
 * @summary
 * 每次candidate取贪心 即越大的异或尽量在高位取1
 */
function findMaximumXOR(nums: number[]): number {
  let res = 0
  const trie = useArrayXORTrie(2 ** 31 - 1)

  for (const num of nums) {
    trie.insert(num)
    res = Math.max(res, trie.search(num))
  }

  return res
}

console.log(findMaximumXOR([3, 10, 5, 25, 2, 8]))
console.log(findMaximumXOR([0]))

console.log(25 ^ 8)
console.log(25 ^ 5)

export {}
