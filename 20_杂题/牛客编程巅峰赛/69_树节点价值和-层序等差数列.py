#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# @param n long长整型 表示标准完全二叉树的结点个数
# @return long长整型
#
from math import floor, log2


MOD = 998244353

# n<=10^9


class Solution:
    # 一个一个加肯定超时
    def tree4TLE(self, n: int) -> int:
        # write code here
        self.res = 0

        def dfs(root: int, depth: int):
            self.res += depth * root
            self.res %= MOD

            if root << 1 <= n:
                dfs(root << 1, depth + 1)
            if (root << 1 | 1) <= n:
                dfs(root << 1 | 1, depth + 1)

        dfs(1, 1)
        return self.res % MOD

    # 分层等差数列统计
    def tree4(self, n: int) -> int:
        res = 0
        TREE_DEPTH = floor(log2(n)) + 1
        depth, left, right = 1, 1, 1
        while depth <= TREE_DEPTH:
            res += depth * (left + right) * (right - left + 1) // 2
            res %= MOD
            depth += 1
            left, right = min(n, left << 1), min(n, right << 1 | 1)
        return res % MOD

