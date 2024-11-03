# abc378-E - Mod Sigma Problem
# https://atcoder.jp/contests/abc378/editorial/11289
# !给定数组A，和模数M。求∑1≤l≤r≤n((∑l≤i≤rA[i])%M)
# 所有子数组之和模M的总和
#
# 子数组之和转为前缀和，即 ∑​(Sr​−Sl​)%M.
# 去掉模，得
# !如果 Sr​>Sl​，则 为Sr​−Sl，否则为 Sr​−Sl+M


class BITMap:

    __slots__ = "n", "total", "_tree"

    def __init__(self, n: int):
        n += 1
        self.n = n
        self.total = 0
        self._tree = dict()

    def add(self, index: int, delta: int) -> None:
        self.total += delta
        index += 1
        while index <= self.n:
            self._tree[index - 1] = self._tree.get(index - 1, 0) + delta
            index += index & -index

    def queryPrefix(self, end: int) -> int:
        if end > self.n:
            end = self.n
        res = 0
        while end > 0:
            res += self._tree.get(end - 1, 0)
            end &= end - 1
        return res

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        if start == 0:
            return self.queryPrefix(end)
        pos, neg = 0, 0
        while end > start:
            pos += self._tree.get(end - 1, 0)
            end &= end - 1
        while start > end:
            neg += self._tree.get(start - 1, 0)
            start &= start - 1
        return pos - neg

    def queryAll(self):
        return self.total


if __name__ == "__main__":
    N, M = map(int, input().split())
    A = list(map(int, input().split()))

    presum = [0] * (N + 1)
    for i in range(N):
        presum[i + 1] = (presum[i] + A[i]) % M

    bit = BITMap(M)  # 值域树状数组，记录前缀和模M的个数
    res = 0
    curSum = 0
    for r in range(N + 1):
        sum1 = presum[r] * r - curSum
        sum2 = bit.queryRange(presum[r] + 1, M) * M
        res += sum1 + sum2
        curSum += presum[r]
        bit.add(presum[r], 1)
    print(res)
