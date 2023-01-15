/* eslint-disable prefer-destructuring */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
// https://leetcode.cn/circle/discuss/1LfDNx/

// 根据上一级目录的编码，整理数据，整理成父子级目录。
// 分享一下对前端这类树形结构转换的心得吧：

// 1.用 interface 清晰定义出每种数据类型，一定要严格定义；
// 2.根据父子关系建立一个有向图，用 Map 保存，键为id，值为原数据结构的数组;
//   有时候父子关系不明确，需要写一个 getParent 函数找到每个数据的父亲;
// 3.找到所有的根节点，存到一个数组里，然后对这些根节点逐个dfs，子节点 push 到根节点的 children 中。
// !下面这段代码可以传入 isRoot/getParent/mapData 函数来实现任意数据结构的转换

/**
 * code != parentCode (no loop)
 * root:parentCode == null
 */
interface Data {
  name: string
  code: string
  parentCode: string | null
}

interface Node {
  name: string
  code: string
  children?: Node[]
}

function transform(dataArray: Data[]): Node[] {
  const graph = new Map<string, Data[]>()
  for (const data of dataArray) {
    const parentCode = data.parentCode
    if (parentCode === null) continue
    if (!graph.has(parentCode)) graph.set(parentCode, [])
    graph.get(parentCode)!.push(data)
  }

  const res = dataArray
    .filter(data => data.parentCode === null)
    .map(data => ({ name: data.name, code: data.code }))
  res.forEach(dfs)
  return res

  function dfs(cur: Node): void {
    for (const next of graph.get(cur.code) || []) {
      const node = { name: next.name, code: next.code }
      if (!cur.children) cur.children = []
      cur.children.push(node)
      dfs(node)
    }
  }
}

console.dir(
  transform([
    { code: 'A', name: 'A', parentCode: null },
    { code: 'A1', name: 'A1', parentCode: 'A' },
    { code: 'A12', name: 'A12', parentCode: 'A1' },
    { code: 'B', name: 'B', parentCode: null },
    { code: 'B2', name: 'B2', parentCode: 'B' },
    { code: 'B22', name: 'B22', parentCode: 'B2' }
  ]),
  { depth: null }
)

// [
//   {
//     name: 'A',
//     code: 'A',
//     children: [
//       {
//         name: 'A1',
//         code: 'A1',
//         children: [ { name: 'A12', code: 'A12' } ]
//       }
//     ]
//   },
//   {
//     name: 'B',
//     code: 'B',
//     children: [
//       {
//         name: 'B2',
//         code: 'B2',
//         children: [ { name: 'B22', code: 'B22' } ]
//       }
//     ]
//   }
// ]
export {}
