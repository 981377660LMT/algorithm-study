# !bad
# !不要在遍历时改变数组长度

nums = [1, 2, 3, 4, 5, 6]
for num in nums:
    if num % 3 != 0:
        nums.remove(num)
print(nums)  # [2, 3, 5, 6]

# nums = [1, 2, 3, 4, 5, 6]
# for i in range(len(nums)):
#     if nums[i] % 3 != 0:  # !IndexError: list index out of range
#         nums[i : i + 1] = []
# print(nums)
#############################################
# Solution1 遍历拷贝 (如果需要修改原数组的话)
nums = [1, 2, 3, 4, 5, 6]
for num in nums[:]:
    if num % 3 != 0:
        nums.remove(num)
print(nums)  # [3, 6]

# Solution2 使用filter/标记删除 (如果不需要修改原数组的话)
nums = [1, 2, 3, 4, 5, 6]
nums = [num for num in nums if num % 3 == 0]
print(nums)  # [3, 6]
#############################################
# SortedList替换i处的元素
from sortedcontainers import SortedList

nums = SortedList([1, 2, 3, 4, 5, 6])
for i, num in enumerate(nums):
    if num % 3 != 0:
        # nums[i] = 0  # !NotImplementedError: use ``del sl[index]`` and ``sl.add(value)`` instead
        nums.remove(num)
        nums.add(0)
print(nums)
