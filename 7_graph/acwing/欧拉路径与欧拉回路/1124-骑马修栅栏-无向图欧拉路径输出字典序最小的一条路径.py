# from collections import defaultdict
# from typing import DefaultDict, List


# def getEulerPath(n: int, adjMap: DefaultDict[int, List[int]], isDirected: bool) -> List[int]:
#     """无向图欧拉路径输出字典序最小的一条路径"""
#     start = 0
#     if isDirected:
#         indegree, outdegree = [0] * n, [0] * n
#         for cur, nexts in adjMap.items():
#             outdegree[cur] += len(nexts)
#             for next in nexts:
#                 indegree[next] += 1
#         for i in range(n):
#             diff = outdegree[i] - indegree[i]
#             if diff == 1:
#                 start = i
#                 break
#     else:
#         for i in sorted(adjMap.keys()):
#             if len(adjMap[i]) & 1:
#                 start = i
#                 break

#     res = []
#     stack = [start]
#     cur = start
#     while stack:
#         if adjMap[cur]:
#             stack.append(cur)
#             next = adjMap[cur].pop()
#             if not isDirected:
#                 adjMap[next].remove(cur)  # 无向图 要删两条边
#             cur = next
#         else:
#             res.append(cur)
#             cur = stack.pop()

#     return res[::-1]


# n = int(input())
# adjMap = defaultdict(list)
# for _ in range(n):
#     x, y = map(int, input().split())
#     x, y = x - 1, y - 1
#     adjMap[x].append(y)
#     adjMap[y].append(x)


# # 要求字典序最小路径，那么就把字典序大的边放在前面，字典序小的边放在后面，早点pop出来
# for nexts in adjMap.values():
#     nexts.sort(reverse=True)


# res = getEulerPath(n, adjMap, isDirected=False)
# for i in res:
#     print(i + 1)
