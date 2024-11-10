# import sys
# from sortedcontainers import SortedList

# sys.setrecursionlimit(int(1e6))
# input = lambda: sys.stdin.readline().rstrip("\r\n")

# Q = int(input())
# heights = SortedList()
# totalIncrease = 0
# outputs = []

# for _ in range(Q):
#     query = input().split()
#     if query[0] == "1":
#         heights.add(0 - totalIncrease)
#     elif query[0] == "2":
#         T = int(query[1])
#         totalIncrease += T
#     elif query[0] == "3":
#         H = int(query[1]) - totalIncrease
#         index = heights.bisect_left(H)
#         harvested = len(heights) - index
#         outputs.append(str(harvested))
#         for i in range(index, len(heights)):
#             heights.pop()

# print("\n".join(outputs))


from sortedcontainers import SortedList


sl = SortedList()


sl.add(-1)
sl.add(4)
sl.add(3)

print(sl[:])
