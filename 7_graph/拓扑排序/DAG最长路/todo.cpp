#include <algorithm>
#include <cassert>
#include <vector>
#include <queue>


class TopologicalSort {
    const int num_nodes;                   // ノード数
    std::vector<std::vector<int>> graph;   // グラフ

public:
    std::vector<int> result; // ソート結果（結果が複数ある場合は辞書順で最小のもの）

public:
    TopologicalSort(int num_nodes) : num_nodes(num_nodes) {
        this->graph.resize(num_nodes);
    }

    // u -> vのdirect edgeを追加
    void add_directed_edge(int u, int v) {
        assert(u != v);
        this->graph[u].emplace_back(v);
    }

    // O(E + V log V)
    // sortができたらtrueを返す
    bool sort() {
        // すべてのノードの入次数を算出
        std::vector<int> indegree(this->num_nodes);
        for (int u = 0; u < this->num_nodes; ++u) {
            for (int v: this->graph[u]) {
                indegree[v]++;
            }
        }

        std::priority_queue<int, std::vector<int>, std::greater<>> que;
        std::vector<bool> used(this->num_nodes);
        for (int u = 0; u < this->num_nodes; ++u) {
            if (indegree[u] == 0) {
                que.push(u);
                used[u] = true;
            }
        }

        while (not que.empty()) {
            const int u = que.top();
            que.pop();
            this->result.emplace_back(u);

            for (int v: this->graph[u]) {
                indegree[v]--;
                if (indegree[v] == 0 and not used[v]) {
                    used[v] = true;
                    que.push(v);
                }
            }
        }

        return (int) this->result.size() == this->num_nodes;
    }

    // グラフの有向パスのうち最長のものの長さ
    // ソート結果が一意なら返り値はN - 1になり，元のグラフはハミルトン路を含むグラフである
    int longest_distance() {
        int ans = 0;
        std::vector<int> distance(this->num_nodes, 0);
        for (int u: this->result) {
            for (int v: this->graph[u]) {
                distance[v] = std::max(distance[v], distance[u] + 1);
                ans = std::max(ans, distance[v]);
            }
        }

        return ans;
    }
};