# 3859. 统计包含 K 个不同整数的子数组
# https://leetcode.cn/problems/count-subarrays-with-k-distinct-integers/
# 给你一个整数数组 nums 和两个整数 k 和 m。
# 返回一个整数，表示满足以下条件的 子数组 的数量：
# 子数组 恰好 包含 ​​​​​​​k 个不同的 整数。
# 在子数组中，每个 不同的 整数 至少 出现 m 次。
# 子数组 是数组中一个连续的、非空 元素序列。


from collections import defaultdict


class Solution:
    def countSubarrays(self, nums: list[int], k: int, m: int) -> int:
        def f(curK: int) -> int:
            """子数组至少包含 curK 个不同整数, 且至少有 k 个整数至少出现m次的子数组数量."""
            counter = defaultdict(int)
            hit, res, left = 0, 0, 0
            for x in nums:
                counter[x] += 1
                hit += counter[x] == m
                while len(counter) >= curK and hit >= k:
                    out = nums[left]
                    hit -= counter[out] == m
                    counter[out] -= 1
                    if counter[out] == 0:
                        del counter[out]
                    left += 1
                res += left
            return res

        return f(k) - f(k + 1)
