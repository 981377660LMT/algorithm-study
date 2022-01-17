n, m = list(map(int, input().split()))

prices = []
for i in range(n):
    base, weight = list(map(int, input().split()))
    prices.append((base + weight * 2, i + 1))

prices.sort(key=lambda x: (-x[0], x[1]))
res = sorted([i for _, i in prices[:m]])
print(*res)
