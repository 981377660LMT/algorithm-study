from typing import List, Optional


class Knapsack01Removable:
    """可撤销01背包,用于求解方案数/可行性."""

    __slots__ = "_dp", "_maxWeight", "_mod"

    def __init__(self, maxWeight: int, mod: Optional[int] = None, dp: Optional[List[int]] = None):
        if dp is not None:
            self._dp = dp
        else:
            self._dp = [0] * (maxWeight + 1)
            self._dp[0] = 1
        self._maxWeight = maxWeight
        self._mod = mod

    def add(self, weight: int) -> None:
        """添加一个重量为weight的物品."""
        dp = self._dp
        if self._mod is None:
            for i in range(self._maxWeight, weight - 1, -1):
                dp[i] += dp[i - weight]
        else:
            mod = self._mod
            for i in range(self._maxWeight, weight - 1, -1):
                dp[i] = (dp[i] + dp[i - weight]) % mod

    def remove(self, weight: int) -> None:
        """移除一个重量为weight的物品.需要保证weight物品存在."""
        dp = self._dp
        if self._mod is None:
            for i in range(weight, self._maxWeight + 1):
                dp[i] -= dp[i - weight]
        else:
            mod = self._mod
            for i in range(weight, self._maxWeight + 1):
                dp[i] = (dp[i] - dp[i - weight]) % mod

    def query(self, weight: int) -> int:
        """
        查询组成重量为weight的物品有多少种方案.
        !注意需要特判重量为0.
        """
        return self._dp[weight] if 0 <= weight <= self._maxWeight else 0

    def copy(self) -> "Knapsack01Removable":
        return Knapsack01Removable(self._maxWeight, self._mod, self._dp[:])

    def __repr__(self) -> str:
        return self._dp.__repr__()


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    # https://atcoder.jp/contests/abc321/tasks/abc321_f
    def solve1():
        MOD = 998244353
        q, maxWeight = map(int, input().split())
        K = Knapsack01Removable(maxWeight, MOD)
        for _ in range(q):
            op, w = input().split()
            w = int(w)
            if op == "+":
                K.add(w)
            else:
                K.remove(w)
            print(K.query(maxWeight))

    # P4141 消失之物
    # https://www.luogu.com.cn/problem/P4141
    # 对每个物品i，在不选择i的情况下输出容量为1-m的方案数.
    # n,m<=2000
    def p4141() -> None:
        def solve(weights: List[int], maxCapacity: int) -> List[List[int]]:
            dp = Knapsack01Removable(maxCapacity, 10)
            for w in weights:
                dp.add(w)
            res = [[0] * (maxCapacity + 1) for _ in range(len(weights))]
            for i, w in enumerate(weights):
                tmp = dp.copy()
                tmp.remove(w)
                row = res[i]
                for j in range(1, maxCapacity + 1):
                    row[j] = tmp.query(j)
            return res

        _, maxCapacity = map(int, input().split())
        weights = list(map(int, input().split()))
        res = solve(weights, maxCapacity)
        for row in res:
            print("".join(map(str, row[1:])))

    p4141()
