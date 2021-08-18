const helper = (s: string, l: number, r: number) => {
  while (l >= 0 && r < s.length && s[l] === s[r]) {
    l--
    r++
  }
  return s.slice(l + 1, r)
}

console.log(helper('abbba', 2, 2))
console.log(helper('abccbd', 2, 3))

export { helper }
