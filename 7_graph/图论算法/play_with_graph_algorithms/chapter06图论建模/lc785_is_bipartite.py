class Solution:

    def is_partite(self, graph):
        V = len(graph)
        visited = [False] * V
        colors = [-1] * V

        for v in range(V):
            if not visited[v]:
                if not self._dfs(v, 0, graph, visited, colors):
                    return False
        return True

    def _dfs(self, v, color, graph, visited, colors):
        visited[v] = True
        colors[v] = color

        for w in graph[v]:
            if not visited[w]:
                if not self._dfs(w, 1 - color, graph, visited, colors):
                    return False
            elif colors[v] == colors[w]:
                return False
        return True


if __name__ == '__main__':
    sol = Solution()
    data = [[1,3], [0,2], [1,3], [0,2]]
    print(sol.is_partite(data))