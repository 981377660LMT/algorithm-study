# 打表、预处理
# 埋め込みによる解法
# https://atcoder.jp/contests/abc363/editorial/10503


embedded = {}  # 打表
embedded[(1, (1,))] = 1


# 需要生成的代码
print(
    f"""
from collections import defaultdict
N, K = map(int, input().split())
S = input()
count = defaultdict(int)
for c in S:
    count[c] += 1
S_key = tuple(reversed(sorted(count.values())))
embedded = {embedded}
print(embedded[(K, S_key)])
"""
)
