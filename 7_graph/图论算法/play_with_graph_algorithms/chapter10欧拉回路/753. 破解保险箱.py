class Solution(object):
    def crackSafe(self, n, k):
        """
        :type n: int
        :type k: int
        :rtype: str
        """

        arr = []
        for i in range(k):
            arr.append(str(i))

        if n == 1:
            return "".join(arr)

        graph = self.build_graph(arr, n - 1)

        seq = self.build_debrujin_sequence(graph, n - 1, arr[0])
        return "".join(seq)

    def build_vertices(self, arr, node_char_count, cnt_index):
        result = []

        if cnt_index >= node_char_count:
            return result

        for i in arr:
            returned_val = self.build_vertices(arr, node_char_count, cnt_index + 1)

            if len(returned_val):
                result.extend([i + val for val in returned_val])
            else:
                result.append(i)

        return result

    # 我们将所有的 n−1 位数作为节点，共有 k**(n-1)个节点，每个节点有 k 条入边和出边
    def build_graph(self, arr, node_char_count):
        vertices = self.build_vertices(arr, node_char_count, 0)

        graph = {}

        for v in vertices:
            graph[v] = []
            for c in arr:
                graph[v].append(v[1:] + c)
        print(graph)
        return graph

    def build_debrujin_sequence(self, graph, node_char_count, start_char):
        start_node = start_char * (node_char_count)
        seq = []
        visited = set([])
        self.dfs(graph, start_node, node_char_count, seq, visited)
        seq.append(start_node[0:-1])
        return seq

    def dfs(self, graph, node, node_char_count, seq, visited):
        vals = graph[node]
        for v in vals:
            edge = (node + v[1:]) if node_char_count > 1 else (node + v)
            if edge not in visited:
                visited.add(edge)
                self.dfs(graph, v, node_char_count, seq, visited)
        seq.append(node[-1])


print(Solution().crackSafe(3, 2))

