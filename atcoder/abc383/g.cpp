#pragma GCC target("arch=skylake-avx512")
#pragma GCC target("avx2")
#pragma GCC optimize("O3")
#pragma GCC target("sse4")
#pragma GCC optimize("unroll-loops")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

#include <bits/stdc++.h>
using namespace std;

int main(){

    int N,K;
    cin >> N >> K;
    vector<long long> A(N);
    for (int i=0; i<N; i++) cin >> A[i];

    if (K == 1) {
        sort(A.begin(), A.end(), greater<long long>());
        vector<long long> prefix(1,0);
        for (auto &x : A) prefix.push_back(prefix.back()+x);
        int M = N;
        for (int i=1; i<=M; i++){
            cout << prefix[i] << (i==M?'\n':' ');
        }
        return 0;
    }

    int length = N-K+1;
    vector<long long> B(length);
    {
        long long cur = 0;
        for (int i=0; i<K; i++) cur += A[i];
        B[0] = cur;
        for (int i=K; i<N; i++){
            cur += A[i]-A[i-K];
            B[i-K+1] = cur;
        }
    }

    int M = N/K;
    vector<long long> dp_prev(length+1, 0);
    vector<long long> dp_curr(length+1, -1000000000000000000LL);

    for (int i=1; i<=length; i++){
        if (i == 1) dp_curr[i] = B[i-1];
        else dp_curr[i] = max(dp_curr[i-1], B[i-1]);
    }
    dp_prev = dp_curr;
    vector<long long> ans(M+1,0);
    ans[1] = dp_prev[length];

    for (int j=2; j<=M; j++){
        for (int i=0; i<=length; i++){
            dp_curr[i] = -1000000000000000000LL;
        }

        deque<int> dq;
        long long prefix_max_curr = -1000000000000000000LL;

        for (int i=1; i<=length; i++){
            if (i>1) prefix_max_curr = max(prefix_max_curr, dp_curr[i-1]);
            else prefix_max_curr = -1000000000000000000LL;

            if (i-K >= 0) {
                long long val_to_add = dp_prev[i-K];
                while(!dq.empty() && dp_prev[dq.back()] <= val_to_add) dq.pop_back();
                dq.push_back(i-K);
            }

            while(!dq.empty() && dq.front() < i-K) dq.pop_front();

            long long candidate = prefix_max_curr;
            if (!dq.empty()) {
                candidate = max(candidate, dp_prev[dq.front()] + B[i-1]);
            }

            dp_curr[i] = candidate;
        }

        dp_prev = dp_curr;
        ans[j] = dp_prev[length];
    }

    for (int i=1; i<=M; i++){
        cout << ans[i] << (i==M?'\n':' ');
    }

    return 0;
}
