# # ModShortestPath
# # “同余最短路”一般是指一种和不定方程相关的建模方式
# # 对于一堆系数 a0, a1, ..., an-1 (0<=a0<=a1<=...<=an-1,最小的非零a0称为base)
# # !我们可以用这个叫同余最短路的方法，来确定他们的`线性组合∑ai*xi`可能取到的值(xi非负)
# # 更确切地说,可以处理出一个数组 dist[0,1,...,base-1]
# # !这里 dist[i]记录的是最小的 x,满足 x=i(mod base)且 x 能被上述系数线性表出
# # 具体的方式是把余 0,余 1,...,余 base-1 看成一个个结点,每个点表示一个剩余类
# # 对每个剩余类,考虑所有的转移:
# # !也就是 x 加上每一个 ai,即连一条从 x 到 (x+ai) mod base 的边，边权为 ai
# # 这样建完图从 0 号点开始跑最短路,每转移一条长度为 d 的边就对应着线性组合加了一个 d
# # 得到最后的 dist 数组

# # !注意:一共有len(A)*min(A)条边

# from typing import List, Sequence, Tuple
# from heapq import heappop, heappush

# INF = int(1e18)


# def modShortestPath(coeffs: List[int]) -> Tuple[int, List[int]]:
#     """确定线性组合∑ai*xi的可能取到的值(ai非负)

#     Args:
#         coeffs (List[int]): 非负整数系数,最小的非零ai称为base

#     Returns:
#         Tuple[int, List[int]]: base, dist
#         base (int): 最小的非零ai
#         dist (List[int]): dist[i]记录的是最小的x,满足x=i(mod base)且x能被系数coeffs线性表出(xi非负)
#         如果不存在这样的x,则dist[i]为INF
#         如果coeff全为0,则返回空数组
#     """
#     coeffs = [v for v in coeffs if v > 0]  # sorted?
#     if not coeffs:
#         return 0, []

#     base = min(coeffs)
#     adjList = [[] for _ in range(base)]
#     for mod in range(base):
#         for v in coeffs:
#             adjList[mod].append(((mod + v) % base, v))
#     dist = dijkstra(base, adjList, 0)
#     return base, dist


# def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
#     dist = [INF] * n
#     dist[start] = 0
#     pq = [(0, start)]

#     while pq:
#         curDist, cur = heappop(pq)
#         if dist[cur] < curDist:
#             continue
#         for next, weight in adjList[cur]:
#             cand = dist[cur] + weight
#             if cand < dist[next]:
#                 dist[next] = cand
#                 heappush(pq, (dist[next], next))
#     return dist
