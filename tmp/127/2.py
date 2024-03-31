from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的二进制数组 possible 。

# 莉叩酱和冬坂五百里正在玩一个有 n 个关卡的游戏，游戏中有一些关卡是 困难 模式，其他的关卡是 简单 模式。如果 possible[i] == 0 ，那么第 i 个关卡是 困难 模式。一个玩家通过一个简单模式的关卡可以获得 1 分，通过困难模式的关卡将失去 1 分。

# 游戏的一开始，莉叩酱将从第 0 级开始 按顺序 完成一些关卡，然后冬坂五百里会完成剩下的所有关卡。

# 假设两名玩家都采取最优策略，目的是 最大化 自己的得分，莉叩酱想知道自己 最少 需要完成多少个关卡，才能获得比冬坂五百里更多的分数。

# 请你返回莉叩酱获得比冬坂五百里更多的分数所需要完成的 最少 关卡数目，如果 无法 达成，那么返回 -1 。


# 注意，每个玩家都至少需要完成 1 个关卡。


class Solution:
    def minimumLevels(self, possible: List[int]) -> int:
        nums = [1 if possible[i] == 1 else -1 for i in range(len(possible))]
        preSum = [0] + list(accumulate(nums))
        n = len(possible)
        for res in range(1, n):  # 完成[1,n-1]个关卡
            sum1 = preSum[res]
            sum2 = preSum[-1] - preSum[res]
            if sum1 > sum2:
                return res
        return -1
