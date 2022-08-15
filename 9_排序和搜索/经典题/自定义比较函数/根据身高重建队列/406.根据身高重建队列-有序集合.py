from typing import List
from sortedcontainers import SortedList

# 根据身高体重重建队列


class Solution:
    def solve1(self, people: List[List[int]]):
        # !逐个插入:高的先发言
        people.sort(key=lambda x: (-x[0], x[1]))
        res = []
        for p in people:
            # 要想使sortedList支持插入到某个位置，需要用id
            res.insert(p[1], p)
        return res

    def solve(self, people: List[List[int]]):
        # !逐个删除:最矮且最惨的先发言 O(nlogn)
        people.sort(key=lambda x: (x[0], -x[1]))
        n = len(people)
        res = [None] * n
        sortedList = SortedList(range(n))
        for height, count in people:
            curIndex = sortedList[count]
            res[curIndex] = (height, count)
            sortedList.pop(count)
        return res


print(Solution().solve(people=[[1, 1], [3, 0], [4, 0]]))
print(Solution().solve1(people=[[1, 1], [3, 0], [4, 0]]))
