// 找出 nums1 中每个元素在 nums2 中的下一个比其大的值。
// stack加map
const nextGreaterElement = (nums1: number[], nums2: number[]) => {
  const stack: number[] = []
  const map = new Map()

  nums2.forEach(num => {
    while (stack.length && num > stack[stack.length - 1]) {
      map.set(stack.pop(), num)
    }
    stack.push(num)
  })

  return nums1.map(num => map.get(num) || -1)
}

console.log(nextGreaterElement([4, 1, 2], [1, 3, 4, 2]))

export {}
