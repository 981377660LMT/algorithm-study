#include <bits/stdc++.h>
using namespace std;

const int MOD=998244353;

int H,W;
vector<string> S;
int dx[4]={1,-1,0,0};
int dy[4]={0,0,1,-1};

int main(){
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    cin>>H>>W;
    S.resize(H);
    for (int i=0;i<H;i++) cin>>S[i];

    int n=H*W;
    auto id=[&](int r,int c){return r*W+c;};

    for (int i=0;i<H;i++){
        for (int j=0;j<W;j++){
            if (S[i][j]=='1' || S[i][j]=='2' || S[i][j]=='3'){
                int c1=S[i][j]-'0';
                for (int k=0;k<4;k++){
                    int ni=i+dx[k],nj=j+dy[k];
                    if(ni<0||ni>=H||nj<0||nj>=W) continue;
                    if (S[ni][nj]=='1' || S[ni][nj]=='2' || S[ni][nj]=='3'){
                        int c2=S[ni][nj]-'0';
                        if (c1==c2) {
                            cout<<0<<"\n";
                            return 0;
                        }
                    }
                }
            }
        }
    }

    vector<vector<int>> g(n);
    for (int i=0;i<H;i++){
        for (int j=0;j<W;j++){
            int u=id(i,j);
            for (int k=0;k<4;k++){
                int ni=i+dx[k],nj=j+dy[k];
                if(ni<0||ni>=H||nj<0||nj>=W) continue;
                int v=id(ni,nj);
                if(u<v){
                    g[u].push_back(v);
                    g[v].push_back(u);
                }
            }
        }
    }

    vector<int> color(n,-1); 

    vector<bool> fixedColor(n,false);
    for (int i=0;i<H;i++){
        for (int j=0;j<W;j++){
            int u=id(i,j);
            if(S[i][j]=='1' || S[i][j]=='2' || S[i][j]=='3'){
                fixedColor[u]=true;
                color[u]= (S[i][j]-'1');
            }
        }
    }

    vector<bool> visited(n,false);
    auto fpow=[&](long long base,long long exp)->long long{
        long long res=1%MOD; long long cur=base%MOD;
        while(exp>0){
            if(exp&1) res=(res*cur)%MOD;
            cur=(cur*cur)%MOD;
            exp>>=1;
        }
        return res;
    };

    long long ans=1;
    for (int start=0;start<n;start++){
        if(!visited[start]){
            vector<int> comp;
            queue<int>q;
            q.push(start);visited[start]=true;
            comp.push_back(start);
            while(!q.empty()){
                int u=q.front();q.pop();
                for (auto &nx:g[u]){
                    if(!visited[nx]){
                        visited[nx]=true;
                        comp.push_back(nx);
                        q.push(nx);
                    }
                }
            }

            int fixedCount=0;
            int freeCount=0;
            for (auto &x:comp){
                if(fixedColor[x]) fixedCount++;
                else freeCount++;
            }

            if(fixedCount==0){
                int sz=(int)comp.size();
                long long ways=3LL*fpow(2,sz-1)%MOD;
                ans=(ans*ways)%MOD;
            } else {


                long long ways=fpow(2,freeCount)%MOD;
                ans=(ans*ways)%MOD;
            }
        }
    }

    cout<<ans%MOD<<"\n";
    return 0;
}
