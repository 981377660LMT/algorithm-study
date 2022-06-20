N = int(input())

points = [(0, 0, 0)]
for _ in range(N):
    time, x, y = map(int, input().split())
    points.append((time, x, y))

for pre, cur in zip(points, points[1:]):
    t1, x1, y1 = pre
    t2, x2, y2 = cur
    diff = t2 - t1
    dist = abs(x1 - x2) + abs(y1 - y2)
    if dist > diff or (diff - dist) & 1:
        print("No")
        exit()
print("Yes")
