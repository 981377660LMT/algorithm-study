function minLength(s: string, numOps: number): number {
  const n = s.length

  const runs: number[] = []
  let currentChar = s[0]
  let currentRun = 1

  for (let i = 1; i < n; i++) {
    if (s[i] === currentChar) {
      currentRun++
    } else {
      runs.push(currentRun)
      currentChar = s[i]
      currentRun = 1
    }
  }
  runs.push(currentRun)

  function check1(): number {
    let flip1 = 0
    let flip2 = 0
    for (let i = 0; i < n; i++) {
      const expected1 = i % 2 === 0 ? '0' : '1'
      const expected2 = i % 2 === 0 ? '1' : '0'
      if (s[i] !== expected1) flip1++
      if (s[i] !== expected2) flip2++
    }
    return Math.min(flip1, flip2)
  }

  function check2(L: number): number {
    let flips = 0
    for (let m of runs) {
      if (m > L) {
        flips += Math.floor(m / (L + 1))
      }
    }
    return flips
  }

  let left = 1
  let right = Math.max(...runs)
  let res = right
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    let requiredFlips: number
    if (mid === 1) {
      requiredFlips = check1()
    } else {
      requiredFlips = check2(mid)
    }
    if (requiredFlips <= numOps) {
      res = mid
      right = mid - 1
    } else {
      left = mid + 1
    }
  }

  return res
}
