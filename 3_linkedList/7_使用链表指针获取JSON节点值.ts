const json = {
  a: { b: { c: 1 } },
  d: { e: 2 },
}

const path = ['a', 'b', 'c']

// 与遍历链表异曲同工
let p: { [k: string]: any } = json
path.forEach(key => {
  p = p[key]
})

console.log(p)
export {}
