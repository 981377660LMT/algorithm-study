import sys

input = sys.stdin.readline

# !N个硬币是否能凑出Y円
# 硬币只有10000/5000/1000 三种
# 1≤N≤2000
# 1e3≤Y≤2e7

YUKICHI = 10000
KAZUYO = 5000
HIDEO = 1000

N, Y = map(int, input().split())

for i in range(N + 1):
    for j in range(N - i + 1):
        k = N - i - j
        if i * YUKICHI + j * KAZUYO + k * HIDEO == Y:
            print(i, j, k)
            exit()

print(-1, -1, -1)
