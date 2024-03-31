from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个整数 numBottles 和 numExchange 。

# numBottles 代表你最初拥有的满水瓶数量。在一次操作中，你可以执行以下操作之一：

# 喝掉任意数量的满水瓶，使它们变成空水瓶。
# 用 numExchange 个空水瓶交换一个满水瓶。然后，将 numExchange 的值增加 1 。
# 注意，你不能使用相同的 numExchange 值交换多批空水瓶。例如，如果 numBottles == 3 并且 numExchange == 1 ，则不能用 3 个空水瓶交换成 3 个满水瓶。


# 返回你 最多 可以喝到多少瓶水。


# !每次交换后numEchange+1
class Solution:
    def maxBottlesDrunk(self, numBottles: int, numExchange: int) -> int:
        res = 0
        emptyBottles = 0
        while numBottles > 0:
            # 喝水
            res += numBottles
            emptyBottles += numBottles
            numBottles = 0
            while emptyBottles >= numExchange:
                # 交换
                emptyBottles -= numExchange
                numBottles += 1
                numExchange += 1
        return res


# numBottles = 13, numExchange = 6
print(Solution().maxBottlesDrunk(13, 6))
