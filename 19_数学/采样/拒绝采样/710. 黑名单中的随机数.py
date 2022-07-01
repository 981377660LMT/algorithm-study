from random import randint
from typing import List


class Solution:
    def __init__(self, n: int, blacklist: List[int]):
        self.b, self.w = len(blacklist), n - len(blacklist)
        self.black = dict()  # 将前半部分的黑名单映射到后半部分的白名单，然后在前半里随机选

        bad = set(blacklist)
        mex = self.w
        for black in blacklist:
            if black < self.w:
                while mex in bad:
                    mex += 1
                self.black[black] = mex
                mex += 1

    def pick(self) -> int:
        """从 [0, n - 1] 范围内的任意整数中选取一个 未加入 黑名单 blacklist 的整数"""
        rand = randint(0, self.w - 1)  # 在或不在都可以选
        return self.black.get(rand, rand)


# Your Solution object will be instantiated and called as such:
# obj = Solution(n, blacklist)
# param_1 = obj.pick()
