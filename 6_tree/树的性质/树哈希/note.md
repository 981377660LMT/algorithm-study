## 有根树的同构

1. 子树结构相同(**子树顺序无影响**) => 树的括号表示法(最小表示法)
   这种做法将子树表示为 **root.value + '(' + subtreeHash +')'**
   注意树的深度不能太大，否则遇到单链会 TLE

   ```TS

     function dfs(root: FileTreeNode): string {
       const subTree: string[] = []
       for (const child of root.children.values()) {
         subTree.push(dfs(child))
       }

       subTree.sort()  // 关键，排除位置的影响
       root.subtreeHash = subTree.join('')

       const res = `${root.value}(${root.subtreeHash})`
       return res
     }
   ```

2. 如果是二叉树，那么结构相同时对应位置也要相同(**子树顺序有影响**) => n 元组
   - 子树 subtreeHash 表示成 n 元组(可以用'#'分隔)时，遇到单链会 TLE (`n*n的字符串拼接`)
   - 优化是使用唯一的哈希 id 来代替哈希值，减少长度

## 无根树的同构

重心+树哈希
`O(m*n^2)`

1. 考虑找树的重心，如果它是唯一的，可以以它为根求最小表示法/树哈希
   如果有两个重心，那么就都求出来，然后取那个较小的
