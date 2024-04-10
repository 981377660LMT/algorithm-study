#include<iostream>
#include<queue>
#include<cstring>
#define int long long
using namespace std;
const int maxn=50010;
const int maxm=1000010;
const int inf=1e9;
inline int read(){
	int x=0,f=1;
	char ch=getchar();
	while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
	while(ch>='0'&&ch<='9'){x=(x<<3)+(x<<1)+(ch-48);ch=getchar();}
	return x*f;
}

int n,m,s,num;
int f[maxn];
int fd(int x){
	if(x==f[x])return x;
	return f[x]=fd(f[x]);
}
struct ask{
	int u,v,uu,vv,w;
}a[maxm];
int cnt;
int head[maxm*3],tot;
struct nd{
	int nxt,to,w;
}e[maxm*12];
void add(int u,int v,int w){
	e[++tot]={head[u],v,w};
	head[u]=tot;
}
int dep[maxn],to[maxn][20];
int in[maxn][20],out[maxn][20];
void dfs(int u,int fa){
	to[u][0]=fa;in[u][0]=out[u][0]=u;
	for(int i=1;(1<<i)<dep[u];i++){
		to[u][i]=to[to[u][i-1]][i-1];
	}
	for(int i=1;(1<<i)<=dep[u];i++){
		in[u][i]=++num;out[u][i]=++num;
		add(out[u][i-1],out[u][i],0);
		add(in[u][i],in[u][i-1],0);
		add(out[to[u][i-1]][i-1],out[u][i],0);
		add(in[u][i],in[to[u][i-1]][i-1],0);
	}
	for(int i=head[u];i;i=e[i].nxt){
		int v=e[i].to;
		if(v!=fa&&v<=n){
			dep[v]=dep[u]+1;
			dfs(v,u);
		}
	}
}
int lg[maxn];
int lca(int u,int v){
	if(u==v)return u;
	if(dep[u]<dep[v])swap(u,v);
	for(int i=lg[dep[u]];i>=0;i--)if(dep[to[u][i]]>=dep[v])u=to[u][i];
	if(u==v)return u;
	for(int i=lg[dep[u]];i>=0;i--)if(to[u][i]!=to[v][i])u=to[u][i],v=to[v][i];
	return to[u][0];
}
int kfa(int u,int k){
	int j=0;
	while(k){
		if(k&1)u=to[u][j];    
		k>>=1;
		++j;
	}
	return u;
}
void build(int u,int v,int w,int t){
	int k=lg[dep[u]-dep[v]+1];
	if(!t)add(out[u][k],num,w);
	else add(num,in[u][k],w);
	u=kfa(u,dep[u]-dep[v]+1-(1<<k));
	if(!t)add(out[u][k],num,w);
	else add(num,in[u][k],w);
}
int dis[maxm*3];
bool vis[maxm*3];
struct Dis{
	int dis,id;
	bool operator <(const Dis&tmp)const{return dis>tmp.dis;}
};
priority_queue<Dis> q;

int T;
signed main(){
	//	freopen(".in","r",stdin);
	//	freopen(".out","w",stdout);
	
	n=read();m=read();s=read();
	lg[1]=0;for(int i=2;i<=n;i++)lg[i]=lg[i>>1]+1;
	for(int i=1;i<=n;i++)f[i]=i;
	while(m--){
		int opt,u,v,w,uu,vv;opt=read();
		if(opt==1){
			u=read();v=read();uu=read();vv=read();w=read();
			if(fd(u)!=fd(v)||fd(uu)!=fd(vv))continue;
			a[++cnt]={u,v,uu,vv,w};
		}
		else{
			u=read();v=read();w=read();
			if(fd(u)==fd(v))continue;
			add(u,v,w);add(v,u,w);
			f[fd(u)]=fd(v);
		}
	}
	num=n;
	for(int i=1;i<=n;i++){
		if(!dep[i]){
			dep[i]=1;
			dfs(i,0);
		}
	}
	for(int i=1;i<=cnt;i++){
		num++;
		int tp1=lca(a[i].u,a[i].v),tp2=lca(a[i].uu,a[i].vv);
		build(a[i].u,tp1,0,0);build(a[i].v,tp1,0,0);
		build(a[i].uu,tp2,a[i].w,1);build(a[i].vv,tp2,a[i].w,1);
	}
	memset(dis,0x3f,sizeof(dis));
	dis[s]=0;q.push({0,s});
	while(!q.empty()){
		int u=q.top().id;q.pop();
		if(vis[u])continue;
		vis[u]=1;
		for(int i=head[u];i;i=e[i].nxt){
			int v=e[i].to;
			if(dis[v]>dis[u]+e[i].w){
				dis[v]=dis[u]+e[i].w;
				q.push({dis[v],v});
			}
		}
	}
	for(int i=1;i<=n;i++){
		if(dis[i]>=0x3f3f3f3f)printf("-1 ");
		else printf("%d ",dis[i]);
	}
}