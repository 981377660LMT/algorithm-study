from collections import defaultdict
from typing import List

# 1 <= nums.length <= 1000
# 0 <= nums[i] < 2^16
# 982. 按位与为零的三元组
# 两数之和 => 哈希表优化


class Solution:
    def countTriplets(self, A: List[int]) -> int:
        memo = defaultdict(int)
        for n1 in A:
            for n2 in A:
                memo[n1 & n2] += 1

        res = 0
        for num in A:
            for key, val in memo.items():
                if num & key == 0:
                    res += val
        return res


# 高维前缀和 & 快速沃尔什变换 更快
# https://leetcode.cn/problems/triples-with-bitwise-and-equal-to-zero/solution/liang-chong-90duo-de-fang-fa-gao-wei-qian-zhui-he-/
