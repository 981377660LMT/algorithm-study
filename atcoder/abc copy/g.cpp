
#include <bits/stdc++.h>
using namespace std;

struct Fenw {
    int n; 
    vector<long long> fenw;
    Fenw(int n=0):n(n),fenw(n+1,0){}
    void init(int n_){n=n_;fenw.assign(n+1,0);}
    void update(int i,long long v){
        for(;i<=n;i+=i&(-i)) fenw[i]+=v;
    }
    long long query(int i){
        long long s=0;for(;i>0;i-=i&(-i)) s+=fenw[i];return s;
    }
    long long rangeQuery(int l,int r){return query(r)-query(l-1);}
};

int main(){

    int N; cin>>N;
    vector<long long>A(N+1),B(N+1);
    for(int i=1;i<=N;i++) cin>>A[i];
    for(int i=1;i<=N;i++) cin>>B[i];
    
    vector<long long>PA(N+1,0),PB(N+1,0);
    for (int i=1;i<=N;i++){
        PA[i]=PA[i-1]+A[i];
        PB[i]=PB[i-1]+B[i];
    }
    
    int K; cin>>K;
    struct Query{int X,Y,id;};
    vector<Query>Q(K);
    for(int i=0;i<K;i++){
        cin>>Q[i].X>>Q[i].Y;
        Q[i].id=i;
    }
    
    vector<long long>vals;
    vals.reserve(2*N);
    for (int i=1;i<=N;i++){
        vals.push_back(A[i]);
        vals.push_back(B[i]);
    }
    sort(vals.begin(), vals.end());
    vals.erase(unique(vals.begin(), vals.end()), vals.end());
    auto compress=[&](long long x){
        return (int)(lower_bound(vals.begin(), vals.end(), x)-vals.begin()+1);
    };
    vector<int>cA(N+1), cB_(N+1);
    for(int i=1;i<=N;i++){
        cA[i]=compress(A[i]);
        cB_[i]=compress(B[i]);
    }
    int M=(int)vals.size();
    
    int block=(int)sqrt(N);
    sort(Q.begin(),Q.end(),[&](auto &a,auto &b){
        int ab=(a.X-1)/block, bb=(b.X-1)/block;
        if(ab!=bb)return ab<bb;
        return a.Y<b.Y;
    });
    
    Fenw fenwA_count(M), fenwA_sum(M), fenwB_count(M), fenwB_sum(M);
    fenwA_count.init(M); fenwA_sum.init(M);
    fenwB_count.init(M); fenwB_sum.init(M);
    
    long long sum_min=0;
    
    auto addA=[&](int i,int curY){
        long long cntL = fenwB_count.query(cA[i]);
        long long smL = fenwB_sum.query(cA[i]);
        long long contrib = smL + A[i]*(curY - cntL);
        sum_min += contrib;
        fenwA_count.update(cA[i],1);
        fenwA_sum.update(cA[i],A[i]);
    };
    auto removeA=[&](int i,int curY){
        long long cntL = fenwB_count.query(cA[i]);
        long long smL = fenwB_sum.query(cA[i]);
        long long contrib = smL + A[i]*(curY - cntL);
        sum_min -= contrib;
        fenwA_count.update(cA[i],-1);
        fenwA_sum.update(cA[i],-A[i]);
    };
    
    auto addB=[&](int j,int curX){
        long long cntL = fenwA_count.query(cB_[j]-1);
        long long smL = fenwA_sum.query(cB_[j]-1);
        long long contrib = smL + B[j]*(curX - cntL);
        sum_min += contrib;
        fenwB_count.update(cB_[j],1);
        fenwB_sum.update(cB_[j],B[j]);
    };
    auto removeB=[&](int j,int curX){
        long long cntL = fenwA_count.query(cB_[j]-1);
        long long smL = fenwA_sum.query(cB_[j]-1);
        long long contrib = smL + B[j]*(curX - cntL);
        sum_min -= contrib;
        fenwB_count.update(cB_[j],-1);
        fenwB_sum.update(cB_[j],-B[j]);
    };
    
    int curX=0, curY=0;
    vector<long long>ans(K);
    for (auto &q:Q) {
        int X=q.X, Y=q.Y;
        while(curY<Y){
            curY++;
            addB(curY,curX);
        }
        while(curY>Y){
            removeB(curY,curX);
            curY--;
        }
        while(curX<X){
            curX++;
            addA(curX,curY);
        }
        while(curX>X){
            removeA(curX,curY);
            curX--;
        }
        long long val = PA[X]*Y + PB[Y]*X - 2*sum_min;
        ans[q.id]=val;
    }
    
    for (auto v:ans) cout<<v<<"\n";
    return 0;
}
