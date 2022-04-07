# 每位用户在1分钟内最多访问u次
# 服务器在1分钟内最多接受g次请求
# DDoS Protection

from collections import defaultdict
from typing import List, Tuple


class Solution:
    def DDoSProtection(self, requests: List[List[int]], u: int, g: int) -> int:
        """求成功访问的次数"""
        if u == 0 or g == 0:
            return 0
        userReq = defaultdict(list)
        allReq = []
        for user, time in sorted(requests, key=lambda x: (x[1], x[0])):
            isUserOk = len(userReq[user]) < u or time - userReq[user][-u] >= 60
            isAllOk = len(allReq) < g or time - allReq[-g] >= 60
            if isUserOk and isAllOk:
                userReq[user].append(time)
                allReq.append(time)
        return len(allReq)


print(Solution().DDoSProtection(requests=[[0, 1], [0, 2]], u=1, g=5))
