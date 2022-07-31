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
# 每一行长度为10，1000行
from collections import defaultdict

INF = int(1e20)
n, m = map(int, input().split())
matrix = [0] * n
for i in range(n):
    states = ["1" if char == "#" else "0" for char in input()][::-1]
    matrix[i] = int("".join(states), 2)


# 1. 很自然想到的解法是 预先处理每行状态+行间暴力dp：二进制枚举每行每个位置是横着贴/竖着贴，
# 预处理出行间合理的转移状态，得到一个有向带权图，再在行间进行暴力枚举
# 这样的时间复杂度是 O(n*4^m) 的，1024*1024*1000 是会超时的

# 2. 可以重新对贴胶带的状态做一个划分：横胶带(默认连续贴)，长度为1的竖胶带，长度大于1的竖胶带。
# `因为贴长度为1的竖胶带不如贴成横胶带，所以要贴竖胶带就要贴长度大于1的竖胶带。`
# 记 dp[row][state] 为 第 row 行贴长度大于1的竖胶带的状态为 state 时的最小总花费，
# 注意到两行同一个位置都贴长度大于1的竖胶带的状态是互斥的，因此可以在上一行从 state 的补集 里枚举子集，
# 看上一行哪些位置贴长度大于1的竖胶带，当前行剩下的1的位置就都贴横胶带；
# 注意枚举时要check状态合法性，因为0的位置是不可以贴胶带的
# 整个算法时间复杂度为 O(n*3^m)


def dfs1(index: int, state: int, cost: int, isPreSelected: bool) -> None:
    """预处理贴横胶带每种状态对应的花费"""
    if index == col:
        horizontal[state] = cost
        return
    dfs1(index + 1, state, cost, False)
    dfs1(index + 1, state | (1 << index), cost + int(not isPreSelected), True)


def dfs2(index: int, state: int, cost: int) -> None:
    """预处理贴竖胶带每种状态对应的花费"""
    if index == col:
        vertical[state] = cost
        return
    dfs2(index + 1, state, cost)
    dfs2(index + 1, state | (1 << index), cost + 1)


row, col = n, m
mask = 1 << m
horizontal = [0] * mask
vertical = [0] * mask
dfs1(0, 0, 0, False)
dfs2(0, 0, 0)

# dp = [defaultdict(lambda: int(1e20)) for _ in range(row)]
dp = [[INF] * mask for _ in range(row)]

for state in range(mask):
    if (state | matrix[0]) == matrix[0]:
        comp = state ^ (mask - 1)  # 贴横的位置
        dp[0][state] = horizontal[comp] + vertical[state]

for r in range(1, row):
    for state in range(mask):
        if not ((state | matrix[r]) == matrix[r]):
            continue
        comp = state ^ (mask - 1)
        g1, g2 = comp, 0  # g1表示上一行贴长度大于1的竖胶带的状态

        while True:
            if (g1 | matrix[r]) == matrix[r]:
                cost1 = dp[r - 1][g1] if r - 1 >= 0 else 0
                cost2 = horizontal[g2 & matrix[r]] + vertical[state]
                dp[r][state] = min(dp[r][state], cost1 + cost2)
            if g1 == 0:
                break
            g1 = comp & (g1 - 1)
            g2 = comp ^ g1


print(min(dp[-1]))
# 从 O(n*4^m) 到 O(n*3^m)
