/**
 * @param {string[]} chars
 * @return {number}
 * 你必须设计并实现一个只使用常量额外空间的算法来解决此问题。
 * @description
 * 快慢指针
 */
const compress = function (chars: string[]): number {
  let [slow, fast] = [0, 0]

  while (fast < chars.length) {
    chars[slow] = chars[fast]
    let count = 1

    while (fast + 1 < chars.length && chars[fast] === chars[fast + 1]) {
      fast++
      count++
    }

    if (count > 1) {
      for (const char of count.toString()) {
        chars[++slow] = char
      }
    }

    fast++
    slow++
  }

  return slow
}

console.log(compress(['a', 'a', 'b', 'b', 'c', 'c', 'c']))
// 输出：返回 6 ，输入数组的前 6 个字符应该是：["a","2","b","2","c","3"]
export default 1
