from sortedcontainers import SortedList

lis = SortedList()
lis.add(1)
lis.add(1)
lis.add(2)
lis.add(4)

print(lis)
lis.discard(1)
print(lis)
print(2 in lis)  # logn
