from typing import List, Optional


class BoundedKnapsackRemovable:
    __slots__ = "_dp", "_mod", "_maxValue", "_countSum"

    def __init__(
        self, maxValue: int, mod: Optional[int] = None, dp: Optional[List[int]] = None
    ) -> None:
        if dp is not None:
            self._dp = dp
        else:
            self._dp = [0] * (maxValue + 1)
            self._dp[0] = 1
        self._mod = mod
        self._maxValue = maxValue
        self._countSum = 0

    def add(self, value: int, count: int) -> None:
        if value <= 0:
            raise ValueError("value must be positive, but got %d" % value)
        dp = self._dp
        self._countSum += count * value
        maxJ = min(self._countSum, self._maxValue)
        if self._mod is None:
            for j in range(value, maxJ + 1):
                dp[j] += dp[j - value]
            for j in range(maxJ, value * (count + 1) - 1, -1):
                dp[j] -= dp[j - value * (count + 1)]
        else:
            mod = self._mod
            for j in range(value, maxJ + 1):
                dp[j] = (dp[j] + dp[j - value]) % mod
            for j in range(maxJ, value * (count + 1) - 1, -1):
                dp[j] = (dp[j] - dp[j - value * (count + 1)]) % mod

    def remove(self, value: int, count: int) -> None:
        maxJ = min(self._countSum, self._maxValue)
        if self._mod is None:
            for i in range(value * (count + 1), maxJ + 1):
                self._dp[i] += self._dp[i - value * (count + 1)]
            for i in range(maxJ, value, -1):
                self._dp[i] -= self._dp[i - value]
        else:
            mod = self._mod
            for i in range(value * (count + 1), maxJ + 1):
                self._dp[i] = (self._dp[i] + self._dp[i - value * (count + 1)]) % mod
            for i in range(maxJ, value, -1):
                self._dp[i] = (self._dp[i] - self._dp[i - value]) % mod

        self._countSum -= count * value

    def query(self, value: int) -> int:
        """!注意需要特判重量为0."""
        return self._dp[value] if 0 <= value <= self._maxValue else 0

    def copy(self) -> "BoundedKnapsackRemovable":
        res = BoundedKnapsackRemovable(self._maxValue, self._mod, self._dp[:])
        res._countSum = self._countSum
        return res

    def __repr__(self) -> str:
        return self._dp.__repr__()
