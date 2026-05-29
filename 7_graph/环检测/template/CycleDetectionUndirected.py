class CycleDetectionUndirected:
    __slots__ = ("n", "m", "edges", "start", "mask", "shift", "elist")

    def __init__(self, n: int, m: int) -> None:
        self.n = n
        self.m = m
        self.edges = []
        self.start = [0] * (n + 1)
        self.shift = 31
        self.mask = (1 << self.shift) - 1

    def add_edge(self, u: int, v: int, edge_id: int) -> None:
        self.edges.append((u, v, edge_id))
        self.edges.append((v, u, edge_id))
        self.start[u + 1] += 1
        self.start[v + 1] += 1

    def find_cycle(self):
        self._csr()
        n = self.n
        m = self.m
        st = self.start
        el = self.elist
        shift = self.shift
        mask = self.mask
        seen = [0] * n
        used = [0] * m
        parent = [0] * n
        parent_edge_id = [0] * n  # 親への辺のID
        for i in range(n):
            if seen[i]:
                continue
            search = [(i, -1, -1)]
            while search:
                v, p, edge_id = search.pop()
                if edge_id != -1 and used[edge_id]:
                    continue
                if seen[v]:
                    parent[v] = p
                    if edge_id != -1:
                        parent_edge_id[v] = edge_id
                    cycle_v = [p]
                    while v != p:
                        p = parent[p]
                        cycle_v.append(p)
                    cycle_v.reverse()
                    cycle_e = [parent_edge_id[p] for p in cycle_v[1:]]
                    cycle_e.append(edge_id)
                    return cycle_v, cycle_e
                seen[v] = 1
                if edge_id != -1:
                    parent_edge_id[v] = edge_id
                    used[edge_id] = 1
                parent[v] = p
                for k in range(st[v], st[v + 1]):
                    uid = el[k]
                    u, edge_id = uid >> shift, uid & mask
                    if used[edge_id]:
                        continue
                    search.append((u, v, edge_id))
        return [], []

    def _csr(self) -> None:
        for i in range(self.n):
            self.start[i + 1] += self.start[i]
        counter = self.start[:]
        self.elist = [0] * len(self.edges)
        shift = self.shift
        for u, v, edge_id in self.edges:
            self.elist[counter[u]] = (v << shift) + edge_id
            counter[u] += 1


if __name__ == "__main__":
    n, m = map(int, input().split())
    cd = CycleDetectionUndirected(n, m)
    for i in range(m):
        u, v = map(int, input().split())
        cd.add_edge(u, v, i)
    cycle_v, cycle_e = cd.find_cycle()
    if len(cycle_v) == 0:
        print(-1)
    else:
        print(len(cycle_v))
        print(" ".join(map(str, cycle_v)))
        print(" ".join(map(str, cycle_e)))
