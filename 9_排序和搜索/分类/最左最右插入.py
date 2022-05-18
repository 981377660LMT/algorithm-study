import bisect

nums = [3, 2, 2, 4, 3, 2, 5]
i = bisect.bisect_left(nums, 2)  # get 1
j = bisect.bisect_right(nums, 2)  # get 4
# j - i 就是 nums 中 2 的个数

print(i, j)

