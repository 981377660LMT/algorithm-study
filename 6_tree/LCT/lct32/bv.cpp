#include <iostream>
#include <vector>
#include <cassert>
#include <cmath>
using namespace std;

// DynamicBitVector
namespace titan23 {

  class DynamicBitVector {
   private:
    static const int BUCKET_MAX = 1000;
    vector<vector<uint8_t>> data;
    vector<int> bucket_data;
    int _size;
    int tot_one;

    void build(const vector<uint8_t> &a) {
      long long s = len();
      int bucket_size = max((int)(s+BUCKET_MAX-1)/BUCKET_MAX, (int)ceil(sqrt(s)));
      data.resize(bucket_size);
      bucket_data.resize(bucket_size);
      for (int i = 0; i < bucket_size; ++i) {
        int start = s*i/bucket_size;
        int stop = min((int)len(), (int)(s*(i+1)/bucket_size));
        vector<uint8_t> d(a.begin()+start, a.begin()+stop);
        int sum = 0;
        for (const uint8_t &x: d) sum += x;
        data[i] = d;
        tot_one += sum;
        bucket_data[i] = sum;
      }
    }

    pair<int, int> get_bucket(int k) const {
      if (k == len()) return {-1, -1};
      if (k < len()/2) {
        for (int i = 0; i < data.size(); ++i) {
          if (k < data[i].size()) return {i, k};
          k -= data[i].size();
        }
      } else {
        int tot = len();
        for (int i = data.size()-1; i >= 0; --i) {
          if (tot-data[i].size() <= k) {
            return {i, k-(tot-data[i].size())};
          }
          tot -= data[i].size();
        }
      }
      assert(false);
    }

   public:
    DynamicBitVector() : _size(0), tot_one(0) {}
    DynamicBitVector(const vector<uint8_t> &a) : _size(a.size()), tot_one(0) {
      build(a);
    }

    void insert(int k, bool key) {
      assert(0 <= k && k <= len());
      if (data.empty()) {
        ++_size;
        tot_one += key;
        bucket_data.emplace_back(key);
        data.push_back({key});
        return;
      }
      auto [bucket_pos, bit_pos] = get_bucket(k);
      if (bucket_pos == -1) {
        bucket_pos = data.size()-1;
        bucket_data.back() += key;
        data.back().emplace_back(key);
      } else {
        bucket_data[bucket_pos] += key;
        data[bucket_pos].insert(data[bucket_pos].begin() + bit_pos, key);
      }
      if (data[bucket_pos].size() > BUCKET_MAX) {
        vector<uint8_t> right(data[bucket_pos].begin() + BUCKET_MAX/2, data[bucket_pos].end());
        data[bucket_pos].erase(data[bucket_pos].begin() + BUCKET_MAX/2, data[bucket_pos].end());
        data.emplace(data.begin() + bucket_pos+1, right);
        bucket_data.insert(bucket_data.begin() + bucket_pos, 0);
        bucket_data[bucket_pos] = 0;
        bucket_data[bucket_pos+1] = 0;
        for (const uint8_t x: data[bucket_pos]) bucket_data[bucket_pos] += x;
        for (const uint8_t x: data[bucket_pos+1]) bucket_data[bucket_pos+1] += x;
      }
      ++_size;
      tot_one += key;
    }

    bool access(int k) const {
      assert(0 <= k && k < len());
      auto [bucket_pos, bit_pos] = get_bucket(k);
      return data[bucket_pos][bit_pos];
    }

    bool pop(int k) {
      assert(0 <= k && k < len());
      auto [bucket_pos, bit_pos] = get_bucket(k);
      bool res = data[bucket_pos][bit_pos];
      bucket_data[bucket_pos] -= res;
      data[bucket_pos].erase(data[bucket_pos].begin() + bit_pos);
      tot_one -= res;
      --_size;
      if (data[bucket_pos].empty()) {
        data.erase(data.begin() + bucket_pos);
        bucket_data.erase(bucket_data.begin() + bucket_pos);
      }
      return res;
    }

    void set(int k, bool v) {
      assert(0 <= k && k < len());
      auto [bucket_pos, bit_pos] = get_bucket(k);
      data[bucket_pos][bit_pos] = v;
    }

    int rank0(int r) const {
      assert(0 <= r && r <= len());
      return r - rank1(r);
    }

    int rank1(int r) const {
      assert(0 <= r && r <= len());
      int s = 0;
      for (int i = 0; i < data.size(); ++i) {
        if (r < data[i].size()) {
          const vector<uint8_t> &d = data[i];
          for (int j = 0; j < r; ++j) {
            if (d[j]) ++s;
          }
          return s;
        }
        s += bucket_data[i];
        r -= data[i].size();
      }
      return s;
      assert(false);
    }

    int rank(int r, bool key) const {
      assert(0 <= r && r <= len());
      return key ? rank1(r) : rank0(r);
    }

