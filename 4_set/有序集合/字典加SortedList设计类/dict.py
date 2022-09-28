d = dict()
d["a"] = 1
d["c"] = 3
d["v"] = 2


# 获取最后一次插入的值(top)
print(next(reversed(d.values())))
# print(list(d.values())[-1])
