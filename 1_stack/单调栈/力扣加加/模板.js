const monotoneStack = function (arr) {
  // 哨兵保证所有元素出栈
  arr.unshift(0)
  arr.push(0)
  // stack一般存数组下标i
  const stack = []
  let res = 0

  for (let i = 0; i < arr.length; i++) {
    while (stack.length && arr[stack[stack.length - 1]] > arr[i]) {
      const tmp = stack.pop()
      // 逻辑...
      // res = ...
      res = Math.max(res, (i - stack[stack.length - 1] - 1) * heights[tmp])
    }
    stack.push(i)
  }

  return res
}
