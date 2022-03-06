# 有 n 个人排队到 1 个水龙头处打水，第 i 个人装满水桶所需的时间是 ti，请问如何安排他们的打水顺序才能使所有人的等待时间之和最小？
n = int(input())
costs = sorted(list(map(int, input().split())))
waits = [0] * n

for i in range(1, n):
    waits[i] = costs[i - 1] + waits[i - 1]

print(sum(waits))
