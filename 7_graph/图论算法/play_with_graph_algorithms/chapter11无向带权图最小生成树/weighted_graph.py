class WeightedGraph:

    """Suppose to use RB tree, but no TreeSet/TreeMap in vanilla Python, using
    dict instead.
    """

    def __init__(self, filename):
        self._filename = filename
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

        # size V list of dictionaries
        self._adj = [dict() for _ in range(self._V)]
        for each_line in lines[1:]:
            a, b, weight = (int(i) for i in each_line.split())
            self.validate_vertex(a)
            self.validate_vertex(b)

            if a == b:
                raise ValueError('Self-Loop is detected!')

            if b in self._adj[a]:
                raise ValueError('Paralles edges are detected!')

            self._adj[a][b] = weight
            self._adj[b][a] = weight

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
        return list(self._adj[v].keys())

    def get_weight(self, v, w):
        if self.has_edge(v, w):
            return self._adj[v][w]
        raise ValueError('No edge {}-{}'.format(v, w))

    def degree(self, v):
        return len(self.adj(v))

    def remove_edge(self, v, w):
        self.validate_vertex(v)
        self.validate_vertex(w)
        if w in self._adj[v]:
            self._adj[v].pop(w)
        if v in self._adj[w]:
            self._adj[w].pop(v)

    def validate_vertex(self, v):
        if v < 0 or v >= self._V:
            raise ValueError('vertex ' + str(v) + ' is invalid')

    def __str__(self):
        res = ['V = {}, E = {}'.format(self._V, self._E)]
        for v in range(self._V):
            res.append(
                '{}: {}'.format(
                    v,
                    ' '.join('{}({})'.format(w, self._adj[v][w]) for w in self._adj[v]),
                ),
            )
        return '\n'.join(res)

    def __repr__(self):
        return self.__str__()

    def __copy__(self):
        return WeightedGraph(self._filename)


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter11/g.txt'
    w_graph = WeightedGraph(filename)
    print(w_graph)