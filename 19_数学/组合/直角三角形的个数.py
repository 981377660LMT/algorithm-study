from typing import List, Tuple


def cal(points: List[Tuple[int, int]]) -> int:
    """求点集组成的不同直角三角形个数"""
    n = len(points)
    res = 0
    for i in range(n):  # 枚举直角顶点
        x0, y0 = points[i]

        for j in range(n):
            if j == i:
                continue
            x1, y1 = points[j]
            for k in range(j + 1, n):
                if k == i or k == j:
                    continue
                x2, y2 = points[k]

                ij = (x1 - x0, y1 - y0)
                ik = (x2 - x0, y2 - y0)
                if ij[0] * ik[0] + ij[1] * ik[1] == 0:  # 垂直
                    res += 1
    return res


n = int(input())
nums = list(map(int, input().split()))
points = []
for i in range(0, 2 * n, 2):
    points.append((nums[i], nums[i + 1]))

print(cal(list(set(points))))
