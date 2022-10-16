# !根据身高体重重建队列
# !1.树状数组维护一个 01 序列的前缀和
# !2.从矮的往高的考虑，因为矮的去除后不影响高的
# !3.对每个人，二分寻找前面恰好有k个人的位置，即找到 第`k+1`个1 所在的位置

from typing import List


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, 'index must be greater than 0'
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos"""
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos"""
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def __repr__(self) -> str:
        preSum = []
        for i in range(self.size):
            preSum.append(self.query(i))
        return str(preSum)

    def __len__(self) -> int:
        return self.size


class Solution:
    def reconstructQueue(self, people: List[List[int]]) -> List[List[int]]:
        n = len(people)
        res = [[] for _ in range(n)]
        bit = BIT1(n + 10)
        for i in range(n):
            bit.add(i + 1, 1)

        # 从矮到高看
        people.sort(key=lambda x: (x[0], -x[1]))
        # for height, preCount in people:
        #     left, right = 1, n
        #     while left <= right:
        #         mid = (left + right) // 2
        #         if bit.query(mid) >= preCount + 1:
        #             right = mid - 1
        #         else:
        #             left = mid + 1

        #     res[left - 1] = [height, preCount]
        #     bit.add(left, -1)
        for height, preCount in people:
            pos = bit.bisectLeft(preCount + 1)
            res[pos - 1] = [height, preCount]
            bit.add(pos, -1)
        return res


print(Solution().reconstructQueue(people=[[7, 0], [4, 4], [7, 1], [5, 0], [6, 1], [5, 2]]))
# 输出：[[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]
