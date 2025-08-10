# abc411-D - Conflict 2
# https://atcoder.jp/contests/abc411/tasks/abc411_d
# 有 1 台服务器和 N 台 PC，每台机器各维护一个字符串，初始均为空。
# 依次处理 Q 条指令：
#
# 1 p  把 PC p 的字符串设为当前服务器字符串
# 2 p s 在 PC p 的字符串末尾追加字符串 s
# 3 p  把服务器字符串设为当前 PC p 的字符串
# 在所有指令执行完后输出服务器字符串。
# 限制
# • 1 ≤ N,Q ≤ 2×10⁵
# • 所有追加字符串 s 的总长度 ≤ 10⁶
# 需要 O(Q+|总追加字符|) 解决。
#
# https://zhuanlan.zhihu.com/p/1919899313892488278
#
# 操作1：将PC的当前节点指向服务器的当前节点（不创建新节点）。
# 操作2：创建新节点，存储追加的字符串，父节点指向PC的当前节点，并更新PC的当前节点。
# 操作3：将服务器的当前节点指向PC的当前节点（不创建新节点）。


from PersistentStack import PersistentStack

N, Q = map(int, input().split())
pc = [PersistentStack() for _ in range(N + 1)]
server = pc[0]  # 服务器字符串（开始为空）

for _ in range(Q):
    cmd = input().split()
    t = cmd[0]
    if t == "1":
        p = int(cmd[1])
        pc[p] = server
    elif t == "2":
        p, s = int(cmd[1]), cmd[2]
        top = pc[p]
        top = top.push(s)
        pc[p] = top
    else:
        p = int(cmd[1])
        server = pc[p]


res = []
cur = server
while cur is not None and cur.value is not None:
    res.append(cur.value)
    cur = cur.pre

print("".join(reversed(res)))
