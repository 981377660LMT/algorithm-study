from sortedcontainers import SortedSet, SortedList, SortedDict

# 充当Java TreeSet
scores = SortedSet()
scores.add(3)
scores.add(3)
scores.add(4)

print(scores)

# 充当Java TreeMap
scores = SortedDict()
scores[3] = 2
scores[1] = 2

# scores.add(3)
# scores.add(4)
print(scores)

# # 充当C++ STL 的multiset(最强)
scores = SortedList(key=lambda x: -x)
scores.add(3)
scores.add(3)
scores.add(4)
scores.add(7)
scores.add(5)
print(*scores.islice(0, 3))  # 取前k项
print(*scores.islice(0, 100))  # 取前k项
# print(*scores.irange(3, 7))  # 取范围值
