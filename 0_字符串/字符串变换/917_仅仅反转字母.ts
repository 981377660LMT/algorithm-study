import { isalpha } from '../../utils/string'

function reverseOnlyLetters(s: string): string {
  const sb = s.split('')
  let [left, right] = [0, s.length - 1]

  while (left < right) {
    while (left < right && !isalpha(sb[left])) left++
    while (left < right && !isalpha(sb[right])) right--
    if (left < right) {
      ;[sb[left], sb[right]] = [sb[right], sb[left]]
      left++
      right--
    }
  }

  return sb.join('')
}
