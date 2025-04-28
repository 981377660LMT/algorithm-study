# 644. 子数组最大平均数 II
# 找出 长度大于等于 k 且含最大平均值的连续子数组 并输出这个最大平均值
# https://leetcode.cn/problems/maximum-average-subarray-ii/solutions/860135/fu-za-du-wei-onde-dan-diao-zhan-fa-by-li-trzz/
#
# 平均数技巧：凸包
# https://atcoder.jp/contests/abc341/tasks/abc341_g
# TODO: cdq分治框架，凸优化框架（abc那道前缀平均数最大值）

from collections import deque
from typing import List


class Solution:
    def findMaxAverage(self, nums: List[int], k: int) -> float:
        # (count, sum)
        stack = deque()
        st = 0
        N = len(nums)
        for i in range(0, k):
            st += nums[i]
        ans = st / k
        ct = k
        for i in range(1, N - k + 1):
            st += nums[i + k - 1]
            ct += 1
            s1 = nums[i - 1]
            c1 = 1
            # S0 / C0 >= S1 / C1
            while stack and stack[-1][1] * c1 >= s1 * stack[-1][0]:
                c0, s0 = stack.pop()
                c1 += c0
                s1 += s0
            stack.append((c1, s1))
            # S0 / C0 <= st / ct
            while stack and stack[0][1] * ct <= st * stack[0][0]:
                c0, s0 = stack.popleft()
                ct -= c0
                st -= s0
            ans = max(ans, st / ct)
        return ans
