type Flat<A = any> = (arr: A[], depth?: number) => A[]

const flat1: Flat = function (arr, depth = 1) {
  // your imeplementation here
  const res: any[] = []

  for (const item of arr) {
    if (Array.isArray(item) && depth) res.push(...flat1(item, depth - 1))
    else res.push(item)
  }

  return res
}

const flat2: Flat = function (arr, depth = 1) {
  // your imeplementation here
  const res: any[] = []
  const queue: any[] = arr.map(item => [item, depth])
  console.log(queue)
  while (queue.length) {
    const [item, depth] = queue.shift()!
    if (Array.isArray(item) && depth) queue.push(...item.map(value => [value, depth - 1]))
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
