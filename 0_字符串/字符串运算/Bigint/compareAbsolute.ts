function compareAbsolute(str1: string, str2: string): -1 | 0 | 1 {
  if (str1.length !== str2.length) return str1.length > str2.length ? 1 : -1

  for (let i = 0; i < str1.length; i++) {
    if (str1[i] !== str2[i]) return str1[i] > str2[i] ? 1 : -1
  }

  return 0
}

export { compareAbsolute }
