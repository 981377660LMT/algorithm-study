// #include <bits/stdc++.h>
// #define il inline
// typedef long long ll;
// const int maxn = 2e5 + 10;
// using namespace std;
// template<class T> il void rd(T& res) {
//     res = 0;char c;bool sign = 0;
//     for(c = getchar();!isdigit(c);c = getchar()) sign |= c == '-';
//     for(;isdigit(c);c = getchar()) res = (res << 1) + (res << 3) + (c ^ 48);
//     (sign) && (res = -res);
//     return;
// }
// struct TreapNode {
//     int ch[2];
//     int rnd,size,val;bool rev;
//     ll sum;
// }tr[maxn << 6];
// int emp[maxn],root[maxn];
// int n,m,i,j,k,q,tot,emp_top,a,b,c,d;ll lans;
// il void _swap(int& x,int& y) {
//     x ^= y ^= x ^= y;
//     return;
// }
// il int new_node(int v) {
//     int id = emp_top ? emp[emp_top--] : ++tot;
//     tr[id].rnd = rand();tr[id].val = v;
//     tr[id].size = 1;tr[id].sum = v;tr[id].rev = 0;
//     tr[id].ch[0] = tr[id].ch[1] = 0;
//     return id;
// }
// il int clone(int o) {
//     int id = emp_top ? emp[emp_top--] : ++tot;
//     tr[id].ch[0] = tr[o].ch[0];tr[id].ch[1] = tr[o].ch[1];
//     tr[id].rnd = tr[o].rnd;tr[id].size = tr[o].size;tr[id].val = tr[o].val;
//     tr[id].rev = tr[o].rev;tr[id].sum = tr[o].sum;
//     return id;
// }
// il void push_down(int o) {
//     if(tr[o].rev) {
//         _swap(tr[o].ch[0],tr[o].ch[1]);
//         if(tr[o].ch[0]) tr[o].ch[0] = clone(tr[o].ch[0]),tr[tr[o].ch[0]].rev ^= 1;
//         if(tr[o].ch[1]) tr[o].ch[1] = clone(tr[o].ch[1]),tr[tr[o].ch[1]].rev ^= 1;
//         tr[o].rev = 0;
//   	}
//   	return;
// }
// il void push_up(int o) {
//     tr[o].size = tr[tr[o].ch[0]].size + tr[tr[o].ch[1]].size + 1;
//     tr[o].sum = tr[tr[o].ch[0]].sum + tr[tr[o].ch[1]].sum + tr[o].val;
//     return;
// }
// void split_k(int now,int k,int& x,int& y) {
//     if(!now) {x = y = 0;return;}
//     push_down(now);
//     if(tr[tr[now].ch[0]].size < k) {
//         x = clone(now);
//         split_k(tr[x].ch[1],k - tr[tr[x].ch[0]].size - 1,tr[x].ch[1],y);
//         push_up(x);
//     } else {
//         y = clone(now);
//         split_k(tr[y].ch[0],k,x,tr[y].ch[0]);
//         push_up(y);
//     }
//     return;
// }
// int merge(int u,int v) {
//     if(!u) return v;if(!v) return u;
//     if(tr[u].rnd < tr[v].rnd) {
//         push_down(u);
//         tr[u].ch[1] = merge(tr[u].ch[1],v);
//         push_up(u);
//         return u;
//     } else {
//         push_down(v);
//         tr[v].ch[0] = merge(u,tr[v].ch[0]);
//         push_up(v);
//         return v;
//     }
// }
// il void insert(int& root,int x,int v) {
//     split_k(root,x,a,b);
//     root = merge(merge(a,new_node(v)),b);
//     return;
// }
// il void erase(int& root,int x) {
//     split_k(root,x,a,c);
//     split_k(a,x - 1,a,b);
//     emp[++emp_top] = b;
//     root = merge(a,c);
//     return;
// }
// il void rever(int& root,int l,int r) {
//     split_k(root,l - 1,a,b);
//     split_k(b,r - l + 1,b,c);
//     tr[b].rev ^= 1;
//     root = merge(a,merge(b,c));
//     return;
// }
// il ll query(int& root,int l,int r) {
//     split_k(root,l - 1,a,b);
//     split_k(b,r - l + 1,b,c);
//     ll res = tr[b].sum;
//     root = merge(a,merge(b,c));
//     return res;
// }
// int main() {
//     srand((unsigned)time(NULL));
//     rd(q);
//     for(int i = 1,x,y,l,r;i <= q;i++) {
//         int his,opt;rd(his);rd(opt);root[i] = root[his];
//         switch(opt) {
//             case 1:rd(x);rd(y);insert(root[i],x ^ lans,y ^ lans);break;
//             case 2:rd(x);erase(root[i],x ^ lans);break;
//             case 3:rd(l);rd(r);rever(root[i],l ^ lans,r ^ lans);break;
//             case 4:rd(l);rd(r);l ^= lans;r ^= lans;printf("%lld\n",lans = query(root[i],l,r));break;
//         }
//     //	cout << a << ' ' << b << ' ' << c << endl;
//     }
// } 
