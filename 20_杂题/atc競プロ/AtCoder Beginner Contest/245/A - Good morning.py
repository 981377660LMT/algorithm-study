import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    a, b, c, d = map(int, input().split())
    # 时分秒 比谁起床早
    if (a, b, 0) < (c, d, 1):
        print("Takahashi")
    else:
        print("Aoki")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
