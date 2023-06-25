from itertools import combinations
from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。如果下标对 i、j 满足 0 ≤ i < j < nums.length ，如果 nums[i] 的 第一个数字 和 nums[j] 的 最后一个数字 互质 ，则认为 nums[i] 和 nums[j] 是一组 美丽下标对 。

# 返回 nums 中 美丽下标对 的总数目。


# 对于两个整数 x 和 y ，如果不存在大于 1 的整数可以整除它们，则认为 x 和 y 互质 。换而言之，如果 gcd(x, y) == 1 ，则认为 x 和 y 互质，其中 gcd(x, y) 是 x 和 k 最大公因数 。
class Solution:
    def countBeautifulPairs(self, nums: List[int]) -> int:
        res = 0
        for a, b in combinations(nums, 2):
            first = int(str(a)[0])
            last = int(str(b)[-1])
            if gcd(first, last) == 1:
                res += 1
        return res
