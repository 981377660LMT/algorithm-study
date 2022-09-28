// 启发式合并

class Solution {
public:
    int n;
    vector<vector<int>> adjList;
    vector<int> val;
    
    int res = 0;
    
    void merge(map<int, int>& a, map<int, int>& b) {
        if (a.size() < b.size()) {
            swap(a, b);
        }
        
        for (auto&& [v, c] : b) {
            if (a.count(v)) {
                res += c * a[v];
            }
            a[v] += c;
        }
    }
    
    map<int, int> dfs(int u, int fa) {
        map<int, int> ma;
        
        int va = val[u];
        ma[va]++;
        res++;
        
        for (int v : adjList[u]) {
            if (v == fa) continue;
            auto son = dfs(v, u);
            auto it = son.begin();
            while (it != son.end() && it->first < va) {
                it = son.erase(it);
            }
            merge(ma, son);
        }
        
        return ma;
    }
    
    int numberOfGoodPaths(vector<int>& vals, vector<vector<int>>& edges) {
        val = vals;
        n = val.size();
        adjList.assign(n, vector<int>());
        
        for (auto&& e : edges) {
            adjList[e[0]].push_back(e[1]);
            adjList[e[1]].push_back(e[0]);
        }
        
        res = 0;
        
        dfs(0, -1);
        
        return res;
    }
};