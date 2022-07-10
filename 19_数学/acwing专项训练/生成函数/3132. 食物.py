# https://www.acwing.com/solution/content/30182/
# https://www.acwing.com/solution/content/47652/
# https://blog.csdn.net/weixin_45925418/article/details/118465962
# 一堆物品，其中每一个的物品数量有不同的限制，求携带n个物品的搭配数？
# n<=1e500

# 广义二项式定理
# (1-x)^(-n) = ∑C(n+k-1,k)*(x^k)
# 即求出生成函数中x^n项的系数
# x*(1/(1-x))^4 里x^n的系数 为 C(n-1+4-1,4-1)=C(n+2,3)= n*(n+1)*(n+2)/6
MOD = 10007
n = int(input())
a = n % MOD
b = (n + 1) % MOD
c = (n + 2) % MOD
d = pow(6, MOD - 2, MOD)
res = a * b * c * d % MOD
print(res)


# import sys
# from functools import lru_cache

# # sys.setrecursionlimit(int(1e9))


# # @lru_cache(None)
# # def fac(n: int) -> int:
# #     """n的阶乘"""
# #     if n == 0:
# #         return 1
# #     return n * fac(n - 1) % MOD


# # @lru_cache(None)
# # def ifac(n: int) -> int:
# #     """n的阶乘的逆元"""
# #     return pow(fac(n), MOD - 2, MOD)


# # def C(n: int, k: int) -> int:
# def C(n: int, k: int) -> int:
#     if n < 0 or k < 0 or n < k:
#         return 0
#     return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD
