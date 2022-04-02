# 求最长的前缀，使得移除一个元素之后，剩下所有元素freq相等

# during iteration
# 1. decrement the biggest count or
# 2. decrement the smallest count
# 移除最多或者最少的


from collections import defaultdict


class Solution:
    def solve(self, nums):
        res = 0
        maxFreq = 0
        counter = defaultdict(int)

        for i, num in enumerate(nums):
            counter[num] += 1
            if counter[num] > maxFreq:
                maxFreq = counter[num]

            charTypes = len(counter)
            if i == (maxFreq - 1) * charTypes:
                res = i + 1
            if i == maxFreq * (charTypes - 1):
                res = i + 1

        return res


print(Solution().solve([1, 1, 1, 2, 2, 3]))
