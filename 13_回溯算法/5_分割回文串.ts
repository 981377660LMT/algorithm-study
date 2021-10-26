// 1 <= s.length <= 16
const partition = (s: string): string[][] => {
  const res: string[][] = []
  const isPalindrome = (str: string) => str === str.split('').reverse().join('')

  const bt = (start: number, path: string[]) => {
    if (start === s.length) {
      res.push(path.slice())
      return
    }

    for (let i = start; i < s.length; i++) {
      const sub = s.slice(start, i + 1)
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
