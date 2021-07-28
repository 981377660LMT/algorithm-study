class Graph:

    """Suppose to use RB tree, but no TreeSet/TreeMap in vanilla Python, using
    set instead.

    Support both directed graph and indirected graph
    """

    def __init__(self, filename, directed=False, reverse=False):
        self._filename = filename
        self._directed = directed
        self._reverse = reverse
        lines = None
        with open(filename, 'r') as f:
            lines = f.readlines()
        if not lines:
            raise ValueError('Expected something from input file!')

        # lines[0] -> V E
        self._V, self._E = (int(i) for i in lines[0].split())

        if self._V < 0:
            raise ValueError('V must be non-negative')

        if self._E < 0:
            raise ValueError('E must be non-negative')

        # size V list of set
        self._adj = [set() for _ in range(self._V)]
        self._indegree = [0] * self._V
        self._outdegree = [0] * self._V

        for each_line in lines[1:]:
            a, b = (int(i) for i in each_line.split())
            if self._reverse:
                a, b = b, a
            self.validate_vertex(a)
            self.validate_vertex(b)

            if a == b:
                raise ValueError('Self-Loop is detected!')

            if b in self._adj[a]:
                raise ValueError('Paralles edges are detected!')

            self._adj[a].add(b)

            if self._directed:
                self._outdegree[a] += 1
                self._indegree[b] += 1

            if not self._directed:
                self._adj[b].add(a)

    def reverse_graph(self):
        return Graph(filename=self._filename, directed=self._directed, reverse=True)

    @property
    def V(self):
        return self._V

    @property
    def E(self):
        return self._E

    def has_edge(self, v, w):
        self.validate_vertex(v)
        self.validate_vertex(w)
        return w in self._adj[v]

    def adj(self, v):
        self.validate_vertex(v)
        return self._adj[v]

    def degree(self, v):
        if self._directed:
            raise ValueError('degree only works on undirected graph')
        return len(self.adj(v))

    def indegree(self, v):
        if not self._directed:
            raise ValueError('indegree only works in directed graph')
        self.validate_vertex(v)
        return self._indegree[v]

    def outdegree(self, v):
        if not self._directed:
            raise ValueError('outdegree only works in directed graph')
        self.validate_vertex(v)
        return self._outdegree[v]

    def remove_edge(self, v, w):
        self.validate_vertex(v)
        self.validate_vertex(w)
        if w in self._adj[v]:
            self._E -= 1
            if self._directed:
                self._outdegree[v] -= 1
                self._indegree[w] -= 1
        self._adj[v].remove(w)
        if not self._directed:
            self._adj[w].remove(v)

    def validate_vertex(self, v):
        if v < 0 or v >= self._V:
            raise ValueError('vertex ' + v + ' is invalid')

    def is_directed(self):
        return self._directed

    def __str__(self):
        res = ['V = {}, E = {}, directed = {}'.format(self._V, self._E, self._directed)]
        for v in range(self._V):
            res.append('{}: {}'.format(v, ' '.join(str(w) for w in self._adj[v])))
        return '\n'.join(res)

    def __repr__(self):
        return self.__str__()

    def __copy__(self):
        return Graph(self._filename, self._directed, self._reverse)


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)

    for v in range(g.V):
        print(g.indegree(v), g.outdegree(v))