const partition = (s: string): string[][] => {
  const res: string[][] = []
  const isPalindrome = (str: string) => str === str.split('').reverse().join('')

  const bt = (start: string, cur: string[], curLen: number) => {
    console.log(curLen, cur)
    if (curLen === s.length) {
      // 注意这里要浅拷贝,不然引用pop之后会没有值
      res.push([...cur])
      return
    }

    for (let i = curLen + 1; i <= start.length; i++) {
      const sub = start.slice(curLen, i)
      if (isPalindrome(sub)) {
        cur.push(sub)
        bt(start, cur, i)
        // 这里需要pop
        cur.pop()
      }
    }
  }
  bt(s, [], 0)

  return res
}

console.dir(partition('aab'), { depth: null })
// [["a","a","b"],["aa","b"]]
export {}
