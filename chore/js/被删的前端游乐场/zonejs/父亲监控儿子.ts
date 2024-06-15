import 'zone.js'

// 创建一个裁判zone，当做teamA和teamB的父zone
const zoneJudgement = Zone.current.fork({
  name: 'judgement',
  properties: {
    // 存放teamA、teamB的排序结果
    result: []
  },

  // 异步任务状态改变时的回调
  onHasTask: (parentZoneDelegate, currentZone, targetZone, hasTaskState) => {
    // setTimeout属于宏任务，!hasTaskState.macroTask标识有宏任务执行完毕
    if (!hasTaskState.macroTask) {
      // 裁判任务执行结束
      switch (targetZone.name) {
        case 'judgement':
          console.log(currentZone.get('result'))
          break
        // A组排序任务执行结束
        case 'teamA':
          currentZone.get('result').push({
            teamA: targetZone.get('team').team
          })
          break
        // B组排序任务执行结束
        case 'teamB':
          currentZone.get('result').push({
            teamB: targetZone.get('team').team
          })
          break
        default:
          break
      }
    }

    // 事件上抛
    parentZoneDelegate.hasTask(targetZone, hasTaskState)
  }
})

const zoneA = zoneJudgement.fork({
  name: 'teamA',
  properties: {
    team: {
      name: 'teamA',
      team: [],
      sort: function () {
        thinking(() => {
          this.team.push(this.team.length + 1)
        })
        thinking(() => {
          this.team.push(this.team.length + 1)
        })
        thinking(() => {
          this.team.push(this.team.length + 1)
        })
      }
    }
  }
})
const zoneB = zoneJudgement.fork({
  name: 'teamB',
  properties: {
    team: {
      name: 'teamB',
      team: [],
      sort: function () {
        thinking(() => {
          this.team.unshift(this.team.length + 1)
        })
        thinking(() => {
          this.team.unshift(this.team.length + 1)
        })
        thinking(() => {
          this.team.unshift(this.team.length + 1)
        })
      }
    }
  }
})

function judgement() {
  zoneA.run(() => {
    const currentZone = Zone.current
    const team = currentZone.get('team')
    team.sort()
  })

  zoneB.run(() => {
    const currentZone = Zone.current
    const team = currentZone.get('team')
    team.sort()
  })
}

zoneJudgement.run(judgement) // [ { teamA: [ 1, 2, 3 ] }, { teamB: [ 3, 2, 1 ] } ]

function thinking(f: () => void): void {
  setTimeout(f, Math.random() * 1000)
}

export {}
