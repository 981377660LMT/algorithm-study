# 给定两个项链的表示，判断他们经过旋转是否可能是一条项链。
from 最小表示法 import findIsomorphic

s1, s2 = input(), input()

min1, min2 = findIsomorphic(s1), findIsomorphic(s2)
if min1 != min2:
    print('No')
else:
    print('Yes')
    print(min1)

