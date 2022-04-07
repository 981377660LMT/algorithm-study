# 求1-n的二进制中1的个数之和
# 关键是不断删除2的n次幂

from functools import lru_cache

# 338. 比特位计数
@lru_cache(None)
def dfs(num: int) -> int:
    if num == 0:
        return 0
    if num == 1:
        return 1
    half = num // 2
    if num & 1:
        return 2 * dfs(half) + half + 1
    else:
        return dfs(half) + dfs(half - 1) + half


class Solution:
    def solve(self, n):
        return dfs(n)


# n ≤ 2 ** 27
print(Solution().solve(n=5))

# 1 到 2^k 里有 2^(k+1) 个1

# i | binary
# ----------
# 1 | 001   1
# 2 | 010   2
# 3 | 011   4
# 4 | 100   5
# 5 | 101   7
# 6 | 110   9
# 7 | 111   12
# 8 | 1000  13
# 9 | 1001  15
# 10 | 1010  17
# 11 | 1011  20
# 12 | 1100  22
# 13 | 1101  25
# 14 | 1110  28

# 0 1 2 3 4 5
# {1,3,5}=2*{0,1,2}+1
# {0,2,4}=2*{0,1,2}
# 因此 f(5)=2*f(2)+3


# 0 1 2 3 4 5 6
# {1,3,5}=2*{0,1,2}+1
# {2,4,6}=2*{1,2,3}
# 因此 f(6)=f(2)+f(3)+3
