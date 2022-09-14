from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 在每次操作中,你可以将当前的数x变成以下六个数中的一个:
# ×*2,×*4，x*8,x/2(如果x被2整除),x/4(如果x被4整除)，x/8(如果x被8整除)
# 例如,如果当前的数x=12，你可以将他变成24、48、96、6、3，你不能将其变成x/8,因为12不能除。
# 现在请问将给定的初始值a通过上述操作变成目标值b需要的最少的操作次数。


def solve(start: int, target: int) -> int:
    if start == target:
        return 0
    elif start > target:
        start, target = target, start
    queue = deque([(start, 0)])
    visited = set([start])
    while queue:
        cur, step = queue.popleft()
        if cur == target:
            return step
        elif cur > target:
            continue

        for next in [cur * 2, cur * 4, cur * 8]:
            if next not in visited and next <= target:
                visited.add(next)
                queue.append((next, step + 1))
    return -1


t = int(input())
for _ in range(t):
    cur, target = map(int, input().split())
    print(solve(cur, target))
