// https://leetcode.cn/problems/smallest-sufficient-team/
// # 1 <= req_skills.length <= 16
// # 1 <= people.length <= 60
// # 你规划了一份需求的技能清单 req_skills，并打算从备选人员名单 people 中选出些人组成一个「必要团队」
// # 请你返回 任一 规模最小的必要团队，团队成员用人员编号表示。
// !1125. 最小的必要团队-状压dp最短路写法

// 1.js的位运算会把数隐式转换成int32类型 比如1<<31就变成-2147483648了
// 2.2**60已经超出js的精度2**53-1了 转成Number会有误差
// 3.状压dp喜欢用python写记忆化 js不方便可以写成bfs(不喜欢写dp递推的写法，怕出错)

function smallestSufficientTeam(reqSkills: string[], people: string[][]): number[] {
  const skill2id = new Map<string, number>()
  reqSkills.forEach((skill, id) => skill2id.set(skill, id))
  const states = new Uint32Array(people.length)
  for (let i = 0; i < people.length; i++) {
    people[i].forEach(skill => {
      states[i] |= 1 << skill2id.get(skill)!
    })
  }

  const dist = new Uint32Array(1 << reqSkills.length).fill(-1) // 记录到达每个状态的最小人数
  const pre = new Int32Array(1 << reqSkills.length).fill(-1) // 记录从哪个状态转移过来
  const select = new Int8Array(1 << reqSkills.length).fill(-1) // 记录每个状态新选择了哪个人
  dist[0] = 0
  const stack = [[0, 0]] // [dist, state]
  while (stack.length) {
    const [curDist, curState] = stack.pop()!
    if (curDist > dist[curState]) continue
    for (let i = 0; i < people.length; i++) {
      const newState = curState | states[i]
      if (dist[newState] > curDist + 1) {
        dist[newState] = curDist + 1
        pre[newState] = curState
        select[newState] = i
        stack.push([dist[newState], newState])
      }
    }
  }

  const res = []
  let cur = (1 << reqSkills.length) - 1
  while (cur > 0) {
    res.push(select[cur])
    cur = pre[cur]
  }
  return res
}

if (require.main === module) {
  console.log(
    smallestSufficientTeam(
      ['hdbxcuzyzhliwv', 'uvwlzkmzgis', 'sdi', 'bztg', 'ylopoifzkacuwp', 'dzsgleocfpl'],
      [
        ['hdbxcuzyzhliwv', 'dzsgleocfpl'],
        ['hdbxcuzyzhliwv', 'sdi', 'ylopoifzkacuwp', 'dzsgleocfpl'],
        ['bztg', 'ylopoifzkacuwp'],
        ['bztg', 'dzsgleocfpl'],
        ['hdbxcuzyzhliwv', 'bztg'],
        ['dzsgleocfpl'],
        ['uvwlzkmzgis'],
        ['dzsgleocfpl'],
        ['hdbxcuzyzhliwv'],
        [],
        ['dzsgleocfpl'],
        ['hdbxcuzyzhliwv'],
        [],
        ['hdbxcuzyzhliwv', 'ylopoifzkacuwp'],
        ['sdi'],
        ['bztg', 'dzsgleocfpl'],
        ['hdbxcuzyzhliwv', 'uvwlzkmzgis', 'sdi', 'bztg', 'ylopoifzkacuwp'],
        ['hdbxcuzyzhliwv', 'sdi'],
        ['hdbxcuzyzhliwv', 'ylopoifzkacuwp'],
        ['sdi', 'bztg', 'ylopoifzkacuwp', 'dzsgleocfpl'],
        ['dzsgleocfpl'],
        ['sdi', 'ylopoifzkacuwp'],
        ['hdbxcuzyzhliwv', 'uvwlzkmzgis', 'sdi'],
        [],
        [],
        ['ylopoifzkacuwp'],
        [],
        ['sdi', 'bztg'],
        ['bztg', 'dzsgleocfpl'],
        ['sdi', 'bztg']
      ]
    )
  )
}
