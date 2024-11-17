# 2847. 数位之积为定值的最小数字
# https://leetcode.cn/problems/smallest-number-with-given-digit-product/
# 给定一个 正 整数 n，返回一个字符串，表示使其各位数字的乘积等于 n 的 最小正整数。
# 如果不存在这样的数字，则返回 "-1" 。


from collections import defaultdict


class Solution:
    def smallestNumber(self, n: int) -> str:
        """特判n=1之后从大到小贪心取数"""
        if n == 1:
            return "1"
        sb = []
        for p in range(9, 1, -1):
            while n % p == 0:
                sb.append(str(p))
                n //= p
        if n != 1:
            return "-1"
        return "".join(sb[::-1])

    def smallestNumber2(self, n: int) -> str:
        if n == 1:
            return "1"
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
    print(Solution().smallestNumber2(105))
    print(Solution().smallestNumber2(7))
    print(Solution().smallestNumber2(44))
