# ConnectedDp (连通性dp)

from typing import List, Tuple


class ConnectedDpSquares:
    """https://github.com/maspypy"""

    @staticmethod
    def connectedDpGraph(
        n: int, mergeReverse: bool
    ) -> Tuple[List[List[int]], List[Tuple[int, int]]]:
        """返回[states,edges]"""
        mp = dict()
        states = [[-1] * n, [-1] * n]
        edges = []
        mp[tuple(states[0])] = 0
        p = -1
        while True:
            p += 1
            if p == len(states):
                break
            if p == 1:
                edges.append((1, 1))
                continue
            now = states[p]
            for nexts, convert in ConnectedDpSquares._getNextStates(now):
                # 当前的连通分量数，消失的连通分量数
                a, b = 0, 0
                for v in range(n):
                    if now[v] == v:
                        a += 1
                        if convert[v] == -1:
                            b += 1
                # 消失的连通分量只有在最终状态时才允许
                if b >= 2:
                    continue
                if b == 1:
                    if max(nexts) != -1:
                        continue
                    edges.append((p, 1))
                    continue
                h = tuple(nexts)
                if mergeReverse:
                    h = min(h, tuple(ConnectedDpSquares._reverseState(nexts)))
                if h not in mp:
                    mp[h] = len(states)
                    states.append(nexts)
                edges.append((p, mp[h]))

        return states, edges

    @staticmethod
    def polygonDpGraph(n: int) -> Tuple[List[List[int]], List[Tuple[int, int]]]:
        """返回[states,edges]"""

        def getOk(now: List[int], nxt: List[int], convert: List[int]) -> bool:
            for i in range(n - 1):
                a1, a2 = now[i] != -1, now[i + 1] != -1
                b1, b2 = nxt[i] != -1, nxt[i + 1] != -1
                if a1 and not a2 and not b1 and b2:
                    return False
                if not a1 and a2 and b1 and not b2:
                    return False
            close = 0
            after = 0
            isNew = [True] * n
            for i in range(n):
                if convert[i] != -1:
                    isNew[convert[i]] = False
            for i in range(n):
                if nxt[i] == i and not isNew[i]:
                    after += 1
            order = [i for i in range(n) if now[i] != -1]
            for k in range(len(order) - 1):
                i, j = order[k], order[k + 1]
                if j == i + 1:
                    continue
                cl = True
                for p in range(i + 1, j):
                    if nxt[p] == -1:
                        cl = False
                        break
                if cl:
                    close += 1
            return a - close == after

        mp = dict()
        states = [[-1] * n, [-1] * n]
        edges = []
        mp[tuple(states[0])] = 0
        p = -1
        while True:
            p += 1
            if p == len(states):
                break
            if p == 1:
                edges.append((1, 1))
                continue
            now = states[p]
            for nexts, convert in ConnectedDpSquares._getNextStates(now):
                a, b = 0, 0
                for v in range(n):
                    if now[v] == v:
                        a += 1
                        if convert[v] == -1:
                            b += 1
                if b >= 2:
                    continue
                if b == 1:
                    if max(nexts) != -1:
                        continue
                    edges.append((p, 1))
                    continue
                ok = getOk(now, nexts, convert)
                if not ok:
                    continue
                h = tuple(nexts)
                if h not in mp:
                    mp[h] = len(states)
                    states.append(nexts)
                edges.append((p, mp[h]))

        return states, edges

    @staticmethod
    def _getNextStates(now: List[int]) -> List[Tuple[List[int], List[int]]]:
        """返回每项是 (新的状态, 从当前状态到新状态的转移) 的数组"""

        def find(x: int) -> int:
            while parent[x] != x:
                parent[x] = parent[parent[x]]
                x = parent[x]
            return x

        def union(a: int, b: int) -> None:
            a, b = find(a), find(b)
            if a == b:
                return
            if a > b:
                a, b = b, a
            parent[b] = a

        n = len(now)
        res = []
        for s in range(1 << n):
            parent = [-1] * (n + n)
            for i in range(n):
                if s & (1 << i):
                    parent[i] = i
            for i in range(n):
                if now[i] != -1:
                    parent[n + i] = n + now[i]
            for i in range(n - 1):
                if parent[i] != -1 and parent[i + 1] != -1:
                    union(i, i + 1)
            for i in range(n):
                if parent[i] != -1 and parent[n + i] != -1:
                    union(i, n + i)
            for i in range(n + n):
                if parent[i] != -1:
                    parent[i] = find(i)
            for i in range(n, n + n):
                if parent[i] >= n:
                    parent[i] = -1
            res.append((parent[:n], parent[n:]))

        return res

    @staticmethod
    def _reverseState(now: List[int]) -> List[int]:
        n = len(now)
        maxI = [-1] * n
        for i in range(n):
            if now[i] != -1:
                maxI[now[i]] = i
        rev = [-1] * n
        for i in range(n):
            if now[i] != -1:
                x = maxI[now[i]]
                rev[n - 1 - i] = n - 1 - x
        return rev


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc296/tasks/abc296_h
    # 形成一个黑色连通块,最少要涂黑多少个格子
    INF = int(1e18)
    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]
    states, edges = ConnectedDpSquares.connectedDpGraph(COL, False)

    S = len(states)
    dp = [INF] * S
    dp[0] = 0
    for r in range(ROW):
        ndp = [INF] * S
        for a, b in edges:
            ok = True
            cost = 0
            for c in range(COL):
                if grid[r][c] == "#" and states[b][c] == -1:
                    ok = False
                if states[b][c] != -1 and grid[r][c] == ".":
                    cost += 1
            if ok:
                ndp[b] = min(ndp[b], dp[a] + cost)
        dp = ndp
    res = INF
    for a, b in edges:
        if b == 1:
            res = min(res, dp[a])
    print(res)
