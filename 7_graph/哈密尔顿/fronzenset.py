# frozenset() 返回一个冻结的集合，冻结后集合不能再添加或删除任何元素。
hashableSet = frozenset([1, 2, 3, 4, 5])
a = {}
a[hashableSet] = 1
print(a[hashableSet])


# frozenset与元组的区别是 可以变化 可以添加 可以删除元素
# fronzenset会把元素哈希
b = hashableSet - {1}
print(isinstance(b, frozenset))

