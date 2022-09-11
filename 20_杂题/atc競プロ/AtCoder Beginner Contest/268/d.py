from collections import deque
from itertools import permutations, zip_longest
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
S = []
for _ in range(n):
    S.append(input())
T = set()
for _ in range(m):
    T.add(input())

if n == 1:
    if S[0] in T:
        print(-1)
        exit(0)
    else:
        print(S[0])
        exit(0)

for perm in permutations(S):
    cur = "_".join(perm)
    if len(cur) > 16:
        print(-1)
        exit(0)

    if cur not in T:
        print(cur)
        exit(0)

    visited = set([tuple([0] * (n - 1))])
    queue = deque([(tuple([0] * (n - 1)), 16 - len(cur))])
    while queue:
        state, remain = queue.popleft()
        if remain < 0:
            continue

        cand = []
        for i in range(n - 1):
            cand.append(perm[i] + "_" * state[i])
        cand.append(perm[-1])
        cand = "_".join(cand)

        if 3 <= len(cand) <= 16 and cand not in T:
            print(cand)
            exit(0)

        if remain > 0:
            for i in range(n - 1):
                new_state = list(state)
                new_state[i] += 1
                new_state = tuple(new_state)
                if new_state not in visited:
                    visited.add(new_state)
                    queue.append((new_state, remain - 1))
print(-1)
