from typing import List
from collections import defaultdict
from functools import lru_cache


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
# O(n*3^2m)  n<=100 m<=5
# 这题数据非常弱 可以预处理行状态+暴力枚举


@lru_cache(None)
def getDigit(num: int, index: int):
    return num // pow(3, index) % 3


# 1. 预处理每行状态
# 这里也可以写dfs(mask)加cache
def dfs(index: int, pre: int, state: int, cost: int) -> None:
    """预处理行状态"""
    if index == 5:
        costByState[state] = cost
        return

    dfs(index + 1, 0, state * 3, cost)
    dfs(index + 1, 1, state * 3 + 1, cost + 1 if pre == 1 else cost + 2)
    dfs(index + 1, 2, state * 3 + 2, cost + 1 if pre == 2 else cost + 3)


costByState = [0] * pow(3, 5)
dfs(0, -1, 0, 0)

# adjMap = defaultdict(set)  # 有向图预处理转移
# for pre in costByState.keys():
#     for cur in costByState.keys():
#         # 屏障不能重叠开
#         if all((getDigit(pre, i), getDigit(cur, i)) not in [(2, 1), (2, 2)] for i in range(5)):
#             adjMap[pre].add(cur)


class Solution:
    def defendSpaceCity(self, time: List[int], position: List[int]) -> int:
        def check(rowIndex: int, pre: int, cur: int) -> bool:
            for c in range(col):
                s1, s2 = getDigit(pre, c), getDigit(cur, c)
                if s2 == 0 and matrix[rowIndex][c] == 1 and s1 != 2:
                    return False
                if (s1, s2) in [(2, 1), (2, 2)]:
                    return False
            return True

        row, col = max(position) + 1, max(time)
        matrix = [[0] * col for _ in range(row)]
        for c, r in zip(time, position):
            matrix[r][c - 1] = 1

        dp = [defaultdict(lambda: int(1e20)) for _ in range(row)]
        for state in range(3 ** col):
            if all(getDigit(state, c) != 0 for c in range(col) if matrix[0][c] == 1):
                dp[0][state] = costByState[state]

        for r in range(1, row):
            for preState, preCost in dp[r - 1].items():
                for curState in range(3 ** col):
                    if check(r, preState, curState):
                        dp[r][curState] = min(dp[r][curState], preCost + costByState[curState])

        return min(dp[-1].values())


print(Solution().defendSpaceCity(time=[1, 2, 1], position=[6, 3, 3]))
print(Solution().defendSpaceCity(time=[1, 1, 1, 2, 2, 3, 5], position=[1, 2, 3, 1, 2, 1, 3]))
print(Solution().defendSpaceCity(time=[2, 1, 2, 2, 1, 1], position=[4, 5, 1, 3, 1, 4]))
print(
    Solution().defendSpaceCity(
        time=[4, 1, 4, 4, 2, 4, 3, 3, 2, 3, 2, 3, 2, 3, 1, 4, 1, 4, 1, 2, 1, 2, 1, 3, 3, 4, 3, 2],
        position=[
            2,
            5,
            8,
            5,
            8,
            6,
            7,
            2,
            3,
            6,
            7,
            3,
            9,
            1,
            9,
            9,
            3,
            3,
            6,
            5,
            8,
            1,
            1,
            4,
            8,
            7,
            5,
            4,
        ],
    )
)
# 5 9 9 28
