// 比较npm版本号
// !加速：转换成一个数字比较，自定义映射规则；优点是更快，缺点是可能超出float64精度
//
// {major}.{minor}.{patch}[-{prerelease}][+{build}]
//
// 1.0.0 < 1.5.2 < 1.5.2-alpha.1 < 1.5.2-alpha.2 < 1.5.2-beta.1 < 1.5.2-beta

function compareSemver(v1: string, v2: string): number {
  v1 = _extractVersion(v1)
  v2 = _extractVersion(v2)
  const [main1, preRelease1] = _splitMainAndPreRelease(v1)
  const [main2, preRelease2] = _splitMainAndPreRelease(v2)
  const mainCompareResult = _compareMain(main1, main2)
  if (mainCompareResult < 0) return -1
  if (mainCompareResult > 0) return 1
  const preReleaseCompareResult = _comparePreRelease(preRelease1, preRelease2)
  if (preReleaseCompareResult < 0) return -1
  if (preReleaseCompareResult > 0) return 1
  return 0
}

function _extractVersion(v: string): string {
  return v.split('+')[0].trim()
}

function _splitMainAndPreRelease(v: string): [string, string] {
  const [main, preRelease] = v.split('-')
  return [main, preRelease || '']
}

function _compareMain(s1: string, s2: string): number {
  const parts1 = s1.split('.')
  const parts2 = s2.split('.')
  const len = Math.max(parts1.length, parts2.length)
  for (let i = 0; i < len; i++) {
    if (i === parts1.length) return -1
    if (i === parts2.length) return 1
    const n1 = Number(parts1[i])
    const n2 = Number(parts2[i])
    if (n1 !== n2) {
      return n1 > n2 ? 1 : -1
    }
  }
  return 0
}

function _comparePreRelease(s1: string, s2: string): number {
  if (!s1 && !s2) return 0
  if (!s1) return 1 // 无预发布版本的版本号更高
  if (!s2) return -1

  const parts1 = s1.split('.')
  const parts2 = s2.split('.')
  const len = Math.max(parts1.length, parts2.length)
  for (let i = 0; i < len; i++) {
    if (i === parts1.length) return -1
    if (i === parts2.length) return 1

    const a = parts1[i]
    const b = parts2[i]
    const aIsNum = /^\d+$/.test(a)
    const bIsNum = /^\d+$/.test(b)

    let compareResult: number
    if (aIsNum && bIsNum) {
      compareResult = Number(a) - Number(b)
    } else if (aIsNum && !bIsNum) {
      compareResult = -1 // 数字标识符优先级更低
    } else if (!aIsNum && bIsNum) {
      compareResult = 1
    } else {
      compareResult = a.localeCompare(b)
    }

    if (compareResult !== 0) {
      return compareResult > 0 ? 1 : -1
    }
  }

  return 0
}

// 基础比较
console.log(compareSemver('1.0.0', '1.0.0')) // 0（相等）
console.log(compareSemver('1.2.3', '1.3.0')) // -1（v1 < v2）
console.log(compareSemver('2.0.0', '1.9.9')) // 1（v1 > v2）

// 预发布版本
console.log(compareSemver('1.0.0-alpha', '1.0.0')) // -1（预发布版本更小）
console.log(compareSemver('1.0.0-beta', '1.0.0-alpha')) // 1（beta > alpha）
console.log(compareSemver('1.0.0-alpha.1', '1.0.0-alpha')) // 1（标识符更多）

// 特殊场景
console.log(compareSemver('1.0.0-1', '1.0.0-alpha')) // -1（数字标识符 < 非数字）
console.log(compareSemver('1.0.0-a', '1.0.0-1')) // 1（非数字 > 数字）
console.log(compareSemver('1.0.0-rc.1', '1.0.0-beta.11')) // 1（rc > beta）

export {}
