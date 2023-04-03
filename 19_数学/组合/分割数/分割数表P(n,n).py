# 分割数表 O(n*sqrt(n))
# 将整数n分拆分成k个非负整数之和的方案数
# dp[n][k]: n个相同的物品放入k个相同的盒子的方案数,每个盒子可以放任意个球.
# dp[n][n]: n个相同的物品分成任意组的方案数


MOD = int(1e9 + 7)


class PartionTableSqrt:
    __slots__ = "_table"

    def __init__(self, max_: int) -> None:
        table = [0] * (max_ + 1)
        table[0] = 1
        for i in range(1, max_ + 1):
            j, sign = 1, 1
            while True:
                tmp1 = i - (j * j * 3 - j) // 2
                if tmp1 < 0:
                    break
                table[i] += table[tmp1] * sign
                table[i] %= MOD
                tmp2 = i - (j * j * 3 + j) // 2
                if tmp2 >= 0:
                    table[i] += table[tmp2] * sign
                    table[i] %= MOD
                j += 1
                sign *= -1
        self._table = table

    def get(self, n: int) -> int:
        if n < 0:
            return 0
        return self._table[n]


table = PartionTableSqrt(int(1e4))
if __name__ == "__main__":
    print(table.get(1000))
