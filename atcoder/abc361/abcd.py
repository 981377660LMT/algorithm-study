from collections import deque, Counter
import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# !TODO: 双向bfs模板
if __name__ == "__main__":
    N = int(input())
    S = input() + ".."
    T = input() + ".."

    if Counter(S) != Counter(T):
        print(-1)
        exit(0)

    def toState(v: str) -> Tuple[int, ...]:
        return tuple(0 if c == "B" else 1 if c == "W" else 2 for c in v)

    queue1 = set([toState(S)])
    queue2 = set([toState(T)])
    visited = set()
    steps = 0
    while queue1 and queue2:
        if len(queue1) > len(queue2):
            queue1, queue2 = queue2, queue1
        nextQueue = set()
        for state in queue1:
            if state in queue2:
                print(steps)
                exit(0)
            if state in visited:
                continue
            visited.add(state)
            for i in range(N + 1):
                if state[i] == 2 and state[i + 1] == 2:
                    for j in range(N + 1):
                        if i == j or i + 1 == j or j + 1 == i:
                            continue
                        newState = list(state)
                        tmpI1, tmpI2 = newState[i], newState[i + 1]
                        tmpJ1, tmpJ2 = newState[j], newState[j + 1]
                        newState[i], newState[i + 1] = tmpJ1, tmpJ2
                        newState[j], newState[j + 1] = tmpI1, tmpI2
                        newState = tuple(newState)
                        nextQueue.add(newState)

        steps += 1
        queue1 = queue2
        queue2 = nextQueue

    print(-1)
