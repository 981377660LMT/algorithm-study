from collections import Counter

w = 'qwertyuqqq'
print(w.replace('q', 'Q', 2))

a = Counter('aass')
b = Counter('as')
# print(a > b)
print(a.subtract(b))
print(a.most_common(), b)

