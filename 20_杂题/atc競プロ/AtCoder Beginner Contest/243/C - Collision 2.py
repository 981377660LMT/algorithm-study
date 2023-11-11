# n 个人站在坐标系的第一象限及其边界，给定他们的坐标和行进方向L/R，判断是否会有人相遇
# !存在 RL就会碰撞

from collections import defaultdict
import sys
import os


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    s = input()
    people = [(*p, d) for p, d in zip(points, s)]
    adjMap = defaultdict(list)  # !按照纵坐标分组
    for x, y, d in people:
        adjMap[y].append((x, y, d))

    for group in adjMap.values():
        group.sort(key=lambda x: x[0])
        directions = "".join(char for _, _, char in group)
        if "RL" in directions:
            print("Yes")
            exit(0)

    print("No")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
