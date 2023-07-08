// !gcd里递归不加判断 会死循环(b站之前lua脚本 入参变成字符串 就是这样挂掉的)

import { lcm } from './gcd'

// 请你统计并返回 nums 的 子数组 中满足 元素最小公倍数为 k 的子数组数目。
function subarrayLCM(nums: number[], k: number): number {
  let res = 0
  for (let i = 0; i < nums.length; i++) {
    let curLcm = 1
    for (let j = i; j < nums.length; j++) {
      curLcm = lcm(curLcm, nums[j])
      res += +(curLcm === k)
    }
  }

  return res
}

export {}
