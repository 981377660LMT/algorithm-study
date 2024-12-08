

n, x, k = map(int, input().split())
n = int(n)
x = int(x)
k = int(k)
items_per_color = dict()
for _ in range(n):
    p, u, c = map(int, input().split())
    if c not in items_per_color:
        items_per_color[c] = []
    items_per_color[c].append((p, u))

dp = [float("-inf")] * (x + 1)
dp[0] = 0

for items in items_per_color.values():
    dp_prev = dp[:]
    for w in range(x + 1):
        dp[w] = dp_prev[w]
    for price, utility in items:
        for w in range(x, price - 1, -1):
            if dp_prev[w - price] != float("-inf"):
                value = dp_prev[w - price] + utility + k
                if dp[w] < value:
                    dp[w] = value

result = max(dp)
print(int(result))
