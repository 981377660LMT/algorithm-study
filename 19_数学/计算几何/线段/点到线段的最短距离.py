# 点到线段的最短距离
# https://tjkendev.github.io/procon-library/python/geometry/segment_line_point_distance.html


from typing import Tuple


# Shortest distance between a line segment and a point.
def segment_line_dist(p: Tuple[int, int], seg: Tuple[int, int, int, int]) -> float:
    def cross2(p, q):
        return p[0] * q[1] - p[1] * q[0]

    def dot2(p, q):
        return p[0] * q[0] + p[1] * q[1]

    def dist2(p):
        return p[0] * p[0] + p[1] * p[1]

    x0, x1 = p
    s1, s2, s3, s4 = seg
    z0 = (s3 - s1, s4 - s2)
    z1 = (x0 - s1, x1 - s2)
    if 0 <= dot2(z0, z1) <= dist2(z0):
        return abs(cross2(z0, z1)) / dist2(z0) ** 0.5
    z2 = (x0 - s3, x1 - s4)
    return min(dist2(z1), dist2(z2)) ** 0.5


# example
print(segment_line_dist((-1, -1), (0, 0, 1, 0)))
# => "1.4142135623730951"
print(segment_line_dist((0.5, 1), (0, 0, 1, 0)))
# => "1.0"
print(segment_line_dist((2, 2), (0, 0, 1, 0)))
# => "2.23606797749979"
