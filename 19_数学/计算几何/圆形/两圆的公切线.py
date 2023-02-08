# 两圆的公切线在圆1上的切点坐标


from math import sqrt
from typing import List, Tuple


def commonTangent1(
    x1: int, y1: int, r1: int, x2: int, y2: int, r2: int
) -> List[Tuple[float, float]]:
    """求出两圆的公切线在圆1上的切点坐标"""
    res = []
    xd = x2 - x1
    yd = y2 - y1

    rr0 = xd * xd + yd * yd
    if (r1 - r2) * (r1 - r2) <= rr0:
        cv = r1 - r2
        if rr0 == (r1 - r2) * (r1 - r2):
            res.append((x1 + r1 * cv * xd / rr0, y1 + r1 * cv * yd / rr0))
        else:
            sv = sqrt(rr0 - cv * cv)
            res.append((x1 + r1 * (cv * xd - sv * yd) / rr0, y1 + r1 * (sv * xd + cv * yd) / rr0))
            res.append((x1 + r1 * (cv * xd + sv * yd) / rr0, y1 + r1 * (-sv * xd + cv * yd) / rr0))
    if (r1 + r2) * (r1 + r2) <= rr0:
        cv = r1 + r2
        if rr0 == (r1 + r2) * (r1 + r2):
            res.append((x1 + r1 * cv * xd / rr0, y1 + r1 * cv * yd / rr0))
        else:
            sv = sqrt(rr0 - cv * cv)
            res.append((x1 + r1 * (cv * xd - sv * yd) / rr0, y1 + r1 * (sv * xd + cv * yd) / rr0))
            res.append((x1 + r1 * (cv * xd + sv * yd) / rr0, y1 + r1 * (-sv * xd + cv * yd) / rr0))
    return sorted(res)


Line = Tuple[Tuple[float, float], Tuple[float, float]]


def commonTangentLines(x1: int, y1: int, r1: int, x2: int, y2: int, r2: int) -> List[Line]:
    """求出两圆的公切线,公切线用两个点表示"""
    res = []
    xd = x2 - x1
    yd = y2 - y1

    rr0 = xd**2 + yd**2
    if (r1 - r2) ** 2 <= rr0:
        cv = r1 - r2
        if rr0 == (r1 - r2) ** 2:
            bx = r1 * cv * xd / rr0
            by = r1 * cv * yd / rr0
            res.append(
                (
                    (x1 + bx, y1 + by),
                    (x1 - yd + bx, y1 + xd + by),
                )
            )
        else:
            sv = (rr0 - cv**2) ** 0.5
            px = cv * xd - sv * yd
            py = sv * xd + cv * yd
            res.append(
                (
                    (x1 + r1 * px / rr0, y1 + r1 * py / rr0),
                    (x2 + r2 * px / rr0, y2 + r2 * py / rr0),
                )
            )
            qx = cv * xd + sv * yd
            qy = -sv * xd + cv * yd
            res.append(
                (
                    (x1 + r1 * qx / rr0, y1 + r1 * qy / rr0),
                    (x2 + r2 * qx / rr0, y2 + r2 * qy / rr0),
                )
            )
    if (r1 + r2) ** 2 <= rr0:
        cv = r1 + r2
        if rr0 == (r1 + r2) ** 2:
            bx = r1 * cv * xd / rr0
            by = r1 * cv * yd / rr0
            res.append(
                (
                    (x1 + bx, y1 + by),
                    (x1 - yd + bx, y1 + xd + by),
                )
            )
        else:
            sv = (rr0 - cv**2) ** 0.5
            px = cv * xd - sv * yd
            py = sv * xd + cv * yd
            res.append(
                (
                    (x1 + r1 * px / rr0, y1 + r1 * py / rr0),
                    (x2 - r2 * px / rr0, y2 - r2 * py / rr0),
                )
            )
            qx = cv * xd + sv * yd
            qy = -sv * xd + cv * yd
            res.append(
                (
                    (x1 + r1 * qx / rr0, y1 + r1 * qy / rr0),
                    (x2 - r2 * qx / rr0, y2 - r2 * qy / rr0),
                )
            )
    return res


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_G&lang=ja
    x1, y1, r1 = map(int, input().split())
    x2, y2, r2 = map(int, input().split())
    res = commonTangent1(x1, y1, r1, x2, y2, r2)
    if res:
        print("\n".join("%.08f %.08f" % p for p in res))
