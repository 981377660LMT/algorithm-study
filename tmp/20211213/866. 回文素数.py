# 答案肯定存在，且小于 2 * 10^8。

# 求出大于或等于 N 的最小回文素数。

# 偶数位回文数的 奇位数字之和与偶位数字之和的差 为 0，故偶数位回文数必然是 11 的倍数。
# 只需要考虑一种镜像方式，直接从小到大遍历回文根即可
class Solution:
    def primePalindrome(self, n: int) -> int:
        def isPrime(n):
            return n >= 2 and all(n % i for i in range(2, int(n ** 0.5) + 1))

        if 8 <= n <= 11:
            return 11

        length = len(str(n))

        # 从回文根开始查找优化
        start = int(str(n)[: (length + 1) // 2])
        for half in range(start, 10 ** 6):
            cand = int(str(half) + str(half)[-2::-1])
            if cand >= n and isPrime(cand):
                return cand


print(Solution().primePalindrome(13))
# 输出：101
