from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)


class BIT1:
    """单点修改
    
    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

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
        if left > right:
            return 0
        return self.query(right) - self.query(left - 1)


class BookMyShow:
    def __init__(self, n: int, m: int):
        self.row = n
        self.col = m
        self.fullRow = -1  # 前多少行满了 注意只能在scatter里维护 gather不能维护整个
        self.tree = BIT1(n + 10)  # 每一行多少人

    def gather(self, k: int, maxRow: int) -> List[int]:
        if k > self.col:
            return []
        """k 个成员中 第一个座位 的排数和座位编号"""

        def check(mid: int) -> bool:
            """mid行能不能再坐k人"""
            if mid <= self.fullRow:
                return False
            count = self.tree.sumRange(mid + 1, mid + 1)
            return self.col - count >= k

        isOk = False
        left, right = 0, maxRow
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                isOk = True
                right = mid - 1
            else:
                left = mid + 1

        if not isOk:
            return []

        resRow = left  # 行
        resCol = self.tree.sumRange(resRow + 1, resRow + 1)
        self.tree.add(resRow + 1, k)
        # if resCol + k == self.col:  # fullRow有问题
        #     self.fullRow = max(self.fullRow, resRow)
        return [resRow, resCol]

    def scatter(self, k: int, maxRow: int) -> bool:
        """如果组里所有 k 个成员 不一定 要坐在一起的前提下，
        都能在第 0 排到第 maxRow 排之间找到座位，那么请返回 true
        """

        def check(upper: int) -> bool:
            # [1,upper+1]行能否坐下k人  发现这里allOk多算了fullCount
            allOk = (upper + 1 - (self.fullRow + 1)) * self.col
            count = self.tree.sumRange((self.fullRow + 1 + 1), upper + 1)
            return (allOk - count) >= k

        isOk = False
        left, right = 0, maxRow
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                isOk = True
                right = mid - 1
            else:
                left = mid + 1

        if not isOk:
            return False

        upper = left  # 要安排到upper行
        preCount = (upper - self.fullRow - 1) * self.col - self.tree.sumRange(
            self.fullRow + 1 + 1, upper + 1 - 1
        )  # 前面可以坐多少人
        curK = k - preCount  # 当行需要安排的座位数

        curCount = self.tree.sumRange(upper + 1, upper + 1)
        if curCount + curK == self.col:
            self.fullRow = upper
        else:
            self.fullRow = max(self.fullRow, upper - 1)
        self.tree.add(upper + 1, curK)

        return True


# 线段树+树状数组
# Your BookMyShow object will be instantiated and called as such:
# obj = BookMyShow(n, m)
# param_1 = obj.gather(k,maxRow)
# param_2 = obj.scatter(k,maxRow)
# bms = BookMyShow(2, 5)
# print(bms.gather(4, 0))
# print(bms.gather(2, 0))
# print(bms.scatter(5, 1))
# print(bms.scatter(5, 1))

# bms = BookMyShow(4, 5)
# print(bms.scatter(6, 2))
# print(bms.gather(6, 3))
# print(bms.scatter(9, 1))
# [null,true,[],false]

bms = BookMyShow(5, 9)
print(bms.gather(10, 1))
print(bms.scatter(3, 3))
print(bms.gather(9, 1))
print(bms.gather(10, 2))
print(bms.gather(2, 0))
# [null,[],true,[1,0],[],[0,3]]
