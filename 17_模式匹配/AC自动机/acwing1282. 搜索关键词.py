# 给定 n 个长度不超过 50 的由小写英文字母组成的单词，以及一篇长为 m 的文章。
# 请问，其中有多少个单词在文章中出现了。
# 注意：每个单词不论在文章中出现多少次，仅累计 1 次。


t = int(input())

# 向Trie中插入节点
def insert(s):
    if s == '':
        return
    global idx
    cur = 0
    for c in s:
        t = ord(c) - ord('a')
        if trie[cur][t] == 0:
            trie[cur][t] = idx
            idx += 1
            trie.append([0] * 26)
            cnt.append(0)
            ne.append(0)
        cur = trie[cur][t]
    cnt[cur] += 1


# 建立AC自动机
def build():
    queue = []
    hh, tt = 0, -1
    for i in range(26):
        if trie[0][i]:
            queue.append(trie[0][i])
            tt += 1
    while hh <= tt:
        t = queue[hh]
        hh += 1  # i - 1.
        for i in range(26):
            c = trie[t][i]  # i
            if c == 0:
                continue
            j = ne[t]
            while j and trie[j][i] == 0:
                j = ne[j]
            if trie[j][i]:
                j = trie[j][i]
            ne[c] = j
            queue.append(c)
            tt += 1


for _ in range(t):
    trie, cnt, ne = [[0] * 26], [0], [0]
    idx = 1
    n = int(input())

    for _ in range(n):
        s = input().strip()
        insert(s)

    build()

    # 读入文章开始查询
    string = input().strip()
    m = len(string)
    j = 0
    res = 0
    for i in range(m):
        t = ord(string[i]) - ord('a')
        while j and trie[j][t] == 0:
            j = ne[j]
        if trie[j][t]:
            j = trie[j][t]
        p = j
        while p:
            res += cnt[p]
            cnt[p] = 0
            p = ne[p]
    print(res)

