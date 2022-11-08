s1 = set(["a", "b", "c", "d", "e"])
s2 = set(["a", "b"])

# !1. 子集 注意Counter不能这么用 Counter判断子集应该`交小并大`
print(s1 >= s2)

# 2. 取差集(删除公共的元素)
print(s1 - s2)
print(s1.difference(s2))
# s1.difference_update(s2)
# print(s1)

# 3. 对称差 (不相交的部分) 删除相同，合并不同
s3 = set(["x", "y", "a", "b"])
print(s1.symmetric_difference(s3))
s1.symmetric_difference_update(s3)
print(s1)
