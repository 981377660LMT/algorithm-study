import itertools

# 请你返回将它变成 0 所需要的步数。 如果当前数字是偶数，你需要把它除以 2 ；否则，减去 1 。
class Solution:
    def numberOfSteps(self, num: int) -> int:
        return num and (bin(num).count('1') + num.bit_length() - 1)

    def numberOfSteps2(self, num: int) -> int:
        for res in itertools.count():
            if not num:
                return res
            num = num - 1 if num & 1 else num >> 1
