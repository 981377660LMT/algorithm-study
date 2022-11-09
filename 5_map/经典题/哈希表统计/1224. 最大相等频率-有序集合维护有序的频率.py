"""最大相等频率

找出最长的前缀 使得
从前缀中 恰好删除一个 元素后，剩下每个数字的出现次数都相同
https://leetcode.cn/problems/maximum-equal-frequency/solution/zui-da-xiang-deng-pin-lu-by-leetcode-sol-5y2m/
"""

# !求最长的前缀，使得移除一个元素之后，剩下所有元素freq相等


from typing import List
from collections import Counter
from sortedcontainers import SortedList


class Solution:
    def maxEqualFreq(self, nums: List[int]) -> int:
        """有序集合维护有序的频率

        during iteration
        1. decrement the biggest count or
        2. decrement the smallest count
        移除最多或者最少的
        """
        n = len(nums)
        counter = Counter(nums)
        freq = SortedList(counter.values())

        for i in range(n - 1, -1, -1):
            # !一种字母
            if len(freq) <= 1:
                return i + 1
            # !删最少的 1 3 3 3
            if freq[0] == 1 and freq[1] == freq[-1]:
                return i + 1
            # !删最多的 3 3 3 4
            if freq[0] == freq[-2] and freq[-1] == freq[-2] + 1:
                return i + 1

            num = nums[i]
            freq.discard(counter[num])
            if counter[num] - 1 > 0:
                freq.add(counter[num] - 1)
            counter[num] -= 1
            if counter[num] == 0:
                del counter[num]

        return 0

    def maxEqualFreq2(self, nums: List[int]) -> int:
        """分三种情况讨论

        1. 最大出现次数 maxFreq == 1 随意删除一个元素
        2. 所有数出现次数都是 maxFreq 或者 maxFreq - 1, 且只有一个数出现次数是 maxFreq
        3. 除开一个数，其他所有数的出现次数都是 maxFreq, 且只有一个数出现次数是 1
        """

        def check(size: int, maxFreq: int, freqCounter: "Counter[int]") -> bool:
            """size个元素, counter是元素出现次数, freqCounter是出现次数的频率"""
            if maxFreq == 1:
                return True
            if (
                freqCounter[maxFreq] * maxFreq + freqCounter[maxFreq - 1] * (maxFreq - 1) == i + 1
                and freqCounter[maxFreq] == 1
            ):
                return True
            if freqCounter[1] == 1 and freqCounter[maxFreq] * maxFreq + 1 == size:
                return True
            return False

        counter, freqCounter = Counter(), Counter()
        res, maxFreq = 0, 0
        for i, num in enumerate(nums):
            counter[num] += 1
            freqCounter[counter[num]] += 1
            if counter[num] != 1:
                freqCounter[counter[num] - 1] -= 1
            maxFreq = max(maxFreq, counter[num])
            if check(i + 1, maxFreq, freqCounter):
                res = i + 1
        return res


print(Solution().maxEqualFreq([1, 1, 1, 2, 2, 3]))
print(Solution().maxEqualFreq(nums=[2, 2, 1, 1, 5, 3, 3, 5]))
print(Solution().maxEqualFreq(nums=[1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5]))
