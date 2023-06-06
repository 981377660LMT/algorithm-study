#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<time.h>
#include<math.h>
#include<iostream>
#include<algorithm>
#include<set>
using namespace std;
template<class Tkey,class Tvalue=bool>
struct Hash{
	struct node{
		Tkey first;Tvalue second;node *next;
		node():next(0){}
		node(const Tkey &x,const Tvalue &y,node *_next=0):first(x),second(y),next(_next){}
	};
	node **v,*key;size_t len,P,max_size;
	void Grow(){
		static double rate=1.7;Hash<Tkey,Tvalue> res(max_size*2,size_t(rate*max_size*2));
		for (size_t i=0;i<P;++i)
			for (node *j=v[i];j;j=j->next)res.insert(j->first,j->second);
		free();*this=res;
	}
	void build(size_t L,size_t p){
		P=p;len=0;max_size=L;++L;key=new node[L];
		v=new node*[p];memset(v,0,sizeof(node*)*p);
	}
	Hash(size_t len=3,size_t p=5){build(len,p);}
	Tvalue& insert(const Tkey &x,const Tvalue &y){
		if (len==max_size)Grow();size_t x1=x%P;
		key[++len]=node(x,y,v[x1]);v[x1]=key+len;
		return key[len].second;
	}
	Tvalue find(const Tkey &x){
		size_t x1=x%P;
		for (node *i=v[x1];i;i=i->next)
			if (i->first==x)return i->second;
		return 0;
	}
	void free(){delete[] key;delete[] v;}
};

struct VanEmdeBoasTree{
	#define M 2147483647
	struct node{
		int min,max,dep;node *aux;Hash<int,node*> son;
		node(int _dep):dep(_dep),min(M),max(-1),aux(0){}
		bool find(int x){
			if (x==min||x==max)return 1;
			if (x<min||x>max||!dep)return 0;
			int i=x>>dep;node *soni=son.find(i);
			return !soni?0:soni->find(x-(i<<dep));
		}
		int pred(int x){
			if (x>=max)return max;if (x<min)return -1;
			int i=x>>dep,hi=i<<dep,lo=x-hi;node *soni=son.find(i);
			if (soni&&lo>=soni->min)return hi+soni->pred(lo);
			int y=aux&&i>0?aux->pred(i-1):-1;
			return y==-1?min:(y<<dep)+son.find(y)->max;
		}
		int succ(int x){
			if (x<=min)return min;if (x>max)return M;
			int i=x>>dep,hi=i<<dep,lo=x-hi;node *soni=son.find(i);
			if (soni&&lo<=soni->max)return hi+soni->succ(lo);
			int y=aux?aux->succ(i+1):M;
			return y==M?max:(y<<dep)+son.find(y)->min;
		}
		void insert(int x){
			if (min>max){min=max=x;return;}
			if (min==max)
				if (x<min){min=x;return;}
				else if (x>max){max=x;return;}
			if (x<min)swap(x,min);if (x>max)swap(x,max);
			int i=x>>dep;node *soni=son.find(i);
			if (!soni)soni=son.insert(i,new node(dep>>1));
			if (soni->empty()){if (!aux)aux=new node(dep>>1);aux->insert(i);}
			soni->insert(x-(i<<dep));
		}
		void erase(int x){
			if (min==x&&max==x){min=M;max=-1;return;}
			if (x==min)
				if (!aux||aux->empty()){min=max;return;}
				else min=x=(aux->min<<dep)+son.find(aux->min)->min;
			if (x==max)
				if (!aux||aux->empty()){max=min;return;}
				else max=x=(aux->max<<dep)+son.find(aux->max)->max;
			int i=x>>dep;node *soni=son.find(i);soni->erase(x-(i<<dep));
			if (soni->empty())aux->erase(i);
		}
		bool empty()const{return min>max;}
	};
	node *root;size_t s;
	VanEmdeBoasTree(){s=0;root=new node(sizeof(int)*8/2);}
	bool find(int x){return root->find(x);}
	int pred(int x){return root->pred(x);}
	int succ(int x){return root->succ(x);}
	void insert(int x){if (!root->find(x))++s,root->insert(x);}
	void erase(int x){if (root->find(x))--s,root->erase(x);}
	int top(){return root->min;}
	void pop(){root->erase(root->min);}
	size_t size()const{return s;}
	#undef M
};
VanEmdeBoasTree a;
int main()
{
	//freopen("1.in","r",stdin);
	//freopen("1.out","w",stdout);
	int t1=clock();
	int n;scanf("%d",&n);
	for (int i=1;i<=n;++i){
		char c;int x;scanf(" %c%d",&c,&x);
		if (c=='I')a.insert(x);
		if (c=='D')a.erase(x);
		if (c=='F')printf("%d\n",a.find(x));
		if (c=='P')printf("%d\n",a.pred(x));
		if (c=='N')printf("%d\n",a.succ(x));
	}
	printf("time=%d\n",clock()-t1);
	system("pause");for (;;);
	return 0;
}


