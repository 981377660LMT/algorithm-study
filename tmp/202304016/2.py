from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个下标从 0 开始的整数数组 nums 和 divisors 。

# divisors[i] 的 可整除性得分 等于满足 nums[j] 能被 divisors[i] 整除的下标 j 的数量。


# 返回 可整除性得分 最大的整数 divisors[i] 。如果有多个整数具有最大得分，则返回数值最小的一个。
class Solution:
    def maxDivScore(self, nums: List[int], divisors: List[int]) -> int:
        scores = [0] * len(divisors)
        for i in range(len(divisors)):
            for j in range(len(nums)):
                if nums[j] % divisors[i] == 0:
                    scores[i] += 1
        maxScores = max(scores)
        res = []
        for i in range(len(scores)):
            if scores[i] == maxScores:
                res.append(divisors[i])
        return min(res)
