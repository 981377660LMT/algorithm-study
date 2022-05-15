class Solution:
    def solve(self, start: int, end: int) -> int:
        """对每一个位考虑,从31到0"""
        diff = end - start + 1
        res = 0
        for bit in range(31, -1, -1):
            # 覆盖不了这个范围
            if (1 << bit) < diff:
                break
            if (1 << bit) & start & end:
                res += 1 << bit

        return res


print(Solution().solve(5, 7))

# 0101 = 5
# 0110 = 6
# 0111 = 7
# ----
# 0100 = 4
