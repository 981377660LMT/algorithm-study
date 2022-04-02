# # 任意排序a,b数组
# # 求a[i]>b[i]的最多对数
# class Solution:
#     def solve(self, a, b):
#         a, b = sorted(a), sorted(b)
#         res = 0
#         while a and b:
#             while b and b[-1] >= a[-1]:
#                 b.pop()
#             if b:
#                 res += 1
#                 a.pop()
#                 b.pop()
#         return res

