// 单调栈要求栈中的元素是单调递减或者单调递减的。
// 这个限制有什么用呢？这个限制（特性）能够解决什么用的问题呢？
// 单调栈适合的题目是求解下一个大于 xxx或者下一个小于 xxx这种题目。
// https://leetcode-cn.com/problems/next-greater-element-i/solution/dong-hua-yan-shi-dan-diao-zhan-496xia-yi-ql65/

// 例子:找出下一个比当前序号大的元素
// 单调栈:栈底元素必须单调(递增)；小的数进来只push,大的数进来一个一个pop,大的数再成为更大的栈底保持单调
const findNextLarge = (nums: number[]) => {
  const monoStack: number[] = []
  const memo: Map<number, number> = new Map()

  nums.forEach(num => {
    // 栈不为空且当前元素大于栈顶元素
    // 说明当前元素是栈顶元素的下一个更大元素
    // while循环表示当前元素是栈中所有已存元素的下一个更大元素
    while (monoStack.length > 0 && num > monoStack[monoStack.length - 1]) {
      memo.set(monoStack.pop()!, num)
    }
    monoStack.push(num)
  })

  console.log(memo, monoStack)
}

console.log(findNextLarge([1, 3, 4, 2, 5]))
export {}
