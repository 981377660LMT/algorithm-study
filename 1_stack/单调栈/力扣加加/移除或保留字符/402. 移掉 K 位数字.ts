/**
 * @param {string} num
 * @param {number} k
 * @return {string}
 * @description 两个相同位数的数字大小关系取决于第一个不同的数的大小。
 * 每次丢弃一次，k 减去 1。当 k 减到 0 ，我们可以提前终止遍历。
   而当遍历完成，如果 k 仍然大于 0。不妨假设最终还剩下 x 个需要丢弃，那么我们需要选择删除末尾 x 个元素。
 */
const removeKdigits = function (num: string, k: number): string {
  const n = num.length
  // 需要保留的长度
  const remain = n - k
  const stack: string[] = []
  for (let i = 0; i < n; i++) {
    while (k && stack.length && Number(stack[stack.length - 1]) > Number(num[i])) {
      stack.pop()!
      k--
    }
    stack.push(num[i])
  }

  // 去除前面的0： 注意parseInt去除数字字符串前面的0时 会出现Infinity的情况
  // 左边消除0需要用replace(/^0+/, '')
  return (stack.slice(0, remain).join('') || 0).toString().replace(/^0+/, '') || '0'
}

console.log(removeKdigits('1432219', 3))
// 输出："1219"
// 解释：移除掉三个数字 4, 3, 和 2 形成一个新的最小的数字 1219 。
console.log(removeKdigits('10200', 1))
console.log(removeKdigits('10', 2))
console.log(removeKdigits('9', 1))
