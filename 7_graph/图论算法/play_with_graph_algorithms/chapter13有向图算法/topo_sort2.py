from play_with_graph_algorithms.chapter13.directed_circle_detection import DirectedCycleDetection
from play_with_graph_algorithms.chapter13.graph_dfs import GraphDFS
from play_with_graph_algorithms.chapter13.graph import Graph


class TopoSort2:

    def __init__(self, g):
        if not g.is_directed():
            raise ValueError('TopoSort only works in directed graph')

        self._g = g
        self._res = []
        self._has_cycle = DirectedCycleDetection(g).has_cycle()

        if self._has_cycle:
            return

        dfs = GraphDFS(g)
        self._res = dfs.post_order[::-1]

        if len(self._res) != self._g.V:
            self._has_cycle = True
            self._res = []

    def has_cycle(self):
        return self._has_cycle
    
    def result(self):
        return self._res


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)

    topo_sort2 = TopoSort2(g)
    print(topo_sort2.result())