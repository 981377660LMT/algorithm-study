/**
 *
 * @param nums digit 数组
 * @description 返回值第二个参数带error 类似go的模式...
 */
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
}

/**
 *
 * @param nums digit 数组
 * @description 返回值第二个参数带error 类似go的模式...
 */
function prePermutation<T>(nums: T[]): [res: T[], ok: boolean] {
  const copy = nums.slice()
  const n = copy.length
  let isExist = false

  loop: for (let left = n - 1; left > -1; left--) {
    for (let right = n - 1; right > left; right--) {
      // 找到了第一对后面小于前面
      if (copy[right] < copy[left]) {
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
}

function reverseRange<T>(nums: T[], i: number, j: number) {
  while (i < j) {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
    i++
    j--
  }
}

export {}