    int select0(int k) const {
      int s = 0;
      for (int i = 0; i < data.size(); ++i) {
        if (k < data[i].size() - bucket_data[i]) {
          for (const uint8_t &x: data[i]) {
            if (!x) --k;
            if (k < 0) return s;
            s++;
          }
          assert(false);
        }
        s += data[i].size();
        k -= data[i].size() - bucket_data[i];
      }
      assert(false);
    }

    int select1(int k) const {
      int s = 0;
      for (int i = 0; i < data.size(); ++i) {
        if (k < bucket_data[i]) {
          for (const uint8_t &x: data[i]) {
            if (x) --k;
            if (k < 0) return s;
            s++;
          }
        }
        s += data[i].size();
        k -= bucket_data[i];
      }
      assert(false);
    }

    int select(int k, bool key) const {
      return key ? select1(k) : select0(k);
    }

    int _insert_and_rank1(int k, bool key) {
      int s = 0;
      int bucket_pos = -1, bit_pos = -1;
      if (k < len()/2) {
        for (int i = 0; i < data.size(); ++i) {
          if (k < data[i].size()) {
            bucket_pos = i;
            bit_pos = k;
            const vector<uint8_t> &d = data[i];
            for (int j = 0; j < k; ++j) {
              s += d[j];
            }
            break;
          }
          s += bucket_data[i];
          k -= data[i].size();
        }
      } else {
        int tot = len();
        s = tot_one;
        for (int i = data.size()-1; i >= 0; --i) {
          if (tot-data[i].size() <= k) {
            bucket_pos = i;
            bit_pos = k-(tot-data[i].size());
            const vector<uint8_t> &d = data[i];
            for (int j = bit_pos; j < d.size(); ++j) {
              s -= d[j];
            }
            break;
          }
          tot -= data[i].size();
          s -= bucket_data[i];
        }
      }

      {
        ++_size;
        tot_one += key;
        if (data.empty()) {
          bucket_data.emplace_back(key);
          data.push_back({{key}});
          return s;
        }
        if (bucket_pos == -1) {
          bucket_pos = data.size()-1;
          bucket_data.back() += key;
          data.back().emplace_back(key);
        } else {
          bucket_data[bucket_pos] += key;
          data[bucket_pos].insert(data[bucket_pos].begin() + bit_pos, key);
        }
        if (data[bucket_pos].size() > BUCKET_MAX) {
          vector<uint8_t> right(data[bucket_pos].begin() + BUCKET_MAX/2, data[bucket_pos].end());
          data[bucket_pos].erase(data[bucket_pos].begin() + BUCKET_MAX/2, data[bucket_pos].end());
          data.emplace(data.begin() + bucket_pos+1, right);
          bucket_data.insert(bucket_data.begin() + bucket_pos, 0);
          bucket_data[bucket_pos] = 0;
          bucket_data[bucket_pos+1] = 0;
          for (const uint8_t &x: data[bucket_pos]) bucket_data[bucket_pos] += x;
          for (const uint8_t &x: data[bucket_pos+1]) bucket_data[bucket_pos+1] += x;
        }
      }
      return s;
    }

    int _access_pop_and_rank1(int k) {
      int prek = k;
      int s = 0;
      int bucket_pos, bit_pos;
      bool res;
      for (int i = 0; i < data.size(); ++i) {
        if (k < data[i].size()) {
          bucket_pos = i;
          bit_pos = k;
          res = data[bucket_pos][bit_pos];
          const vector<uint8_t> &d = data[i];
          for (int j = 0; j < k; ++j) {
            if (d[j]) ++s;
          }
          break;
        }
        s += bucket_data[i];
        k -= data[i].size();
      }
      bucket_data[bucket_pos] -= res;
      data[bucket_pos].erase(data[bucket_pos].begin() + bit_pos);
      tot_one -= res;
      --_size;
      if (data[bucket_pos].empty()) {
        data.erase(data.begin() + bucket_pos);
        bucket_data.erase(bucket_data.begin() + bucket_pos);
      }
      return s << 1 | res;
    }

    pair<bool, int> _access_ans_rank1(int k) const {
      assert(0 <= k && k < len());
      int s = 0;
      for (int i = 0; i < data.size(); ++i) {
        if (k < data[i].size()) {
          const vector<uint8_t> &d = data[i];
          for (int j = 0; j < k; ++j) {
            s += d[j];
          }
          return {data[i][k], s};
        }
        s += bucket_data[i];
        k -= data[i].size();
      }
      assert(false);
    }

    vector<uint8_t> tovector() const {
      vector<uint8_t> a(len());
      int ptr = 0;
      for (const vector<uint8_t> &d: data) for (const uint8_t &x: d) {
        a[ptr++] = x;
      }
      return a;
    }

    void print() const {
      vector<uint8_t> a = tovector();
      int n = (int)a.size();
      assert(n == len());
      cout << "[";
      for (int i = 0; i < n-1; ++i) {
        cout << a[i] << ", ";
      }
      if (n > 0) {
        cout << a.back();
      }
      cout << "]";
      cout << endl;
    }

    bool empty() const { return _size == 0; }
    int len() const { return _size; }
  };
} // name space titan23

