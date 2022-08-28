# from collections import deque
# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# INF = int(1e18)

# a, b = map(int, input().split())


# queue = deque([(a, b, 0)])
# visited = set([(a, b)])
# while queue:
#     num1, num2, curDist = queue.popleft()
#     if num1 % num2 == 0 or num2 % num1 == 0:
#         print(curDist)
#         exit(0)
#     s1, s2 = str(num1), str(num2)

#     if len(s1) > 1:
#         for i in range(len(s1)):
#             nextNum1 = int(s1[:i] + s1[i + 1 :])
#             if (nextNum1, num2) not in visited:
#                 visited.add((nextNum1, num2))
#                 queue.append((nextNum1, num2, curDist + 1))
#     if len(s2) > 1:
#         for i in range(len(s2)):
#             nextNum2 = int(s2[:i] + s2[i + 1 :])
#             if (num1, nextNum2) not in visited:
#                 visited.add((num1, nextNum2))
#                 queue.append((num1, nextNum2, curDist + 1))

# print(-1)

########################################
# from collections import Counter
# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# INF = int(1e18)

# n = int(input())
# nums = list(map(int, input().split()))
# if n == 1:
#     print(0)
#     exit(0)
# if n == 2:
#     print(1 if nums[0] == nums[1] else 0)
#     exit(0)

# evens = [num for i, num in enumerate(nums) if not i & 1]
# odds = [num for i, num in enumerate(nums) if i & 1]
# counter1 = Counter(evens)
# counter2 = Counter(odds)
# max1 = max(evens)
# max2 = max(odds)

# res1 = sum(v * (max1 - k) for k, v in counter1.items())
# res2 = sum(v * (max2 - k) for k, v in counter2.items())
# res = res1 + res2
# if max1 == max2:
#     res += min(len(evens), len(odds))
# print(res)
# # 先加到最大值
# # 再加小的那个
########################################
# # 奇数长度时答案只有两种情况
# # dere循环
# # rede循环

# from itertools import cycle, groupby, zip_longest
# from math import log2
# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# INF = int(1e18)

# TARGET1 = "dere" * int(1e5)
# TARGET2 = "rede" * int(1e5)

# s = input()
# n = len(s)
# if n <= 2:
#     print(0)
#     exit(0)

# if n & 1:
#     res = INF

#     cand = 0
#     for c1, c2 in zip(s, TARGET1):
#         cand += int(c1 != c2)
#     res = min(cand, res)

#     cand = 0
#     for c1, c2 in zip(s, TARGET2):
#         cand += int(c1 != c2)
#     res = min(cand, res)

#     print(res)
# else:
#     res = INF
#     target = (n - 1) // 2  # 这么多个好e
#     # 只能有一个冗余字符
#     # 相同字符分组
#     groups = [(char, len(list(group))) for char, group in groupby(s)]

#     cand = 0
#     for g, c2 in zip_longest(groups, TARGET1[:n], fillvalue=("", 0)):
#         char, count = g
#         if char != c2:
#             cand += count
#     res = min(cand, res)

#     cand = 0
#     for g, c2 in zip_longest(groups, TARGET2[:n], fillvalue=("", 0)):
#         char, count = g
#         if char != c2:
#             cand += count
#     res = min(cand, res)

#     print(max(0, res))
# ########################################
# from collections import defaultdict
# from itertools import combinations
# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# INF = int(1e18)


# n = int(input())
# nums = list(map(int, input().split()))


# adjMap = defaultdict(list)
# for i, num in enumerate(nums):
#     adjMap[num].append(i)

# res = 0
# for key, indexes in adjMap.items():
#     for i, (pre, cur) in enumerate(zip(indexes, indexes[1:])):
#         if cur - pre <= 1:
#             continue
#         left, right = i + 1, len(indexes) - (i + 1)
#         # 求区间[pre+1:cur]中小于key的数的个数 需要某种数据结构优化
#         count = sum(num < key for num in nums[pre + 1 : cur])
#         res += count * left * right
# print(res)


# # test = 0
# # for i, j, k in combinations(range(n), 3):
# #     if nums[i] == nums[k] > nums[j]:
# #         test += 1
# # print(test)
########################################
# #
# # 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
# #
# # 输出中国队可能获得的最高分数。
# # @param americaTeam int整型一维数组 美国队队员的速度列表
# # @param chinaTeam int整型一维数组 中国队队员的速度列表
# # @return int整型
# #
# from bisect import bisect_right
# from collections import defaultdict
# from itertools import permutations
# from typing import List


# # 注意运动员奔跑的速度区间为[1,20]


# class Solution:
#     # def process(self, americaTeam: List[int], chinaTeam: List[int]) -> int:
#     # !田忌赛马不对
#     #     chinaTeam.sort()
#     #     res = 0
#     #     for a in americaTeam:
#     #         pos = bisect_right(chinaTeam, a)
#     #         cur = 0
#     #         if pos == len(chinaTeam):
#     #             cur = chinaTeam.pop(0)
#     #         else:
#     #             cur = chinaTeam.pop(pos)
#     #         if a == cur:
#     #             res += 1
#     #         elif a < cur:
#     #             res += 3

#     #     return res

#     def process(self, americaTeam: List[int], chinaTeam: List[int]) -> int:
#         res = 0
#         for perm in permutations(chinaTeam):
#             cur = 0
#             for a, b in zip(americaTeam, perm):
#                 if a == b:
#                     cur += 1
#                 elif a < b:
#                     cur += 3
#             res = max(res, cur)
#         return res
########################################
#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 能完成组队返回1，不能完成组队返回0
# @param colors string字符串 T恤颜色字符串
# @return int整型
#

# R必须与B组队 且R在B前
# Y自由插入
class Solution:
    def processData(self, colors: str) -> int:
        while True:
            cur = colors
            cur = cur.replace("RB", "")
            cur = cur.replace("YB", "")
            cur = cur.replace("RY", "")
            if cur == colors:
                break
            colors = cur
        colors = colors.replace("Y", "")
        if not colors:
            return 1

        return 0
