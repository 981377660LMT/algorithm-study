// 维护图中桥的个数，支持增加边的操作
// https://ei1333.github.io/library/graph/connected-components/incremental-bridge-connectivity.hpp
// 概要
// 辺の追加クエリのみ存在するとき, 二重辺連結成分を効率的に管理するデータ構造.

// 使い方
// IncrementalBridgeConnectivity(sz): sz 頂点で初期化する.
// Find(k): 頂点 k が属する二重辺連結成分(の代表元)を求める.
// GetBridgeSize(): 現在の橋の個数を返す.
// AddEdge(x, y): 頂点 x と y との間に無向辺を追加する.
// TODO
