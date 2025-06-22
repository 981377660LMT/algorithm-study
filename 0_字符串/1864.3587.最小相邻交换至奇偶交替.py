# 3587. 最小相邻交换至奇偶交替
# https://leetcode.cn/problems/minimum-adjacent-swaps-to-alternate-parity/solutions/3705565/ba-ou-shu-shi-zuo-kong-wei-qi-shu-shi-zu-oqbi/
# 给你一个由互不相同的整数组成的数组 nums 。
# 在一次操作中，你可以交换任意两个 `相邻` 元素。
# 在一个排列中，当所有相邻元素的奇偶性交替出现，我们认为该排列是 有效排列。这意味着每对相邻元素中一个是偶数，一个是奇数。
# 请返回将 nums 变成任意一种 有效排列 所需的最小相邻交换次数。
# 如果无法重排 nums 来获得有效排列，则返回 -1。
#
# !注意和1864的区别，这里必须交换`相邻`元素，而不是任意两个元素。


from typing import List


INF = int(1e18)


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        pos1 = [i for i, v in enumerate(nums) if v & 1]
        n, m = len(nums), len(pos1)

        def calc(start: int) -> int:
            """start=0/1 => '1'在偶数/奇数位置."""
            if (n - start + 1) // 2 != m:
                return INF
            return sum(abs(i - j) for i, j in zip(range(start, n, 2), pos1))

        res = min(calc(0), calc(1))
        return res if res < INF else -1
