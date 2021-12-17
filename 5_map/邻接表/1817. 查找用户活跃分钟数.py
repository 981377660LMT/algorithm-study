from typing import List
from collections import defaultdict


class Solution:
    def findingUsersActiveMinutes(self, logs: List[List[int]], k: int) -> List[int]:
        res = [0] * k
        counter = defaultdict(set)
        for id, time in logs:
            counter[id].add(time)
        for countSet in counter.values():
            res[len(countSet) - 1] += 1
        return res


print(Solution().findingUsersActiveMinutes(logs=[[0, 5], [1, 2], [0, 2], [0, 5], [1, 3]], k=5))
# 输出：[0,2,0,0,0]
# 解释：
# ID=0 的用户执行操作的分钟分别是：5 、2 和 5 。因此，该用户的用户活跃分钟数为 2（分钟 5 只计数一次）
# ID=1 的用户执行操作的分钟分别是：2 和 3 。因此，该用户的用户活跃分钟数为 2
# 2 个用户的用户活跃分钟数都是 2 ，answer[2] 为 2 ，其余 answer[j] 的值都是 0
