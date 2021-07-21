from play_with_graph_algorithms.chapter13.graph import Graph
from play_with_graph_algorithms.chapter03.graph_dfs import GraphDFS


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)
    graph_dfs = GraphDFS(g)

    print('DFS pre order: ', graph_dfs.pre_order)
    print('DFS post order: ', graph_dfs.post_order)