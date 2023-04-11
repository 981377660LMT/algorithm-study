from collections import defaultdict
from functools import reduce


trie = lambda: defaultdict(trie)
trie = trie()

# Flatten 2D array (matrix)
a = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
print(sum(a, []))


# Array index trick
# 检查回文
str = "abcba"
print(all(str[i] == str[~i] for i in range(len(str))))

# Chained boolean comparisons
if "a" < "c" > "b":
    print(1)
