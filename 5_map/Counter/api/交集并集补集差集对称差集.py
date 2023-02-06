# https://blog.csdn.net/wizblack/article/details/78796557
# 集合/Counter里的五种操作
#
# 交集
# 并集
# 补集
# 差集
# 对称差集


from collections import Counter


s1, s2 = set("abaac"), set("bcd")
counter1, counter2 = Counter("abaac"), Counter("bcd")


# 1. 交集
print("1.交集:")
print(s1 & s2, counter1 & counter2)

# 2. 并集
print("2.并集:")
print(s1 | s2, counter1 | counter2)  # !取并后counter的频率是两者中较大的那个
print(counter1 + counter2)  # !counter特有加法

# 3. 差集、补集(多出来的元素/缺少的元素)
print("3.差集、补集:")
print(s1 - s2, counter1 - counter2)
print(s2 - s1, counter2 - counter1)

# 4. 对称差集(不同的元素)
print("4.对称差集:")
print(s1 ^ s2, (counter1 - counter2) + (counter2 - counter1))
