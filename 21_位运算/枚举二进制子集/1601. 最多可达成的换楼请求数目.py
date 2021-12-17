from typing import List
from functools import lru_cache

# 1 <= n <= 20
# 1 <= requests.length <= 16

# 从请求列表中选出的若干个请求是可行的需要满足 每栋楼员工净变化为 0
# 请你从原请求列表中选出若干个请求，
# 使得它们是一个可行的请求列表，
# 并返回所有可行列表中最大请求数目。
class Solution:
    def maximumRequests(self, n: int, requests: List[List[int]]) -> int:
        def get_bit1(x: int) -> int:
            res = 0
            while x:
                res += 1
                x &= x - 1
            return res

        res = 0
        for state in range(1 << len(requests)):
            buildings = [0] * n
            for i in range(len(requests)):
                if state & (1 << i):
                    u, v = requests[i]
                    buildings[u] -= 1
                    buildings[v] += 1
            if all(v == 0 for v in buildings):
                res = max(res, get_bit1(state))
        return res


print(Solution().maximumRequests(n=5, requests=[[0, 1], [1, 0], [0, 1], [1, 2], [2, 0], [3, 4]]))
