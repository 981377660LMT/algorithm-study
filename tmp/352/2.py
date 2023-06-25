from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个整数：num1 和 num2 。

# 在一步操作中，你需要从范围 [0, 60] 中选出一个整数 i ，并从 num1 减去 2i + num2 。

# 请你计算，要想使 num1 等于 0 需要执行的最少操作数，并以整数形式返回。


# 如果无法使 num1 等于 0 ，返回 -1 。


# !对二进制数的理解、结论
# 类似之前的某场t2 ceilingLog/floorLog
class Solution:
    def makeTheIntegerZero(self, num1: int, num2: int) -> int:
        if num1 == 1 + num2:
            return 1
        if num1 < 1 + num2:
            return -1

        def check(mid: int) -> bool:
            num2Sum = mid * num2
            target = num1 - num2Sum
            if target < mid:
                return False
            # 选mid个数能否凑出remain
            bitCount = target.bit_count()
            return bitCount <= mid

        for res in range(100):
            if check(res):
                return res
        return -1


# num1 = 3, num2 = -2
print(Solution().makeTheIntegerZero(3, -2))
