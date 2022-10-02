"""问最多能看多少章的漫画
两本漫画可以换任意一章，问最多能看多少章漫画

二分答案比较方便
"""

from collections import Counter
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# print to file
# print("a" * 4000)
with open("out.txt", "w") as f:
    print("a" * 4000, file=f)
# 現在持っている単行本が 1 冊以下の場合、何もしない。
# そうでない場合、現在持っている単行本から 2 冊を選んで売り、代わりに好きな巻を選んで 1 冊買う
# その後、高橋君は『すぬけ君』を 1 巻、2 巻、3 巻、… と順番に読みます。ただし、次に読むべき巻を持っていない状態になった場合、(他の巻を持っているかどうかに関わらず)その時点で『すぬけ君』を読むのをやめます。
# 初め、高橋君は『すぬけ君』の単行本を N 冊持っています。i 番目の単行本は ai巻です。
if __name__ == "__main__":
    n = int(input())
    nums = sorted(map(int, input().split()))

    def check(mid: int) -> bool:
        counter = Counter(nums)
        money = 0
        for key in counter:
            if key > mid:
                money += counter[key]
            else:
                money += counter[key] - 1

        for i in range(1, mid + 1):
            if counter[i] == 0:
                if money >= 2:
                    money -= 2
                else:
                    return False
        return True

    left, right = 1, int(1e10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    print(right)

# 4000*'a'
# print to file
with open("out.txt", "w") as f:
    print(4000 * "a", file=f)
