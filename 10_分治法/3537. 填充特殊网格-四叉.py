# 3537. 填充特殊网格
# https://leetcode.cn/problems/fill-a-special-grid/description/
# 给你一个非负整数 N，表示一个 2N x 2N 的网格。你需要用从 0 到 22N - 1 的整数填充网格，使其成为一个 特殊 网格。一个网格当且仅当满足以下 所有 条件时，才能称之为 特殊 网格：
#
# 右上角象限中的所有数字都小于右下角象限中的所有数字。
# 右下角象限中的所有数字都小于左下角象限中的所有数字。
# 左下角象限中的所有数字都小于左上角象限中的所有数字。
# 每个象限也都是一个特殊网格。
# 返回一个 2N x 2N 的特殊网格。

from typing import List


class Solution:
    def specialGrid(self, n: int) -> List[List[int]]:
        res = [[0] * (1 << n) for _ in range(1 << n)]
        count = 0

        def dfs(u: int, d: int, l: int, r: int) -> None:
            if d - u == 1:
                nonlocal count
                res[u][l] = count
                count += 1
                return
            mid = (d - u) // 2
            dfs(u, u + mid, l + mid, r)  # 右上角
            dfs(u + mid, d, l + mid, r)  # 右下角
            dfs(u + mid, d, l, l + mid)  # 左下角
            dfs(u, u + mid, l, l + mid)  # 左上角

        dfs(0, 1 << n, 0, 1 << n)
        return res
