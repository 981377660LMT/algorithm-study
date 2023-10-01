# bisect 使用案例
from bisect import bisect_left, bisect_right

empty = []
print(bisect_left(empty, 0))
print(bisect_right(empty, 0))

# py3.10以上的版本支持key
print(bisect_left([(1, 2), (2, 3), (3, 4)], 2, key=lambda x: x[0]))


nums = [1, 3, 4, 4, 5, 6, 7, 7, 7, 8, 8, 9, 10]
# 小于等于8的个数
print(bisect_right(nums, 8))
# 小于8的个数
print(bisect_left(nums, 8))

# 大于等于8的个数
print(len(nums) - bisect_left(nums, 8))
# 大于8的个数
print(len(nums) - bisect_right(nums, 8))
