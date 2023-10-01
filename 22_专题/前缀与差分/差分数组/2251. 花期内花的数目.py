# 给你一个下标从 0 开始的二维整数数组 flowers ，
# 其中 flowers[i] = [starti, endi] 表示第 i 朵花的 花期 从 starti 到 endi （都 包含）。
# 同时给你一个下标从 0 开始大小为 n 的整数数组 people ，people[i] 是第 i 个人来看花的时间。
# 请你返回一个大小为 n 的整数数组 answer ，其中 answer[i]是第 i 个人到达时在花期内花的 数目 。


# https://leetcode.cn/problems/number-of-flowers-in-full-bloom/solutions/1447392/python-san-chong-jing-dian-jie-fa-by-981-5tt2/
# 解题思路
# 离散化+差分
# 离线查询+小根堆
# 动态开点线段树/树状数组


from bisect import bisect_right
from collections import defaultdict
from heapq import heappop, heappush
from itertools import accumulate
from typing import List


class Solution:
    def fullBloomFlowers(self, flowers: List[List[int]], people: List[int]) -> List[int]:
        diff = defaultdict(int)
        for left, right in flowers:
            diff[left] += 1
            diff[right + 1] -= 1

        # !离散化的keys、原数组前缀和
        keys = sorted(diff)
        preSum = [0] + list(accumulate(diff[key] for key in keys))
        return [preSum[bisect_right(keys, p)] for p in people]

    def fullBloomFlowers2(self, flowers: List[List[int]], people: List[int]) -> List[int]:
        queries = sorted([(person, index) for index, person in enumerate(people)])
        flowers = sorted(flowers)

        fi, res = 0, [0] * len(queries)
        pq = []
        for qi in range(len(queries)):
            while fi < len(flowers) and flowers[fi][0] <= queries[qi][0]:
                heappush(pq, flowers[fi][1])
                fi += 1
            while pq and pq[0] < queries[qi][0]:
                heappop(pq)
            res[queries[qi][1]] = len(pq)
        return res
