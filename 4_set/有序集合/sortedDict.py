from sortedcollections import SortedDict

info = SortedDict({1: [1, 2, 3], 3: [7, 4], 7: [0, 9], 2: [4]})

print(*info.irange(2, 6))
print(list(info.islice(2, 6)))

