class AlphaPresum:
    def __init__(self, s: str) -> None:
        n = len(s)
        preSum = [[0] * 26 for _ in range(n + 1)]

        for i in range(1, n + 1):
            preSum[i][ord(s[i - 1]) - ord('a')] += 1
            for j in range(26):
                preSum[i][j] += preSum[i - 1][j]

        self._preSum = preSum

    def getCountOfSlice(self, char: str, left: int, right: int) -> int:
        """
        >>> preSum = AlphaPresum("abcdabcd")
        >>> print(preSum.getCountOfSlice('a', 0, 2)) # s[0:2]间'a'个数为1
        >>> print(preSum.getCountOfSlice('a', 0, 8)) # s[0:8]间'a'个数为2
        """
        assert 0 <= left <= right < len(self._preSum)
        ord_ = ord(char) - ord('a')
        return self._preSum[right][ord_] - self._preSum[left][ord_]


preSum = AlphaPresum("abcdabcd")
# s[0:2]间'a'个数为1
print(preSum.getCountOfSlice('a', 0, 2))
print(preSum.getCountOfSlice('a', 0, 8))
