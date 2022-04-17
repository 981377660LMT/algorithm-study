# 海岸是一条无限长的直线，陆地位于海岸的一侧，海洋位于另外一侧
# 每个小岛都位于海洋一侧的某个点上。
# 雷达装置均位于海岸线上，且雷达的监测范围为 d，
# 当小岛与某雷达的距离不超过 d 时，该小岛可以被雷达覆盖。


def main():
    n, d = map(int, input().split())
    intervals = []
    for _ in range(n):
        x, y = map(int, input().split())
        if y > d:
            print(-1)
            exit()

        delta = (d ** 2 - y ** 2) ** 0.5
        left, right = x - delta, x + delta
        intervals.append((left, right))  # 覆盖这个岛屿需要的雷达的范围

    # 用最少的点覆盖区间
    intervals.sort(key=lambda x: x[1])

    res = 0
    preEnd = -int(1e20)
    for start, end in intervals:
        if start > preEnd:
            res += 1
            preEnd = end
    print(res)


main()
