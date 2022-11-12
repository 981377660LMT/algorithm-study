// 从 1 加到 n 的长度之和
const preSum = [0]
for (let i = 1; i <= 20000; i++) {
  preSum.push(preSum[i - 1] + String(i).length)
}

function splitMessage(message: string, limit: number): string[] {
  const n = message.length
  for (let group = 1; group <= n; group++) {
    const res = check(group)
    if (res.length) return res
  }
  return []

  function check(group: number) {
    const L = String(group).length + 3
    const maxSufLen = String(group - 1).length + L
    if (maxSufLen > limit) return []
    const len1 = L * (group - 1)
    const len2 = preSum[group - 1]
    const remainWord = n - (limit * (group - 1) - len1 - len2)
    if (remainWord < 0 || remainWord > limit) return []
    if (remainWord + L + String(group).length > limit) return []

    // 合理的
    const res = []
    let preI = 0
    for (let i = 0; i < group - 1; i++) {
      const sufLen = String(i + 1).length + L
      const wordLen = limit - sufLen
      res.push(`${message.slice(preI, preI + wordLen)}<${i + 1}/${group}>`)
      preI += wordLen
    }

    // 最后一组
    res.push(`${message.slice(preI)}<${group}/${group}>`)
    return res
  }
}

console.log(splitMessage('this is really a very awesome message', 9))

export {}
