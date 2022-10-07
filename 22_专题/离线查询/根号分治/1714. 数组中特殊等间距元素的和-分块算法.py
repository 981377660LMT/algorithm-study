from collections import defaultdict
from functools import lru_cache
from math import sqrt
from typing import List

MOD = int(1e9 + 7)
# 1 <= n <= 5 * 1e4
# 其中 queries[i] = [xi, yi]。 第 i 个查询指令的答案是 nums[j] 中满足该条件的所有元素的和：
# xi <= j < n 且 (j - xi) 能被 yi 整除。 (即分段点的和)


# 每个查询要计算nums[start:n:step]的和
# 分块思想:
# 1. 如果step比较大(大于根号n，那么查询只需根号n次运算，没问题)
# 2. 如果step比较小(小于根号n,那么只需要在时间复杂度O(n*根号n)内预处理答案，然后O(1)查询)
class Solution:
    def solve1(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        """每个查询要计算nums[start:n:step]的和
        
        注意分块算法O(n*n^1/2)非常卡常 不适合python写
        """
        n = len(nums)
        sqrt_ = int(sqrt(n))
        memo = [[0] * (sqrt_ + 1) for _ in range(n + 1)]
        for step in range(1, sqrt_ + 1):
            for start in range(n - 1, -1, -1):
                memo[start][step] = memo[min(n, start + step)][step] + nums[start]
                memo[start][step] %= MOD

        res = [0] * len(queries)
        for i, (start, step) in enumerate(queries):
            if step <= sqrt_:
                res[i] = memo[start][step]
            else:
                res[i] = sum(nums[start:n:step]) % MOD
        return res

    def solve(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        """每个查询要计算nums[start:n:step]的和

        记忆化的形式代替分块思想
        注意到不能全部dfs,这样的复杂度是O(n*step)的
        将step限制在根号n以内,那么只需要在时间复杂度O(n*根号n)内预处理答案,然后O(1)查询
        如果step大于根号n,那么只需要在时间复杂度O(根号n)内暴力查询即可

        dfs超时了
        """

        @lru_cache(None)
        def dfs(start: int, step: int) -> int:
            if start >= n:
                return 0
            return (nums[start] + dfs(start + step, step) % MOD) % MOD

        n = len(nums)
        res = [0] * len(queries)
        sqrt_ = int(sqrt(n))
        for i, (start, step) in enumerate(queries):
            if step <= sqrt_:
                res[i] = dfs(start, step)
            else:
                res[i] = sum(nums[start:n:step]) % MOD
        dfs.cache_clear()
        return res


print(Solution().solve1(nums=[0, 1, 2, 3, 4, 5, 6, 7], queries=[[0, 3], [5, 1], [4, 2]]))
# 输出: [9,18,10]
# 解释: 每次查询的答案如下：
# 1) 符合查询条件的索引 j 有 0、 3 和 6。 nums[0] + nums[3] + nums[6] = 9
# 2) 符合查询条件的索引 j 有 5、 6 和 7。 nums[5] + nums[6] + nums[7] = 18
# 3) 符合查询条件的索引 j 有 4 和 6。 nums[4] + nums[6] = 10
