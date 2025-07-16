export {}

// 3606. 优惠券校验器
//
// https://leetcode.cn/problems/coupon-code-validator/description/
//
// 给你三个长度为 n 的数组，分别描述 n 个优惠券的属性：code、businessLine 和 isActive。其中，第 i 个优惠券具有以下属性：
//
// code[i]：一个 字符串，表示优惠券的标识符。
// businessLine[i]：一个 字符串，表示优惠券所属的业务类别。
// isActive[i]：一个 布尔值，表示优惠券是否当前有效。
// 当以下所有条件都满足时，优惠券被认为是 有效的 ：
//
// code[i] 不能为空，并且仅由字母数字字符（a-z、A-Z、0-9）和下划线（_）组成。
// businessLine[i] 必须是以下四个类别之一："electronics"、"grocery"、"pharmacy"、"restaurant"。
// isActive[i] 为 true 。
// 返回所有 有效优惠券的标识符 组成的数组，按照以下规则排序：
//
// 先按照其 businessLine 的顺序排序："electronics"、"grocery"、"pharmacy"、"restaurant"。
// 在每个类别内，再按照 标识符的字典序（升序）排序。
//
// 技巧：
// - fail fast (尽早失败) 的思路可以提高性能，避免不必要的计算。
// - 谓词下推.

const CODE_REGEX = /^[a-zA-Z0-9_]+$/
const BUSINESS_LINES = ['electronics', 'grocery', 'pharmacy', 'restaurant']

const compareString = (a: string, b: string): 0 | -1 | 1 => {
  if (a === b) return 0
  return a < b ? -1 : 1
}

function validateCoupons(code: string[], businessLine: string[], isActive: boolean[]): string[] {
  const n = code.length
  const indexes: number[] = Array(n)
  for (let i = 0; i < n; i++) indexes[i] = i

  const check = (i: number): boolean => {
    if (!isActive[i]) return false
    if (code[i].length === 0) return false
    if (!BUSINESS_LINES.includes(businessLine[i])) return false
    if (!CODE_REGEX.test(code[i])) return false
    return true
  }

  const compare = (i: number, j: number): number => {
    return compareString(businessLine[i], businessLine[j]) || compareString(code[i], code[j])
  }

  const transform = (i: number): string => {
    return code[i]
  }

  return indexes
    .filter(i => check(i))
    .sort((i, j) => compare(i, j))
    .map(i => transform(i))
}

export {}
