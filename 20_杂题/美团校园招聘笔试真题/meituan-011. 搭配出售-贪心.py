from operator import ne


a, b, c, d, e, f, g = list(map(int, input().split()))

plans = [(a, e), (b, f), (c, g)]
plans.sort(key=lambda x: -x[1])

res = 0
for need, price in plans:
    have = min(d, need)
    res += have * price
    d -= have
    if d == 0:
        break

print(res)
