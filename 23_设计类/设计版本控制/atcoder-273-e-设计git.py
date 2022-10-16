from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# ADD x : 整数 x を A の末尾に追加する。
# DELETE : A の末尾の要素を削除する。ただし、A が空である場合は何もしない。
# SAVE y : ノートの y ページ目に書かれている数列を消し、代わりに現在の A を y ページ目に書き込む。
# LOAD z : A をノートの z ページ目に書かれている数列で置き換える。
# Q 個のクエリを与えられる順に実行し、各クエリの実行直後における A の末尾の要素を出力してください。
# !持久化数组???

if __name__ == "__main__":
    note = defaultdict(int)  # 时间戳(版本)
    q = int(input())
    history = [[] for _ in range(q + 1)]
    book = defaultdict(int)
    index = -1
    curTime = 0
    for i in range(q):
        t, *args = input().split()
        if t == "ADD":
            x = int(args[0])
            history[index].append((i + 1, x))
            index += 1
            curTime += 1
        elif t == "DELETE":
            index = max(0, index - 1)
            curTime += 1
        elif t == "SAVE":
            page = int(args[0]) - 1
            book[page] = i
        elif t == "LOAD":
            page = int(args[0]) - 1
            curTime = book[page]


# 版本控制
# 快照数组 和那个好像不一样
