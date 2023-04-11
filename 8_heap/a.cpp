template<class T> struct PriorityQueue {
    vector<T> d;

    PriorityQueue() {}

    PriorityQueue(const vector<T> &d_) : d(d_) {
	make_heap();
    }

    template<class Iter> PriorityQueue(Iter first, Iter last) : d(first, last) {
	make_heap();
    }

    void push(const T &x) {
	int k = d.size();
	d.push_back(x);
	up(k);
    }

    void pop_min() {
	if (d.size() < 3u) {
	    d.pop_back(); 
	} else {
	    swap(d[1], d.back()); d.pop_back();
	    int k = down(1);
	    up(k);
	}
    }

    void pop_max() {
	if (d.size() < 2u) { 
	    d.pop_back();
	} else {
	    swap(d[0], d.back()); d.pop_back();
	    int k = down(0);
	    up(k);
	}
    }

    const T& get_min() const {
	return d.size() < 2u ? d[0] : d[1];
    }

    const T& get_max() const {
	return d[0];
    }

    int size() const { return d.size(); }

    bool empty() const { return d.empty(); }

    void make_heap() {
	for (int i=d.size(); i--; ) {
	    if (i & 1 && d[i-1] < d[i]) swap(d[i-1], d[i]);
	    int k = down(i);
	    up(k, i);
	}
    }

    inline int parent(int k) const {
	return ((k>>1)-1)&~1;
    }

    int down(int k) {
	int n = d.size();
	if (k & 1) { // min heap
	    while (2*k+1 < n) {
		int c = 2*k+3;
		if (n <= c || d[c-2] < d[c]) c -= 2;
		if (c < n && d[c] < d[k]) { swap(d[k], d[c]); k = c; }
		else break;
	    }
	} else { // max heap
	    while (2*k+2 < n) {
		int c = 2*k+4;
		if (n <= c || d[c] < d[c-2]) c -= 2;
		if (c < n && d[k] < d[c]) { swap(d[k], d[c]); k = c; }
		else break;
	    }
	}
	return k;
    }

    int up(int k, int root=1) {
	if ((k|1) < (int)d.size() && d[k&~1] < d[k|1]) {
	    swap(d[k&~1], d[k|1]);
	    k ^= 1;
	}

	int p;
	while (root < k && d[p=parent(k)] < d[k]) { // max heap
	    swap(d[p], d[k]);
	    k = p;
	}
	while (root < k && d[k] < d[p=parent(k)|1]) { // min heap
	    swap(d[p], d[k]);
	    k = p;
	}
	return k;
    }
};