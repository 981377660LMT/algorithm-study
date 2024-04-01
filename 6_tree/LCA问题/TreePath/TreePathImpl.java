package template.graph;

public class TreePathImpl implements TreePath {
  DepthOnTree dot;
  KthAncestor ancestor;
  LcaOnTree lca;
  int a;
  int b;
  int c;

  public TreePathImpl(DepthOnTree dot, KthAncestor kthAncestor, LcaOnTree lca) {
    this.dot = dot;
    this.ancestor = kthAncestor;
    this.lca = lca;
  }

  public void init(int a, int b) {
    this.a = a;
    this.b = b;
    c = lca.lca(a, b);
  }

  public int length() {
    return dot.depth(a) + dot.depth(b) - 2 * dot.depth(c);
  }

  /**
   * a is 0-th, k <= length()
   * <p>
   * O(log_2n)
   */
  public int kthNodeOnPath(int k) {
    if (k <= dot.depth(a) - dot.depth(c)) {
      return ancestor.kthAncestor(a, k);
    }
    return ancestor.kthAncestor(b, length() - k);
  }

  // 某个点是否在路径上.
  @Override
  public boolean onPath(int u) {
    return lca.lca(u, c) == c && (lca.lca(u, a) == u || lca.lca(u, b) == u);
  }

  // TODO
  // 在一棵树上，对于路径 (x,y) 和路径 (u,v)，判断它们相交，等价于判断是否满足：
  // **lca(u,v) 在路径 (x,y) 上，或者 lca(x,y) 在路径 (u,v) 上**
  public boolean intersect(TreePath other) {
    return other.onPath(c) || onPath(other.c);
  }
}
#include <bits/stdc++.h>
#define __ ios::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define rep(i,a,b) for(int i = a; i <= b; i++)
#define LOG1(x1,x2) cout << x1 << ": " << x2 << endl;
#define LOG2(x1,x2,y1,y2) cout << x1 << ": " << x2 << " , " << y1 << ": " << y2 << endl;
#define LOG3(x1,x2,y1,y2,z1,z2) cout << x1 << ": " << x2 << " , " << y1 << ": " << y2 << " , " << z1 << ": " << z2 << endl;
typedef long long ll;
typedef double db;
const int N = 1e5+100;
const db EPS = 1e-9;
using namespace std;

int n,q,t,tot,head[N],f[N][25],d[N];
struct Edge{
	int to,next;
}e[2*N];

void add(int x,int y){
	e[++tot].to = y, e[tot].next = head[x], head[x] = tot;
}

void dfs(int x,int fa){
	f[x][0] = fa;
	for(int i = 1; (1<<i) <= d[x]; i++)
		f[x][i] = f[f[x][i-1]][i-1];
	for(int i = head[x]; i; i = e[i].next){
		int y = e[i].to;
		if(y == fa) continue;
		d[y] = d[x]+1; dfs(y,x);
	}
}

int lca(int x,int y){
	if(d[x] > d[y]) swap(x,y);
	for(int i = t; i >= 0; i--)
		if(d[f[y][i]] >= d[x]) y = f[y][i];
	if(x == y) return x;
	for(int i = t; i >= 0; i--)
		if(f[x][i] != f[y][i]) x = f[x][i], y = f[y][i];
	return f[x][0];
}

int main()
{
	scanf("%d%d",&n,&q);
	t = (int)(log(n)/log(2))+1;
	tot = 1;
	rep(i,1,n-1){
		int a,b; scanf("%d%d",&a,&b);
		add(a,b); add(b,a);
	}
	dfs(1,0);
	rep(i,1,q){
		int a,b,c,D; scanf("%d%d%d%d",&a,&b,&c,&D);
		int x[5];
		x[1] = lca(a,c), x[2] = lca(a,D), x[3] = lca(b,c), x[4] = lca(b,D);
		int p1 = 0,p2 = 0;
		rep(j,1,4)
			if(d[x[j]] > d[p1]) p2 = p1, p1 = x[j];
			else if(d[x[j]] > d[p2]) p2 = x[j];
		int h1 = lca(a,b), h2 = lca(c,D);
		if(p1 == p2){
			if(d[p1] < d[h1] || d[p1] < d[h2]) printf("0\n");
			else printf("1\n");
		}
		else{
			int ans = d[p1]+d[p2]-2*d[lca(p1,p2)]+1;
		 	printf("%d\n",ans);
		}
	}
	return 0;
}
