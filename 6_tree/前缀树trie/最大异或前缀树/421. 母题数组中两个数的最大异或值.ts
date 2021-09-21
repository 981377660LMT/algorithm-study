/**
 * @param {number[]} nums  0 <= nums[i] <= 2**31 - 1
 * @return {number}
 * 你可以在 O(n) 的时间解决这个问题吗？
 * @summary
 * 每次candidate取贪心 即越大的异或尽量在高位取1
 */
const findMaximumXOR1 = function (nums: number[]): number {
  let res = 0

  for (let i = 31; ~i; i--) {
    // 当前最大值candidate
    // 判断是否存在两个数异或为candidate
    res <<= 1
    const candidate = res + 1

    const prefix = new Set<number>()
    for (const num of nums) {
      prefix.add(num >> i)
    }

    for (const pre of prefix) {
      if (prefix.has(pre ^ candidate)) {
        res = candidate
        break
      }
    }
  }

  return res
}

// class TrieNode {
//   children: Map<number, TrieNode>
//   constructor() {
//     this.children = new Map()
//   }
// }

// // 每个数的二进制从高位到低位依次插入字典树中，最上层0为根节点
// class Trie {
//   private root: TrieNode

//   constructor() {
//     this.root = new TrieNode()
//   }

//   insert(num: number) {
//     let root = this.root
//     for (let i = 31; ~i; i--) {
//       const bit = (num >> i) & 1
//       !root.children.has(bit) && root.children.set(bit, new TrieNode())
//       root = root.children.get(bit)!
//     }
//   }

//   search(num: number) {
//     console.dir(this.root, { depth: null })
//     let root = this.root
//     let res = 0
//     for (let i = 31; ~i; i--) {
//       const bit = (num >> i) & 1
//       const needBit = 1 - bit
//       if (root.children.has(needBit)) {
//         res = (res << 1) + 1
//         root = root.children.get(needBit)!
//       } else if (root.children.has(bit)) {
//         res = res << 1
//         root = root.children.get(bit)!
//       }
//     }

//     return res
//   }
// }

class Trie {
  private static N = 32 * 10000
  private static trieNode = Array.from<any, [number, number]>({ length: Trie.N }, () => [0, 0])
  private static index = 0

  constructor() {
    for (let i = 0; i < Trie.index; i++) {
      Trie.trieNode[i] = [0, 0]
    }
    Trie.index = 0
  }

  insert(num: number) {
    let root = 0
    for (let i = 31; ~i; i--) {
      const bit = (num >> i) & 1
      if (!Trie.trieNode[root][bit]) {
        Trie.index++
        Trie.trieNode[root][bit] = Trie.index
      }
      root = Trie.trieNode[root][bit]
    }
  }

  search(num: number) {
    let root = 0
    let res = 0

    for (let i = 31; ~i; i--) {
      if (!root) break
      const bit = (num >> i) & 1
      const needBit = 1 - bit
      if (Trie.trieNode[root][needBit]) {
        res = (res << 1) + 1
        root = Trie.trieNode[root][needBit]
      } else {
        res = res << 1
        root = Trie.trieNode[root][bit]
      }
    }

    return res
  }

  static main() {
    const trie = new Trie()
    trie.insert(1)
    trie.insert(4)
    console.log(Trie.trieNode)
  }
}

// 前缀树解法
// https://leetcode-cn.com/problems/maximum-xor-of-two-numbers-in-an-array/solution/python3-ha-xi-biao-he-er-wei-shu-zu-shi-7nh3b/
const findMaximumXOR = function (nums: number[]): number {
  let res = 0
  const trie = new Trie()

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
