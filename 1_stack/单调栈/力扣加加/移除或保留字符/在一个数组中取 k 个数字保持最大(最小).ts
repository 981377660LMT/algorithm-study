const pickMaxKFromArray = (arr: number[], k: number) => {
  let remain = arr.length - k
  const stack: number[] = []
  for (let i = 0; i < arr.length; i++) {
    while (remain && stack.length && arr[stack[stack.length - 1]] < arr[i]) {
      stack.pop()
      remain--
    }
    stack.push(i)
  }
  return stack.map(index => arr[index]).join('')
}

console.log(pickMaxKFromArray([2, 3, 1, 4, 5, 2, 9, 1], 3))
