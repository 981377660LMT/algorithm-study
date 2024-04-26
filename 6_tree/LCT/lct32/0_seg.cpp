#include <vector>
#include <unordered_map>
using namespace std;

namespace titan23 {

  template <class T,
            T (*_op)(T, T),
            T (*_e)()>
  struct DynamicSegmentTree {

    int _u, _size, _log;
    unordered_map<int, T> _data;

    DynamicSegmentTree(const int u) {
      _build(u);
    }

    void _build(const int u) {
      this->_u = u;
      this->_log = 32 - __builtin_clz(_u-1);
      this->_size = 1 << _log;
    }

    T _get(const int k) const {
      auto it = _data.find(k);
      return it == _data.end()? _e : it->second;
    }

    T get(int k) const {
      if (k < 0) k += _u;
      return _get(k+_size);
    }

    void set(int k, const T v) {
      if (k < 0) k += _u;
      k += _size;
      _data[k] = v;
      for (int i = 0; i < _log; ++i) {
        k >>= 1;
        _data[k] = _op(_get(k<<1), _get(k<<1|1));
      }
    }

    T prod(int l, int r) {
      l += _size;
      r += _size;
      T lres = _e(), rres = _e();
      while (l < r) {
        if (l & 1) {
          lres = _op(lres, _get(l++));
        }
        if (r & 1) {
          rres = _op(_get(r^1), rres);
        }
        l >>= 1;
        r >>= 1;
      }
      return _op(lres, rres);
    }

    T all_prod() {
      return _get(1);
    }

    template<typename F>  // F: function<bool (T)> f
    int max_right(int l, F &&f) {
      if (l == _u) return _u;
      l += _size;
      T s = _e();
      while (1) {
        while ((l & 1) == 0) {
          l >>= 1;
        }
        if (!f(_op(s, _get(l)))) {
          while (l < _size) {
            l <<= 1;
            if (f(_op(s, _get(l)))) {
              s = _op(s, _get(l));
              l |= 1;
            }
          }
          return l - _size;
        }
        s = _op(s, _get(l));
        ++l;
        if ((l & (-l)) == l) break;
      }
      return _u;
    }

    template<typename F>  // F: function<bool (T)> f
    int min_left(int r, F &&f) {
      if (r == 0) return 0;
      r += _size;
      T s = _e();
      while (1) {
        --r;
        while ((r > 1) && (r & 1)) {
          r >>= 1;
        }
        if (!f(_op(_get(r), s))) {
          while (r < _size) {
            r = r << 1 | 1;
            if (f(_op(_get(r), s))) {
              s = _op(_get(r), s);
              r ^= 1;
            }
          }
          return r + 1 - _size;
        }
        s = _op(_get(r), s);
        if ((r & (-r)) == r) break;
      }
      return 0;
    }

    vector<T> tovector() {
      vector<T> res(_u);
      for (int i = 0; i < _u; ++i) {
        res[i] = get(i);
      }
      return res;
    }

    void print() {}
  };
}  // namespace titan23
