# 给定一个正整数，找出与其二进制表达式中1的个数相同且大小最接近的那两个数（一个略大，一个略小）。
# num的范围在[1, 2147483647]之间；
# 如果找不到前一个或者后一个满足条件的正数，那么输出 -1。
from typing import List
from nextPermutation import nextPermutation, prePermutation


class Solution:
    def findClosedNumbers(self, num: int) -> List[int]:
        bin_ = list(bin(num)[2:].zfill(32))

        ok1, nextPerm = nextPermutation(bin_)
        cand1 = int("".join(nextPerm), 2) if ok1 else -1
        if cand1 > 1 << 31:
            cand1 = -1

        ok2, prePerm = prePermutation(bin_)
        cand2 = int("".join(prePerm), 2) if ok2 else -1
        if cand2 > 1 << 31:
            cand2 = -1

        return [cand1, cand2]
