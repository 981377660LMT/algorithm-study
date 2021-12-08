# 你需要用 红，黄，绿 三种颜色之一给每一个格子上色，且确保相邻格子颜色不同（
# 给你网格图的行数 n 。


# 总结：Two pattern for each row, 121 and 123.
# We consider the state of the first row,
# pattern 121: 121, 131, 212, 232, 313, 323.
# pattern 123: 123, 132, 213, 231, 312, 321.
# So we initialize a121 = 6, a123 = 6.

# We consider the next possible for each pattern:
# Patter 121 can be followed by: 212, 213, 232, 312, 313
# Patter 123 can be followed by: 212, 231, 312, 232

# 121 => three 121, two 123
# 123 => two 121, two 123

# So we can write this dynamic programming transform equation:
# b121 = a121 * 3 + a123 * 2
# b123 = a121 * 2 + a123 * 2
class Solution:
    def numOfWays(self, n: int) -> int:
        a123, a121, mod = 6, 6, int(1e9 + 7)
        for _ in range(n - 1):
            a121, a123 = a121 * 3 + a123 * 2, a121 * 2 + a123 * 2
        return (a123 + a121) % mod


print(Solution().numOfWays(n=2))
# 54
