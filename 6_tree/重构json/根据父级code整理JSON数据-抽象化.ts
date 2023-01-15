/* eslint-disable no-param-reassign */
/* eslint-disable eqeqeq */
/* eslint-disable prefer-destructuring */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// https://leetcode.cn/circle/discuss/1LfDNx/

// 根据上一级目录的编码，整理数据，整理成父子级目录。
// 分享一下对前端这类树形结构转换的心得吧：

// 1.用 interface 清晰定义出每种数据类型，一定要严格定义；
// 2.根据父子关系建立一个有向图，用 Map 保存，键为id，值为原数据结构的数组;
//   有时候父子关系不明确，需要写一个 getParent 函数找到每个数据的父亲;
// 3.找到所有的根节点，存到一个数组里，然后对这些根节点逐个dfs，子节点 push 到根节点的 children 中。
// !下面这段代码可以传入 getId/getParent/mapData 函数来实现任意数据结构的转换

interface Operation<From, To, Id> {
  getId: (from: From) => Id
  getParentId: (from: From) => Id | null | undefined
  mapFrom: (from: From) => To
}

function transform<
  From extends object,
  To extends object & { children?: To[] | null | undefined },
  Id extends PropertyKey = PropertyKey
>(from: Iterable<From>, operation: Operation<From, To, Id>): To[] {
  const { getId, getParentId, mapFrom } = operation
  const ROOT_SYMBOL = Symbol('ROOT')
  const ID_SYMBOL = Symbol('ID')

  const tree = new Map<PropertyKey, Array<From & { [ID_SYMBOL]: Id }>>()
  for (const data of from) {
    let parentId: ReturnType<typeof getParentId> | typeof ROOT_SYMBOL = getParentId(data)
    if (parentId == undefined) parentId = ROOT_SYMBOL
    if (!tree.has(parentId)) tree.set(parentId, [])
    tree.get(parentId)!.push({ ...data, [ID_SYMBOL]: getId(data) })
  }

  const res = { children: [] }
  dfs(res, ROOT_SYMBOL)
  return res.children

  function dfs(cur: { children?: To[] | null | undefined }, curId: Id | typeof ROOT_SYMBOL): void {
    for (const next of tree.get(curId) || []) {
      const node = mapFrom(next)
      if (!cur.children) cur.children = []
      cur.children.push(node)
      dfs(node, next[ID_SYMBOL])
    }
  }
}

if (require.main === module) {
  // usage
  const res = transform(
    [
      { code: 'A', name: 'A', parentCode: null },
      { code: 'A1', name: 'A1', parentCode: 'A' },
      { code: 'A12', name: 'A12', parentCode: 'A1' },
      { code: 'B', name: 'B', parentCode: null },
      { code: 'B2', name: 'B2', parentCode: 'B' },
      { code: 'B22', name: 'B22', parentCode: 'B2' }
    ],
    {
      getId: data => data.code,
      getParentId: data => data.parentCode,
      mapFrom: data => ({ name: data.name, code: data.code })
    }
  )

  console.dir(res, { depth: null })
}

export {}
