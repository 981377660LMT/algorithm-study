n = int(input())
nums = list(map(int, input().split()))
up = [0] * n
down = [0] * n
# 数组中求一个最长的山脉的长度，满足山脉中没有两个相邻的相同高度的点，且只有一个峰


# 预处理+查表
for i in range(n):
    for j in range(i):
        if nums[i] > nums[j]:
            up[i] = max(up[i], up[j] + 1)

nums = nums[::-1]
for i in range(n):
    for j in range(i):
        if nums[i] > nums[j]:
            down[i] = max(down[i], down[j] + 1)

# 注意这里
down = down[::-1]

res = []
for i in range(n):
    res.append(down[i] + up[i] + 1)
print(max(res))

