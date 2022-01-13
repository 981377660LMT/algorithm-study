from sortedcontainers import SortedDict

sd = SortedDict({1: 2, 3: 0, 2: 6, 7: 9, 4: 33, 9: 88})


# irange是根据key的范围,返回迭代器切片
# islice是根索引,返回迭代器切片
print(sd, *sd.irange(3, 8), next(sd.islice(3, 5)))


# peekitem/keys/values 用于访问元素
higher = sd.bisect_right(6)
print(higher)
print(sd.peekitem(higher))
# print(sd.keys()[higher])
# print(sd.values()[higher])


# 删除
sd.pop(1)
sd.popitem(index=-1)
