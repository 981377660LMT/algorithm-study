var monostoneStack = function (T) {
  let stack = []
  let result = []
  for (let i = 0; i < T.length; i++) {
    result[i] = 0
    while (stack.length > 0 && T[stack[stack.length - 1]] < T[i]) {
      let peek = stack.pop()
      result[peek] = i - peek
    }
    stack.push(i)
  }
  return result
}
