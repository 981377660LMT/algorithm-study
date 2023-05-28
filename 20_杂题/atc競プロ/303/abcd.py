from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# あなたのパソコンのキーボードには、a キー・Shift キー・CapsLock キーの
# 3 種類のキーがあります。また、CapsLock キーにはランプが付いています。 初め、CapsLock キーのランプは OFF であり、パソコンの画面には空文字列が表示されています。

# あなたは、以下の
# 3 種類の操作のうち
# 1 つを選んで実行するということを
# 0 回以上何度でも行うことができます。

# X ミリ秒かけて a キーのみを押す。CapsLock キーのランプが OFF ならば画面の文字列の末尾に a が付け足され、ON ならば A が付け足される。
# Y ミリ秒かけて Shift キーと a キーを同時に押す。CapsLock キーのランプが OFF ならば画面の文字列の末尾に A が付け足され、 ON ならば a が付け足される。
# Z ミリ秒かけて CapsLock キーを押す。CapsLock キーのランプが OFF ならば ON に、ON ならば OFF に切り替わる。
# A と a からなる文字列
# S が与えられます。画面の文字列を
# S に一致させるのに必要な最短の時間は何ミリ秒かを求めてください。

dist = [[INF] * 2 for _ in range(len(s) + 1)]
dist[0][0] = 0
queue = deque()
queue.append((0, 0, 0))  # dist, caps, pos
while queue:
    curDist, curCaps, curPos = queue.popleft()
    if curPos == len(s):
        continue
    if dist[curPos][curCaps] < curDist:
        continue

    cand1 = curDist + z
    if dist[curPos][1 ^ curCaps] > cand1:
        dist[curPos][1 ^ curCaps] = cand1
        queue.append((cand1, 1 ^ curCaps, curPos))

    upper = int(s[curPos].isupper())
    cand2 = curDist + x
    if upper == curCaps:
        if dist[curPos + 1][curCaps] > cand2:
            dist[curPos + 1][curCaps] = cand2
            queue.append((cand2, curCaps, curPos + 1))

    cand3 = curDist + y
    if upper != curCaps:
        if dist[curPos + 1][curCaps] > cand3:
            dist[curPos + 1][curCaps] = cand3
            queue.append((cand3, curCaps, curPos + 1))

print(min(dist[-1]))


if __name__ == "__main__":
    x, y, z = map(int, input().split())
    s = input()

    dp = [0, INF]
    for c in s:
        target = 1 if c.isupper() else 0
        ndp = [INF, INF]
        for i in range(2):
            for j in range(2):
                # 按相同
                if j == target:
                    ndp[j] = min(ndp[j], dp[i] + x)
                else:
                    ndp[j] = min(ndp[j], dp[1^i] + z+

                # 按shift

                # 按caps
        dp = ndp
    print(min(dp))
