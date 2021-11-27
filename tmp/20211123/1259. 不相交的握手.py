from math import factorial

# 将握手的人之间连线，请你返回连线不会相交的握手方案数。
class Solution:
    def numberOfWays(self, num_people: int) -> int:
        n = num_people // 2
        return (factorial(2 * n) // factorial(n) // factorial(n) // (n + 1)) % (10 ** 9 + 7)


print(Solution().numberOfWays(2))
print(Solution().numberOfWays(4))
print(Solution().numberOfWays(6))
print(Solution().numberOfWays(8))
# 1 2 5 14 42 142  卡特兰数
