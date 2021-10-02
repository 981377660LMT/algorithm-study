function nextGreatestLetter(letters: string[], target: string): string {
  let l = 0
  let r = letters.length - 1

  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = letters[mid]
    if (midElement === target) l++
    else if (midElement < target) l = mid + 1
    else if (midElement > target) r = mid - 1
  }

  return letters[l] || letters[0]
}

console.log(nextGreatestLetter(['c', 'f', 'j'], 'a'))
console.log(nextGreatestLetter(['c', 'f', 'j'], 'k'))
// 输出: "c"
// 在比较时，字母是依序循环出现的。举个例子：
// 如果目标字母 target = 'z' 并且字符列表为 letters = ['a', 'b']，则答案返回 'a'
