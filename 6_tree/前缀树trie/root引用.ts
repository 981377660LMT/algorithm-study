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
