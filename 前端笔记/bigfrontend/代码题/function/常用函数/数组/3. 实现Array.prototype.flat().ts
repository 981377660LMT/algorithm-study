type Func = (arr: Array<any>, depth?: number) => Array<any>

const flat: Func = function (arr, depth = 1) {
  // your imeplementation here
  const res: any[] = []

  for (const item of arr) {
    if (Array.isArray(item) && depth) res.push(...flat(item, depth - 1))
    else res.push(item)
  }

  return res
}

const flat2: Func = function (arr, depth = 1) {
  // your imeplementation here
  const res: any[] = []
  const stack: any[] = arr.map(item => [item, depth])
  console.log(stack)
  while (stack.length) {
    const [item, depth] = stack.shift()!
    if (Array.isArray(item) && depth) stack.push(...item.map(v => [v, depth - 1]))
    else res.push(item)
  }

  return res
}

export {}

const arr = [1, [2], [3, [4]]]
console.log(flat2(arr))
console.log(flat2(arr, 1))
console.log(flat2(arr, 2))
// 追问

// 能否不用递归而用迭代的方式实现？
