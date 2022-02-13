from collections import defaultdict


# tree直接用dict 省去离散化步骤
class BIT:
    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

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

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index


#
# 求解合法的(i,j)对的数量
# @param n int整型 n个人
# @param m int整型 m个窗口
# @param a int整型一维数组 长度为n的vector,顺序表示1-n号客人的办理业务所需时间
# @return long长整型
#
from heapq import heappop, heappush


class Solution:
    def getNumValidPairs(self, n: int, m: int, need: list[int]) -> int:
        """银行预约排队模型"""
        window = [0] * m
        finish = [0] * n
        for i in range(n):
            cost = need[i]
            start = heappop(window)
            end = start + cost
            finish[i] = end
            heappush(window, end)

        # 求finish逆序对数
        finish = finish[::-1]
        bit = BIT(int(1e12))
        res = 0
        for num in finish:
            bit.add(num, 1)
            res += bit.query(num - 1)
        return res


print(Solution().getNumValidPairs(5, 2, [1, 3, 2, 5, 4]))
