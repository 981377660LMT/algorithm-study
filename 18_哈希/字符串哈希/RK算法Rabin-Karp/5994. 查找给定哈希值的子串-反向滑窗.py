MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


# 注意这题不是rabin-karp模板，是滑窗滚动哈希
# 1. 因为 pow(power, -1, modulo) 不一定存在(power与mod不一定互质)，所以不能正向滑窗
# 2. python 里求次方优先使用pow函数
# 3. 不记录切片，而是记录切片端点，最后再取切片
class Solution:
    # 错误解法
    def subStrHash1(self, s: str, power: int, modulo: int, k: int, hashValue: int) -> str:
        curSum = sum((ord(s[i]) - 96) * pow(power, i, modulo) for i in range(k)) % modulo
        if curSum == hashValue:
            return s[:k]
        for right in range(k, len(s)):
            curSum -= ord(s[right - k]) - 96
            curSum *= pow(power, -1, modulo)
            curSum += (ord(s[right]) - 96) * pow(power, k - 1, modulo)
            curSum %= modulo
            if curSum == hashValue:
                return s[right - k + 1 : right + 1]

    # pow(power, -1, modulo)正向不行，求反向pow(power, k-1, modulo)
    # 滚动哈希，不可直接用pow计算，否则超时；要用快速幂
    def subStrHash(self, s: str, power: int, modulo: int, k: int, hashValue: int) -> str:
        n = len(s)
        leftCand = 0
        hash = 0
        for i in range(n - 1, n - 1 - k, -1):
            ord_ = ord(s[i]) - 96
            hash = hash * power + ord_
            hash %= modulo

        if hash == hashValue:
            leftCand = n - k

        for left in range(n - k - 1, -1, -1):
            right = left + k
            # 注意不能用一般的幂，要用带mod的快速幂；也可以使用预处理
            hash -= (ord(s[right]) - 96) * (pow(power, k - 1, modulo))
            hash %= modulo
            hash *= power
            hash %= modulo
            hash += ord(s[left]) - 96
            hash %= modulo
            if hash == hashValue:
                leftCand = left

        return s[leftCand : leftCand + k]


print(Solution().subStrHash(s="leetcode", power=7, modulo=20, k=2, hashValue=0))
print(Solution().subStrHash(s="fbxzaad", power=31, modulo=100, k=3, hashValue=32))

