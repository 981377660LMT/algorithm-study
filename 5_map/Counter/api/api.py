from collections import Counter

c = Counter("aabbbccdddddeefggggg")

# 1. 子集关系用交/并集表示 (交小并大)
sub_c = Counter("a")
print(sub_c & c == sub_c)
print(sub_c | c == c)

# 2. 合并取最大频率 相交取最小频率
or_a = Counter("abbbbb") | c
print(or_a)

# 3.减(差集/补集) 小于等于0直接消除
diff_b = Counter("aa")
print(c - diff_b)

# 4. 空counter
print(not Counter())

# 5.freq总和
print(len(list(c.elements())))
print(sum(c.values()))
