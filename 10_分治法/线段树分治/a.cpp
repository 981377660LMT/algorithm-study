#pragma once
#include <algorithm>
#include <limits>
#include <set>
#include <utility>
#include <vector>

enum class DyConOperation {
    Begins = 1,
    Ends = 2,
    Event = 3,
};

template <class Time = int> struct offline_dynamic_connectivity {

    std::vector<std::pair<Time, int>> queries;

    struct Edge {
        Time since;
        Time until;
        int edge_id;
    };
    std::vector<Edge> edges;

    offline_dynamic_connectivity() = default;

    void add_observation(Time clk, int event_id) { queries.emplace_back(clk, event_id); }

    void apply_time_range(Time since, Time until, int edge_id) {
        edges.push_back(Edge{since, until, edge_id});
    }

    struct Procedure {
        DyConOperation op;
        int id_;
    };

    std::vector<std::vector<Procedure>> nodes;
    std::vector<Procedure> seg;

    void rec(int now) {
        seg.insert(seg.end(), nodes[now].cbegin(), nodes[now].cend());
        if (now * 2 < int(nodes.size())) rec(now * 2);
        if (now * 2 + 1 < int(nodes.size())) rec(now * 2 + 1);

        for (auto itr = nodes[now].rbegin(); itr != nodes[now].rend(); ++itr) {
            if (itr->op == DyConOperation::Begins) {
                seg.push_back(Procedure{DyConOperation::Ends, itr->id_});
            }
        }
    }

    std::vector<Procedure> generate() {
        if (queries.empty()) return {};

        std::vector<Time> query_ts;
        {
            query_ts.reserve(queries.size());
            for (const auto &p : queries) query_ts.push_back(p.first);
            std::sort(query_ts.begin(), query_ts.end());
            query_ts.erase(std::unique(query_ts.begin(), query_ts.end()), query_ts.end());

            std::vector<int> t_use(query_ts.size() + 1);
            t_use.at(0) = 1;

            for (const Edge &e : edges) {
                t_use[std::lower_bound(query_ts.begin(), query_ts.end(), e.since) - query_ts.begin()] =
                    1;
                t_use[std::lower_bound(query_ts.begin(), query_ts.end(), e.until) - query_ts.begin()] =
                    1;
            }
            for (int i = 1; i < int(query_ts.size()); ++i) {
                if (!t_use[i]) query_ts[i] = query_ts[i - 1];
            }

            query_ts.erase(std::unique(query_ts.begin(), query_ts.end()), query_ts.end());
        }

        const int N = query_ts.size();
        int D = 1;
        while (D < int(query_ts.size())) D *= 2;

        nodes.assign(D + N, {});

        for (const Edge &e : edges) {
            int l =
                D + (std::lower_bound(query_ts.begin(), query_ts.end(), e.since) - query_ts.begin());
            int r =
                D + (std::lower_bound(query_ts.begin(), query_ts.end(), e.until) - query_ts.begin());

            while (l < r) {
                if (l & 1) nodes[l++].push_back(Procedure{DyConOperation::Begins, e.edge_id});
                if (r & 1) nodes[--r].push_back(Procedure{DyConOperation::Begins, e.edge_id});
                l >>= 1, r >>= 1;
            }
        }

        for (const auto &op : queries) {
            int t = std::upper_bound(query_ts.begin(), query_ts.end(), op.first) - query_ts.begin();
            nodes.at(t + D - 1).push_back(Procedure{DyConOperation::Event, op.second});
        }
        seg.clear();
        rec(1);
        return seg;
    }
};