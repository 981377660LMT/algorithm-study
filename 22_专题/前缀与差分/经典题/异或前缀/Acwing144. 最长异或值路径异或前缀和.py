# 144. 最长异或值路径
# 给定一个树，树上的边都具有权值。
# 树中一条路径的异或长度被定义为路径上所有边的权值的异或和：
# 给定上述的具有 n 个节点的树，你能找到异或长度最大的路径吗？


# 树的任意两个点的异或路径和,找LCA
# 建立树的异或路径sum,先求左右节点left/right的公共父节点root，
# （求公共父节点也有原题，用两个变量存是否存在子节点dfs即可）。

# preXor表示从根节点到当前节点的异或和
# 答案就是(preXor[left] ^ preXor[lca]) ^ (preXor[right] ^ preXor[lca]) ^ lca->val

# 最大异或路径 - 异或前缀和

# https://www.acwing.com/problem/content/description/146/
