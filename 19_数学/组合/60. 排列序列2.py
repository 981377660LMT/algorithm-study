from math import factorial

# 给定 n 和 k，返回k他是第几个排列。
class Solution:
    def getRank(self, n: int, k: int) -> int:
        res = 0
        word = str(k)
        for i in range(n):
            afterSmaller = 0
            for j in range(i + 1, n):
                if word[j] < word[i]:
                    afterSmaller += 1

            # 比多少个大
            res += factorial(n - 1 - i) * afterSmaller

        return res + 1


print(Solution().getRank(n=3, k=213))
# n = 3, k = 3
