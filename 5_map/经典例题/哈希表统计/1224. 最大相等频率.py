from typing import List
from collections import Counter

# 请你帮忙从该数组中找出能满足下面要求的 最长 前缀
# 从前缀中 删除一个 元素后，使得所剩下的每个数字的出现次数相同。


class Solution:
    def maxEqualFreq(self, nums: List[int]) -> int:
        # 便于循环时处理不删的情况
        nums.append(0x3F3F3F3F)
        counter = Counter(nums)
        for i in range(len(nums) - 1, -1, -1):
            counter[nums[i]] -= 1
            # 除去频率为0的数字
            freq = sorted(filter(lambda f: f > 0, counter.values()))
            print(freq)

            # 三种情况

            # 一个数
            if len(freq) == 1:
                return i

            # [1, 2, 2, 2]
            if freq[0] == 1 and freq[1] == freq[-1]:
                return i

            # [2, 2, 2, 3]
            if freq[0] == freq[-2] == freq[-1] - 1:
                return i


print(Solution().maxEqualFreq(nums=[2, 2, 1, 1, 5, 3, 3, 5]))
# 输出：7
# 解释：对于长度为 7 的子数组 [2,2,1,1,5,3,3]，如果我们从中删去 nums[4]=5，
# 就可以得到 [2,2,1,1,3,3]，里面每个数字都出现了两次。
