# a^(p-1)模p与1同余(费马小定理)
# 因此在模p的意义下 a^-1 为 a^(p-2)

print(pow(4, 3 - 1, mod=3))


n = int(input())

for i in range(n):
    a, p = map(int, input().split())
    res = pow(a, p - 2, p)  # 快速幂
    print(res if a % p != 0 else 'impossible')
