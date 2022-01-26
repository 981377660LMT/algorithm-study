from math import factorial
from itertools import permutations

n = int(input())
arr = input().split()

res = []
for p in permutations(arr):
    cand = list(p)
    # b后不能有a
    firstB = next((i for i, char in enumerate(cand) if char == 'B'), None)
    if firstB is not None and 'A' in cand[firstB:]:
        continue
    res.append(cand)

for item in res:
    print('-'.join(item), end=' ')
print(factorial(n) // 2)

