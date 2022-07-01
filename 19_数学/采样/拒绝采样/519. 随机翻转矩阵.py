from random import randint
from typing import List


class Solution:
    """
    随机选取一个满足 matrix[i][j] == 0 的下标 (i, j) ，并将它的值变为 1 。
    所有满足 matrix[i][j] == 0 的下标 (i, j) 被选取的概率应当均等。
    尽量最少调用内置的随机函数，并且优化时间和空间复杂度。
    """

    def __init__(self, m: int, n: int):
        self.black = dict()
        self.row, self.col = m, n
        self.remain = m * n

    def flip(self) -> List[int]:
        """返回一个满足 matrix[i][j] == 0 的随机下标 [i, j] ，并将其对应格子中的值变为 1

        每次flip后 都将flip的值记录到黑名单 并映射到另一个不是黑名单的下标(这里指最后一个数)
        """
        if self.remain == 0:
            return []
        self.remain -= 1
        rand = randint(0, self.remain)
        res = self.black.get(rand, rand)  # rand在或不在都可以选
        self.black[rand] = self.black.get(self.remain, self.remain)
        return [res // self.col, res % self.col]

    def reset(self) -> None:
        """将矩阵中所有的值重置为 0"""
        self.remain = self.row * self.col
        self.black.clear()


# Your Solution object will be instantiated and called as such:
# obj = Solution(m, n)
# param_1 = obj.flip()
# obj.reset()
