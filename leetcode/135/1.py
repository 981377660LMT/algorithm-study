from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个 正 整数 x 和 y ，分别表示价值为 75 和 10 的硬币的数目。

# Alice 和 Bob 正在玩一个游戏。每一轮中，Alice 先进行操作，Bob 后操作。每次操作中，玩家需要拿出价值 总和 为 115 的硬币。如果一名玩家无法执行此操作，那么这名玩家 输掉 游戏。


# 两名玩家都采取 最优 策略，请你返回游戏的赢家。
# !每次必须1个75和4个10.
class Solution:
    def losingPlayer(self, x: int, y: int) -> str:
        failTurn = 0
        while True:
            if x < 1 or y < 4:
                break
            x -= 1
            y -= 4
            failTurn ^= 1
        return "Alice" if failTurn == 1 else "Bob"
