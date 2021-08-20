const reverseString = function (s: string[]) {
  let l = -1,
    r = s.length
  while (++l < --r) [s[l], s[r]] = [s[r], s[l]]
  return s
}

export default 1
