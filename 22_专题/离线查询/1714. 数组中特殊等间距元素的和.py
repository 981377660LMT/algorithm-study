from functools import lru_cache
from typing import List

MOD = int(1e9 + 7)
# 1 <= n <= 5 * 104
# 其中 queries[i] = [xi, yi]。 第 i 个查询指令的答案是 nums[j] 中满足该条件的所有元素的和：
# xi <= j < n 且 (j - xi) 能被 yi 整除。 (即分段点的和)

# 离线查询：先map后排序再处理，优先处理范围小的query，后面再扩大范围时直接去cache里取


class Solution:
    def solve(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        m, n = len(nums), len(queries)
        queries = sorted(([x, y, i] for i, (x, y) in enumerate(queries)), reverse=True)
        res = [0] * n
        # 保存前一次的start 和 res ,key为start%step 与step
        memo = dict()
        for start, step, qi in queries:
            preStart, preRes = memo.get((start % step, step), (m, 0))
            curRes = (sum(nums[start:preStart:step]) + preRes) % MOD
            res[qi] = curRes
            memo[(start % step, step)] = (start, curRes)
        return res


print(Solution().solve(nums=[0, 1, 2, 3, 4, 5, 6, 7], queries=[[0, 3], [5, 1], [4, 2]]))
# 输出: [9,18,10]
# 解释: 每次查询的答案如下：
# 1) 符合查询条件的索引 j 有 0、 3 和 6。 nums[0] + nums[3] + nums[6] = 9
# 2) 符合查询条件的索引 j 有 5、 6 和 7。 nums[5] + nums[6] + nums[7] = 18
# 3) 符合查询条件的索引 j 有 4 和 6。 nums[4] + nums[6] = 10
