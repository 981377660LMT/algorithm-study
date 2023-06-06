#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<time.h>
#include<math.h>
#include<iostream>
#include<algorithm>
#include<vector>
using namespace std;
const int N=1000005;
const double eps=1e-6;
struct point{
	double x,y;
	point(double _x=0,double _y=0):x(_x),y(_y){}
};
struct line{
	point x,y;int id1,id2;
	line(){}
	line(const point &_x,const point &_y,int _id1,int _id2):x(_x),y(_y),id1(_id1),id2(_id2){}
}a[N];
inline double cha(double x1,double y1,double x2,double y2){return x1*y2-x2*y1;}
point cross(const line &a,const line &b){
	double s1=cha(b.y.x-b.x.x,b.y.y-b.x.y,a.x.x-b.x.x,a.x.y-b.x.y),
	s2=cha(b.y.x-b.x.x,b.y.y-b.x.y,a.y.x-b.x.x,a.y.y-b.x.y);
	if (fabs(s1-s2)<eps)return point(1e100,1e100);
	return point((a.x.x*s2-a.y.x*s1)/(s2-s1),(a.x.y*s2-a.y.y*s1)/(s2-s1));
}
inline int on_segment(const point &p,const line &l){
	if (fabs(cha(l.y.x-l.x.x,l.y.y-l.x.y,p.x-l.x.x,p.y-l.x.y))/sqrt((l.y.x-l.x.x)*(l.y.x-l.x.x)+(l.y.y-l.x.y)*(l.y.y-l.x.y))<eps&&
	min(l.x.x,l.y.x)+eps<=p.x&&max(l.x.x,l.y.x)-eps>=p.x&&
	min(l.x.y,l.y.y)+eps<=p.y&&max(l.x.y,l.y.y)-eps>=p.y)return 1;
	if (fabs(l.y.x-p.x)<eps&&fabs(l.y.x-l.x.x)<eps&&min(l.x.y,l.y.y)+eps<=p.y&&max(l.x.y,l.y.y)-eps>=p.y)return 1;
	if (fabs(l.y.y-p.y)<eps&&fabs(l.y.y-l.x.y)<eps&&min(l.x.x,l.y.x)+eps<=p.x&&max(l.x.x,l.y.x)-eps>=p.x)return 1;
	return 0;
}
inline bool in(const point &p,const line &l){
	return cha(l.y.x-l.x.x,l.y.y-l.x.y,p.x-l.x.x,p.y-l.x.y)>-eps;
}
struct node{
	line p;
	node *l,*r;
}c[N],*root;
int n,m,c1=0,s=0;
node *build(int l,int r){
	if (l>r)return 0;
	node *res=c+(++c1);res->p=a[l];
	int r1=r;
	for (int i=l+1;i<=r;++i){++s;
		point p=cross(a[l],a[i]);
		if (on_segment(p,a[i]))
			if (in(a[i].x,a[l]))a[++r1]=line(a[i].x,p,a[i].id1,a[i].id2);
			else a[++r1]=line(p,a[i].y,a[i].id1,a[i].id2);
		else if (in(a[i].x,a[l])&&in(a[i].y,a[l]))a[++r1]=a[i];
	}
	res->l=build(r+1,r1);r1=r;
	for (int i=l+1;i<=r;++i){++s;
		point p=cross(a[l],a[i]);
		if (on_segment(p,a[i]))
			if (!in(a[i].x,a[l]))a[++r1]=line(a[i].x,p,a[i].id1,a[i].id2);
			else a[++r1]=line(p,a[i].y,a[i].id1,a[i].id2);
		else if (!(in(a[i].x,a[l])&&in(a[i].y,a[l])))a[++r1]=a[i];
	}
	res->r=build(r+1,r1);
	return res;
}
int locate(node *x,const point &p){
	return in(p,x->p)?(x->l?locate(x->l,p):x->p.id1):(x->r?locate(x->r,p):x->p.id2);
}
int main()
{
	//freopen("1.in","r",stdin);
	//freopen("1.out","w",stdout);
	int t1=clock();
	scanf("%d",&n);
	for (int i=1;i<=n;++i)scanf("%lf%lf%lf%lf",&a[i].x.x,&a[i].x.y,&a[i].y.x,&a[i].y.y);
	random_shuffle(a+1,a+1+n);
	root=build(1,n);
	scanf("%d",&m);
	while (m--){
		double x,y;scanf("%lf%lf",&x,&y);
		printf("%d\n",locate(root,point(x,y)));
	}
	printf("time=%d\n",clock()-t1);
	printf("c1=%d s=%d\n",c1,s);
	system("pause");for (;;);
	return 0;
}


