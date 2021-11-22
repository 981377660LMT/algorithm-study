function hashCode(str: string) {
  let hash = 0

  for (const char of str) {
    // hash = 31 * hash + char.codePointAt(0)!
    // 与乘以31相同
    hash = (hash << 5) - hash + char.codePointAt(0)!
  }

  return hash | 0 // Convert to 32bit integer
}
