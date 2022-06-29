# n个整数 选5个数 乘积模p余q的方案数
# n<=100 ai<=1e9
# !注意n^5/5! 不会超时

import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline


n, p, q = map(int, input().split())
nums = list(map(int, input().split()))

# # !1.TLE
# res = 0
# for sub in combinations(nums, 5):
#     mod_ = reduce(lambda x, y: x * y % p, sub, 1)
#     res += mod_ % p == q
# print(res)

# # !2. 普通的枚举+每一步都不让数字超过1<<63-1   842 ms ac
# res = 0
# for i in range(n):
#     for j in range(i + 1, n):
#         for k in range(j + 1, n):
#             for l in range(k + 1, n):
#                 for m in range(l + 1, n):
#                     mod_ = ((nums[i] * nums[j] % p) * (nums[k] * nums[l] % p)) % p * nums[m] % p
#                     res += mod_ == q
# print(res)


# # !3. 2的基础上保存中间结果   489 ms ac
res = 0
for i in range(n):
    m1 = nums[i]
    for j in range(i + 1, n):
        m2 = m1 * nums[j] % p
        for k in range(j + 1, n):
            m3 = m2 * nums[k] % p
            for l in range(k + 1, n):
                m4 = m3 * nums[l] % p
                for m in range(l + 1, n):
                    m5 = m4 * nums[m] % p
                    res += m5 == q
print(res)


# 优化：
# 1. combinations会比直接遍历慢很多
# 2. 计算过程及时取余 不超过 1<<63-1
# 3. 遍历时保存中间计算结果
