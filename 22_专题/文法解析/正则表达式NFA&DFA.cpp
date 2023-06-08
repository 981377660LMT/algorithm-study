#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<time.h>
#include<math.h>
#include<iostream>
#include<algorithm>
using namespace std;
#define N 505
#define nn 16  //N/32
#define P 300007
#define CH 10
typedef unsigned int uint;
struct node{
	node *next;int to;char c;
}*son[N],edge[N*5];
struct hash_node{
	uint cur[nn],nxt[nn];
	hash_node *next;
}*head[CH][P],buf[P];
uint next[N][CH][nn],mask[CH][nn],now[nn];
int q[N],n,tot;bool vis[N];char reg_s[N];
void ins_edge(int a,int b,int c){
	static int top=0;node *p=edge+(top++);
	p->to=b;p->c=c;p->next=son[a];son[a]=p;
}
inline int getlevel(char c){
	if (c=='+')return 0;
	else if (c=='*')return 2;
	else return 1;
}
void build_NFA(char *s,int l,int r,int a,int b){
	while (s[l]=='('){
		int top=0;bool flag=0;
		for (int i=l;i<r;++i){
			if (s[i]=='(')++top;
			if (s[i]==')')--top;
			if (!top){flag=1;break;}
		}
		if (flag)break;
		build_NFA(s,l+1,r-1,a,b);return;
	}
	if (l==r){ins_edge(a,b,s[l]-'0');return;}
	int top=0,loc=l;
	for (int i=l;i<=r;++i){
		if (top==0&&(loc==l||getlevel(s[i])<getlevel(s[loc])))loc=i;
		if (s[i]=='(')++top;
		if (s[i]==')')--top;
	}
	if (s[loc]=='+'){
		build_NFA(s,l,loc-1,a,b);
		build_NFA(s,loc+1,r,a,b);
		return;
	}
	if (s[loc]=='*'){
		int t1=tot++,t2=tot++;
		build_NFA(s,l,loc-1,t1,t2);
		ins_edge(a,t1,-1);ins_edge(t2,b,-1);
		ins_edge(t2,t1,-1);ins_edge(t1,t2,-1);
		return;
	}
	int tmp=tot++;
	build_NFA(s,l,loc-1,a,tmp);
	build_NFA(s,loc,r,tmp,b);
}
void insert(int c,uint _cur[],uint _nxt[]){
	uint code=0;for (int i=0;i<nn;++i)code^=_cur[i];code%=P;
	static int top=0;hash_node *p=buf+(top++);
	memcpy(p->cur,_cur,nn*4); memcpy(p->nxt,_nxt,nn*4);
	p->next=head[c][code];head[c][code]=p;
}
inline void read(char &c){for (c=getchar();c!=EOF&&(c<'0'||c>'9');c=getchar());}
int main()
{
	//freopen("1.in","r",stdin);
	//freopen("1.out","w",stdout);
	memset(son,0,sizeof(son));tot=0;
	scanf("%d%s",&n,&reg_s);int start=tot++,end=tot++;
	build_NFA(reg_s,0,strlen(reg_s)-1,start,end);
	for (int i=0;i<tot;++i)
		for (int j=0;j<n;++j){
			int head=0,tail=0;
			for (node *p=son[i];p;p=p->next)
				if (p->c==j)q[tail++]=p->to;
			if (!tail){
				mask[j][i/32]+=1u<<(i%32);
				continue;
			}
			memset(vis,0,sizeof(vis));vis[q[0]]=1;
			while (head<tail){
				int cur=q[head++];
				for (node *p=son[cur];p;p=p->next)
					if (!vis[p->to]&&p->c==-1){
						vis[p->to]=1;q[tail++]=p->to;
					}
			}
			uint *cur=next[i][j];
			for (int k=0;k<tail;++k){
				int a=q[k]/32,b=q[k]%32;
				cur[a]|=1u<<b;
			}
		}
	memset(vis,0,sizeof(vis));
	memset(now,0,sizeof(now));
	q[0]=start;vis[start]=1;int head=0,tail=1;
	while (head<tail){
		int cur=q[head++];
		for (node *p=son[cur];p;p=p->next)
			if (!vis[p->to]&&p->c==-1){
				vis[p->to]=1;q[tail++]=p->to;
			}
	}
	for (int i=0;i<tail;++i){
		int a=q[i]/32,b=q[i]%32;
		now[a]|=1u<<b;
	}
	int pos=1;char ch;
	for (read(ch);ch>='0'&&ch<='9';read(ch),++pos){
		for (int i=0;i<nn;++i)now[i]&=~mask[ch-'0'][i];
		uint code=0;for (int i=0;i<nn;++i)code^=now[i];code%=P;bool flag=0;
		for (hash_node *p=::head[ch-'0'][code];p;p=p->next){
			bool have=1;
			for (int i=0;i<nn;++i)
				if (p->cur[i]!=now[i]){have=0;break;}
			if (have){
				memcpy(now,p->nxt,nn*4);
				flag=1;break;
			}
		}
		if (!flag){
			uint _next[nn];memset(_next,0,sizeof(_next));
			for (int i=0;i<nn;++i){
				uint tmp=now[i];
				while (tmp>0){
					int loc=__builtin_ctz(tmp);uint *cur=next[i*32+loc][ch-'0'];
					for (int j=0;j<nn;++j)_next[j]|=cur[j];
					tmp-=1u<<loc;
				}
			}
			_next[start/32]|=1u<<start%32;insert(ch-'0',now,_next);memcpy(now,_next,nn*4);
		}
		if (now[end/32]&1u<<end%32)printf("%d ",pos);
	}
	printf("\n");
	system("pause");for (;;);
	return 0;
}


