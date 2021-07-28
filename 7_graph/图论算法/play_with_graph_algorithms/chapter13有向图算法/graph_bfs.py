from collections import deque
from play_with_graph_algorithms.chapter13.graph import Graph
from play_with_graph_algorithms.chapter05.graph_bfs import GraphBFS


if __name__ == '__main__':
    filename = 'play_with_graph_algorithms/chapter13/ug.txt'
    g = Graph(filename, directed=True)
    graph_bfs = GraphBFS(g)
    print('BFS order : {}'.format(graph_bfs.order()))
