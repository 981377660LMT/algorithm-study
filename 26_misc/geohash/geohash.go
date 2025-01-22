// https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484507&idx=1&sn=ba2c821a75660b7aa6ec9e150b93f06e
//
// geohash 服务模块需要对外提供的几个 API 整理如下：
//
// • Hash：将用户输入的经纬度 lon、lat 转为 geohash 字符串
// • Get：通过传入的 geohash 字符串，获取到对应于矩形区域块的 GEOEntry 实例
// • Add：通过用户传入的经纬度 lon、lat，构造出 point 实例并添加到对应的矩形区域中
// • ListByPrefix：通过用户输入的 geohash 字符串，获取到对应矩形区域块内所有子矩形区域块的 GEOEntry 实例（包含本身）
// • Rem：通过用户输入的 geohash 字符串，删除对应矩形区域块的 GEOEntry
// • ListByRadiusM：通过用户输入的中心点 lon、lat，以及对应的距离范围 radius，返回范围内所有的点集合

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

func main() {
	S := NewGEOService[int]()
	fmt.Println(S.Hash(116.404, 39.915))
}

var distRank = []int{0, 20, 150, 600, 5000, 20000, 160000}

// 将十进制数值转为 base32 编码
var Base32 = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'M', 'N',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func base32ToIndex(base32 byte) int {
	if base32 >= '0' && base32 <= '9' {
		return int(base32 - '0')
	}
	if base32 >= 'B' && base32 <= 'H' {
		return int(base32 - 'B' + 26)
	}
	if base32 >= 'J' && base32 <= 'K' {
		return int(base32 - 'J' + 33)
	}
	if base32 >= 'M' && base32 <= 'N' {
		return int(base32 - 'J' + 35)
	}
	if base32 >= 'P' && base32 <= 'Z' {
		return int(base32 - 'J' + 37)
	}
	return -1
}

type GEOService[V any] struct {
	mux  sync.RWMutex
	root *geoTrieNode[V]
}

func NewGEOService[V any]() *GEOService[V] {
	return &GEOService[V]{
		root: &geoTrieNode[V]{},
	}
}

// 通过二分的方式，把经纬度转换成对应的 hash 字符串，固定为 40 位二进制数字组成
// 每 5 个 bit 位由一个 base32 映射得到，因此总共由 8 位 base32 字符组成
func (g *GEOService[V]) Hash(lon, lat float64) string {
	// lon 通过二分转为 20 个二进制 bit 位
	lonBits := g.getBinaryBits(&strings.Builder{}, lon, -180, 180)

	// lat 通过二分转为 20 个二进制 bit 位
	latBits := g.getBinaryBits(&strings.Builder{}, lat, -90, 90)

	// 经纬度交错在一次，每 5 个 bit 一组，转换成 base32 字符串
	var geoHash strings.Builder
	var fiveBitsBuffer strings.Builder
	for i := 1; i <= 40; i++ {
		if i&1 == 1 {
			fiveBitsBuffer.WriteByte(lonBits[(i-1)>>1])
		} else {
			fiveBitsBuffer.WriteByte(latBits[(i-1)>>1])
		}
		// 位数不足 5 位，则继续累积位数
		if i%5 != 0 {
			continue
		}

		// 将每个 bit 位组成的二进制数字转为十进制数值
		val, _ := strconv.ParseInt(fiveBitsBuffer.String(), 2, 64)
		// 转为 base32 编码，并写入到 stringBuilder 中
		geoHash.WriteByte(Base32[val])

		fiveBitsBuffer.Reset()
	}

	return geoHash.String()
}

func (g *GEOService[V]) Get(geoHash string) (*GEOEntry[V], bool) {
	// 加读锁，和写操作隔离，保证并发安全
	g.mux.RLock()
	defer g.mux.RUnlock()

	target := g.get(geoHash)
	if target == nil || !target.end {
		return nil, false
	}
	return &target.GEOEntry, true
}

