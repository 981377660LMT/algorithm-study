from play_with_graph_algorithms.chapter13.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter12.bellman_ford import BellmanFord


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/wg.txt'
    g = WeightedGraph(filename, directed=True)
    bf = BellmanFord(g, s=0)

    if not bf.has_neg():
        strings = []
        for v in range(g.V):
            strings.append(bf.dist_to(v))
        print(' '.join(str(i) for i in strings))
        print(bf.path(1))
    else:
        print('exist negative circle')