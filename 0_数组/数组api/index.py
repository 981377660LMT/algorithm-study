# python 的 list 没有 lastIndexOf 方法
# 可以用 next 代替

print("as".rfind("s"))
print("as".rindex("s"))

nums = [1, 7, 3, 4, 5, 6, 7, 8, 9, 10]
print([121, 11, 2].index(11))
print(next((i for i in range(len(nums) - 1, -1, -1) if nums[i] == 7)))  # next 代替 lastIndexOf
