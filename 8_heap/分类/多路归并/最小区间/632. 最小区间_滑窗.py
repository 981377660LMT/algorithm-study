# https://leetcode.cn/problems/smallest-range-covering-elements-from-k-lists/solutions/2982588/liang-chong-fang-fa-dui-pai-xu-hua-dong-luih5/
# 你有 k 个 非递减排列 的整数列表。找到一个 最小 区间，使得 k 个列表中的每个列表至少有一个数包含在其中。
# 我们定义如果 b-a < d-c 或者在 b-a == d-c 时 a < c，则区间 [a,b] 比 [c,d] 小。
#
# !把所有元素都合在一起排序
# !合法区间等价于 pairs 的一个连续子数组，满足列表编号 0,1,2,…,k−1 都在这个子数组中

from typing import List

INF = int(1e18)


class Solution:
    def smallestRange(self, nums: List[List[int]]) -> List[int]:
        pairs = sorted((v, i) for i, vs in enumerate(nums) for v in vs)
        resL, resR = -INF, INF
        counter = [0] * len(nums)
        lack = len(nums)
        left = 0
        for r, i in pairs:
            if counter[i] == 0:
                lack -= 1
            counter[i] += 1
            while lack == 0:
                l, j = pairs[left]
                if r - l < resR - resL:
                    resL, resR = l, r
                counter[j] -= 1
                if counter[j] == 0:
                    lack += 1
                left += 1

        return [resL, resR]
