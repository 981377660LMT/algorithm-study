# 3669. K 因数分解
# https://leetcode.cn/problems/balanced-k-factor-decomposition/
# 给你两个整数 n 和 k，将数字 n 恰好分割成 k 个正整数，使得这些整数的 乘积 等于 n。
# 返回一个分割方案，使得这些数字中 最大值 和 最小值 之间的 差值 最小化。结果可以以 任意顺序 返回。
# 4 <= n <= 1e5
# 2 <= k <= 5
# k 严格小于 n 的正因数的总数。


from typing import List, Optional


def getAllFactors(n: int) -> List[List[int]]:
    res = [[] for _ in range(n + 1)]
    for f in range(1, n + 1):
        for m in range(f, n + 1, f):
            res[m].append(f)
    return res


INF = int(1e18)


allFactors = getAllFactors(int(1e5 + 10))


class Solution:
    def minDifference(self, n: int, k: int) -> List[int]:
        minDiff = INF
        path = [0] * k
        res: Optional[List[int]] = None

        def dfs(pos: int, cur: int, min_: int, max_: int) -> None:
            if pos == k - 1:
                nonlocal minDiff, res
                d = max(max_, cur) - min(min_, cur)
                if d < minDiff:
                    minDiff = d
                    path[pos] = cur
                    res = path[:]
                return
            for f in allFactors[cur]:
                path[pos] = f
                dfs(pos + 1, cur // f, min(min_, f), max(max_, f))

        dfs(0, n, INF, 0)
        return res  # type: ignore
