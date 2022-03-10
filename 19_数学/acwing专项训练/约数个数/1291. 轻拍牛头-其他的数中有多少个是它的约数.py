# 共有 N 个整数 A1,A2,…,AN，对于每一个数 Ai，求其他的数中有多少个是它的约数。

# 因数筛

# 这题会卡counter 需要用数组计数

MAX = int(1e6 + 10)
n = int(input())
nums = [0] * n
counter = [0] * MAX
multiCounter = [0] * MAX
for i in range(n):
    nums[i] = int(input())
    counter[nums[i]] += 1

for factor in range(1, MAX):
    if not counter[factor]:
        continue
    for multi in range(factor, MAX, factor):
        multiCounter[multi] += counter[factor]

for num in nums:
    print(multiCounter[num] - 1)

