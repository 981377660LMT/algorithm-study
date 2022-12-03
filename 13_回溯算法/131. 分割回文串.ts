// 1 <= s.length <= 16
const partition = (s: string): string[][] => {
  const res: string[][] = []
  const isPalindrome = (str: string) => str === str.split('').reverse().join('')

  const bt = (index: number, path: string[]) => {
    if (index === s.length) {
      res.push(path.slice())
      return
    }

    for (let i = index; i < s.length; i++) {
      const sub = s.slice(index, i + 1)
      if (isPalindrome(sub)) {
        path.push(sub)
        bt(i + 1, path)
        path.pop()
      }
    }
  }
  bt(0, [])

  return res
}

console.dir(partition('aab'), { depth: null })
// [["a","a","b"],["aa","b"]]
export {}
