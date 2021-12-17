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

  // 返回值第二个参数带error 类似go的模式...
  function nextPermutation<T>(nums: T[]): [res: T[], ok: boolean] {
    const copy = nums.slice()
    const n = copy.length
    let isExist = false

    loop: for (let left = n - 1; left > -1; left--) {
      for (let right = n - 1; right > left; right--) {
        // 找到了第一对后面大于前面
        if (copy[right] > copy[left]) {
          // 交换完排序
          ;[copy[left], copy[right]] = [copy[right], copy[left]]
          reverseRange(copy, left + 1, n - 1)
          isExist = true
          break loop
        }
      }
    }

    if (isExist) return [copy, true]
    else return [[], false]

    function reverseRange<T>(nums: T[], i: number, j: number) {
      while (i < j) {
        ;[nums[i], nums[j]] = [nums[j], nums[i]]
        i++
        j--
      }
    }
  }
}

// console.log(nextPalindrome('1221'))
// console.log(nextPalindrome('32123'))
console.log(nextPalindrome('23143034132'))
export {}
