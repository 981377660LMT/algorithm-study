// 如果某个整数是 num 中各位数字的一个 排列 且它的 值大于 num ，则称这个整数为 妙数
// 暴力解法
// 2 <= num.length <= 1000
// 1 <= k <= 1000
function getMinSwaps(num: string, k: number): number {
  const nums = num.split('').map(Number)
  const perm = nextPermutation(nums.slice(), k)

  let res = 0
  for (let i = 0; i < num.length; i++) {
    let cur = i
    // 从后面找第一个相同的位置搬过来
    while (nums[i] !== perm[i]) {
      res++
      cur++
      ;[nums[i], nums[cur]] = [nums[cur], nums[i]]
    }
  }

  return res

  function nextPermutation(nums: number[], k = 1) {
    const n = nums.length

    while (k--) {
      next()
    }

    return nums

    function next() {
      loop: for (let left = n - 1; ~left; left--) {
        for (let right = n - 1; right > left; right--) {
          // 找到了第一对后面大于前面
          if (nums[right] > nums[left]) {
            // 交换玩排序
            ;[nums[left], nums[right]] = [nums[right], nums[left]]

            // 后面都比前面小
            reverseRange(nums, left + 1, n - 1)
            break loop
          }
        }
      }
    }

    function reverseRange(nums: number[], i: number, j: number) {
      while (i < j) {
        ;[nums[i], nums[j]] = [nums[j], nums[i]]
        i++
        j--
      }
    }
  }
}

if (require.main === module) {
  console.log(getMinSwaps('5489355142', 4))
}

export {}
// # 输出：2
// # 解释：第 4 个最小妙数是 "5489355421" ，要想得到这个数字：
// # - 交换下标 7 和下标 8 对应的位："5489355142" -> "5489355412"
// # - 交换下标 8 和下标 9 对应的位："5489355412" -> "5489355421"

// 例如，num = "5489355142" ：
// 第 1 个最小妙数是 "5489355214"
// 第 2 个最小妙数是 "5489355241"
// 第 3 个最小妙数是 "5489355412"
// 第 4 个最小妙数是 "5489355421"
