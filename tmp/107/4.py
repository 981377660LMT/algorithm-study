from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n ，表示服务器的总数目，再给你一个下标从 0 开始的 二维 整数数组 logs ，其中 logs[i] = [server_id, time] 表示 id 为 server_id 的服务器在 time 时收到了一个请求。

# 同时给你一个整数 x 和一个下标从 0 开始的整数数组 queries  。

# 请你返回一个长度等于 queries.length 的数组 arr ，其中 arr[i] 表示在时间区间 [queries[i] - x, queries[i]] 内没有收到请求的服务器数目。


# 注意时间区间是个闭区间。


# !对滑动窗口更加本质的理解


class Solution:
    def countServers(self, n: int, logs: List[List[int]], x: int, queries: List[int]) -> List[int]:
        #  滑动窗口
        record = defaultdict(list)  # 记录每个时间点的服务器
        for serverId, time in logs:
            record[time].append(serverId)
        keys = sorted(record)
        qWithId = sorted([(time, i) for i, time in enumerate(queries)])
        counter = defaultdict(int)
        res = [0] * len(queries)
        left = 0
        right = 0
        for time, qId in qWithId:
            # add
            start, end = time - x, time
            while right < len(keys) and keys[right] <= end:
                for serverId in record[keys[right]]:
                    counter[serverId] += 1
                right += 1
            # remove
            while left < len(keys) and keys[left] < start:
                for serverId in record[keys[left]]:
                    counter[serverId] -= 1
                    if counter[serverId] == 0:
                        counter.pop(serverId)
                left += 1
            res[qId] = n - len(counter)
        return res


# [2,2,3,1,2,2,1,3,2,2]
# 3
# [[1,35],[1,32],[1,11],[1,39],[2,29]]
# 8
# [38,30,23,33,15,31,34,22,11,14]
print(
    Solution().countServers(
        3,
        [[1, 35], [1, 32], [1, 11], [1, 39], [2, 29]],
        8,
        [38, 30, 23, 33, 15, 31, 34, 22, 11, 14],
    )
)
