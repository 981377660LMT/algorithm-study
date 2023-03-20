from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 value 。

# 在一步操作中，你可以对 nums 中的任一元素加上或减去 value 。

# 例如，如果 nums = [1,2,3] 且 value = 2 ，你可以选择 nums[0] 减去 value ，得到 nums = [-1,2,3] 。
# 数组的 MEX (minimum excluded) 是指其中数组中缺失的最小非负整数。

# 例如，[-1,2,3] 的 MEX 是 0 ，而 [1,0,3] 的 MEX 是 2 。
# 返回在执行上述操作 任意次 后，nums 的最大 MEX 。


class Solution:
    def findSmallestInteger(self, nums: List[int], value: int) -> int:
        # mod
        mods = [num % value for num in nums]
        mp = defaultdict(list)
        for mod in mods:
            if len(mp[mod]) == 0:
                mp[mod].append(mod)
            else:
                mp[mod].append(mp[mod][-1] + value)
        allNums = set()
        for v in mp.values():
            allNums.update(v)
        s = set(allNums)
        res = 0
        while res in s:
            res += 1
        return res


# nums = [1,-10,7,13,6,8], value = 5
print(Solution().findSmallestInteger(nums=[1, -10, 7, 13, 6, 8], value=5))
# [3,0,3,2,4,2,1,1,0,4]
# 5
print(Solution().findSmallestInteger(nums=[3, 0, 3, 2, 4, 2, 1, 1, 0, 4], value=5))
