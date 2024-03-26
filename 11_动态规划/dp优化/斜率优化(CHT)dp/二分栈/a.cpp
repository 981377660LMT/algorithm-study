const int N=1e5+5;

int n,h[N],s[N],p[N];
ll f[N],buc[N];
vector <int> st[N];
ll cal(int p,ll l){return f[p-1]+l*l*s[p];}

#define tp st[c][st[c].size()-1]
#define se st[c][st[c].size()-2]

int chk(int i,int j){
	int l=p[j],r=buc[s[j]]+1;
	while(l<r){
		int m=l+r>>1;
		if(cal(i,m-p[i]+1)<cal(j,m-p[j]+1))l=m+1;
		else r=m;
	} return l;
}

int main(){
	cin>>n;
	for(int i=1;i<=n;i++)s[i]=read(),p[i]=++buc[s[i]];
	for(int i=1;i<=n;i++){
		int c=s[i];
		while(st[c].size()>1&&chk(se,tp)<=chk(tp,i))st[c].pop_back();
		st[c].push_back(i);
		while(st[c].size()>1&&chk(se,tp)<=p[i])st[c].pop_back();
		f[i]=cal(tp,p[i]-p[tp]+1);
	} cout<<f[n]<<endl;
	return 0;
}