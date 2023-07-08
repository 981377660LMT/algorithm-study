# 2747. 统计没有收到请求的服务器数目
# 给你一个整数 n ，表示服务器的总数目，再给你一个下标从 0 开始的 二维 整数数组 logs ，
# 其中 logs[i] = [server_id, time] 表示 id 为 server_id 的服务器在 time 时收到了一个请求。
# 同时给你一个整数 x 和一个下标从 0 开始的整数数组 queries  。
# 请你返回一个长度等于 queries.length 的数组 arr ，
# 其中 arr[i] 表示在时间区间 [queries[i] - x, queries[i]] 内没有收到请求的服务器数目。
# 注意时间区间是个闭区间。

# 离线查询(排序)+滑动窗口

from collections import defaultdict
from typing import List


class Solution:
    def countServers(self, n: int, logs: List[List[int]], x: int, queries: List[int]) -> List[int]:
        qWithId = [(end, qid) for qid, end in enumerate(queries)]
        qWithId.sort()
        logs.sort(key=lambda x: x[1])
        res = [0] * len(queries)

        ql, qr = 0, 0  # 模拟队列
        counter = defaultdict(int)
        for end, qid in qWithId:
            start = end - x
            while qr < len(logs) and logs[qr][1] <= end:
                id, _ = logs[qr]
                counter[id] += 1
                qr += 1
            while ql < len(logs) and logs[ql][1] < start:
                id, _ = logs[ql]
                counter[id] -= 1
                if counter[id] == 0:
                    del counter[id]
                ql += 1

            res[qid] = n - len(counter)

        return res


# 3
# [[2,4],[2,1],[1,2],[3,1]]
# 2
# [3,4]

print(Solution().countServers(3, [[2, 4], [2, 1], [1, 2], [3, 1]], 2, [3, 4]))
