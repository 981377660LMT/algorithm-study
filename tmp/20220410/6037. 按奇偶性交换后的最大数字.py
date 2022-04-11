from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def largestInteger(self, num: int) -> int:
        nums = list(map(int, str(num)))
        odds = sorted([num for num in nums if num % 2 == 1])
        evens = sorted([num for num in nums if num % 2 == 0])
        res = []
        for i in range(len(nums)):
            cur = nums[i]
            if cur & 1:
                res.append(str(odds.pop()))
            else:
                res.append(str(evens.pop()))
        return int(''.join(res))
