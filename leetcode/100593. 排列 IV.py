# 100593. 排列 IV
# https://leetcode.cn/contest/biweekly-contest-151/problems/permutations-iv/
# 给你两个整数 n 和 k，一个 交替排列 是前 n 个正整数的排列，且任意相邻 两个 元素不都为奇数或都为偶数。
# 返回第 k 个 交替排列 ，并按 字典序 排序。如果有效的 交替排列 少于 k 个，则返回一个空列表。
# 1 <= n <= 100
# 1 <= k <= 1e15

from typing import List
from functools import lru_cache


class Solution:
    def permute(self, n: int, k: int) -> List[int]:
        @lru_cache(None)
        def dfs(remainOdd: int, remainEven: int, parity: int) -> int:
            if remainOdd == 0 and remainEven == 0:
                return 1
            res = 0
            if parity == -1:
                if remainOdd > 0:
                    res += remainOdd * dfs(remainOdd - 1, remainEven, 0)
                if remainEven > 0:
                    res += remainEven * dfs(remainOdd, remainEven - 1, 1)
            elif parity == 1:
                if remainOdd > 0:
                    res += remainOdd * dfs(remainOdd - 1, remainEven, 0)
            else:
                if remainEven > 0:
                    res += remainEven * dfs(remainOdd, remainEven - 1, 1)
            return res if res <= k else k + 1

        odds = (n + 1) // 2
        evens = n // 2
        if k > dfs(odds, evens, -1):
            return []

        choose = list(range(1, n + 1))
        res = []
        parity = -1

        for _ in range(n):
            ok = False
            for num in choose:
                if parity != -1 and num & 1 != parity:
                    continue

                if num & 1 == 1:
                    if odds <= 0:
                        continue
                    count = dfs(odds - 1, evens, 0)
                else:
                    if evens <= 0:
                        continue
                    count = dfs(odds, evens - 1, 1)

                if count >= k:
                    res.append(num)
                    choose.remove(num)
                    if num & 1 == 1:
                        odds -= 1
                        parity = 0
                    else:
                        evens -= 1
                        parity = 1
                    ok = True
                    break
                else:
                    k -= count

            if not ok:
                return []

        return res


print(Solution().permute(30, int(1e15)))
