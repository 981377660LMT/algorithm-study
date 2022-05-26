from typing import List


class Solution:
    def readBinaryWatch(self, num: int) -> List[str]:
        return [
            str(a) + ":" + str(b).rjust(2, '0')
            for a in range(12)
            for b in range(60)
            if (bin(a) + bin(b)).count('1') == num
        ]


print(Solution().readBinaryWatch(1))
print(bin(12).count('1'))
print(int(12).bit_count())
