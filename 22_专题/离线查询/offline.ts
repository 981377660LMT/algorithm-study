const query = [77, 4, 2, 5, 1]
const res = Array(query.length).fill(0)

const ids = Array.from({ length: query.length }, (_, i) => i).sort(
  // 自定义处理query的顺序(搜索范围是从'小'到'大')
  (i1, i2) => query[i1] - query[i2]
)

for (const id of ids) {
  console.log('离线查询顺序' + id, query[id])
  res[id] = id + query[id]
}

console.log(res)
export {}
