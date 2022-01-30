MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


# 注意这题不是rabin-karp哈希模板
# 这题需要反向滑窗
class Solution:
    def subStrHash(self, s: str, power: int, modulo: int, k: int, hashValue: int) -> str:
        curSum = sum((ord(s[i]) - 96) * (power ** i) for i in range(k)) % modulo
        if curSum == hashValue:
            return s[:k]
        print(curSum)
        for right in range(k, len(s)):
            curSum -= ord(s[right - k]) - 96
            curSum *= pow(power, -1, modulo)
            curSum += (ord(s[right]) - 96) * power ** (k - 1)
            curSum %= modulo
            if curSum == hashValue:
                return s[right - k + 1 : right + 1]


print(Solution().subStrHash(s="leetcode", power=7, modulo=20, k=2, hashValue=0))
# print(Solution().subStrHash(s="fbxzaad", power=31, modulo=100, k=3, hashValue=32))

