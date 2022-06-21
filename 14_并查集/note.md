关键词：无向图、分组的`传递性`

<!-- 解决连接问题于路径问题 -->

<!-- 两点之间是否可以通过路径连接起来？ -->
<!-- 哪些点属于同一个集合？ -->

<!-- 网络中node的连接状态 -->

“是否存在一条从入口到出口的路径”，那么这就是一个简单的联通问题
并查集还可以在人工智能中用作图像人脸识别。比如将同一个人的不同角度，不同表情的面部数据进行联通。这样就可以很容易地回答两张图片是否是同一个人，无论拍摄角度和面部表情如何。
并查集只能回答“联通与否”，而不能回答诸如“具体的联通路径是什么”

并查集（Union-find Algorithm）定义了两个用于此数据结构的操作：
Find：确定元素属于哪一个子集。它可以被用来确定两个元素是否属于同一子集。
Union：将两个子集合并成同一个集合。

```Python
检测无向图中是否有环
uf = UF()
for a, b in edges:
    if uf.connected(a, b): return False
    uf.union(a, b)
return True
```

主要 API connected 和 union 中的复杂度都是 find 函数造成的，所以说它们的**复杂度和 find 一样**。
find 主要功能就是从某个节点向上遍历到树根，其时间复杂度就是树的高度
问题的关键在于，如何想办法避免树的不平衡

优化:

1. 平衡性优化:rank 优化(小一些的树接到大一些的树下面，这样就能避免头重脚轻，更平衡一些)
   加一个权重数组：通过比较树的重量，就可以保证树的生长相对平衡而不会退化成链表，树的高度大致在 logN 这个数量级
   简单的做法是

   ```JS

    const union = (key1: number, key2: number) => {
        const root1 = find(key1)
        const root2 = find(key2)
        // rank优化:总是让大的根指向小的根
        parent[Math.max(root1, root2)] = Math.min(root1, root2)
    }
   ```

   **实际上应该维护一个 rank 数组**

2. 路径压缩(**非常关键!!!**)：进一步压缩每棵树的高度，使树高始终保持为常数,使得 union 和 connected API 时间复杂度为 O(1)。
   ```JS
   const find = (val: number) => {
   while (parent[val] !== val) {
     parent[val] = parent[parent[val]] // 这一句进行路径压缩
     val = parent[val]
   }
   return val
   }
   ```
   优化后高度不超过 3

进行路径压缩后 rank 可以不要(此时 rank 对优化作用很小)
