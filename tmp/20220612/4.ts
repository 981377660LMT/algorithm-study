function distributeCookies(cookies: number[], k: number): number {
  const n = cookies.length
  let res = Infinity
  cookies.sort((a, b) => -(a - b))
  dfs(0, 0, new Uint32Array(k))
  return res

  function dfs(index: number, curMax: number, groups: Uint32Array): void {
    if (curMax > res) return
    if (index === n) {
      res = Math.min(res, curMax)
      return
    }

    for (let i = 0; i < k; i++) {
      if (i > 0 && groups[i] === groups[i - 1]) continue
      groups[i] += cookies[index]
      dfs(index + 1, Math.max(curMax, groups[i]), groups)
      groups[i] -= cookies[index]
    }
  }
}

export {}
