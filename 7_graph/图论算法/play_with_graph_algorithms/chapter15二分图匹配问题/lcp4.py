from play_with_graph_algorithms.chapter15.bipartite_matching import BipartiteMatching
from play_with_graph_algorithms.chapter13.weighted_graph import WeightedGraph
from play_with_graph_algorithms.chapter15.hungarian import Hungarian


def domino(n, m, broken):
    board = [[0] * m for _ in range(n)]
    for i, j in broken:
        board[i][j] = 1
    # 由于Graph并没有实现重载
    # 这里使用WeightedGraph实现的重载
    g = WeightedGraph(empty_graph=True, V=n * m)
    for i in range(n):
        for j in range(m):
            if j + 1 < m and board[i][j] == 0 and board[i][j + 1] == 0:
                g.add_edge(i * m + j, i * m + (j + 1), 1)
            if i + 1 < n and board[i][j] == 0 and board[i + 1][j] == 0:
                g.add_edge(i * m + j, (i + 1) * m + j, 1)

    bm = BipartiteMatching(g)
    return bm.max_matching()


def domino_hungarian(n, m, broken, use_dfs=False):
    board = [[0] * m for _ in range(n)]
    for i, j in broken:
        board[i][j] = 1
    # 由于Graph并没有实现重载
    # 这里使用WeightedGraph实现的重载
    g = WeightedGraph(empty_graph=True, V=n * m)
    for i in range(n):
        for j in range(m):
            if j + 1 < m and board[i][j] == 0 and board[i][j + 1] == 0:
                g.add_edge(i * m + j, i * m + (j + 1), 1)
            if i + 1 < n and board[i][j] == 0 and board[i + 1][j] == 0:
                g.add_edge(i * m + j, (i + 1) * m + j, 1)

    hungarian = Hungarian(g, use_dfs=use_dfs)
    return hungarian.maxmatching()


if __name__ == "__main__":
    print('Maxflow:')
    print(domino(2, 3, broken=[[1, 0], [1, 1]]))
    print(domino(3, 3, broken=[]))

    print('Hungarian BFS:')
    print(domino_hungarian(2, 3, broken=[[1, 0], [1, 1]]))
    print(domino_hungarian(3, 3, broken=[]))

    print('Hungarian DFS:')
    print(domino_hungarian(2, 3, broken=[[1, 0], [1, 1]], use_dfs=True))
    print(domino_hungarian(3, 3, broken=[], use_dfs=True))