// 单调栈要求栈中的元素是单调递减或者单调递减的。
// 这个限制有什么用呢？这个限制（特性）能够解决什么用的问题呢？
// 单调栈适合的题目是求解下一个大于 xxx或者下一个小于 xxx这种题目。
// 以及维护好数组的一个区间内的「最大值」「次大值」等等。
// https://leetcode-cn.com/problems/next-greater-element-i/solution/dong-hua-yan-shi-dan-diao-zhan-496xia-yi-ql65/

// 例子:找出下一个比当前序号大的元素
// 单调栈:栈底元素必须单调(递增)；小的数进来只push,大的数进来一个一个pop,大的数再成为更大的栈底保持单调
const findNextLarge = (nums: number[]): number[] => {
  const stack: number[] = []
  const res = Array<number>(nums.length).fill(-1)

  for (let i = 0; i < nums.length; i++) {
    while (stack.length > 0 && nums[stack[stack.length - 1]] < nums[i]) {
      res[stack.pop()!] = i
    }
    stack.push(i)
  }

  return res
}
console.log(findNextLarge([1, 3, 4, 2, 5]))

export {}
