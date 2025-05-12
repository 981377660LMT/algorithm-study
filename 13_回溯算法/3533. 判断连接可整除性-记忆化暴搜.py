# 3533. 判断连接可整除性(记忆化暴搜)
# https://leetcode.cn/problems/concatenated-divisibility/description/
#
# 给你一个正整数数组 nums 和一个正整数 k。
# 当 nums 的一个 排列 中的所有数字，按照排列顺序 连接其十进制表示 后形成的数可以 被 k  整除时，我们称该排列形成了一个 可整除连接 。
# 返回能够形成 可整除连接 且 字典序 最小 的排列（按整数列表的形式表示）。如果不存在这样的排列，返回一个空列表。
#
# !按照顺序搜索，碰到的第一个可整除连接就是字典序最小的排列。

from functools import lru_cache
from typing import List


class Solution:
    def concatenatedDivisibility(self, nums: List[int], k: int) -> List[int]:
        nums.sort()
        pow10 = [pow(10, len(str(v))) for v in nums]

        res = []

        @lru_cache(None)
        def dfs(remain: int, mod: int) -> bool:
            if remain == 0:
                return mod == 0
            for i, (p10, v) in enumerate(zip(pow10, nums)):
                if not (remain >> i) & 1:
                    continue
                if dfs(remain ^ (1 << i), (mod * p10 + v) % k):
                    res.append(v)
                    return True
            return False

        ok = dfs((1 << len(nums)) - 1, 0)
        dfs.cache_clear()
        if not ok:
            return []
        return res[::-1]
