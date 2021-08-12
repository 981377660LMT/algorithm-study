const partition = (s: string): string[][] => {
  const res: string[][] = []
  const isPalindrome = (str: string) => str === str.split('').reverse().join('')

  const bt = (remain: string, path: string[]) => {
    // console.log(curLen, path)
    if (remain.length === 0) {
      res.push(path.slice())
      return
    }

    for (let i = 0; i < remain.length; i++) {
      const sub = remain.slice(0, i + 1)
      if (isPalindrome(sub)) {
        path.push(sub)
        bt(remain.slice(i + 1), path)
        path.pop()
      }
    }
  }
  bt(s, [])

  return res
}

console.dir(partition('aab'), { depth: null })
// [["a","a","b"],["aa","b"]]
export {}
