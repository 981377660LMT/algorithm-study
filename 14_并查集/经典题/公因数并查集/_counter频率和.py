from collections import Counter

c = Counter({1: 2, 2: 3})

# 2 5
print(len(c), len(list(c.elements())), c.total())
