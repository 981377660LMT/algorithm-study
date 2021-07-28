class WeightedGraph:

    """Suppose to use RB tree, but no TreeSet/TreeMap in vanilla Python, using
    dict instead.

    Support both directed graph and indirected graph
    """

    def __init__(
        self,
        filename=None,
        directed=False,
        is_redisual=False,
        empty_graph=False,
        V=None,
    ):
        if empty_graph is False:
            self._generate_graph_from_file(
                filename=filename,
                directed=directed,
                is_redisual=is_redisual,
                empty_graph=empty_graph,
            )
        else:
            self._generate_empty_graph(directed, V)

    def _generate_empty_graph(self, directed, V):
        self._V = V
        self._directed = directed
        self._E = 0
        self._adj = [dict() for _ in range(self._V)]

        self._filename = None

    def _generate_graph_from_file(
        self,
        filename,
        directed=False,
        is_redisual=False,
        empty_graph=False,
    ):
        self._filename = filename
        self._directed = directed
        self._is_redisual = is_redisual
        lines = None
        with open(filename, 'r') as f:
            lines = f.readlines()
        if not lines:
            raise ValueError('Expected something from input file!')

        # lines[0] -> V e
        self._E = 0
        self._V, self._e = (int(i) for i in lines[0].split())

        if self._V < 0:
            raise ValueError('V must be non-negative')

        if self._e < 0:
            raise ValueError('E must be non-negative')

        # size V list of dictionaries
        self._adj = [dict() for _ in range(self._V)]
        for each_line in lines[1:]:
            a, b, weight = (int(i) for i in each_line.split())
            self.add_edge(a, b, weight)
            if self._is_redisual:
                self.add_edge(b, a, 0)

    def add_edge(self, a, b, weight):
        self.validate_vertex(a)
        self.validate_vertex(b)

        if a == b:
            raise ValueError('Self-Loop is detected!')

        # if b in self._adj[a]:
        #     raise ValueError('Paralles edges are detected!')

        self._adj[a][b] = weight
        if not self._directed:
            self._adj[b][a] = weight

        self._E += 1

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

    def set_weight(self, v, w, net_weight):
        if not self.has_edge(v, w):
            raise ValueError('No edge {}-{}'.format(v, w))
        self._adj[v][w] = net_weight
        if not self._directed:
            self._adj[w][v] = net_weight

    # def degree(self, v):
    #     return len(self.adj(v))

    def remove_edge(self, v, w):
        self.validate_vertex(v)
        self.validate_vertex(w)
        if w in self._adj[v]:
            self._adj[v].pop(w)
        if not self._directed and v in self._adj[w]:
            self._adj[w].pop(v)

    def validate_vertex(self, v):
        if v < 0 or v >= self._V:
            raise ValueError('vertex ' + str(v) + ' is invalid')

    def is_directed(self):
        return self._directed

    def generate_redisual_graph(self, empty_graph=False):
        # redisual graph is definitely a directed graph
        if not empty_graph:
            return WeightedGraph(filename=self._filename, directed=True, is_redisual=True)
        return WeightedGraph(empty_graph=True, directed=True, is_redisual=True)

    def __str__(self):
        res = ['V = {}, E = {}, directed = {}'.format(self._V, self._E, self._directed)]
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
    filename = 'play_with_graph_algorithms/chapter13/wg.txt'
    w_graph = WeightedGraph(filename, directed=True)
    print(w_graph)
    