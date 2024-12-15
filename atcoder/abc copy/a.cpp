#include <bits/stdc++.h>
using namespace std;

inline int v2(long long x) {
    int c=0;
    while((x & 1LL)==0) {
        x >>=1;
        c++;
    }
    return c;
}

int main(){
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    int N; cin >> N;
    vector<long long> A(N);
    for (int i=0; i<N; i++) cin >> A[i];

    const int MAX_T=25;
    const int MAX_R=24;
    vector<vector<long long>> groups(MAX_T);
    for (auto &x: A) {
        int t=v2(x);
        long long O = x>>t;
        groups[t].push_back(O);
    }

    vector<long long> sumO(MAX_T,0), cntO(MAX_T,0);
    __int128 diag=0;
    for (int t=0; t<MAX_T; t++) {
        cntO[t]= (long long)groups[t].size();
        for (auto &o:groups[t]) {
            sumO[t]+=o;
            diag+=o;
        }
    }

    __int128 full_matrix=0;

    for (int t=0; t<MAX_T; t++){
        auto &g=groups[t];
        int m=(int)g.size();
        if (m==0) continue;
        if (m==1) {
            full_matrix += sumO[t];
            continue;
        }

        sort(g.begin(), g.end());
        vector<__int128> C(MAX_R+2,0), S(MAX_R+2,0);

        for (int r=1; r<=MAX_R; r++){
            int modBase=(1<<r);
            vector<long long> freq(modBase,0), sumVal(modBase,0);
            for (int i=0; i<m; i++){
                long long o=g[i];
                int rem=(int)(o&(modBase-1));
                int need=(modBase - rem)&(modBase-1);
                long long pairCount=freq[need];
                long long pairSum=sumVal[need];
                if (pairCount>0) {
                    __int128 addC=pairCount;
                    __int128 addS= pairSum + (__int128)o*pairCount;
                    C[r]+=addC;
                    S[r]+=addS;
                }
                freq[rem]++;
                sumVal[rem]+=o;
            }
        }

        for (int r=MAX_R; r>=1; r--){
            C[r]-=C[r+1];
            S[r]-=S[r+1];
        }

        __int128 group_i_j_sum=0; 
        for (int r=1; r<=MAX_R; r++){
            __int128 addVal = S[r] >> r;
            group_i_j_sum += addVal;
        }

        full_matrix += sumO[t] + 2*group_i_j_sum;
    }

    for (int ti=0; ti<MAX_T; ti++){
        if (cntO[ti]==0) continue;
        for (int tj=ti+1; tj<MAX_T; tj++){
            if (cntO[tj]==0) continue;
            __int128 val = (__int128)cntO[tj]*sumO[ti] 
                         + (__int128)((long long)1<<(tj-ti))*cntO[ti]*sumO[tj];
            full_matrix += 2*val;
        }
    }

    __int128 tmp = full_matrix + diag;
    long long ans = (long long)(tmp>>1);
    cout << ans << "\n";

    return 0;
}
