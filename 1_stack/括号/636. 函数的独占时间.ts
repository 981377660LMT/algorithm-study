type ID = number
/**
 * @param {number} n  有一个 单线程 CPU 正在运行一个含有 n 道函数的程序  1 <= n <= 100
 * @param {string[]} logs
 * @return {number[]}
 * "0:start:3" 意味着标识符为 0 的函数调用在时间戳 3 的 起始开始执行 ；
 * 而 "1:end:2" 意味着标识符为 1 的函数调用在时间戳 2 的 末尾结束执行
 * 函数的 独占时间 定义是在这个函数在程序所有函数调用中执行时间的总和，
 * 调用其他函数花费的时间不算该函数的独占时间
 * 两个开始事件不会在同一时间戳发生
   两个结束事件不会在同一时间戳发生
   @summary 
   执行时间:end_timestamp - start_timestamp + 1
   end必定是栈顶元素的end
 */
const exclusiveTime = function (n: number, logs: string[]): number[] {
  const res = Array<number>(n).fill(0)
  const stack: ID[] = []

  let pre = 0
  for (const log of logs) {
    const details = log.split(':')
    const id = parseInt(details[0])
    const state = details[1]
    const time = parseInt(details[2])
    if (state === 'start') {
      if (stack.length) res[stack[stack.length - 1]] += time - pre
      stack.push(id)
      pre = time
    } else {
      res[stack.pop()!] += time - pre + 1
      // 注意这个时间 时间以左边为准 所以要加1
      pre = time + 1
    }
  }

  return res
}

console.log(exclusiveTime(2, ['0:start:0', '1:start:2', '1:end:5', '0:end:6']))
// 输出：[3,4]
console.log(
  exclusiveTime(1, ['0:start:0', '0:start:2', '0:end:5', '0:start:6', '0:end:6', '0:end:7'])
)
// 输出：[8]
export {}
