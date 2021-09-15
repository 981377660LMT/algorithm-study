// 你能否用 O(n) 时间复杂度完成此题?
// 记录状态
// 第一个数肯定符合
// 改造成switch更加具有可读性
// const wiggleMaxLength = (nums: number[]): number => {
//   let res = 1
//   let isUp: boolean | undefined

//   for (let i = 0; i < nums.length - 1; i++) {
//     if (nums[i] === nums[i + 1]) {
//       continue
//     } else if (nums[i] < nums[i + 1]) {
//       if (isUp === undefined) {
//         isUp = false
//         res++
//       }
//       if (isUp) {
//         isUp = !isUp
//         res++
//       }
//     } else {
//       if (isUp === undefined) {
//         isUp = true
//         res++
//       }
//       if (!isUp) {
//         isUp = !isUp
//         res++
//       }
//     }
//   }

//   return res
// }
/**
 *
 * @param nums
 * @returns
 * 给你一个整数数组 nums ，返回 nums 中作为 摆动序列 的 最长子序列的长度
 * @summary
 *  贪心:只有真正的V型（波峰或者波谷位置）翻转才会增加摆动序列的长度。
 */
const wiggleMaxLength = (nums: number[]): number => {
  let res = 1
  let state: 'initial' | 'needUp' | 'needDown' = 'initial'

  for (let i = 0; i < nums.length - 1; i++) {
    switch (state) {
      case 'initial':
        if (nums[i] < nums[i + 1]) {
          state = 'needDown'
          res++
        } else if (nums[i] > nums[i + 1]) {
          state = 'needUp'
          res++
        }
        break
      case 'needUp':
        if (nums[i] < nums[i + 1]) {
          state = 'needDown'
          res++
        }
        break
      case 'needDown':
        if (nums[i] > nums[i + 1]) {
          state = 'needUp'
          res++
        }
        break
      default:
        break
    }
  }

  return res
}

console.log(wiggleMaxLength([1, 7, 4, 9, 2, 5]))
console.log(wiggleMaxLength([1, 2, 3]))

export {}
