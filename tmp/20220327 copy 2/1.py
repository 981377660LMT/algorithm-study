from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def convertTime(self, current: str, correct: str) -> int:
        def toSecond(time: str) -> int:
            h, m = map(int, time.split(':'))
            return h * 60 + m

        res = 0
        num1, num2 = toSecond(current), toSecond(correct)
        diff = num2 - num1

        for time in (60, 15, 5, 1):
            res += diff // time
            diff %= time

        return res


print(Solution().convertTime(current="02:30", correct="04:35"))
