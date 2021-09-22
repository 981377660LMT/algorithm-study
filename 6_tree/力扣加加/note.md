1. DFS 细分为前中后序遍历， BFS 细分为带层的和不带层的。

```JS
前序遍历
function dfs(root) {
    if (满足特定条件）{
        // 返回结果 or 退出搜索空间
    }
     // 主要逻辑
    dfs(root.left)
    dfs(root.right)
}

后序遍历
function dfs(root) {
    if (满足特定条件）{
        // 返回结果 or 退出搜索空间
    }
    dfs(root.left)
    dfs(root.right)
     // 主要逻辑
}
```

2. BFS
   **不带层信息**

```JS
const visited = {}
function bfs() {
    let q = new Queue()
    q.push(初始状态)
    while(q.length) {
        let i = q.pop()
        if (visited[i]) continue
        if (i 是我们要找的目标) return 结果
        for (i的可抵达状态j) {
            if (j 合法) {
                q.push(j)
            }
        }
    }
    return 没找到
}
```

**带层信息**
while 循环控制一层一层往下走，for 循环利用长度变量控制从左到右遍历每一层二叉树节点。

```TS
const main = (root: BinaryTree | null) => {
  if (!root) return
  const queue: [BinaryTree, number][] = [[root, 0]]
  let depth=1

  while (queue.length) {
    const levelLength = queue.length
    // 遍历当前层的所有节点
    for (let i = 0; i < levelLength; i++) {
      const [head, level] = queue.shift()!
      queue.left&&queue.push([head.left, level + 1])
      queue.right&&queue.push([head.right, level + 1])
    }
    depth++;
  }

  return res
}
```

3. 树的题目就三种类型，分别是：搜索类，构建类和修改类，而这三类题型的比例也是逐渐降低的

- 搜索类

```python
# 其中 path 是树的路径， 如果需要就带上，不需要就不带
def dfs(root, path):
    # 空节点
    if not root: return
    # 叶子节点
    if not root.left and not root.right: return
    path.append(root)
    # 逻辑可以写这里，此时是前序遍历
    dfs(root.left)
    dfs(root.right)
    # 需要弹出，不然会错误计算。
    # 比如对于如下树：
    """
              5
             / \
            4   8
           /   / \
          11  13  4
         /  \    / \
        7    2  5   1
    """
    # 如果不 pop，那么 5 -> 4 -> 11 -> 2 这条路径会变成 5 -> 4 -> 11 -> 7 -> 2，其 7 被错误地添加到了 path

    path.pop()
    # 逻辑也可以写这里，此时是后序遍历

    return 你想返回的数据
```

表示状态 几个状态就模几 例如 3 个状态 0 1 2 每次+1 再模 3

- 构建类
  普通二叉树的构建和二叉搜索树的构建

  1. 给你一个 BFS 的遍历的结果数组，让你构建出原始的树结构。
  2. 如果是二叉搜索树，那么就有可能根据一种遍历序列构造出来

1008. 前序遍历构造二叉搜索树

- 修改类

  116.  填充每个节点的下一个右侧节点指针，

4.  **二叉搜索树**的中序遍历的结果是一个有序数组。如果碰到二叉搜索树的搜索类题目，一定先想下能不能利用这个性质来做。
    给**完全二叉树**编号，这样父子之间就可以通过编号轻松求出。已知一个节点的编号是 i，那么其左子节点就是 2 i，右子节点就是 2 1 + 1，父节点就是 (i + 1) / 2。
    将空节点当成普通节点
    **距离**与**路径**

5.  七个技巧
    - 我们在写 dfs 函数的时候，要将函数中表示当前节点的形参也写成 root 即 dfs(root)
    - 双递归(如果题目有类似，任意节点开始 xxxx 或者所有 xxx 这样的说法，就可以考虑使用双递归)
      一个主递归函数和一个内部递归函数。主递归函数负责计算以某一个节点开始的 xxxx，内部递归函数负责计算 xxxx，这样就实现了以所有节点开始的 xxxx。

```Python
def dfs_inner(root):
    # 这里写你的逻辑，就是前序遍历
    dfs_inner(root.left)
    dfs_inner(root.right)
    # 或者在这里写你的逻辑，那就是后序遍历
def dfs_main(root):
    return dfs_inner(root) + dfs_main(root.left) + dfs_main(root.right)
```

