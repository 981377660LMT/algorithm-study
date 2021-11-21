import itertools


class Solution:
    def numberOfSteps(self, num: int) -> int:
        return num and (bin(num).count('1') + num.bit_length() - 1)

    def numberOfSteps2(self, num: int) -> int:
        for res in itertools.count():
            if not num:
                return res
            num = num - 1 if num & 1 else num >> 1
