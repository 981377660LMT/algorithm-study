// 给你一些区域列表 regions ，`每个列表的第一个区域都包含这个列表内所有其他区域。`
// 如果区域 X 包含区域 Y ，那么区域 X  比区域 Y 大。
// 给定两个区域 region1 和 region2 ，找到同时包含这两个区域的 最小 区域。

function findSmallestRegion(regions: string[][], region1: string, region2: string): string {
  const parent = new Map<string, string>()
  for (const [p, ...children] of regions) {
    for (const child of children) {
      parent.set(child, p)
    }
  }

  const visited = new Set<string>()
  while (parent.has(region1)) {
    visited.add(region1)
    region1 = parent.get(region1)!
  }

  while (parent.has(region2)) {
    if (visited.has(region2)) return region2
    region2 = parent.get(region2)!
  }

  return region2
}

console.log(
  findSmallestRegion(
    [
      ['Earth', 'North America', 'South America'],
      ['North America', 'United States', 'Canada'],
      ['United States', 'New York', 'Boston'],
      ['Canada', 'Ontario', 'Quebec'],
      ['South America', 'Brazil'],
    ],
    'Quebec',
    'New York'
  )
)
