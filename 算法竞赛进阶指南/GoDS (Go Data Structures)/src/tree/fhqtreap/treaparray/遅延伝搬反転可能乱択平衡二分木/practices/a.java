import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.InputStreamReader;
import java.io.OutputStreamWriter;

public class Solution {
  static class Node {
    int l, r;
    int key, val, minv;
    int size, add, rev;
  }
  static class Pair {
    int x, y;

    Pair(int x, int y) {
      this.x = x;
      this.y = y;
    }
  }



  final static int N = 200010, INF = (int) 1e9;

  static int n, m, root, idx;
  static int[] h = new int[N];
  static Node[] tr = new Node[N];


  static int x, y, z, k;
  static Pair temp;

  static int get_min(int l, int r) {
    temp = split(root, l - 1, x, y);
    x = temp.x;
    y = temp.y;
    temp = split(y, r - l + 1, y, z);
    y = temp.x;
    z = temp.y;
    int ans = tr[y].minv;
    root = merge(merge(x, y), z);
    return ans;
  }

  static void insert(int l, int t) {
    temp = split(root, l, x, y);
    x = temp.x;
    y = temp.y;
    root = merge(merge(x, get_node(t)), y);
  }

  static void remove(int l) {
    temp = split(root, l - 1, x, y);
    x = temp.x;
    y = temp.y;
    temp = split(y, 1, y, z);
    y = temp.x;
    z = temp.y;
    root = merge(x, z);
  }

  static void modify(int l, int r, int v) {
    temp = split(root, l - 1, x, y);
    x = temp.x;
    y = temp.y;
    temp = split(y, r - l + 1, y, z);
    y = temp.x;
    z = temp.y;
    tr[y].val += v;
    tr[y].minv += v;
    tr[y].add += v;
    root = merge(merge(x, y), z);
  }

  static void revolve(int l, int r, int c) {
    c = c % (r - l + 1);
    temp = split(root, l - 1, x, y);
    x = temp.x;
    y = temp.y;
    temp = split(y, r - l + 1 - c, y, z);
    y = temp.x;
    z = temp.y;
    temp = split(z, c, z, k);
    z = temp.x;
    k = temp.y;
    z = merge(z, y);
    x = merge(x, z);
    root = merge(x, k);
  }

  static void reverse(int l, int r) {
    temp = split(root, l - 1, x, y);
    x = temp.x;
    y = temp.y;
    temp = split(y, r - l + 1, y, z);
    y = temp.x;
    z = temp.y;
    swap(y);
    tr[y].rev ^= 1;
    root = merge(merge(x, y), z);
  }

  static Pair split(int p, int siz, int x, int y) {
    if (p == 0)
      return new Pair(0, 0);

    pushdown(p);
    if (tr[tr[p].l].size < siz) {
      x = p;
      temp = split(tr[p].r, siz - tr[tr[p].l].size - 1, tr[p].r, y);
      tr[p].r = temp.x;
      y = temp.y;
    } else {
      y = p;
      temp = split(tr[p].l, siz, x, tr[p].l);
      x = temp.x;
      tr[p].l = temp.y;
    }
    pushup(p);
    return new Pair(x, y);
  }

  static int merge(int x, int y) {
    if (x == 0 || y == 0)
      return x + y;
    if (tr[x].key < tr[y].key) {
      pushdown(x);
      tr[x].r = merge(tr[x].r, y);
      pushup(x);
      return x;
    } else {
      pushdown(y);
      tr[y].l = merge(x, tr[y].l);
      pushup(y);
      return y;
    }
  }

  static int build(int l, int r) {
    if (l > r)
      return 0;
    int mid = l + r >> 1;
    int u = get_node(h[mid]);
    tr[u].l = build(l, mid - 1);
    tr[u].r = build(mid + 1, r);
    pushup(u);
    return u;
  }

  static void pushup(int u) {
    tr[u].size = tr[tr[u].l].size + tr[tr[u].r].size + 1;
    int t = Math.min(tr[u].l == 0 ? INF : tr[tr[u].l].minv, tr[u].r == 0 ? INF : tr[tr[u].r].minv);
    tr[u].minv = Math.min(tr[u].val, t);
  }

  static void pushdown(int u) {
    if (tr[u].rev != 0) {
      swap(tr[u].l);
      swap(tr[u].r);
      tr[tr[u].l].rev ^= 1;
      tr[tr[u].r].rev ^= 1;
      tr[u].rev = 0;
    }
    if (tr[u].add != 0) {
      int l = tr[u].l, r = tr[u].r;
      tr[l].minv += tr[u].add;
      tr[l].add += tr[u].add;
      tr[l].val += tr[u].add;
      tr[r].minv += tr[u].add;
      tr[r].add += tr[u].add;
      tr[r].val += tr[u].add;
      tr[u].add = 0;
    }
  }


  static int get_node(int val) {
    ++idx;
    tr[idx].val = tr[idx].minv = val;
    tr[idx].key = (int) (Math.random() * INF) + 1;
    tr[idx].size = 1;
    tr[idx].add = tr[idx].rev = 0;
    tr[idx].l = tr[idx].r = 0;
    return idx;
  }

  static void swap(int u) {


    int temp = tr[u].l;
    tr[u].l = tr[u].r;
    tr[u].r = temp;
  }



}
