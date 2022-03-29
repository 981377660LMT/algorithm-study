from typing import List
from collections import Counter
from functools import lru_cache


# 如果 nums 的一个子集中，所有元素的乘积可以用若干个 `互不相同的质数` 相乘得到

# 那么我们称它为 好子集 。
# 1 <= nums.length <= 105
# 1 <= nums[i] <= 30，小于等于30的质数正好是10个，暗示状压

# 每个质数p只能在好子集中出现0或1次，对应着选或不选
# 遍历可能作为好子集元素的数(一定是0-30的素数)，乘积做组合，求每种情况的频率
# 看到限制里面数字最大也不会超过 30 ，立刻想到暴力。把`组合全部弄出来`，每种组合的次数就是数字出现次数的乘积。
# 需要小心的是，[1, 1, 1] 这种是不算的，乘起来要大于 1 的才算

# 预处理
primes = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29]
composites = set([4, 8, 9, 12, 16, 18, 20, 24, 25, 27, 28])

# 0-30每个数包含的质数
factors = [sum(1 << i for i, p in enumerate(primes) if x % p == 0) for x in range(31)]
MOD = 10 ** 9 + 7
available = (1 << len(primes)) - 1


class Solution:
    def numberOfGoodSubsets(self, nums: List[int]) -> int:
        freq = Counter(nums)

        @lru_cache(None)
        def dfs(cur: int, state: int) -> int:
            if cur == 1:
                return 1

            # 不取当前
            res = dfs(cur - 1, state)

            # 取当前；交小并大
            if cur not in composites and state | factors[cur] == state:
                res += freq[cur] * dfs(cur - 1, state ^ (factors[cur]))

            return res

        # 减1表示减去空集；答案为可选的数目*1的子集数
        return (dfs(30, available) - 1) * pow(2, freq[1], MOD) % MOD


print(Solution().numberOfGoodSubsets(nums=[1, 2, 3, 4]))
# 输出：6
# 解释：好子集为：
# - [1,2]：乘积为 2 ，可以表示为质数 2 的乘积。
# - [1,2,3]：乘积为 6 ，可以表示为互不相同的质数 2 和 3 的乘积。
# - [1,3]：乘积为 3 ，可以表示为质数 3 的乘积。
# - [2]：乘积为 2 ，可以表示为质数 2 的乘积。
# - [2,3]：乘积为 6 ，可以表示为互不相同的质数 2 和 3 的乘积。
# - [3]：乘积为 3 ，可以表示为质数 3 的乘积。
