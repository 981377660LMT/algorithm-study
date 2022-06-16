'''
给你一个整数数组 nums​​​ 和一个整数 k 。你需要将这个数组划分到 k 个相同大小的子集中，
使得同一个子集里面没有两个相同的元素。
请你返回将数组分成 k 个子集后，各子集 不兼容性(是该子集里面最大值和最小值的差) 的 和 的 最小值 ，如果无法分成分成 k 个子集，返回 -1 。
'''

from typing import List
from functools import lru_cache
from itertools import combinations

# 1 <= k <= nums.length <= 16
# nums.length 能被 k 整除。

# !引入优化：每组的第一个元素必定是剩下元素中最小的那个。
# 有了这个优化后等于每组减少了一个元素，排列组合取4个变成了取3个就能大幅缩短时间。
INF = int(1e20)


class Solution:
    def minimumIncompatibility(self, nums: List[int], k: int) -> float:
        @lru_cache(None)
        def dfs(state: int) -> int:
            # 还没有被选的数
            available = [i for i in range(len(nums)) if not state & (1 << i)]
            if not available:
                return 0

            res = INF
            # 还没有被选的数里取size-1个，尝试每一种组合
            for group in combinations(available[1:], size - 1):
                group = (available[0],) + group
                # 不能有相同的数
                if len(set([nums[i] for i in group])) < size:
                    continue
                nextState = state
                for i in group:
                    nextState |= 1 << i
                res = min(res, dfs(nextState) + nums[group[-1]] - nums[group[0]])

            return res

        size = len(nums) // k
        if size == 1:
            return 0
        nums.sort()
        res = dfs(0)
        dfs.cache_clear()
        # 没找到答案就返回-1
        return res if res != INF else -1


print(Solution().minimumIncompatibility([6, 3, 8, 1, 3, 1, 2, 2], 4))
# 输出：6
# 解释：最优的子集分配为 [1,2]，[2,3]，[6,8] 和 [1,3] 。
# 不兼容性和为 (2-1) + (3-2) + (8-6) + (3-1) = 6 。

