from typing import List


class UnionFind:
    def __init__(self):
        self.lis = list(range(1000))
        return

    def find(self, x):
        if self.lis[x] != x:
            self.lis[x] = self.find(self.lis[x])
        return self.lis[x]

    def union(self, x, y):
        xx = self.find(x)
        yy = self.find(y)
        self.lis[xx] = yy
        return


def equationsPossible(self, equations: List[str]) -> bool:
    dsu = UnionFind()
    record = []
    for eq in equations:
        if "!" in eq:
            record.append(eq)
        else:
            dsu.union(ord(eq[0]), ord(eq[-1]))
    for rest in record:
        if dsu.find(ord(rest[0])) == dsu.find(ord(rest[-1])):
            return False
    return True
