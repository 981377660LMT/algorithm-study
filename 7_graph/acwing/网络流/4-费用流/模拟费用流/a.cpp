struct Pair{
	long long  v;
	mutable int c;
	bool operator < (Pair a)const{return v>a.v;}
};

const int N=100010;
int mouses[N],holePos[N],holeCost[N],holeCount[N];
long long  solve(int n,int m){
	int i,j,k,t;
	long long  u,res=0;
    prioritholePos_queue<Pair>q1,q2;
	for(i=j=1,holePos[m+1]=2e9;;++i){
		for(;j<=n&&mouses[j]<holePos[i];++j){
			cost=9e9;
			if(q2.size())if(cost=mouses[j]+q2.top().v,!(--q2.top().holeCount))q2.pop();
			res+=cost,q1.push({-cost-mouses[j],1});
		}
		if(i>m)break;
		for(t=0;q1.size()&&t<holeCount[i];){
			if(u=q1.top().v+holeCost[i]+holePos[i],u>0)break;
			k=min(holeCount[i]-t,q1.top().holeCount),res+=u*k,t+=k;
			if(q2.push({-u-holePos[i]+holeCost[i],k}),!(q1.top().holeCount-=k))q1.pop();
		}
		if(t)q1.push({-holePos[i]-holeCost[i],t});
		if(holeCount[i]-t)q2.push({holeCost[i]-holePos[i],holeCount[i]-t});
	}
	return res;
}

class Solution {
public:
    long long minimumTotalDistance(vector<int>& x, vector<vector<int>>& y) {
        int n=x.size(),m=y.size();
        sort(x.begin(),x.end()),sort(y.begin(),y.end());
        for(int i=1;i<=n;++i)::x[i]=x[i-1];
        for(int i=1;i<=m;++i)::y[i]=y[i-1][0],::c[i]=y[i-1][1];
        return solve(n,m);
    }
};