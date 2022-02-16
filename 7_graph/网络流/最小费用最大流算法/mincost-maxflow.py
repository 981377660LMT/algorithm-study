# class MincostMaxflow:
#     def __init__(self, n):
#         self.n = n
#         self.graph = [[] for _ in range(n)]
#         self.mincost = 0
#         self.minflow = 0

#     def addEdge(self, cur, next, cap, cost):
#         self.graph[cur].append([next, cap, cost, len(self.graph[next])])
#         self.graph[next].append([cur, 0, -cost, len(self.graph[cur]) - 1])

#     def findPath(self, cur, target, pre):
#         if cur == target:
#             return True
#         for i in range(len(self.graph[cur])):
#             v, c, w, p = self.graph[cur][i]
#             if pre[v] == -1 and c > 0:
#                 pre[v] = i
#                 if self.findPath(v, target, pre):
#                     return True
#         return False

#     def minCostMaxFlow(self, s, t):
#         pre = [-1] * self.n
#         flow = 0
#         cost = 0
#         while True:
#             for i in range(self.n):
#                 pre[i] = -1
#             pre[s] = 0
#             if not self.findPath(s, t, pre):
#                 break
#             minflow = int(1e20)
#             for i in range(self.n):
#                 if pre[i] != -1:
#                     v, c, w, p = self.graph[i][pre[i]]
#                     minflow = min(minflow, c)
#             flow += minflow
#             cost += minflow * self.graph[s][pre[s]][2]
#             for i in range(self.n):
#                 if pre[i] != -1:
#                     v, c, w, p = self.graph[i][pre[i]]
#                     self.graph[i][pre[i]][1] -= minflow
#                     self.graph[v][p][1] += minflow
#             self.minflow = flow
#             self.mincost = cost
#         return self.minflow, self.mincost


# if __name__ == '__main__':
#     flow = MincostMaxflow(4)
#     flow.addEdge(1, 1, 3, 1)
#     flow.addEdge(0, 3, 4, 3)
#     flow.addEdge(3, 2, 3, 1)
#     flow.addEdge(1, 2, 2, 2)
#     print(flow.minCostMaxFlow(2, 2))

# todo
