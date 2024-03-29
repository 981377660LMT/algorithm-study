# 有区别的n个球放入k个无区别的盒子(每个盒子至少放一个球)

MOD = int(1e9 + 7)


class Stirling2Table:
    """第二类斯特林数表."""

    __slots__ = ("_n", "_k", "_table")

    def __init__(self, n: int, k: int):
        self._n = n
        self._k = k
        table = [[0] * (k + 1) for _ in range(n + 1)]
        table[0][0] = 1
        for i in range(1, n + 1):
            for j in range(1, k + 1):
                table[i][j] = table[i - 1][j - 1] + table[i - 1][j] * j
                table[i][j] %= MOD
        self._table = table

    def get(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        return self._table[n][k]

    def __getitem__(self, index: tuple[int, int]) -> int:
        return self.get(*index)


if __name__ == "__main__":
    n, k = map(int, input().split())
    table = Stirling2Table(n, k)
    print(table.get(n, k))
    print(table[n, k])
