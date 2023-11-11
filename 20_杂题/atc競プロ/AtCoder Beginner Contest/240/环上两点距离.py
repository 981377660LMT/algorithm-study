import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    a, b = map(int, input().split())
    dist = min((a - b) % 10, (b - a) % 10)  # 环上距离
    print(["No", "Yes"][dist == 1])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
