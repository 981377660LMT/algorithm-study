/**
 *
 * @param arr
 * @param k
 * @returns
 * @description
 * 在一个数组中取 k 个数字保持最大(不改变相对顺序)
 */
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
  // 注意slice
  return stack
    .slice(0, k)
    .map(index => arr[index])
    .join('')
}

console.log(pickMaxKFromArray([2, 3, 1, 4, 5, 2, 9, 1], 3))
