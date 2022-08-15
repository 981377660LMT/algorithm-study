# from heapq import heappop, heappush


# # 污染减到至少一半需要多少filter
# def solution(A):
#     sum_ = sum(A)

#     pq = []
#     for num in A:
#         heappush(pq, -num)

#     target = sum_ / 2
#     cur = 0
#     res = 0
#     while cur < target:
#         big = -heappop(pq)
#         cur += big / 2
#         res += 1
#         heappush(pq, -big / 2)
#     return res


# print(solution([10, 10]))
# print(solution([5, 19, 8, 1]))
# print(solution([3, 0, 5]))
# print(solution([1, 1, 1]))

# !RuntimeError: dictionary changed size during iteration
# fraction sum to 1
# X:molecule  Y:denominator
# ways to choose a pair of fractions that sum up to 1

# from collections import defaultdict
# from math import comb, gcd
# from typing import List

# MOD = int(1e9 + 7)


# def solution(X: List[int], Y: List[int]):
#     n = len(X)
#     if n == 0:
#         return 0
#     adjMap = defaultdict(lambda: defaultdict(int))
#     for a, b in zip(X, Y):
#         if a >= b:
#             continue
#         gcd_ = gcd(a, b)
#         a //= gcd_
#         b //= gcd_
#         adjMap[b][a] += 1

#     res = 0
#     for sum_ in adjMap:
#         for num1 in adjMap[sum_]:
#             count1 = adjMap[sum_][num1]
#             num2 = sum_ - num1
#             if num1 == num2:
#                 res += count1 * (count1 - 1)
#             else:
#                 count2 = adjMap[sum_][num2]
#                 res += count1 * count2
#     return (res // 2) % MOD


# print(solution([1, 1, 1], [1, 1, 1]))
# print(solution([1, 2, 3, 1, 2, 12, 8, 4], [5, 10, 15, 2, 4, 15, 10, 5]))


# 选x个数 求最小和 每个


INF = int(1e18)


def solution(A, X, Y):
    preSum = [[0] for _ in range(Y + 10)]
    for i, num in enumerate(A):
        mod_ = i % Y
        preSum[mod_].append(preSum[mod_][-1] + num)

    res = INF
    for p in preSum:
        for start in range(len(p)):
            end = start + X
            if end >= len(p):
                break
            res = min(res, p[end] - p[start])

    return res


# print(solution([4, 2, 5, 4, 3, 5, 1, 4, 2, 7], 3, 2))
