// Given a string s, return the minimum number of characters needed to be inserted so that the string becomes a palindrome.

class Solution {
  constructor(private memo: Map<string, number> = new Map<string, number>()) {}

  solve(s: string): number {
    const helper = (left: number, right: number): number => {
      if (left >= right) return 0
      const key = `${left}#${right}`
      if (this.memo.has(key)) return this.memo.get(key)!

      if (s[left] === s[right]) {
        const tmp = helper(left + 1, right - 1)
        this.memo.set(key, tmp)
        return tmp
      } else {
        const tmp = 1 + Math.min(helper(left, right - 1), helper(left + 1, right))
        this.memo.set(key, tmp)
        return tmp
      }
    }

    return helper(0, s.length - 1)
  }

  test() {
    console.log(this.memo)
  }
}

const big1 = new Solution()
console.log(big1.solve('radr'))
big1.test()
