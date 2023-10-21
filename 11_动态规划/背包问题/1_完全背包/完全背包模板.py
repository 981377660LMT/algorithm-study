# 描述
# 你有一个背包，最多能容纳的体积是V。
# 现在有n个物品，每种物品有任意多个，第i个物品的体积为vi ,价值为wi。
# (1)求这个背包至多能装多大价值的物品?
# (2)若背包恰好装满，求至多能装多大价值的物品?

n, capacity = list(map(int, input().split()))
# n, capacity = 3, 5

size, worth = [0 for _ in range(n)], [0 for _ in range(n)]
for i in range(n):
    size[i], worth[i] = map(int, input().split(" "))
# size, worth = [2, 4, 1], [10, 5, 4]

dp1 = [0] * (capacity + 1)
dp2 = [float("-inf")] * (capacity + 1)
dp2[0] = 0

for i in range(n):
    for j in range(capacity + 1):
        if j >= size[i]:
            dp1[j] = max(dp1[j], dp1[j - size[i]] + worth[i])
            dp2[j] = max(dp2[j], dp2[j - size[i]] + worth[i])

# 至多装多大价值
print(dp1[-1])
# 恰好装满，至多装多大价值
print(0 if dp2[-1] == float("-inf") else dp2[-1])
