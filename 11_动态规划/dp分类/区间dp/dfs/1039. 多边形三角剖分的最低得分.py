from typing import List
from functools import lru_cache

# 返回多边形进行三角剖分后可以得到的最低分。
# 3 <= A.length <= 50
# f(0,n-1) = f(0,j) + f(j,n-1) + A[0]*A[k]*A[n-1]
class Solution:
    def minScoreTriangulation(self, values: List[int]) -> int:
        @lru_cache(None)
        def dfs(l, r) -> int:
            # 注意此处边界为range(l + 1, r)没有值
            if l + 1 >= r:
                return 0

            res = 0x7FFFFFFF
            for mid in range(l + 1, r):
                res = min(res, dfs(l, mid) + dfs(mid, r) + values[l] * values[mid] * values[r])
            return res

        return dfs(0, len(values) - 1)


print(Solution().minScoreTriangulation([3, 7, 4, 5]))
# 输出：144
# 解释：有两种三角剖分，可能得分分别为：3*7*5 + 4*5*7 = 245，或 3*4*5 + 3*4*7 = 144。最低分数为 144。
