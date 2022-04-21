from collections import defaultdict
from functools import lru_cache
from pprint import pprint
from typing import Generator, List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)

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
# 这题数据非常弱


@lru_cache(None)
def getDigit(num: int, index: int):
    """返回 `3` 进制下 `num` 的 `index` 位的数字，`index` 最低位(最右)为 0 """
    prefix = num // pow(3, index)
    return prefix % 3


class Solution:
    def defendSpaceCity(self, time: List[int], position: List[int]) -> int:
        def dfs(
            i: int, curRow: int, pre: int, cur: int, cost: int,
        ) -> Generator[Tuple[int, int], None, None]:
            """返回(行的状态，这一行的花费)
        
            0:不主动开启
            1:只主动开自己
            2:主动开自己和下一个
            """
            if i == col:
                yield cur, cost
                return

            isLeftOpen1 = i > 0 and getDigit(cur, i - 1) == 1  # 左边开了一个单独的
            isLeftOpen2 = i > 0 and getDigit(cur, i - 1) == 2  # 左边开了一个联合的
            isCoveredByUp = curRow > 0 and getDigit(pre, i) == 2  # 当前上方开了一个联合的

            if isCoveredByUp:
                yield from dfs(i + 1, curRow, pre, cur * 3, cost)  # 不能重叠开
            else:
                if matrix[curRow][i] == 0:
                    yield from dfs(i + 1, curRow, pre, cur * 3, cost)  # 可以不开
                yield from dfs(
                    i + 1, curRow, pre, cur * 3 + 1, cost + (1 if isLeftOpen1 else 2)
                )  # 单独开
                if curRow < row - 1:
                    yield from dfs(
                        i + 1, curRow, pre, cur * 3 + 2, cost + (1 if isLeftOpen2 else 3)  # 联合开
                    )

        row, col = max(position) + 1, max(time)
        matrix = [[0] * col for _ in range(row)]
        for c, r in zip(time, position):
            matrix[r][c - 1] = 1

        dp = [defaultdict(lambda: int(1e20)) for _ in range(row)]
        startStates = dfs(0, 0, 0, 0, 0)
        for state, cost in startStates:
            dp[0][state] = cost

        for r in range(1, row):
            for preState, preCost in dp[r - 1].items():
                # 这里应该要预处理，每个位置是不贴/单独开/主动联合开，计算出对应的花费，然后行间转移看是否合法(可以预处理处一个有向图的合法状态转移，时间复杂度O(n*3^2m))
                for curState, curCost in dfs(0, r, preState, 0, 0):
                    dp[r][curState] = min(dp[r][curState], preCost + curCost)

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
