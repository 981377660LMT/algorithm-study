# 建立树的异或路径sum,先求左右节点left/right的公共父节点root，
# （求公共父节点也有原题，用两个变量存是否存在子节点dfs即可）。
# 答案就是(sum[left] ^ sum[root]) ^ (sum[right] ^ sum[root]) ^ root->val
# 时间复杂度O(n^2)

# 树上差分
# 启示：任意两个点一般找LCA
