from collections import defaultdict
from typing import List
import gc

gc.disable()

# 1 <= time.length == position.length <= 500
# 1 <= time[i] <= 5
# 0 <= position[i] <= 100
# 矩阵为5行100列
# 太空城中的一些舱室将要受到陨石雨的冲击，这些舱室按照编号 0 ~ N 的顺序依次排列。
# 为了阻挡陨石损毁舱室，太空城可以使用能量展开防护屏障，具体消耗如下：
# 选择一个舱室开启屏障，能量消耗为 2
# 选择相邻两个舱室开启联合屏障，能量消耗为 3
# 对于已开启的一个屏障，多维持一时刻，能量消耗为 1
# 已知陨石雨的影响范围和到达时刻，time[i] 和 position[i] 分别表示该陨石的到达时刻和冲击位置。
# 请返回太空舱能够守护所有舱室所需要的最少能量。


# 5的话是暗示三进制状压dp
# O(n*3^m) 子集状压 DP + 贪心预处理


# 1.预处理每行的状态的花费
def dfs1(index: int, state: int, cost: int, isPreSelected: bool) -> None:
    """预处理单独开的花费"""
    if index == 5:
        one[state] = cost
        return
    dfs1(index + 1, state, cost, False)
    dfs1(index + 1, state | (1 << index), cost + (1 if isPreSelected else 2), True)


def dfs2(index: int, state: int, cost: int, isPreSelected: bool) -> None:
    """预处理联合开的花费"""
    if index == 5:
        two[state] = cost
        return
    dfs2(index + 1, state, cost, False)
    dfs2(index + 1, state | (1 << index), cost + (1 if isPreSelected else 3), True)


one = [0] * (1 << 5)
two = [0] * (1 << 5)
dfs1(0, 0, 0, False)
dfs2(0, 0, 0, False)

# 也可以这么写，判断左边一位是否与最低位相邻
for i in range(1, 1 << 5):
    lowbit = i & -i
    one[i] = one[i - lowbit] + (2 if ((lowbit << 1) & i) == 0 else 1)
    two[i] = two[i - lowbit] + (3 if ((lowbit << 1) & i) == 0 else 1)


class Solution:
    def defendSpaceCity(self, time: List[int], position: List[int]) -> int:
        row, col = max(position) + 1, max(time)
        mask = 1 << col
        matrix = [0 for _ in range(row)]
        for c, r in zip(time, position):
            matrix[r] |= 1 << (c - 1)

        dp = [defaultdict(lambda: int(1e20)) for _ in range(row)]  # dp[i][j]表示第i行开联合盾的状态为j时的最小总花费

        for r in range(row):
            # 这个利用屏障不能重叠的性质把时间复杂度从 O(n*3^2m) 降到了 O(n*3^m) ; O(n*3^2m)的行间暴力枚举解法会多计算了两行相同位置都开联合屏障的情况
            for state in range(mask):  # state 表示 当前行哪些位置开联合
                comp = state ^ (mask - 1)  # 枚举上一行哪些位置开联合
                g1, g2 = comp, 0  # g1表示上一行实际开联合的位置

                while True:
                    cost1 = dp[r - 1][g1] if r - 1 >= 0 else 0
                    cost2 = (
                        one[g2 & matrix[r]] + two[state]
                    )  # 每一行哪些开1 哪些开2； g2 & matrix[r] 表示这一行必须单独开的位置
                    dp[r][state] = min(dp[r][state], cost1 + cost2)
                    if g1 == 0:
                        break

                    g1 = comp & (g1 - 1)
                    g2 = comp ^ g1

        return min(dp[-1].values())  # return dp[-1][0]  最后一行不能开联合盾


print(Solution().defendSpaceCity(time=[1, 2, 1], position=[6, 3, 3]))
print(Solution().defendSpaceCity(time=[1, 1, 1, 2, 2, 3, 5], position=[1, 2, 3, 1, 2, 1, 3]))
