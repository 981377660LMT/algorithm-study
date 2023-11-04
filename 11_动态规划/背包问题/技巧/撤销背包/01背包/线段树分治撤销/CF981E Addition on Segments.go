// https://www.luogu.com.cn/problem/solution/CF981E

// /**********/
// #define MAXN 10011
// ll n,q;
// std::bitset<MAXN>ans;
// struct Segment_Tree
// {
//     std::vector<ll>t[MAXN<<2|1];
//     #define rt t[num]
//     #define tl t[num<<1]
//     #define tr t[num<<1|1]
//     void modify(un ql,un qr,ll k,un l=1,un r=n,un num=1)
//     {
//         if(ql<=l&&r<=qr)
//         {
//             rt.push_back(k);return;
//         }
//         un mid=(l+r)>>1;
//         if(ql<=mid)modify(ql,qr,k,l,mid,num<<1);
//         if(qr>mid)modify(ql,qr,k,mid+1,r,num<<1|1);
//     }
//     void dfs(std::bitset<MAXN>pre,un l=1,un r=n,un num=1)
//     {
//         std::bitset<MAXN>cur=pre;
//         for(auto x:rt)cur|=(cur<<x);
//         un mid=(l+r)>>1;
//         if(l==r){ans|=cur;return;}
//         dfs(cur,l,mid,num<<1),dfs(cur,mid+1,r,num<<1|1);
//     }
// }sgt;
// std::bitset<MAXN>s;
// int main()
// {
//     n=read(),q=read();
//     while(q--)
//     {
//         ll l=read(),r=read(),k=read();
//         sgt.modify(l,r,k);
//     }
//     s[0]=1;
//     sgt.dfs(s);
//     ll cnt=0;
//     for(ll i=1;i<=n;++i)
//         if(ans[i])++cnt;
//     printf("%lld\n",cnt);
//     for(ll i=1;i<=n;++i)
//         if(ans[i])printf("%lld ",i);
//     return 0;
// }