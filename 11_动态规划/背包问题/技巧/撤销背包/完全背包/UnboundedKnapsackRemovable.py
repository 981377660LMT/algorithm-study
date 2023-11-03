from typing import List, Optional


class UnboundedKnapsackRemovable:
    """可撤销完全背包，用于求解方案数/可行性."""

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
            for i in range(weight, self._maxWeight + 1):
                dp[i] += dp[i - weight]
        else:
            mod = self._mod
            for i in range(weight, self._maxWeight + 1):
                dp[i] = (dp[i] + dp[i - weight]) % mod

    def remove(self, weight: int) -> None:
        """移除一个重量为weight的物品.需要保证weight物品存在."""
        dp = self._dp
        if self._mod is None:
            for i in range(self._maxWeight, weight - 1, -1):
                dp[i] -= dp[i - weight]
        else:
            mod = self._mod
            for i in range(self._maxWeight, weight - 1, -1):
                dp[i] = (dp[i] - dp[i - weight]) % mod

    def query(self, weight: int) -> int:
        """查询组成重量为weight的物品有多少种方案."""
        return self._dp[weight] if 0 <= weight <= self._maxWeight else 0

    def copy(self) -> "UnboundedKnapsackRemovable":
        return UnboundedKnapsackRemovable(self._maxWeight, self._mod, self._dp[:])

    def __repr__(self) -> str:
        return self._dp.__repr__()


if __name__ == "__main__":
    # https://atcoder.jp/contests/agc049/tasks/agc049_d
    ...
