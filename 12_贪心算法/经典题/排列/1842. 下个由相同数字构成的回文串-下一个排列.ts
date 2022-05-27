import { nextPermutation } from './api/nextPermutation'

function nextPalindrome(num: string): string {
  const n = num.length
  const preHalf = num.split('').slice(0, n >> 1)
  const [nextPerm, ok] = nextPermutation(preHalf)

  if (!ok) return ''

  let prefix = nextPerm.join('')
  //  长度为奇数的情况，要单独加中间的字符
  if ((num.length & 1) === 1) prefix += num[n >> 1]
  const suffix = nextPerm.reverse().join('')
  return prefix + suffix
}

if (require.main === module) {
  console.log(nextPalindrome('1221'))
  console.log(nextPalindrome('32123'))
  console.log(nextPalindrome('23143034132'))
}
export {}
