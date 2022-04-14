from typing import List
from collections import defaultdict

# 根据身高体重重建队列


class BIT:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class Solution:
    def reconstructQueue(self, people: List[List[int]]) -> List[List[int]]:
        n = len(people)
        res = [[] for _ in range(n)]
        bit = BIT(n + 10)
        for i in range(n):
            bit.add(i + 1, 1)

        # 从矮到高看
        people.sort(key=lambda x: (x[0], -x[1]))
        for height, preCount in people:
            left, right = 1, n
            while left <= right:
                mid = (left + right) // 2
                if bit.query(mid) >= preCount + 1:
                    right = mid - 1
                else:
                    left = mid + 1

            res[left - 1] = [height, preCount]
            bit.add(left, -1)

        return res


print(Solution().reconstructQueue(people=[[7, 0], [4, 4], [7, 1], [5, 0], [6, 1], [5, 2]]))
# 输出：[[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]
