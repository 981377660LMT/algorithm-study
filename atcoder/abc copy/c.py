import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    a, b, c, d, e = map(int, input().split())

    problems = [("A", a), ("B", b), ("C", c), ("D", d), ("E", e)]

    from itertools import combinations

    names = []
    for length in range(1, 6):
        for combo in combinations(problems, length):
            name = "".join(p[0] for p in combo)
            score = sum(p[1] for p in combo)
            names.append((name, score))

    names.sort(key=lambda x: (-x[1], x[0]))

    for n, s in names:
        print(n)
