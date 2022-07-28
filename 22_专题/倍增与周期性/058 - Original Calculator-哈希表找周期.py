# n<=1e5,k<=1e18

# !.周期性解法(鸽巢原理)
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e5)


def move(x: int) -> int:
    y = sum(int(char) for char in str(x))
    return (x + y) % MOD


state, k = map(int, input().split())
visited = dict()
while k:
    visited[state] = k
    state = move(state)
    k -= 1
    if state in visited:
        period = visited[state] - k
        k %= period

print(state)

