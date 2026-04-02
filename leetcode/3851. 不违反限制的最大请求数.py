# 3851. 不违反限制的最大请求数
# https://leetcode.cn/problems/maximum-requests-without-violating-the-limit/description/
# 给定一个二维整数数组 requests，其中 requests[i] = [useri, timei] 表示 useri 在 timei 进行了一次请求。
# 同时给定两个整数 k 和 window。
# 如果存在一个整数 t，使得某个用户在闭区间 [t, t + window] 内的请求次数严格大于 k，则用户违反了限制。
# 可以丢弃任意数量的请求。
# 返回一个整数，表示没有用户违反限制的可 保留 的 最大 请求数。


from collections import defaultdict, deque
from typing import List


class Solution:
    def maxRequests(self, requests: List[List[int]], k: int, window: int) -> int:
        mp = defaultdict(list)
        for user, time in requests:
            mp[user].append(time)

        res = 0
        for times in mp.values():
            times.sort()
            queue = deque()
            for t in times:
                while queue and t - queue[0] > window:
                    queue.popleft()
                if len(queue) < k:
                    res += 1
                    queue.append(t)
        return res
