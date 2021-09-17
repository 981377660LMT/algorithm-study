const root = {
  val: 1,
  children: {
    val: 2,
    children: {
      val: 3,
      children: null,
    },
  },
}
let rootP = root as any

rootP = root.children
rootP.val = 666
console.log(root)

// 谨慎操作,小心引用

const arr = Array(3)
arr[1] = 1
console.log(arr)
export default 1
