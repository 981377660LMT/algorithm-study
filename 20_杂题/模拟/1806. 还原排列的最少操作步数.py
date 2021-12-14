# 给你一个偶数 n​​​​​​ ，已知存在一个长度为 n 的排列 perm ，其中 perm[i] == i​（下标 从 0 开始 计数）。
# 一步操作中，你将创建一个新数组 arr ，对于每个 i ：

# 如果 i % 2 == 0 ，那么 arr[i] = perm[i / 2]
# 如果 i % 2 == 1 ，那么 arr[i] = perm[n / 2 + (i - 1) / 2]

# 要想使 perm 回到排列初始值，至少需要执行多少步操作？返回最小的 非零 操作步数。

# 元素从1出发，最终需要回到1：等同于把1换回原来的位置
class Solution:
    def reinitializePermutation(self, n: int) -> int:
        i = n / 2
        res = 1
        while i != 1:
            if i % 2 == 0:
                i = i / 2
            else:
                i = n / 2 + (i - 1) / 2
            res += 1
        return res


print(Solution().reinitializePermutation(n=4))
# 输出：2
# 解释：最初，perm = [0,1,2,3]
# 第 1 步操作后，perm = [0,2,1,3]
# 第 2 步操作后，perm = [0,1,2,3]
# 所以，仅需执行 2 步操作
