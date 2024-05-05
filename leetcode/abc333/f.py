import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 人の人が一列に並んでおり、人
# i は先頭から
# i 番目に並んでいます。

# 以下の操作を、列に並んでいる人が
# 1 人になるまで繰り返します。

# 先頭に並んでいる人を
# 2
# 1
# ​
#   の確率で列から取り除き、そうでない場合は列の末尾に移す。
# 人
# i=1,2,…,N それぞれについて、人
# i が最後まで列に並んでいる
# 1 人になる確率を
# mod 998244353 で求めて下さい。(取り除くかどうかの選択はすべてランダムかつ独立です。)

INV2 = (MOD + 1) // 2
if __name__ == "__main__":
    n = int(input())
    # void rd() {
    #   cin >> n >> m;
    # }

    # cin>>n>>m;
    # f[1][1]=1;
    # for(int i=2;i<=n;i++){
    # 	double s=0.5,t=0;
    # 	for(int j=2;j<=i;j++){
    # 		s=s/2;
    # 		t=t/2+f[i-1][j-1]/3;
    # 	}
    # 	f[i][i]=t/(1-s);
    # 	f[i][1]=f[i][i]/2;
    # 	for(int j=2;j<i;j++)
    # 		f[i][j]=f[i][j-1]/2+f[i-1][j-1]/2;
    # }printf("%.9lf",f[n][m]);

    dp = [[0] * (n + 1) for _ in range(n + 1)]
    dp[1][1] = 1
    for i in range(2, n + 1):
        dp[i][1] = 1
        s = INV2
        for j in range(2, i + 1):
            s = s * INV2 % MOD
            dp[i][j] = (dp[i - 1][j - 1] * INV2 + dp[i][j - 1] * INV2) % MOD
        dp[i][i] = dp[i][i - 1] * INV2 % MOD
        for j in range(2, i):
            dp[i][j] = (dp[i][j - 1] * INV2 + dp[i - 1][j - 1] * INV2) % MOD

    for i in range(1, n + 1):
        print(dp[n][i], end=" ")
