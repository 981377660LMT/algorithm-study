
#pragma once
#include <algorithm>
#include <cassert>
#include <limits>
#include <vector>

template <typename M, typename O, std::pair<typename M::T, bool> (*act)(typename M::T, typename O::T)>
class SegmentTreeBeats {
    using T = typename M::T;
    using E = typename O::T;

public:
    SegmentTreeBeats() = default;
    explicit SegmentTreeBeats(int n) : SegmentTreeBeats(std::vector<T>(n, M::id())) {}
    explicit SegmentTreeBeats(const std::vector<T>& v) {
        size = 1;
        while (size < (int) v.size()) size <<= 1;
        node.resize(2 * size, M::id());
        lazy.resize(2 * size, O::id());
        std::copy(v.begin(), v.end(), node.begin() + size);
        for (int i = size - 1; i > 0; --i) node[i] = M::op(node[2 * i], node[2 * i + 1]);
    }

    T operator[](int k) {
        return fold(k, k + 1);
    }

    void update(int l, int r, const E& x) { update(l, r, x, 1, 0, size); }

    T fold(int l, int r) { return fold(l, r, 1, 0, size); }

private:
    int size;
    std::vector<T> node;
    std::vector<E> lazy;

    void push(int k) {
        if (lazy[k] == O::id()) return;
        if (k < size) {
            lazy[2 * k] = O::op(lazy[2 * k], lazy[k]);
            lazy[2 * k + 1] = O::op(lazy[2 * k + 1], lazy[k]);
        }
        bool success;
        std::tie(node[k], success) = act(node[k], lazy[k]);
        if (!success) {
            assert(k < size);
            push(2 * k);
            push(2 * k + 1);
            node[k] = M::op(node[2 * k], node[2 * k + 1]);
        }
        lazy[k] = O::id();
    }

    void update(int a, int b, const E& x, int k, int l, int r) {
        push(k);
        if (r <= a || b <= l) return;
        if (a <= l && r <= b) {
            lazy[k] = O::op(lazy[k], x);
            push(k);
            return;
        }
        int m = (l + r) / 2;
        update(a, b, x, 2 * k, l, m);
        update(a, b, x, 2 * k + 1, m, r);
        node[k] = M::op(node[2 * k], node[2 * k + 1]);
    }

    T fold(int a, int b, int k, int l, int r) {
        push(k);
        if (r <= a || b <= l) return M::id();
        if (a <= l && r <= b) return node[k];
        int m = (l + r) / 2;
        return M::op(fold(a, b, 2 * k, l, m),
                     fold(a, b, 2 * k + 1, m, r));
    }
};


// the monoid for range chmin/chmax/add range sum query

using ll = long long;
constexpr ll INF = 1e18;

struct S {
    ll max_val, smax_val;
    ll min_val, smin_val;
    ll sum;
    int max_cnt, min_cnt, len;
    S() : max_val(-INF), smax_val(-INF), min_val(INF), smin_val(INF),
          sum(0), max_cnt(0), min_cnt(0), len(0) {}
    S(ll x, int len) : max_val(x), smax_val(-INF), min_val(x), smin_val(INF),
                       sum(x * len), max_cnt(len), min_cnt(len), len(len) {}
};

struct MinMaxSumMonoid {
    using T = S;
    static T id() { return S(); }
    static T op(T a, T b) {
        T c;
        c.sum = a.sum + b.sum;
        c.len = a.len + b.len;
        if (a.min_val < b.min_val) {
            c.min_val = a.min_val;
            c.min_cnt = a.min_cnt;
            c.smin_val = std::min(a.smin_val, b.min_val);
        } else if (a.min_val > b.min_val) {
            c.min_val = b.min_val;
            c.min_cnt = b.min_cnt;
            c.smin_val = std::min(a.min_val, b.smin_val);
        } else {
            c.min_val = a.min_val;
            c.min_cnt = a.min_cnt + b.min_cnt;
            c.smin_val = std::min(a.smin_val, b.smin_val);
        }
        if (a.max_val > b.max_val) {
            c.max_val = a.max_val;
            c.max_cnt = a.max_cnt;
            c.smax_val = std::max(a.smax_val, b.max_val);
        } else if (a.max_val < b.max_val) {
            c.max_val = b.max_val;
            c.max_cnt = b.max_cnt;
            c.smax_val = std::max(a.max_val, b.smax_val);
        } else {
            c.max_val = a.max_val;
            c.max_cnt = a.max_cnt + b.max_cnt;
            c.smax_val = std::max(a.smax_val, b.smax_val);
        }
        return c;
    }
};

struct F {
    ll lb, ub, diff;
    F(ll lb = -INF, ll ub = INF, ll diff = 0) : lb(lb), ub(ub), diff(diff) {}
    static F chmin(ll x) { return F(-INF, x, 0); }
    static F chmax(ll x) { return F(x, INF, 0); }
    static F add(ll x) { return F(-INF, INF, x); }
    bool operator==(const F& rhs) const { return lb == rhs.lb && ub == rhs.ub && diff == rhs.diff; }
};

struct ChminChmaxAddMonoid {
    using T = F;
    static T id() { return F(); }
    static T op(T a, T b) {
        F c;
        c.lb = std::clamp(a.lb + a.diff, b.lb, b.ub) - a.diff;
        c.ub = std::clamp(a.ub + a.diff, b.lb, b.ub) - a.diff;
        c.diff = a.diff + b.diff;
        return c;
    }
};

std::pair<S, bool> act(S a, F b) {
    if (a.len == 0) return {a, true};
    if (a.min_val == a.max_val || b.lb == b.ub || a.max_val <= b.lb || b.ub < a.min_val) {
        return {S(std::clamp(a.min_val, b.lb, b.ub) + b.diff, a.len), true};
    }
    if (a.smin_val == a.max_val) {
        a.min_val = a.smax_val = std::max(a.min_val, b.lb) + b.diff;
        a.max_val = a.smin_val = std::min(a.max_val, b.ub) + b.diff;
        a.sum = a.min_val * a.min_cnt + a.max_val * a.max_cnt;
        return {a, true};
    }
    if (b.lb < a.smin_val && a.smax_val < b.ub) {
        ll min_nxt = std::max(a.min_val, b.lb);
        ll max_nxt = std::min(a.max_val, b.ub);
        a.sum += (min_nxt - a.min_val) * a.min_cnt - (a.max_val - max_nxt) * a.max_cnt + b.diff * a.len;
        a.min_val = min_nxt + b.diff;
        a.max_val = max_nxt + b.diff;
        a.smin_val += b.diff;
        a.smax_val += b.diff;
        return {a, true};
    }
    return {a, false};
}