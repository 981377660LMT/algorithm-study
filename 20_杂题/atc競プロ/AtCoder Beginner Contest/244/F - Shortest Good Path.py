"""
一条长度为k路径表示 A1,A2,A3.... ..Ak,需要点Ai到点Ai+1有一条直接相连的边。
我们可以用一个01字符串来表示路径 S=s1s2s3...sn
如果si=0 那么路径中i出现了偶数次 如果si=1 那么路径中i出现了奇数次 (即每次需要将状态flip/异或)
可以看到一个01字符串能表示多条路径。
现在需要你求出,`所有的01字符串`所代表的路径中,最短的那一条路径的长度的总和。
"""
# 与えられるグラフは単純かつ連結
# 状压dp
# !n<=17 可以 O(n^2*2^n)
# !总之就是bfs求最短路
# !dp表示字符串奇偶状态为state 结尾为i (每条路径只与字符串奇偶状态和结束点有关)


from collections import defaultdict, deque
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u].add(v)
        adjMap[v].add(u)

    dist = [[int(1e18)] * (1 << n) for _ in range(n)]
    queue = deque()
    for i in range(n):
        dist[i][1 << i] = 1  # 一个点算长度为1
        queue.append((i, 1 << i))

    while queue:
        cur, state = queue.popleft()
        for next in adjMap[cur]:
            nextState = state ^ (1 << next)  # !奇偶性
            cand = dist[cur][state] + 1
            if cand < dist[next][nextState]:
                dist[next][nextState] = cand
                queue.append((next, nextState))

    res = 0
    for state in range(1, 1 << n):
        cur = int(1e18)
        for i in range(n):
            cur = min(cur, dist[i][state])
        res += cur
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
