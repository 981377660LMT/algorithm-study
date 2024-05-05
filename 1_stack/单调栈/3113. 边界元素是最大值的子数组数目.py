# 3113. 边界元素是最大值的子数组数目
# https://leetcode.cn/problems/find-the-number-of-subarrays-where-boundary-elements-are-maximum/solutions/
#
# 给你一个 正 整数数组 nums 。
# 请你求出 nums 中有多少个子数组，
# !满足子数组中第一个和最后一个元素都是这个子数组中的最大值。
# 树上版本 https://leetcode.cn/problems/number-of-good-paths/

from bisect import bisect_left, bisect_right
from collections import defaultdict
from typing import List
from 每个元素作为最值的影响范围 import getRange


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def numberOfSubarrays(self, nums: List[int]) -> int:
        """
        第一个数和最后一个数都是区间最大值.
        对每个左端点, 求出它可以作为最大值的区间范围(左严格右非严格), 然后在这个范围里看有哪些右端点合法.
        合法的右端点用区间求频率检测.
        """
        ranges = getRange(nums, isMax=True, isLeftStrict=True, isRightStrict=False)
        mp = defaultdict(list)
        for i, v in enumerate(nums):
            mp[v].append(i)

        def query(start: int, end: int, value: int) -> int:
            return bisect_left(mp[value], end) - bisect_left(mp[value], start)

        res = 0
        for i, (_, right) in enumerate(ranges):
            cur = nums[i]
            res += query(i, right + 1, cur)  # 左端点开始
        return res

    def numberOfSubarrays2(self, nums: List[int]) -> int:
        """
        第一个数是区间最小值，最后一个数是区间最大值.
        维护两个单调栈.
        第一个单调栈求右边第一个>a[i]的下标lg[i].
        然后第二个栈维护的是候选的左端点，淘汰那些不可能作为左端点的下标.
        然后在第二个栈上二分找a[i]对应可能的左边.
        """
        res = 0
        rightBigger = [-1]
        leftCand = []  # 从栈底到栈顶递增的栈, 考虑所有可能的左边界
        for i, v in enumerate(nums):
            while len(rightBigger) > 1 and nums[rightBigger[-1]] <= v:
                rightBigger.pop()
            top = rightBigger[-1]
            rightBigger.append(i)
            res += len(leftCand) - bisect_right(leftCand, top) + 1
            while leftCand and nums[leftCand[-1]] > v:
                leftCand.pop()
            leftCand.append(i)
        return res


if __name__ == "__main__":
    # nums = [1,4,3,3,2]
    print(Solution().numberOfSubarrays([1, 4, 3, 3, 2]))  # 6

    def check2():
        from random import randint, seed

        seed(0)

        def bf(arr: List[int]) -> int:
            n = len(arr)
            res = 0
            for i in range(n):
                for j in range(i, n):
                    first, last = arr[i], arr[j]
                    min_, max_ = min(arr[i : j + 1]), max(arr[i : j + 1])
                    if first == min_ and last == max_:
                        res += 1
            return res

        for _ in range(1000):
            n = randint(1, 100)
            nums = [randint(1, 100) for _ in range(n)]
            res1 = Solution().numberOfSubarrays2(nums)
            res2 = bf(nums)
            if res1 != res2:
                print(nums, res1, res2)
                raise ValueError("error")
        print("check2 pass")

    check2()
