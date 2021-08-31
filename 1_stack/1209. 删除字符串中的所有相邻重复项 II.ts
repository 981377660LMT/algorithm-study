/**
 * @param {string} s
 * @param {number} k
 * @return {string}
 */
// var removeDuplicates = function (
//   s: string,
//   k: number,
//   reg = new RegExp('([a-z])\\1{' + String(k - 1) + '}', 'g')
// ): string {
//   const removed = s.replace(reg, '')
//   return s.length === removed.length ? s : removeDuplicates(removed, k, reg)
// }
// 堆内存溢出
// <--- Last few GCs --->
// [41:0x4d9d750]     3184 ms: Scavenge 108.0 (146.2) -> 107.9 (146.7) MB, 0.7 / 0.0 ms  (average mu = 0.971, current mu = 0.975) allocation failure
// [41:0x4d9d750]     3187 ms: Scavenge 108.4 (146.7) -> 108.4 (146.7) MB, 0.5 / 0.0 ms  (average mu = 0.971, current mu = 0.975) allocation failure
// [41:0x4d9d750]     3188 ms: Scavenge 108.4 (146.7) -> 108.4 (146.7) MB, 0.5 / 0.0 ms  (average mu = 0.971, current mu = 0.975) allocation failure
// <--- JS stacktrace --->
// FATAL ERROR: MarkCompactCollector: young object promotion failed Allocation failed - JavaScript heap out of memory
var removeDuplicates = function (s: string, k: number): string {
  const stack: [string, number][] = []
  for (const char of s) {
    if (!stack.length || stack[stack.length - 1][0] !== char) stack.push([char, 1])
    else if (stack[stack.length - 1][1] + 1 < k) stack[stack.length - 1][1]++
    else stack.pop()
  }

  return stack.map(item => item[0].repeat(item[1])).join('')
}

console.log(removeDuplicates('deeedbbcccbdaa', 3))

export {}
