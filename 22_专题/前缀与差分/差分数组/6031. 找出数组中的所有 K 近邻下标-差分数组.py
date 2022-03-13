from typing import List

MOD = int(1e9 + 7)


class Solution:
    def findKDistantIndices(self, nums: List[int], key: int, k: int) -> List[int]:
        """
        K 近邻下标 是 nums 中的一个下标 i ，
        并满足至少存在一个下标 j 使得 |i - j| <= k 且 nums[j] == key 。

        因为要范围update 所以用差分数组
        """
        n = len(nums)
        diff = [0] * (n + 10)
        for i, num in enumerate(nums):
            if num != key:
                continue
            left, right = max(0, i - k), min(n - 1, i + k)
            diff[left] += 1
            diff[right + 1] -= 1
        for i in range(1, len(diff)):
            diff[i] += diff[i - 1]

        return [i for i, num in enumerate(diff) if num]


print(Solution().findKDistantIndices(nums=[3, 4, 9, 1, 3, 9, 5], key=9, k=1))
print(
    Solution().findKDistantIndices(
        nums=[
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
        key=939,
        k=34,
    )
)

