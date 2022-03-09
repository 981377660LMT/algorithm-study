# 整个城市所有的道路都是双向车道,道路的两个方向均需要铲雪
# 铲雪车只能把它开过的地方（车道）的雪铲干净，无论哪儿有雪，铲雪车都得从停放的地方出发，游历整个城市的街道。
# 现在的问题是：最少要花多少时间去铲掉所有道路上的雪呢？
'''
按照题意，所有的道路都是双向车道,图中必然存在欧拉回路，把所有路长度加起来，除以铲雪
速度即可
'''

import sys
import math

s = sys.stdin.readline()
total_len = 0
speed = 20000 / 60
while True:
    try:
        s = sys.stdin.readline()
    except Exception:
        break
    if s == '':
        break

    a1, b1, a2, b2 = map(int, s.split())
    total_len += math.sqrt((a2 - a1) ** 2 + (b2 - b1) ** 2)

min_val = int(total_len * 2 / speed + 0.5)
h, m = min_val // 60, min_val % 60
print('{:d}:{:02d}'.format(h, m))

