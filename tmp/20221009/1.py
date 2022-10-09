from typing import List


INF = int(1e20)

# 返回处理用时最长的那个任务的员工的 id 。如果存在两个或多个员工同时满足，则返回几人中 最小 的 id 。


class Solution:
    def hardestWorker(self, n: int, logs: List[List[int]]) -> int:
        res, resId = 0, INF
        n = len(logs)
        for i in range(n):
            pre = logs[i - 1][1] if i > 0 else 0
            time = logs[i][1] - pre
            if time > res or (time == res and logs[i][0] < resId):
                res = time
                resId = logs[i][0]
        return resId


print(
    Solution().hardestWorker(70, [[36, 3], [1, 5], [12, 8], [25, 9], [53, 11], [29, 12], [52, 14]])
)
