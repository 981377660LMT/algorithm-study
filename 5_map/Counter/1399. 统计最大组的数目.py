from collections import Counter

# 返回数字数目并列最多的组有多少个。


class Solution:
    def countLargestGroup(self, n: int) -> int:
        nums = [sum(int(char) for char in str(num)) for num in range(1, n + 1)]
        freqCounter = Counter(Counter(nums).values())
        return max(freqCounter.items())[1]


print(Solution().countLargestGroup(13))


# 输入：n = 13
# 输出：4
# 解释：总共有 9 个组，将 1 到 13 按数位求和后这些组分别是：
# [1,10]，[2,11]，[3,12]，[4,13]，[5]，[6]，[7]，[8]，[9]。总共有 4 个组拥有的数字并列最多。

