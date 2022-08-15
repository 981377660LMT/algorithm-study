nums = [1, 2, 3, 4, 5]

# 1. all + check 函数
print(all(num > 0 for num in nums))

# 2. for ... break + else (适用于all函数太长的情况)
for num in nums:
    if num <= 0:
        break
else:
    print("all positive")


# !注意break/continue只能打破一层循环，写all的逻辑需要用 `for else`/`或者用一个标志位flag`
# 其实all里可以传check函数的
