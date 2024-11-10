import sys

N, M = map(int, sys.stdin.readline().split())
X = list(map(int, sys.stdin.readline().split()))
A = list(map(int, sys.stdin.readline().split()))

total_stones = sum(A)
if total_stones != N:
    print(-1)
    exit(0)

positions = dict(zip(X, A))
sorted_positions = sorted(positions.items())

operations = 0
surplus = 0
last_pos = 0

for pos, count in sorted_positions:
    gap = pos - last_pos - 1

    if surplus < gap:
        print(-1)
        exit(0)

    operations += gap * surplus - gap * (gap + 1) // 2
    surplus -= gap

    stones = count + surplus
    if stones < 1:
        print(-1)
        exit(0)

    surplus = stones - 1
    operations += surplus
    last_pos = pos

gap = N - last_pos

if surplus < gap:
    print(-1)
    exit(0)

operations += gap * surplus - gap * (gap + 1) // 2

print(operations)
