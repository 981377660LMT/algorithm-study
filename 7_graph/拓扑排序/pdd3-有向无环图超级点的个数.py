# pdd的T3
# !给一个有向无环图，超级点定义为图中任意点都可以从自己到超级点或者超级点到自己；问有多少个超级点.
# 在一个有向无环图中，如果一个点u是超级点，那么它必须满足对于图中的任意其他点v，
# 要么存在一条从u到v的路径，要么存在一条从v到u的路径，求图中超级点的个数。
# 数据范围2<=节点数<=3e5，1<=边数<=3e5
#

# 如果不是DAG，缩点即可


from heapq import heapify, heappop, heappush
from typing import List, Tuple


# // get min and max time for each node
#     vector<int> min_time(n, 0), max_time(n, 0);
#     int tot = 0;
#     for (int i = 0; i < n; i++) {
#         for (int j : adj[i]) {
#             min_time[j] = max(min_time[j], min_time[i] + 1);
#             tot = max(tot, min_time[j]);
#         }
#     }
#     max_time = vector<int>(n, tot);
#     for (int i = n - 1; i >= 0; i--) {
#         for (int j : adj[i]) {
#             max_time[i] = min(max_time[i], max_time[j] - 1);
#         }
#     }
#     vector<int> my_ans(n, 1);
#     for (int i = 0; i < n; i++) {
#         bool feas = true;
#         if (min_time[i] != max_time[i]) feas = false;
#         int t = min_time[i];
#         for (int j = 0; j < n; j++) {
#             if (j != i && (max_time[j] >= t && min_time[j] <= t)) feas = false;
#         }
#         my_ans[i] = feas;
#     }
# TODO: verify
# 4个点，边为[(0,1),(0,2),(0,3)]，2这个点会被误判为关键点
def superPoints(n: int, edges: List[Tuple[int, int]]) -> List[bool]:
    # adjList = [[] for _ in range(n)]
    # rAdjList = [[] for _ in range(n)]
    # for u, v in edges:
    #     adjList[u].append(v)
    #     rAdjList[v].append(u)
    # order = topoSortByHeap(n, adjList, minFirst=True)[0]
    # rOrder = topoSortByHeap(n, rAdjList, minFirst=True)[0]
    # print(order, rOrder)
    # rOrder.reverse()
    # return [a == b for a, b in zip(order, rOrder)]
    ...


if __name__ == "__main__":
    n = 4
    edges = [(0, 1), (0, 2), (0, 3)]
    print(superPoints(n, edges))  # [True, True, True, True, True]
