/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */
/* eslint-disable class-methods-use-this */

// #pragma once
// #include <algorithm>
// #include <limits>
// #include <set>
// #include <utility>
// #include <vector>

import { UnionFindArrayWithUndo } from '../../14_并查集/UnionFindWithUndo'
import { UnionFindArrayWithUndoAndWeight } from '../../14_并查集/UnionFindWithUndoAndWeight'

// enum class DyConOperation {
//     Begins = 1,
//     Ends = 2,
//     Event = 3,
// };

// template <class Time = int> struct offline_dynamic_connectivity {

//     std::vector<std::pair<Time, int>> queries;

//     struct Edge {
//         Time since;
//         Time until;
//         int edge_id;
//     };
//     std::vector<Edge> edges;

//     offline_dynamic_connectivity() = default;

//     void add_observation(Time clk, int event_id) { queries.emplace_back(clk, event_id); }

//     void apply_time_range(Time since, Time until, int edge_id) {
//         edges.push_back(Edge{since, until, edge_id});
//     }

//     struct Procedure {
//         DyConOperation op;
//         int id_;
//     };

//     std::vector<std::vector<Procedure>> nodes;
//     std::vector<Procedure> seg;

//     void rec(int now) {
//         seg.insert(seg.end(), nodes[now].cbegin(), nodes[now].cend());
//         if (now * 2 < int(nodes.size())) rec(now * 2);
//         if (now * 2 + 1 < int(nodes.size())) rec(now * 2 + 1);

//         for (auto itr = nodes[now].rbegin(); itr != nodes[now].rend(); ++itr) {
//             if (itr->op == DyConOperation::Begins) {
//                 seg.push_back(Procedure{DyConOperation::Ends, itr->id_});
//             }
//         }
//     }

//     std::vector<Procedure> generate() {
//         if (queries.empty()) return {};

//         std::vector<Time> query_ts;
//         {
//             query_ts.reserve(queries.size());
//             for (const auto &p : queries) query_ts.push_back(p.first);
//             std::sort(query_ts.begin(), query_ts.end());
//             query_ts.erase(std::unique(query_ts.begin(), query_ts.end()), query_ts.end());

//             std::vector<int> t_use(query_ts.size() + 1);
//             t_use.at(0) = 1;

//             for (const Edge &e : edges) {
//                 t_use[std::lower_bound(query_ts.begin(), query_ts.end(), e.since) - query_ts.begin()] =
//                     1;
//                 t_use[std::lower_bound(query_ts.begin(), query_ts.end(), e.until) - query_ts.begin()] =
//                     1;
//             }
//             for (int i = 1; i < int(query_ts.size()); ++i) {
//                 if (!t_use[i]) query_ts[i] = query_ts[i - 1];
//             }

//             query_ts.erase(std::unique(query_ts.begin(), query_ts.end()), query_ts.end());
//         }

//         const int N = query_ts.size();
//         int D = 1;
//         while (D < int(query_ts.size())) D *= 2;

//         nodes.assign(D + N, {});

//         for (const Edge &e : edges) {
//             int l =
//                 D + (std::lower_bound(query_ts.begin(), query_ts.end(), e.since) - query_ts.begin());
//             int r =
//                 D + (std::lower_bound(query_ts.begin(), query_ts.end(), e.until) - query_ts.begin());

//             while (l < r) {
//                 if (l & 1) nodes[l++].push_back(Procedure{DyConOperation::Begins, e.edge_id});
//                 if (r & 1) nodes[--r].push_back(Procedure{DyConOperation::Begins, e.edge_id});
//                 l >>= 1, r >>= 1;
//             }
//         }

//         for (const auto &op : queries) {
//             int t = std::upper_bound(query_ts.begin(), query_ts.end(), op.first) - query_ts.begin();
//             nodes.at(t + D - 1).push_back(Procedure{DyConOperation::Event, op.second});
//         }
//         seg.clear();
//         rec(1);
//         return seg;
//     }
// };

/**
 * 线段树分治.
 * 线段树分治是一种处理动态修改和询问的离线算法.
 * 通过将某一元素的出现时间段在线段树上保存，我们可以 dfs 遍历整棵线段树，
 * 运用可撤销数据结构维护来得到每个时间点的答案.
 */
class OfflineDynamicConnectivity {
  private readonly _add: (edgeId: number) => void
  private readonly _undo: () => void
  private readonly _query: (queryId: number) => void
  private readonly _edges: { start: number; end: number; id: number }[] = []
  private readonly _queries: { time: number; id: number }[] = []
  private _nodes: number[][] = []
  private _edgeId = 0
  private _queryId = 0

  /**
   * dfs 遍历整棵线段树来得到每个时间点的答案.
   * @param add 添加编号为`edgeId`的边后产生的副作用.
   * @param undo 撤销一次`add`操作.
   * @param query 响应编号为`queryId`的查询.
   */
  constructor(options: { add: (edgeId: number) => void; undo: () => void; query: (queryId: number) => void } & ThisType<void>)
  constructor(add: (edgeId: number) => void, undo: () => void, query: (queryId: number) => void)
  constructor(arg1: any, arg2?: any, arg3?: any) {
    if (typeof arg1 === 'object') {
      this._add = arg1.add
      this._undo = arg1.undo
      this._query = arg1.query
    } else {
      this._add = arg1
      this._undo = arg2
      this._query = arg3
    }
  }

