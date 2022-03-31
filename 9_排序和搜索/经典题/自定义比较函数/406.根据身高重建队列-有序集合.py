from typing import List
from sortedcontainers import SortedList

# 根据身高体重重建队列
class Solution:
    def solve1(self, queue: List[List[int]]):
        # 高的先发言
        queue.sort(key=lambda x: (-x[0], x[1]))
        res = []
        for p in queue:
            # 要想使sortedList支持插入到某个位置，需要用id
            res.insert(p[1], p)
        return res

    def solve(self, queue: List[List[int]]):
        # 最矮且最惨的先发言
        queue.sort(key=lambda x: (x[0], -x[1]))
        n = len(queue)
        res = [None] * n
        sortedList = SortedList(range(n))
        for height, count in queue:
            curIndex = sortedList[count]
            res[curIndex] = (height, count)
            sortedList.pop(count)
        return res


print(Solution().solve(queue=[[1, 1], [3, 0], [4, 0]]))
