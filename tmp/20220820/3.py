from typing import List
from itertools import accumulate
from string import ascii_lowercase


class Solution:
    def shiftingLetters(self, s: str, shifts: List[List[int]]) -> str:
        diff = [0] * (len(s) + 1)
        for left, right, di in shifts:
            diff[left] += 1 if di == 1 else -1
            diff[right + 1] -= 1 if di == 1 else -1
        diff = list(accumulate(diff))
        res = []
        for char, offset in zip(s, diff):
            index = (ord(char) - 97 + offset) % 26
            res.append(ascii_lowercase[index])
        return "".join(res)


print(Solution().shiftingLetters(s="dztz", shifts=[[0, 0, 0], [1, 1, 1]]))