  /**
   * 在时间范围`[startTime, endTime)`内添加一条编号为`edgeId`的边.
   */
  addEdge(startTime: number, endTime: number, edgeId?: number): void {
    if (startTime >= endTime) return
    if (edgeId == undefined) edgeId = this._edgeId++
    this._edges.push({ start: startTime, end: endTime, id: edgeId })
  }

  /**
   * 在时间`time`时添加一个编号为`queryId`的查询.
   */
  addQuery(time: number, queryId?: number): void {
    if (queryId == undefined) queryId = this._queryId++
    this._queries.push({ time, id: queryId })
  }

  run(): void {
    if (!this._queries.length) return
    const times: number[] = Array(this._queries.length)
    for (let i = 0; i < this._queries.length; i++) times[i] = this._queries[i].time
    times.sort((a, b) => a - b)
    this._uniqueInplace(times)
    const usedTimes = new Uint8Array(times.length + 1)
    usedTimes[0] = 1
    for (let i = 0; i < this._edges.length; i++) {
      const e = this._edges[i]
      usedTimes[this._lowerBound(times, e.start)] = 1
      usedTimes[this._lowerBound(times, e.end)] = 1
    }
    for (let i = 1; i < times.length; i++) {
      if (!usedTimes[i]) times[i] = times[i - 1]
    }
    this._uniqueInplace(times)

    const n = times.length
    let offset = 1
    while (offset < n) offset <<= 1
    this._nodes = Array(offset + n)
    for (let i = 0; i < this._nodes.length; i++) this._nodes[i] = []
    for (let i = 0; i < this._edges.length; i++) {
      const e = this._edges[i]
      let left = offset + this._lowerBound(times, e.start)
      let right = offset + this._lowerBound(times, e.end)
      const eid = e.id * 2
      while (left < right) {
        if (left & 1) this._nodes[left++].push(eid) // add
        if (right & 1) this._nodes[--right].push(eid)
        left >>= 1
        right >>= 1
      }
    }

    for (let i = 0; i < this._queries.length; i++) {
      const q = this._queries[i]
      this._nodes[offset + this._upperBound(times, q.time) - 1].push(q.id * 2 + 1) // query
    }
    this._dfs(1)
  }

  private _dfs(now: number): void {
    const curNodes = this._nodes[now]
    for (let i = 0; i < curNodes.length; i++) {
      const id = curNodes[i]
      if (id & 1) {
        // query
        this._query((id - 1) / 2)
      } else {
        // add
        this._add(id / 2)
      }
    }
    if (now << 1 < this._nodes.length) this._dfs(now << 1)
    if (((now << 1) | 1) < this._nodes.length) this._dfs((now << 1) | 1)
    for (let i = curNodes.length - 1; ~i; i--) {
      if (!(curNodes[i] & 1)) {
        this._undo()
      }
    }
  }

  private _uniqueInplace(sorted: number[]): void {
    let slow = 0
    for (let fast = 0; fast < sorted.length; fast++) {
      if (sorted[fast] !== sorted[slow]) sorted[++slow] = sorted[fast]
    }
    sorted.length = slow + 1
  }

  private _lowerBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private _upperBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

export { OfflineDynamicConnectivity }

if (require.main === module) {
  // let add = 0
  // let undo = 0
  // let query = 0
  // const dc = new OfflineDynamicConnectivity({
  //   add(edgeId) {
  //     // console.log('add', edgeId)
  //     add++
  //   },
  //   undo() {
  //     // console.log('undo')
  //     undo++
  //   },
  //   query(queryId) {
  //     // console.log('query', queryId)
  //     query++
  //   }
  // })

  // // dc.addEdge(0, 100, 1) // 在时间范围[0, 100)内添加一条编号为1的边
  // // dc.addQuery(50, 2) // 在时间50时添加一个编号为2的查询
  // // dc.run()
  // const n = 1e5
  // const m = 1e5
  // const q = 1e5
  // for (let i = 0; i < m; i++) {
  //   dc.addEdge(Math.floor(Math.random() * n), Math.floor(Math.random() * n), i)
  // }
  // for (let i = 0; i < q; i++) {
  //   dc.addQuery(Math.floor(Math.random() * n), i)
  // }
  // console.time('offline_dynamic_connectivity')
  // dc.run()
  // console.timeEnd('offline_dynamic_connectivity')

  // console.log(add, undo, query)
  test()
  function test(): void {
    let add = 0
    let undo = 0
    let query = 0
    const dc = new OfflineDynamicConnectivity({
      add(edgeId) {
        // console.log('add', edgeId)
        add++
      },
      undo() {
        // console.log('undo')
        undo++
      },
      query(queryId) {
        // console.log('query', queryId)
        query++
      }
    })

    const n = 1e5
    for (let i = 0; i < 1e5; i++) {
      dc.addEdge(0, i)

      dc.addQuery(i)
    }
    dc.run()
    console.log(add, undo, query)
  }
}
