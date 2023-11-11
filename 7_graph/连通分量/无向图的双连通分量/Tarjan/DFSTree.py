"""非递归版的Tarjan算法,用于求解无向图的双连通分量"""

from typing import List, Tuple


class DFSTree:
    # cf: https://codeforces.com/blog/entry/68138
    def __init__(self, N: int):
        self.N = N
        self.edges = []

        self.E = [[] for _ in range(self.N)]
        # span-edge and back-edge (directed)
        self.span_edge = [[] for _ in range(self.N)]
        self.back_edge = [[] for _ in range(self.N)]

        self.ord = [-1] * self.N
        self.low = [-1] * self.N
        self.par = [None] * self.N

        self.is_art = [False] * self.N
        self.is_bridge = []

        self.built = False

    def add_edge(self, u: int, v: int) -> None:
        """add edge"""
        eid = len(self.edges)
        self.edges.append((u, v))
        self.E[u].append((v, eid))
        self.E[v].append((u, eid))
        self.is_bridge.append(False)

    def build(self) -> None:
        """build dfs tree"""
        cnt = 0
        self.tvcc_id = [-1] * len(self.edges)
        self.bcc_num = 0
        for i in range(self.N):
            if ~self.ord[i]:
                continue
            stack = [(i, -1, -1)]
            estack = []
            while stack:
                v, p, p_eid = stack.pop()
                if v < 0:
                    v = ~v
                    for d, i in self.span_edge[v]:
                        if d == p:
                            continue
                        self.low[v] = min(self.low[v], self.low[d])
                    if ~p and self.ord[p] <= self.low[v]:
                        while True:
                            eid = estack.pop()
                            self.tvcc_id[eid] = self.bcc_num
                            if eid == p_eid:
                                break
                        self.bcc_num += 1
                else:
                    if ~self.ord[v]:
                        continue
                    self.ord[v] = cnt
                    self.low[v] = cnt
                    cnt += 1
                    # p -> v is span-edge.
                    if ~p:
                        self.par[v] = (p, p_eid)
                        self.span_edge[p].append((v, p_eid))
                        estack.append(p_eid)
                    stack.append((~v, p, p_eid))

                    for d, eid in self.E[v][::-1]:
                        if eid == p_eid:
                            continue
                        if ~self.ord[d]:
                            # v -> d is back-edge since v is already visited.
                            self.back_edge[v].append((d, eid))
                            self.low[v] = min(self.low[v], self.ord[d])
                            estack.append(eid)
                            continue
                        stack.append((d, v, eid))

        self._search_bridge()
        self._search_articulation_points()
        self.built = True

    def bridges(self) -> List[Tuple[int, int]]:
        """return list of edges that are bridges"""
        assert self.built
        return [e for i, e in enumerate(self.edges) if self.is_bridge[i]]

    def articulation_points(self) -> List[int]:
        """return list of vertices that are articulation points"""
        assert self.built
        return [i for i in range(self.N) if self.is_art[i]]

    def two_edge_connected_components(self) -> List[List[int]]:
        """边双连通分量"""
        assert self.built
        tecc_id = [-1] * self.N
        cnt = 0
        for i in range(self.N):
            if ~tecc_id[i]:
                continue
            stack = [i]
            while stack:
                v = stack.pop()
                if ~tecc_id[v]:
                    continue
                tecc_id[v] = cnt
                for d, eid in self.E[v]:
                    if ~tecc_id[d] or self.is_bridge[eid]:
                        continue
                    stack.append(d)
            cnt += 1
        ret = [[] for _ in range(cnt)]
        for i, tid in enumerate(tecc_id):
            ret[tid].append(i)
        return ret

    def biconnected_components(self) -> List[List[int]]:
        """点双连通分量"""
        assert self.built
        ret = [set() for _ in range(self.bcc_num)]
        for eid, tid in enumerate(self.tvcc_id):
            assert ~eid
            u, v = self.edges[eid]
            ret[tid].add(u)
            ret[tid].add(v)
        ret = [list(s) for s in ret]
        for i in range(self.N):
            if not self.E[i]:
                ret.append([i])
        return ret

    def _search_bridge(self) -> None:
        for u in range(self.N):
            for v, i in self.span_edge[u]:
                # (u, v) is bridge if vertex u has child v
                # that does not have lowlink to pass over its parent
                self.is_bridge[i] = self.ord[u] < self.low[v]

    def _search_articulation_points(self) -> None:
        for u in range(self.N):
            if self.par[u] is None:
                self.is_art[u] = len(self.span_edge[u]) >= 2
            else:
                for v, _ in self.span_edge[u]:
                    if self.ord[u] <= self.low[v]:
                        self.is_art[u] = True
                        break


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    tree = DFSTree(n)
    for _ in range(m):
        u, v = map(int, input().split())
        tree.add_edge(u, v)
    tree.build()

    # groups = tree.biconnected_components()
    groups = tree.two_edge_connected_components()
    print(len(groups))
    for g in groups:
        print(len(g), *g)
