# import re

# arr = re.split(r'[a-z]', input())
# print(len([v for v in arr if v]))
# 将字符串中的所有字母都看作是空白字符，统计出现的数的个数

res = 0
i = 0
s = input()
while i < len(s):
    if s[i].isdigit():
        while i < len(s) and s[i].isdigit():
            i += 1
        res += 1
    else:
        i += 1

print(res)