func (g *GEOService[V]) Add(lon, lat float64, val V) {
	geoHash := g.Hash(lon, lat)

	// 加写锁
	g.mux.Lock()
	defer g.mux.Unlock()

	target := g.get(geoHash)
	if target != nil && target.end {
		// 往目标节点对应的 GEOEntry 中追加这个新的节点
		target.add(lon, lat, val)
		return
	}

	//  目标节点不存在，则执行遍历
	node := g.root
	for i := 0; i < len(geoHash); i++ {
		index := base32ToIndex(geoHash[i])
		if node.children[index] == nil {
			node.children[index] = &geoTrieNode[V]{}
		}
		node.children[index].passCnt++
		node = node.children[index]
	}
	node.end = true
	node.add(lon, lat, val)
	node.Hash = geoHash
}

func (g *GEOService[V]) ListByPrefix(prefix string) []*GEOEntry[V] {
	// 加读锁
	g.mux.RLock()
	defer g.mux.RUnlock()

	target := g.get(prefix)
	if target == nil {
		return nil
	}
	return target.dfs()
}

func (g *GEOService[V]) Rem(geoHash string) bool {
	// 加写锁
	g.mux.Lock()
	defer g.mux.Unlock()

	target := g.get(geoHash)
	if target == nil || !target.end {
		return false
	}

	node := g.root
	for i := 0; i < len(geoHash); i++ {
		index := base32ToIndex(geoHash[i])
		node.children[index].passCnt--
		// 如果某个 child passCnt 减至 0，则直接丢弃整个 child 返回
		if node.children[index].passCnt == 0 {
			node.children[index] = nil
			return true
		}
		node = node.children[index]
	}
	node.end = false
	return true
}

// 查询指定范围范围内
// radius 单位为 m.
// m 转为对应的经纬度 经度1度≈111m；纬度1度≈111m
func (g *GEOService[V]) ListByRadiusM(lon, lat float64, radiusM int) ([]*GEOPoint[V], error) {
	// 加一把读锁
	g.mux.RLock()
	defer g.mux.RUnlock()

	// 1 根据用户指定的查询范围，确定所需要的 geohash 字符串的长度，保证对应的矩形区域长度大于等于 radiusM
	bitsLen, err := g.getBitsLengthByRadiusM(2 * radiusM)
	if err != nil {
		return nil, err
	}

	// 2 针对于传入的 lon、lat 中心点，沿着上、下、左、右方向进行偏移，获取到包含自身在内的总共 9 个点的点矩阵
	// 核心是为了保证通过 9 个点获取到的矩形区域一定能完全把检索范围包含在内
	points := g.getCenterPoints(lon, lat, radiusM)

	// 3. 针对这9个点，通过 ListByPrefix 方法，分别取出区域内的所有子 GEOEntry
	var rawEntries []*GEOEntry[V]
	for i := 0; i < len(points); i++ {
		geoHash := g.Hash(points[i][0], points[i][1])[:bitsLen]
		rawEntries = append(rawEntries, g.ListByPrefix(geoHash)...)
	}

	// 4. 针对所有 entry，取出其中包含的所有 point
	// 取出 point 之后，计算其与 center point 的相对距离，如果超过范围则进行过滤
	var geoPoints []*GEOPoint[V]
	// 遍历所有的 entry
	for _, rawEntry := range rawEntries {
		// 遍历每个 entry 中所有的 point
		for _, rawPoint := range rawEntry.GetPoints() {
			// 计算一个 point 与中心点 lon、lat 的相对距离
			dist := g.calDistance(lon, lat, rawPoint.Lon, rawPoint.Lat)
			// 如果相对距离大于 radiusM，则进行过滤
			if dist > float64(radiusM) {
				continue
			}
			// 相对距离满足条件，则将 point 追加到 list 中(可能会有重复的 point)
			geoPoints = append(geoPoints, rawPoint)
		}
	}

	return geoPoints, nil
}