- 前后遍历
  自顶向下:递归调用函数时将这些值传递到子节点，一般是通过参数传到子树中(例如已知根节点性质需要递推到子树，dfs 前序遍历比较常规)
  自底向上:首先对所有子节点递归地调用函数，然后根据返回值和根节点本身的值得到答案(例如**子树的和**,**二叉树的坡度**)
  自底向上通常用后序遍历
  如果你能使用参数和节点本身的值来决定什么应该是传递给它子节点的参数，那就用前序遍历。
  如果对于树中的任意一个节点，如果你知道它子节点的答案，你能计算出当前节点的答案，那就用后序遍历。
  如果遇到二叉搜索树则考虑中序遍历
- 虚拟节点(根节点被修改,新建一个虚拟节点当做新的根节点)

```JS
返回移除了所有不包含 1 的子树的原二叉树
var pruneTree = function (root) {
  function dfs(root) {
    if (!root) return 0;
    const r = dfs(root.right);
    if (l == 0) root.left = null;
    if (r == 0) root.right = null;
    return root.val + l + r;
  }
  ans = new TreeNode(-1);
    const l = dfs(root.left);
  ans.left = root;
  dfs(ans);
  return ans.left;
};
```

- 计算子树和

  ```JS
  function dfs(root) {
  if (!root) return 0;
  const l = dfs(root.left);
  const r = dfs(root.right);
  return root.val + l + r;
  }
  ```

- 计算子树高度

```Python
def dfs(node):
  if not node: return 0
  l = dfs(node.left)
  r = dfs(node.right)
  return max(l, r) + 1
```

- 边界
  空节点
  叶子节点

- dfs 携带更多的参数
  pre 节点/路径和/路径
  参数扩展经常用于 dfs 前序遍历
- 返回元组/列表
  例如层序遍历存储 level 与 root

树的题目一种中心点就是**遍历**，这是搜索问题和修改问题的基础。
而遍历从大的方向分为广度优先遍历和深度优先遍历，这就是我们的**两个基本点**。两个基本点可以进一步细分，比如广度优先遍历有带层信息的和不带层信息的（其实只要会带层信息的就够了）。深度优先遍历常见的是前序和后序，中序多用于二叉搜索树，因为二叉搜索树的中序遍历是严格递增的数组。
树的题目从大的方向上来看就**三种**，一种是搜索类，这类题目最多，这种题目牢牢把握**开始点，结束点 和 目标**即可。构建类型的题目我之前的专题以及讲过了，一句话概括就是根据一种遍历结果确定根节点位置，**根据另外一种遍历结果（如果是二叉搜索树就不需要了）确定左右子树**。修改类题目不多，这种问题边界需要特殊考虑，这是和搜索问题的本质区别，可以使用虚拟节点技巧。另外搜索问题，如果返回值不是根节点也可以考虑虚拟节点。
树有四个比较重要的对做题帮助很大的概念，分别是完全二叉树，二叉搜索树，路径和距离，这里面相关的题目推荐大家好好做一下，都很经典。
七种干货技巧，很多技巧都说明了在什么情况下可以使用

游程编码和 Huffman 都是**无损压缩算法**，即解压缩过程不会损失原数据任何内容
实际情况，我们先用游程编码一遍，然后再用 Huffman 再次编码一次。几乎所有的无损压缩格式都用到了它们，比如 PNG，GIF，PDF，ZIP 等。
Huffman 编码
用短的编码表示出现频率高的字符，用长的编码来表示出现频率低的字符
run-length encode(游程编码)
将重复且连续出现多次的字符使用（连续出现次数，某个字符）来描述。
比如一个字符串：
使用游程编码可以将其描述为：5A4B3C

平衡二叉搜索树的查找和有序数组的二分查找本质都是一样的，只是数据的存储方式不同罢了。那为什么有了有序数组二分，还需要二叉搜索树呢？原因在于树的结构对于**动态数据**比较友好，比如数据是频繁变动的，比如经常添加和删除，那么就可以使用二叉搜索树。

常用判断边界条件
**检查两个 root 都不为 空**

```JS
const helper = (root1: BinaryTree | null, root2: BinaryTree | null): BinaryTree | null => {
    if (!root1) return root2
    if (!root2) return root1
    ...
  }
```

这等价于写法

```JS
if(root1&&root1){


return ...
}
return root1||root1
```

但是前者更加优雅

满二叉树的性质：
最小值是 2 ** (level - 1)，最大值是 2 ** level - 1，其中 level 是树的深度。
假如父节点的索引为 i，那么左子节点就是 2*i， 右边子节点就是 2*i + 1。
假如子节点的索引是 i，那么父节点的索引就是 i // 2。
先思考一般情况（不是之字形）， 然后通过观察找出规律
