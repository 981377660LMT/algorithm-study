import * as jsondiffpatch from 'jsondiffpatch'
// import 'jsondiffpatch/dist/formatters-styles/html.css'
// import 'jsondiffpatch/dist/formatters-styles/annotated.css'

const diffpatcher = jsondiffpatch.create()

{
  // 原始对象
  const left = {
    name: 'John Doe',
    age: 30,
    tasks: ['write code', 'fix bugs'],
    contact: {
      email: 'john.doe@example.com'
    }
  }

  // 修改后的对象
  const right = {
    name: 'Johnathan Doe', // 修改
    age: 31, // 修改
    tasks: ['write code', 'review code', 'fix bugs'], // 数组中添加元素
    // contact 属性被删除
    status: 'active' // 新增属性
  }

  // 生成 delta
  const delta = diffpatcher.diff(left, right)

  // 打印 delta
  console.log(JSON.stringify(delta, null, 2))

  // 使用 delta 来更新 left 对象
  // 注意：patch 会直接修改原始对象。如果不想修改，可以先克隆。
  const patchedLeft = JSON.parse(JSON.stringify(left)) // 创建一个克隆
  diffpatcher.patch(patchedLeft, delta)

  // 验证 patchedLeft 是否和 right 对象深度相等
  console.log(patchedLeft)
  // 输出会和 right 对象的内容完全一样
  /*
{
  name: 'Johnathan Doe',
  age: 31,
  tasks: [ 'write code', 'review code', 'fix bugs' ],
  status: 'active'
}
*/

  // 使用 delta 来从 right 对象恢复到 left 对象的状态
  const unpatchedRight = JSON.parse(JSON.stringify(right)) // 创建一个克隆
  diffpatcher.unpatch(unpatchedRight, delta)

  // 验证 unpatchedRight 是否和 left 对象深度相等
  console.log(unpatchedRight)
  // 输出会和 left 对象的内容完全一样
  /*
{
  name: 'John Doe',
  age: 30,
  tasks: [ 'write code', 'fix bugs' ],
  contact: { email: 'john.doe@example.com' }
}
*/
}

{
  const diffpatcherWithObjectHash = jsondiffpatch.create({
    // @ts-ignore
    objectHash: obj => obj.id
  })

  const listBefore = [
    { id: 'a', value: 1 },
    { id: 'b', value: 2 }
  ]

  const listAfter = [
    { id: 'b', value: 3 }, // value 改变
    { id: 'a', value: 1 } // 顺序改变
  ]

  const delta = diffpatcherWithObjectHash.diff(listBefore, listAfter)

  console.log(JSON.stringify(delta, null, 2))
}
