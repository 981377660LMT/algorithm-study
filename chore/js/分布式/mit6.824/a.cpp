#include <bits/stdc++.h>
using namespace std;

const int MOD = 998244353;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    int N,M; cin >> N >> M;
    vector<int> B(N);
    for (int i=0; i<N; i++) cin >> B[i];

    vector<int> dp(M+1,0), sum_runs(M+1,0);
    vector<int> dp_prev(M+1,0), sum_runs_prev(M+1,0);

    if (B[0] == -1) {
        for (int x=1; x<=M; x++) {
            dp[x] = 1;
            sum_runs[x] = 1;
        }
    } else {
        int c = B[0];
        dp[c] = 1;
        sum_runs[c] = 1;
    }

    for (int i=2; i<=N; i++) {
        dp_prev = dp;
        sum_runs_prev = sum_runs;
        for (int x=1; x<=M; x++) {
            dp[x]=0; sum_runs[x]=0;
        }

        vector<int> prefix_dp(M+1,0), prefix_sum_runs(M+1,0);
        for (int x=1; x<=M; x++){
            prefix_dp[x] = (prefix_dp[x-1] + dp_prev[x]) % MOD;
            prefix_sum_runs[x] = (prefix_sum_runs[x-1] + sum_runs_prev[x]) % MOD;
        }

        int total_dp = prefix_dp[M];
        int total_sum_runs = prefix_sum_runs[M];

        if (B[i-1] == -1) {
            for (int x=1; x<=M; x++){
                int ways = total_dp;
                long long val = (long long)total_sum_runs + total_dp - prefix_dp[x];
                val %= MOD;
                if (val<0) val += MOD;

                dp[x] = ways;
                sum_runs[x] = (int)val;
            }
        } else {
            int c = B[i-1];
            int ways = total_dp;
            long long val = (long long)total_sum_runs + total_dp - prefix_dp[c];
            val %= MOD;
            if (val<0) val += MOD;

            dp[c] = ways;
            sum_runs[c] = (int)val;
        }
    }

    int ans=0;
    if (B[N-1] == -1) {
        for (int x=1; x<=M; x++){
            ans += sum_runs[x];
            if (ans>=MOD) ans-=MOD;
        }
    } else {
        ans = sum_runs[B[N-1]] % MOD;
    }

    cout << ans << "\n";
    return 0;
}
