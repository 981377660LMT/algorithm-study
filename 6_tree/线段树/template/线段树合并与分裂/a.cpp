struct {
    int ls, rs, v;
} T[M];
int cnt, nn, R[N], P[N];
namespace segt {
int query(int x, int p, int l = 1, int r = nn) {
    if (l == r) return T[p].v;
    int mid = l + r >> 1;
    if (x <= mid)
        return query(x, T[p].ls, l, mid);
    else
        return query(x, T[p].rs, mid + 1, r);
}
int kth(int k, int p, int l = 1, int r = nn) {
    if (l == r) return T[p].v >= k ? l : 0;
    int mid = l + r >> 1;
    if (T[T[p].ls].v >= k)
        return kth(k, T[p].ls, l, mid);
    else
        return kth(k - T[T[p].ls].v, T[p].rs, mid + 1, r);
}
void pushup(int p) {
    T[p].v = T[T[p].ls].v + T[T[p].rs].v;
}
void add(int x, int d, int &p, int l = 1, int r = nn) {
    if (!p) p = ++cnt;
    if (l == r) return (void)(T[p].v += d);
    int mid = l + r >> 1;
    if (x <= mid)
        add(x, d, T[p].ls, l, mid);
    else
        add(x, d, T[p].rs, mid + 1, r);
    pushup(p);
}
int merge(int p, int q, int l = 1, int r = nn) {
    if (!p || !q) return p + q;
    if (l == r) return T[p].v += T[q].v, p;
    int mid = l + r >> 1;
    T[p].ls = merge(T[p].ls, T[q].ls, l, mid);
    T[p].rs = merge(T[p].rs, T[q].rs, mid + 1, r);
    pushup(p);
    return p;
}
} // namespace segt
    nn = n;
    dsu::init(n);
    P[0] = -1;
    for (int i = 1, x; i <= n; ++i) {
        cin >> x;
        P[x] = i;
        segt::add(x, 1, R[i]);
    }
    for (int i = 0, x, y; i < m; ++i) {
        cin >> x >> y;
        int fx = dsu::find(x), fy = dsu::find(y);
        fa[fy] = fx;
        segt::merge(R[fx], R[fy]);
    }
    cin >> q;
    while (q--) {
        char o;
        int x, y;
        cin >> o >> x >> y;
        if (o == 'Q') {
            int fx = dsu::find(x);
            int ans = segt::kth(y, R[fx]);
            cout << P[ans] << '\n';
        } else {
            int fx = dsu::find(x), fy = dsu::find(y);
            fa[fy] = fx;
            R[fy] = segt::merge(R[fx], R[fy]);
        }
    }
    return 0;