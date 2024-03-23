#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <algorithm>
#include <string>
using namespace std;
const int maxn = 1110000;
const double alpha = 0.6;
struct SuffixBalancedTree {
	string data;
	double weight[maxn], left, right;
	int size[maxn], children[maxn][2], root;
	int* scapegoat, length, collector[maxn];
	void init() {
		data.resize(1);
		root = 0;
	}
	void pushup(int x) {
		size[x] = size[children[x][0]] + size[children[x][1]] + 1;
	}
	void dfs(int x) {
		if (!x) return;
		dfs(children[x][0]);
		collector[++length] = x;
		dfs(children[x][1]);
	}
	int build(int a, int b, double L, double R) {
		if (a > b) return 0;
		int mid = (a + b) >> 1;
		int x = collector[mid];
		size[x] = b - a + 1;
		weight[x] = (L + R) / 2;
		children[x][0] = build(a, mid - 1, L, weight[x]);
		children[x][1] = build(mid + 1, b, weight[x], R);
		return x;
	}
	void add(int& x, double L, double R) {
		if (!x) {
			x = data.size() - 1;
			children[x][0] = children[x][1] = 0;
			weight[x] = (L + R) / 2;
			size[x] = 1;
			return;
		}
		int y = data.size() - 1;
		if (data[y] < data[x] || data[y] == data[x] && weight[y - 1] < weight[x - 1])
			add(children[x][0], L, weight[x]);
		else
			add(children[x][1], weight[x], R);
		pushup(x);
		if (size[children[x][0]] > size[x] * alpha || size[children[x][1]] > size[x] * alpha)
			scapegoat = &x, left = L, right = R;
	}
	void push_front(char c) { //向开头添加一个新的字符之后会改变其他位置的下标
		data.push_back(c);
		scapegoat = nullptr;
		add(root, 0, 1);
		if (scapegoat) {
			length = 0;
			dfs(*scapegoat);
			*scapegoat = build(1, length, left, right);
		}
	}
	int merge(int x, int y) {
		if (!x || !y) 
			return x | y;
		if (size[x] > size[y]) {
			children[x][1] = merge(children[x][1], y);
			pushup(x);
			return x;
		}
		else {
			children[y][0] = merge(x, children[y][0]);
			pushup(y);
			return y;
		}
	}
	void del(int& x) {
		const int y = data.size() - 1;
		size[x] -= 1;
		if (x == y)
			x = merge(children[x][0], children[x][1]);
		else if (weight[y] < weight[x])
			del(children[x][0]);
		else
			del(children[x][1]);
	}
	void pop_front() {
		del(root);
		data.pop_back();
	}
	double weight(int index) { //返回下标index的权值，权值越小，字典序越小（下标从1开始）
		return weight[data.size() - index];
	}
	int rank(int index) { //返回下标index的排名（下标从1开始）
		double key = weight(index);
		int ret = 0, x = root;
		while (x) {
			if (key < weight[x])
				x = children[x][0];
			else
				ret += size[children[x][0]] + 1, x = children[x][1];
		}
		return ret;
	}
	int rank(const char *s) { //返回字典序小于s的后缀的个数
		int ret = 0, x = root, n = strlen(s);
		while (x) {
			int L = min(x, n) + 1;
			int flag = 0;
			for (int i = 0; i < L; ++i) {
				if (s[i] != data[x - i]) {
					flag = s[i] - data[x - i];
					break;
				}
			}
			if (flag <= 0)
				x = children[x][0];
			else
				ret += size[children[x][0]] + 1, x = children[x][1];
		}
		return ret;
	}
}tree;
int main() {

	return 0;
}