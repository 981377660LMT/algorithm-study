# from functools import lru_cache
# import sys

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# MOD = 998244353
# INF = int(4e18)
# # 長さ N の正整数列 A=(A
# # 1
# # ​
# #  ,A
# # 2
# # ​
# #  ,…,A
# # N
# # ​
# #  ) および整数 x,y が与えられます。
# # 次の条件をすべて満たすように、xy 座標平面上に N+1 個の点 p
# # 1
# # ​
# #  ,p
# # 2
# # ​
# #  ,…,p
# # N
# # ​
# #  ,p
# # N+1
# # ​
# #   を配置することができるか判定してください。(同じ座標に 2 個以上の点を配置してもよいです。)

# # p
# # 1
# # ​
# #  =(0,0)
# # p
# # 2
# # ​
# #  =(A
# # 1
# # ​
# #  ,0)
# # p
# # N+1
# # ​
# #  =(x,y)
# # 点 p
# # i
# # ​
# #   と点 p
# # i+1
# # ​
# #   の距離は A
# # i
# # ​
# #   (1≤i≤N)
# # 線分 p
# # i
# # ​
# #  p
# # i+1
# # ​
# #   と線分 p
# # i+1
# # ​
# #  p
# # i+2
# # ​
# #   のなす角は 90 度 (1≤i≤N−1)
# # p
# # 1
# # ​
# #  =(0,0)
# # p
# # 2
# # ​
# #  =(A
# # 1
# # ​
# #  ,0)
# # p
# # N+1
# # ​
# #  =(x,y)
# # 点 p
# # i
# # ​
# #   と点 p
# # i+1
# # ​
# #   の距離は A
# # i
# # ​
# #   (1≤i≤N)
# # 線分 p
# # i
# # ​
# #  p
# # i+1
# # ​
# #   と線分 p
# # i+1
# # ​
# #  p
# # i+2
# # ​
# #   のなす角は 90 度 (1≤i≤N−1)
# DIR4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]
# if __name__ == "__main__":
#     n, x, y = map(int, input().split())
#     nums = list(map(int, input().split()))

#     # @lru_cache(None)
#     # def dfs(index: int, row: int, col: int, horizol: bool) -> bool:
#     #     if index == n:
#     #         return row == y and col == x
#     #     if horizol:  # 横向き
#     #         nextRow1, nextCol1 = row + nums[index], col
#     #         if dfs(index + 1, nextRow1, nextCol1, False):
#     #             return True
#     #         nextRow2, nextCol2 = row - nums[index], col
#     #         if dfs(index + 1, nextRow2, nextCol2, False):
#     #             return True
#     #         return False
#     #     else:
#     #         nextRow1, nextCol1 = row, col + nums[index]
#     #         if dfs(index + 1, nextRow1, nextCol1, True):
#     #             return True
#     #         nextRow2, nextCol2 = row, col - nums[index]
#     #         if dfs(index + 1, nextRow2, nextCol2, True):
#     #             return True
#     #         return False

#     # res = dfs(1, 0, nums[0], True)
#     # print("Yes" if res else "No")
#     # dp = [set(), set([(nums[0], 0)])]  # vertial , horizol
#     # for i in range(1, n):
#     #     ndp = [set(), set()]
#     #     for row, col in dp[0]:
#     #         ndp[1].add((row + nums[i], col))
#     #         ndp[1].add((row - nums[i], col))
#     #     for row, col in dp[1]:
#     #         ndp[0].add((row, col + nums[i]))
#     #         ndp[0].add((row, col - nums[i]))
#     #     dp = ndp

#     # print("Yes" if ((x, y) in dp[0] or (x, y) in dp[1]) else "No")

#     # !背包,两个维度
#     goods1 = [nums[i] for i in range(n) if i % 2 == 0]  # 横向
#     goods2 = [nums[i] for i in range(n) if i % 2 == 1]  # 竖向
#     goods1 = goods1[1:]

#     # 横向是否可以
#     dp1 = set([nums[0]])
#     for i in range(len(goods1)):
#         ndp1 = set()
#         for j in dp1:
#             ndp1.add(j + goods1[i])
#             ndp1.add(j - goods1[i])
#         dp1 = ndp1

#     if x not in dp1:
#         print("No")
#         exit(0)

#     dp2 = set([0])
#     for i in range(len(goods2)):
#         ndp2 = set()
#         for j in dp2:
#             ndp2.add(j + goods2[i])
#             ndp2.add(j - goods2[i])
#         dp2 = ndp2

#     if y not in dp2:
#         print("No")
#         exit(0)

#     print("Yes")
print(1 << 10000)
