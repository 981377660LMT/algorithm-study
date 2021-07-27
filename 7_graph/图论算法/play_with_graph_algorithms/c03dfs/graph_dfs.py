from play_with_graph_algorithms.chapter02.adj_set import AdjSet as Graph


class GraphDFS:

    def __init__(self, G, recursive=True):
        self._pre_order = []
        self._post_order = []
        self._G = G
        self._visited = [False] * G.V

        # 遍历所有的点，相当于遍历图中所有可能存在的联通块
        for v in range(G.V):
            if not self._visited[v]:
                if recursive:
                    self._dfs_recursive(v)
                else:
                    self._dfs_iteration(v)

    def _dfs_recursive(self, v):
        self._visited[v] = True
        self._pre_order.append(v)
        for w in self._G.adj(v):
            if not self._visited[w]:
                self._dfs_recursive(w)
        self._post_order.append(v)

    def _dfs_iteration(self, v):
        """For pre-order that's straight-forward by using one stack,
        but for post-order we need an augmented stack2
        """
        stack1 = [v]
        self._visited[v] = True
        stack2 = []
        while stack1:
            curr = stack1.pop()
            self._pre_order.append(curr)
            stack2.append(curr)
            for w in self._G.adj(curr):
                if not self._visited[w]:
                    stack1.append(w)
                    self._visited[w] = True
        self._post_order += stack2[::-1]
    
    @property
    def pre_order(self):
        return self._pre_order

    @property
    def post_order(self):
        return self._post_order


if __name__ == '__main__':
    print('For one block, recursive:')
    filename = 'play_with_graph_algorithms/chapter03/g1.txt'
    g = Graph(filename)
    graph_dfs = GraphDFS(g)
    print(graph_dfs.pre_order)
    print(graph_dfs.post_order)

    print('*' * 40)

    print('For two blocks, recursive:')
    filename = 'play_with_graph_algorithms/chapter03/g2.txt'
    g = Graph(filename)
    graph_dfs = GraphDFS(g)
    print(graph_dfs.pre_order)
    print(graph_dfs.post_order)

    print('*' * 40)

    print('For one block, iteration:')
    filename = 'play_with_graph_algorithms/chapter03/g1.txt'
    g = Graph(filename)
    graph_dfs = GraphDFS(g, recursive=False)
    print(graph_dfs.pre_order)
    print(graph_dfs.post_order)

    print('*' * 40)

    print('For two blocks, iteration:')
    filename = 'play_with_graph_algorithms/chapter03/g2.txt'
    g = Graph(filename)
    graph_dfs = GraphDFS(g, recursive=False)
    print(graph_dfs.pre_order)
    print(graph_dfs.post_order)
