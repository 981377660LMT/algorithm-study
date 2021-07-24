from play_with_graph_algorithms.chapter13.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter12.floyd import Floyd


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/wg.txt'
    g = WeightedGraph(filename, directed=True)
    floyd = Floyd(g)

    if not floyd.has_neg_cycle():
        for v in range(g.V):
            strings = []
            for w in range(g.V):
                strings.append(str(floyd.dist_to(v, w)))
            print(' '.join(strings))
    else:
        print('exist negative cycle')
