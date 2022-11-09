# def solve(
#     n: int,
#     m: int,
#     mousePos: List[int],
#     holePos: List[int],
#     holeWage: List[int],
#     holeCount: List[int],
# ) -> int:
#     events: List["Event"] = []
#     Q0: List["Pair"] = []
#     Q1: List["Pair"] = []
#     for i in range(1, n + 1):
#         events.append(Event(mousePos[i], 0, 1, 0))
#     for i in range(1, m + 1):
#         events.append(Event(holePos[i], 1, holeCount[i], holeWage[i]))

#     events.sort(key=lambda x: x.pos)
#     heappush(Q0, Pair(int(1e12), int(1e5)))
#     heappush(Q1, Pair(int(1e12), int(1e5)))
#     res = 0
#     for event in events:
#         if event.kind == 0:
#             tmp = heappop(Q1)
#             res += tmp.cost + event.pos
#             tmp.count -= 1
#             heappush(Q0, Pair(-event.pos * 2 - tmp.cost, 1))
#             if tmp.count:
#                 heappush(Q1, tmp)
#         else:
#             flow = 0
#             while flow < event.count and Q0[0].cost + event.pos + event.wage < 0:
#                 tmp = heappop(Q0)
#                 f = min(tmp.count, event.count - flow)
#                 res += f * (tmp.cost + event.pos + event.wage)
#                 flow += f
#                 tmp.count -= f
#                 heappush(Q1, Pair(-event.pos * 2 - tmp.cost, f))
#                 if tmp.count:
#                     heappush(Q0, tmp)
#             if flow:
#                 heappush(Q0, Pair(-event.pos - event.wage, flow))
#             if flow < event.count:
#                 heappush(Q1, Pair(-event.pos + event.wage, event.count - flow))
#     return res
