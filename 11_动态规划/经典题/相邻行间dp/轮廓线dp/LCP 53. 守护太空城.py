from functools import cache
from typing import List, Tuple

# 选择一个舱室开启屏障，能量消耗为 2
# 选择相邻两个舱室开启联合屏障，能量消耗为 3
# 对于已开启的一个屏障，多维持一时刻，能量消耗为 1
# 1 <= time.length == position.length <= 500
# 1 <= time[i] <= 5
# 0 <= position[i] <= 100


# class Solution:
#     def defendSpaceCity(self, time: List[int], position: List[int]) -> int:
#         @cache
#         def dfs(r: int, c: int, cState: Tuple[int, ...]) -> int:
#             """轮廓线+三进制状压dp
#             Args:
#                 cState (Tuple[int, ...]): 前COL个格子开盾的状态

#             规定开盾是从下往上，从右往左开
#             """
#             if r == ROW:
#                 return 0
#             if c == COL:
#                 return dfs(r + 1, 0, cState)

#             mustUnion = r > 0 and grid[r - 1][c] == 1 and cState[0] == 0  # 必须开联合盾
#             canUnion = r > 0 and cState[0] == 0  # 可以开联合盾
#             canKeepOne = not mustUnion and c > 0 and cState[-1] == 1  # 可以维持单独盾
#             canKeepUnion = not mustUnion and r > 0 and c > 0 and cState[-1] == 0  # 可以维持联合盾

#             res = int(1e20)
#             if mustUnion:
#                 res = min(res, 3 + dfs(r, c + 1, cState[1:] + (2,)))
#             else:
#                 # 开联合盾 3
#                 if canUnion:
#                     res = min(res, 3 + dfs(r, c + 1, cState[1:] + (2,)))

#                 # 新开单独盾 2
#                 res = min(res, 2 + dfs(r, c + 1, cState[1:] + (1,)))

#                 # 多维持 1 时刻
#                 if canKeep:
#                     res = min(res, 1 + dfs(r, c + 1, cState[1:] + (1,)))

#                 # 不开盾
#                 res = min(res, dfs(r, c + 1, cState[1:] + (0,)))

#             return res

#         ROW, COL = max(position) + 2, max(time)  # 保证最后一行都会被覆盖 因此加一个虚拟行全为0
#         grid = [[0] * COL for _ in range(ROW)]
#         for r, c in set(zip(position, time)):
#             grid[r][c] = 1

#         res = dfs(0, 0, tuple([0] * COL))
#         dfs.cache_clear()
#         return res


# 细节太多了...
