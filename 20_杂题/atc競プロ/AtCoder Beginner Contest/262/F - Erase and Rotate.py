# n<=2e5
# k<=n-1
# !最多可以进行k次erase和rotate操作 求全排列最小字典序 (其实就是进行k次 因为可以一直删末尾)
# erase:删除P的一个字母
# rotate:将P的末尾字母移到P的开头


# !1.回転を行わない場合
# 首部字母为p0-pk中的某个最小值 pi
# 第二个字母为pi+1-pk+1中的某个最小值 pj
# 第三个字母...
# 暴力查找与更新O(n^2) 需要线段树维护

# !2.回転を 1 回以上行う場合
# 不能删除rotate过来的项 因为这会多耗费一步操作
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def main() -> None:
    n, k = map(int, input().split())
    perm = list(map(int, input().split()))  # 1-n 的某个排列


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()


# TODO
