# 0 ≤ n < 2 ** 31
# Try iterating from backward. Keep track of the largest number found so far and its index.

# 注意到n很小 |s|^3解法


class Solution:
    def solve2(self, n: int) -> int:
        chars = list(str(n))
        res = n
        for i in range(len(chars)):
            for j in range(i + 1, len(chars)):
                if chars[j] > chars[i]:
                    chars[i], chars[j] = chars[j], chars[i]
                    res = max(res, int("".join(chars)))
                    chars[i], chars[j] = chars[j], chars[i]
        return res

    # |S|解法
    def solve(self, n: int) -> int:
        chars = list(str(n))
        big = len(chars) - 1
        swap = [-1, -1]

        for i in range(big, -1, -1):
            if chars[big] > chars[i]:
                swap = [big, i]
            elif chars[big] < chars[i]:
                big = i
        chars[swap[0]], chars[swap[1]] = chars[swap[1]], chars[swap[0]]

        return int("".join(chars))


print(Solution().solve(n=1929))
print(Solution().solve(n=4343))
# 9921