func (g *GEOService[V]) getBitsLengthByRadiusM(radiusM int) (int, error) {
	if radiusM > 160*1000 || radiusM < 0 {
		return -1, fmt.Errorf("invalid radius: %d", radiusM)
	}

	var i int
	for {
		if radiusM <= distRank[i+1] {
			return 8 - i, nil
		}
		i++
	}
}

func (g *GEOService[V]) calDistance(lon1, lat1, lon2, lat2 float64) float64 {
	return 111 * (math.Pow(lon1-lon2, 2) + math.Pow(lat1-lat2, 2))
}

// 获取中心点的矩阵.
func (g *GEOService[V]) getCenterPoints(lon, lat float64, radiusM int) [9][2]float64 {
	dif := float64(radiusM) / 111
	left := lon - dif
	if left < -180 {
		left += 360
	}
	right := lon + dif
	if right > 180 {
		right -= 360
	}
	bot := lat - dif
	if bot < -90 {
		bot += 180
	}
	top := lat + dif
	if top > 90 {
		top -= 180
	}

	return [9][2]float64{
		{
			left, top,
		},
		{
			lon, top,
		},
		{
			right, top,
		},
		{
			left, lat,
		},
		{
			lon, lat,
		},
		{
			right, lat,
		},
		{
			left, bot,
		},
		{
			lon, bot,
		},
		{
			right, bot,
		},
	}
}

func (gn *geoTrieNode[V]) dfs() []*GEOEntry[V] {
	var entries []*GEOEntry[V]
	if gn.end {
		entries = append(entries, &gn.GEOEntry)
	}

	for i := 0; i < len(gn.children); i++ {
		if gn.children[i] == nil {
			continue
		}
		entries = append(entries, gn.children[i].dfs()...)
	}

	return entries
}

func (g *GEOService[V]) get(geoHash string) *geoTrieNode[V] {
	node := g.root
	for i := 0; i < len(geoHash); i++ {
		index := base32ToIndex(geoHash[i])
		if index == -1 || node.children[index] == nil {
			return nil
		}
		node = node.children[index]
	}
	return node
}

func (g *GEOService[V]) getBinaryBits(bits *strings.Builder, val, start, end float64) string {
	for i := 0; i < 20; i++ {
		mid := (start + end) / 2
		if val < mid {
			bits.WriteByte('0')
			end = mid
		} else {
			bits.WriteByte('1')
			start = mid
		}
	}
	return bits.String()
}

type geoTrieNode[V any] struct {
	// 子节点列表，顺序与 base32 一致
	children [32]*geoTrieNode[V]
	passCnt  int
	end      bool

	// geohash 字符串对应的矩形区域
	GEOEntry[V]
}

type GEOEntry[V any] struct {
	Points map[string]V
	// 矩形区域对应的 geohash 字符串
	Hash string
}

// 获取到矩形区域内所有的点集合
func (g *GEOEntry[V]) GetPoints() []*GEOPoint[V] {
	points := make([]*GEOPoint[V], 0, len(g.Points))
	for key, val := range g.Points {
		lon, lat := g.lonlat(key)
		points = append(points, &GEOPoint[V]{
			Lon: lon,
			Lat: lat,
			Val: val,
		})
	}
	return points
}

func (g *GEOEntry[V]) add(lon, lat float64, val V) {
	if g.Points == nil {
		g.Points = make(map[string]V)
	}
	g.Points[g.key(lon, lat)] = val
}

func (g *GEOEntry[V]) key(lon, lat float64) string {
	return fmt.Sprintf("%v_%v", lon, lat)
}

func (g *GEOEntry[V]) lonlat(key string) (lon, lat float64) {
	info := strings.Split(key, "_")
	lon, _ = strconv.ParseFloat(info[0], 64)
	lat, _ = strconv.ParseFloat(info[1], 64)
	return
}

type GEOPoint[V any] struct {
	Lon, Lat float64
	Val      V
}
