# 2847. 数位之积为定值的最小数字
# https://leetcode.cn/problems/smallest-number-with-given-digit-product/

from collections import defaultdict


class Solution:
    def smallestNumber(self, n: int) -> str:
        counter = defaultdict(int)
        for f in [2, 3, 5, 7]:
            while n % f == 0:
                n //= f
                counter[f] += 1

        if n > 1:
            return "-1"

        while counter[3] >= 2:
            counter[3] -= 2
            counter[9] += 1
        while counter[2] >= 3:
            counter[2] -= 3
            counter[8] += 1
        while counter[2] >= 1 and counter[3] >= 1:
            counter[2] -= 1
            counter[3] -= 1
            counter[6] += 1
        while counter[2] >= 2:
            counter[2] -= 2
            counter[4] += 1

        return "".join(str(k) * v for k, v in sorted(counter.items()))


if __name__ == "__main__":
    print(Solution().smallestNumber(105))
    print(Solution().smallestNumber(7))
    print(Solution().smallestNumber(44))
