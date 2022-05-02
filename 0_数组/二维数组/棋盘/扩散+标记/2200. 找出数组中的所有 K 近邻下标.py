# 给你一个下标从 0 开始的整数数组 nums 和两个整数 key 和 k 。
# K 近邻下标 是 nums 中的一个下标 i ，
# 并满足至少存在一个下标 j 使得 |i - j| <= k 且 nums[j] == key 。
# 以列表形式返回按 递增顺序 排序的所有 K 近邻下标。


from typing import List


class Solution:
    def findKDistantIndices(self, nums: List[int], key: int, k: int) -> List[int]:
        """扩散+标记，碰到相同的守卫(点)就停下，保证每个数只被扫到两次
        """
        n = len(nums)
        visited = [False] * n
        for i, num in enumerate(nums):
            if num != key:
                continue

            visited[i] = True
            # 标记左右
            for di in (-1, 1):
                ni = i + di
                while 0 <= ni < n and abs(ni - i) <= k and nums[ni] != key:  # 碰到相同的守卫(点)就停下
                    visited[ni] = True
                    ni += di

        return [i for i in range(n) if visited[i]]


print(
    Solution().findKDistantIndices(
        [
            734,
            228,
            636,
            204,
            552,
            732,
            686,
            461,
            973,
            874,
            90,
            537,
            939,
            986,
            855,
            387,
            344,
            939,
            552,
            389,
            116,
            93,
            545,
            805,
            572,
            306,
            157,
            899,
            276,
            479,
            337,
            219,
            936,
            416,
            457,
            612,
            795,
            221,
            51,
            363,
            667,
            112,
            686,
            21,
            416,
            264,
            942,
            2,
            127,
            47,
            151,
            277,
            603,
            842,
            586,
            630,
            508,
            147,
            866,
            434,
            973,
            216,
            656,
            413,
            504,
            360,
            990,
            228,
            22,
            368,
            660,
            945,
            99,
            685,
            28,
            725,
            673,
            545,
            918,
            733,
            158,
            254,
            207,
            742,
            705,
            432,
            771,
            578,
            549,
            228,
            766,
            998,
            782,
            757,
            561,
            444,
            426,
            625,
            706,
            946,
        ],
        939,
        34,
    )
)
