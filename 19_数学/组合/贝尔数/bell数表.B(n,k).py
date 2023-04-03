MOD = int(1e9 + 7)


class BellTable1:
    """贝尔数表.计算Bell(n, k)."""

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
        self._table = table  # 第二类斯特林数

    def get(self, n: int, k: int) -> int:
        res = 0
        for i in range(k + 1):
            res += self._table[n][i]
            res %= MOD
        return res

    def __getitem__(self, index: tuple[int, int]) -> int:
        return self.get(*index)


if __name__ == "__main__":
    n, k = map(int, input().split())
    table = BellTable1(n, k)
    print(table.get(n, k))
