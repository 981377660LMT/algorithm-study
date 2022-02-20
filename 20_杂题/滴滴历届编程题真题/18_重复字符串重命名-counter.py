from collections import Counter


chars = input().split(",")
counter = Counter()
res = []
for char in chars:
    if char in counter:
        res.append(f"{char}_{counter[char] - 1}")
        counter[char] += 1
    else:
        res.append(char)
        counter[char] = 1
print(res)

