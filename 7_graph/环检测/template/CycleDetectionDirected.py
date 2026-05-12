class CycleDetectionDirected:
    __slots__ = ("n", "edges", "start", "mask", "shift", "elist")

    def __init__(self, n: int):
        self.n = n
        self.edges = []
        self.start = [0] * (n + 1)
        self.shift = 31
        self.mask = (1 << self.shift) - 1

    def add_edge(self, u: int, v: int, edge_id: int) -> None:
        self.edges.append((u, v, edge_id))
        self.start[u + 1] += 1

    def find_cycle(self):
        self._csr()
        n = self.n
        st = self.start
        el = self.elist
        shift = self.shift
        mask = self.mask
        seen = [0] * n
        finished = [0] * n
        stack = []
        for i in range(n):
            if seen[i]:
                continue
            search = [
                ((i << shift) + mask) * 2 + 1,
                ((i << shift) + mask) * 2,
            ]  # (v = i, edge_id = mask, t = 1, 0)
            seen[i] = True
            while search:
                vidt = search.pop()
                vid, t = vidt >> 1, vidt & 1
                v, edge_id = vid >> shift, vid & mask
                if finished[v]:
                    continue
                if t == 0:
                    seen[v] = 1
                    stack.append((v, edge_id))
                    for k in range(st[v], st[v + 1]):
                        uid = el[k]
                        u, id = uid >> shift, uid & mask
                        if seen[u] and finished[u] == 0:
                            cycle = [id]
                            while stack:
                                v, id = stack.pop()
                                if v == u:
                                    break
                                cycle.append(id)
                            return cycle[::-1]

                        elif seen[u] == 0:
                            search.append(uid * 2 + 1)
                            search.append(uid * 2)

                else:
                    stack.pop()
                    finished[v] = 1
        return []

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
    cd = CycleDetectionDirected(n)
    for i in range(m):
        u, v = map(int, input().split())
        cd.add_edge(u, v, i)
    cycle = cd.find_cycle()
    if len(cycle) == 0:
        print(-1)
    else:
        print(len(cycle))
        print("\n".join(map(str, cycle)))
