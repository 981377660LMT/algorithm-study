from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def countEven(self, num: int) -> int:
        res = 0
        for n in range(1, num + 1):
            s = list(str(n))
            nums = [int(num) for num in s]
            if sum(nums) % 2 == 0:
                res += 1
        return res
