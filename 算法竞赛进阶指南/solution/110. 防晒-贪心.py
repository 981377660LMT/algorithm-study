# 有 C 头奶牛进行日光浴，
# 第 i 头奶牛需要 minSPF[i] 到 maxSPF[i] 单位强度之间的阳光。

# 每头奶牛在日光浴前必须涂防晒霜，
# 防晒霜有 L 种，涂上第 i 种之后，
# 身体接收到的阳光强度就会稳定为 SPF[i]，第 i 种防晒霜有 cover[i] 瓶。

# 求最多可以满足多少头奶牛进行日光浴。
# 1≤C,L≤2500


c, l = map(int, input().split())
cows = []
for _ in range(c):
    cows.append(list(map(int, input().split())))  # minSPF  和 maxSPF

# sun-protection factor 防晒因子
spf = []
for _ in range(l):
    spf.append(list(map(int, input().split())))  # SPF 与 瓶数


# 奶牛按照minSPF降序排序, 然后依次选择防晒霜.
# 每头奶牛选择在它区间内SPF值尽可能大的防晒霜.小的spf要给别人用
cows.sort(reverse=True)
spf.sort(reverse=True)
res = 0
for i in range(c):
    for j in range(l):
        if spf[j][0] > cows[i][1]:
            continue
        if cows[i][0] > spf[j][0]:
            break

        # 找到了
        if spf[j][1]:
            spf[j][1] -= 1
            res += 1
            break

print(res)
