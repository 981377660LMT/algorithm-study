// # 请你返回可在 一次交换（交换两数字 A[i] 和 A[j] 的位置）后得到的、按字典序排列小于 A 的最大可能排列。
// # The same technic to solve 31. Next Permutation

// 注意这里只交换一次
function prevPermOpt1(arr: number[]): number[] {
  const n = arr.length

  loop: for (let left = n - 1; left > -1; left--) {
    for (let right = n - 1; right > left; right--) {
      // 注意这里 相等的要放到前面  3113 => 交换0和1位置 1313
      if (right >= 1 && arr[right] == arr[right - 1]) continue
      if (arr[right] < arr[left]) {
        ;[arr[right], arr[left]] = [arr[left], arr[right]]
        break loop
      }
    }
  }

  return arr
}

// print(Solution().prevPermOpt1(arr=[3, 2, 1]))
// # 输出：[3,1,2]
// # 解释：交换 2 和 1
// print(Solution().prevPermOpt1([1, 1, 5]))
// # 输出：[1,1,5]
// # 解释：已经是最小排列
// print(Solution().prevPermOpt1([1, 9, 4, 6, 7]))
// # 输出：[1,7,4,6,9]
// # 解释：交换 9 和 7
console.log(prevPermOpt1([1, 9, 4, 6, 7]))
console.log(prevPermOpt1([3, 1, 1, 3]))
// 注意答案是1313
