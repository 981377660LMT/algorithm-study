sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

MAPPING = {"A": "BC", "B": "CA", "C": "AB"}


def main() -> None:
    def dfs(t: int, k: int) -> str:
        if t == 0:
            return s[k - 1]
        if k == 0:
            mod_ = t % 3
            res = s[0]
            for _ in range(mod_):
                res = MAPPING[res][0]
            return res
        else:
            pre = dfs(t - 1, k >> 1)
            return MAPPING[pre][not k & 1]

    s = input()
    q = int(input())
    for _ in range(q):
        t, k = map(int, input().split())
        # print(["a", "b", "c"][dfs(t, k)])
        print(dfs(t, k))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
