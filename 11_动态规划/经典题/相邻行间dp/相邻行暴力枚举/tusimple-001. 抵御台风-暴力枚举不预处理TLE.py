# Algorithm for finding the fewest rectangles to cover a set of rectangles without overlapping
# 字符串由字符 '.' 和 '#' 构成，其中 '.' 表示关闭的窗户，'#' 表示打开的窗户。
# TuTu 希望他尽可能少的切割胶带。
# 所以 TuTu 每次可以切割下一条胶带并用它封住一些连续的（相邻）水平窗户、
# 或者连续的（相邻）垂直窗户，使其不备台风侵袭。
# 现在 TuTu 想知道他最少需要切割下多少条胶带。
# 已经贴上的窗户不能再贴胶带


# 输入:
# 3 5

# .#.#
#####
# .#.#

# 输出:
# 5

# 说明:
# A.C.E
# ABCDE
# A.C.E


# 1<= n <= 1000, 1<= m <= 10
from collections import defaultdict
from typing import Generator, List, Tuple

# 生成器比append到数组快很多


def dfs(
    index: int, curRow: int, preState: int, curState: int, curCost: int,
) -> Generator[Tuple[int, int], None, None]:
    """返回(行的状态，这一行的花费);curState 记录哪几个位置向下 默认横向连接"""
    if index == col:
        yield curState, curCost
        return

    # 和等于1有关系吗
    if matrix[curRow][index] == 1:
        # 这个位置向下连接
        isUpVertical = (preState >> index) & 1
        yield from dfs(
            index + 1,
            curRow,
            preState,
            curState | (1 << index),
            curCost + (0 if isUpVertical else 1),
        )

        # 这个位置横向连接，看前一个是否横向连接
        isLeftHorizontal = (
            index > 0 and matrix[curRow][index - 1] == 1 and not ((curState >> (index - 1)) & 1)
        )

        yield from dfs(
            index + 1, curRow, preState, curState, curCost + (0 if isLeftHorizontal else 1)
        )
    else:
        yield from dfs(index + 1, curRow, preState, curState, curCost)


n, m = map(int, input().split())

matrix = []
for _ in range(n):
    matrix.append([1 if char == '#' else 0 for char in input()])
matrix = [list(col) for col in zip(*matrix)]

row, col = len(matrix), len(matrix[0])


# 每一行长度为10，1000行
dp = [defaultdict(lambda: int(1e20)) for _ in range(row)]
startStates = dfs(0, 0, 0, 0, 0)
for state, cost in startStates:
    dp[0][state] = cost

for r in range(1, row):
    # 这里可以改一下，可以在外面预处理出前后两行的state，计算当前行的newCost
    for preState, preCost in dp[r - 1].items():
        curStates = dfs(0, r, preState, 0, 0)
        for curState, curCost in curStates:
            dp[r][curState] = min(dp[r][curState], preCost + curCost)

print(min(dp[-1].values()))
