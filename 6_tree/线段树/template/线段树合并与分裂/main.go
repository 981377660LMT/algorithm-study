package main

func main() {
	CF600E()
}

// Lomsat gelral
// https://www.luogu.com.cn/problem/CF600E
func CF600E() {}

type name struct {
}

// il int merge(int now1,int now2,int l,int r){
// 	if(!now1 || !now2) return now1|now2;
// 	if(l==r){sum[now1]+=sum[now2]; return now1;}
// 	ls[now1]=merge(ls[now1],ls[now2],l,mid);
// 	rs[now1]=merge(rs[now1],rs[now2],mid+1,r);
// 	pu(now1); return now1;
// }

// func
// 合并后子树的线段树的信息丢失，必须离线问题，自下而上回答

// 子树的线段树信息不丢失，可以在线回答问题。
