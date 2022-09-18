_ = calCircle1(x1, y1, x2, y2, r)
                if x0 is None or y0 is None:
                    continue
                res = max(
                    res,
                    sum(
                        (x - x0) * (x - x0) + (y - y0) * (y - y0) <= r * r + EPS for x, y in points
                    ),
                )