n, m = map(int, input().split())
cmax = 1 << m
mask = cmax - 1
horizontal = [0]
vertical = [0]

# dp处理行的状态，不用dfs
for i in range(1, cmax):
    lowbit = i & -i
    horizontal.append(horizontal[i - lowbit] + (1 if (lowbit << 1) & i == 0 else 0))
    vertical.append(vertical[i - lowbit] + 1)

dp = [10001] * cmax
dp[0] = 0

for i in range(n):
    cs = int(''.join(['1' if char == '#' else '0' for char in input()]), 2)
    for ci in range(mask, -1, -1):
        if ci & cs == ci:
            cans = dp[0] + vertical[ci]
            pre = ci
            while pre:
                tmp = dp[pre] + vertical[pre ^ ci]
                if tmp < cans:
                    cans = tmp
                pre = (pre - 1) & ci
            tmp = cans + horizontal[ci ^ cs]
            if tmp < 10001:
                cans = tmp

            cj = ci ^ mask
            while cj:
                lowbit = cj & -cj
                cj ^= lowbit
                desu = dp[ci | lowbit]
                if desu < cans:
                    cans = desu
            dp[ci] = cans
        else:
            dp[ci] = 10001

print(dp[0])
