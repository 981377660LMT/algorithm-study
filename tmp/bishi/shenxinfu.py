# # # !abs(nums[i]-v)<=x 时可以入侵主机
# # # 需要按顺序入侵主机
# # # n<=1e5 nums[i]<=1e9


# # # 5 4
# # # 4 9 6 7 8


# from typing import Optional, Tuple


# def getIntersect(
#     interval1: Tuple[int, int], interval2: Tuple[int, int]
# ) -> Optional[Tuple[int, int]]:
#     """获取两个区间交集"""
#     if interval1[0] > interval2[1] or interval2[0] > interval1[1]:
#         return None
#     return (max(interval1[0], interval2[0]), min(interval1[1], interval2[1]))


# n, x = map(int, input().split())  # 主机数 主机辨识精度
# nums = list(map(int, input().split()))  # 主机辨识度

# intervals = []
# for num in nums:
#     intervals.append((num - x, num + x))

# res = 0
# curInter = intervals[0]
# for i in range(1, n):
#     nextInter = intervals[i]
#     intersect = getIntersect(curInter, nextInter)
#     if intersect is None:
#         res += 1
#         curInter = nextInter
#     else:
#         curInter = intersect
# print(res)
