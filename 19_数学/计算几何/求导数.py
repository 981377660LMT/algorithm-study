# # 求导数


from typing import Callable

EPS = 1e-8
INF = 1e18
N = 10  # batch size
ALPHA = 1.4  # step size


def differentiate(f: Callable[[float], float], x: float) -> float:
    h = 1e-2
    a = [[0.0] * N for _ in range(N)]
    res = INF
    err = INF
    a[0][0] = (f(x + h) - f(x - h)) / (2 * h)
    for i in range(1, N):
        h /= ALPHA
        a[0][i] = (f(x + h) - f(x - h)) / (2 * h)
        fac = ALPHA * ALPHA
        for j in range(1, i + 1):
            a[j][i] = (a[j - 1][i] * fac - a[j - 1][i - 1]) / (fac - 1.0)
            fac *= ALPHA * ALPHA
            errt = max(abs(a[j][i] - a[j - 1][i]), abs(a[j][i] - a[j - 1][i - 1]))
            if errt <= err:
                err = errt
                res = a[j][i]
                if err < EPS:
                    return res
        if abs(a[i][i] - a[i - 1][i - 1]) >= 2 * err:
            break
    return res


if __name__ == "__main__":

    def f(x: float) -> float:
        return -(x**2)

    def df(x: float) -> float:
        return -2 * x

    print(differentiate(lambda x: -(x**2), 1.0), df(1.0))
