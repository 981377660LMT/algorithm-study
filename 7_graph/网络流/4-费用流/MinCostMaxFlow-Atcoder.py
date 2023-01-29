from heapq import heappush, heappop

# https://atcoder.jp/contests/abc247/submissions/33589096
class MinCostFlow:
    """
    https://github.com/atcoder/ac-library/blob/master/atcoder/internal_csr.hpp
    https://github.com/atcoder/ac-library/blob/master/atcoder/mincostflow.hpp
    https://github.com/atcoder/ac-library/blob/master/document_en/mincostflow.md
    https://github.com/atcoder/ac-library/blob/master/document_ja/mincostflow.md
    """

    def __init__(self, n):
        self.n = n
        self._edges = []

    def add_edge(self, fr, to, cap, cost):
        assert 0 <= fr < self.n
        assert 0 <= to < self.n
        assert 0 <= cap
        assert 0 <= cost
        self._edges.append(self.edge(fr, to, cap, cost))
        return len(self._edges) - 1

    class edge:
        def __init__(self, fr, to, cap, cost):
            self.fr = fr
            self.to = to
            self.cap = cap
            self.flow = 0
            self.cost = cost

    def get_edge(self, i):
        assert 0 <= i < len(self._edges)
        return self._edges[i]

    def edges(self):
        return self._edges

    def flow(self, s, t, flow_limit=1 << 60):
        return self.slope(s, t, flow_limit)[-1]

    def __csr(self, edges):
        # Compressed Sparse Row
        self.start = [0] * (self.n + 1)
        for fr, _ in edges:
            self.start[fr + 1] += 1
        for i in range(self.n):
            self.start[i + 1] += self.start[i]
        counter = self.start.copy()
        self.elist = [0] * len(edges)
        for fr, e in edges:
            self.elist[counter[fr]] = e
            counter[fr] += 1

    class _edge:
        def __init__(self, to, rev, cap, cost):
            self.to = to
            self.rev = rev
            self.cap = cap
            self.cost = cost

    def __g(self):
        degree = [0] * self.n
        redge_idx = [0] * self.m
        elist = [(0, None)] * (2 * self.m)
        now_elist = 0
        for i in range(self.m):
            e = self._edges[i]
            self.edge_idx[i] = degree[e.fr]
            degree[e.fr] += 1
            redge_idx[i] = degree[e.to]
            degree[e.to] += 1
            elist[now_elist] = (e.fr, self._edge(e.to, -1, e.cap - e.flow, e.cost))
            now_elist += 1
            elist[now_elist] = (e.to, self._edge(e.fr, -1, e.flow, -e.cost))
            now_elist += 1
        self.__csr(elist)
        for i in range(self.m):
            e = self._edges[i]
            self.edge_idx[i] += self.start[e.fr]
            redge_idx[i] += self.start[e.to]
            self.elist[self.edge_idx[i]].rev = redge_idx[i]
            self.elist[redge_idx[i]].rev = self.edge_idx[i]

    def slope(self, s, t, flow_limit=1 << 60):
        assert 0 <= s < self.n
        assert 0 <= t < self.n
        assert s != t
        self.m = len(self._edges)
        self.edge_idx = [0] * self.m
        self.__g()
        result = self.__slope(s, t, flow_limit)
        for i in range(self.m):
            e = self.elist[self.edge_idx[i]]
            self._edges[i].flow = self._edges[i].cap - e.cap
        return result

    def __dual_ref(self, s, t):
        log = self.n.bit_length()
        mask = (1 << log) - 1
        dist = [1 << 60] * self.n
        vis = [0] * self.n
        que_min = []
        que = []
        dist[s] = 0
        que_min.append(s)
        while que_min or que:
            if que_min:
                v = que_min.pop()
            else:
                v = heappop(que) & mask
            if vis[v]:
                continue
            vis[v] = 1
            if v == t:
                break
            dual_v = self.dual[v]
            dist_v = dist[v]
            for i in range(self.start[v], self.start[v + 1]):
                e = self.elist[i]
                if not e.cap:
                    continue
                cost = e.cost - self.dual[e.to] + dual_v
                if dist[e.to] - dist_v > cost:
                    dist_to = dist_v + cost
                    dist[e.to] = dist_to
                    self.prev_e[e.to] = e.rev
                    if dist_to == dist_v:
                        que_min.append(e.to)
                    else:
                        heappush(que, dist_to << log | e.to)
        if not vis[t]:
            return False
        for v in range(self.n):
            if not vis[v]:
                continue
            self.dual[v] -= dist[t] - dist[v]
        return True

    def __slope(self, s, t, flow_limit):
        self.dual = [0] * self.n
        self.prev_e = [0] * self.n
        flow = 0
        cost = 0
        prev_cost_per_flow = -1
        result = [(0, 0)]
        while flow < flow_limit:
            if not self.__dual_ref(s, t):
                break
            c = flow_limit - flow
            v = t
            while v != s:
                c = min(c, self.elist[self.elist[self.prev_e[v]].rev].cap)
                v = self.elist[self.prev_e[v]].to
            v = t
            while v != s:
                e = self.elist[self.prev_e[v]]
                e.cap += c
                self.elist[e.rev].cap -= c
                v = self.elist[self.prev_e[v]].to
            d = -self.dual[s]
            flow += c
            cost += c * d
            if prev_cost_per_flow == d:
                result.pop()
            result.append((flow, cost))
            prev_cost_per_flow = d
        return result


if __name__ == "__main__":

    from collections import defaultdict
    from sys import stdin

    input = lambda: stdin.readline().strip()
    big = 10**9

    n = int(input())
    dic = defaultdict(int)
    for _ in range(n):
        a, b, c = map(int, input().split())
        dic[(a - 1, b - 1)] = max(dic[(a - 1, b - 1)], c)

    mcf = MinCostFlow(302)
    s = 300
    t = 301
    for i in range(150):
        mcf.add_edge(s, i, 1, 0)
        mcf.add_edge(150 + i, t, 1, 0)

    for (a, b), c in dic.items():
        mcf.add_edge(a, b + 150, 1, big - c)

    res = mcf.slope(s, t, n)

    print(res[-1][0])
    ans = [0] * res[-1][0]
    flow = 0
    cost = 0
    for new_flow, new_cost in res:
        for i in range(flow + 1, new_flow + 1):
            ans[i - 1] = i * big - cost - (new_cost - cost) * (i - flow) // (new_flow - flow)
        flow = new_flow
        cost = new_cost
    print(*ans, sep="\n")
