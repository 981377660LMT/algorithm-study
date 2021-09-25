type PassengerId = number
type StartStationName = string
type StartTime = number
type TwoStationNames = string
type Count = number
type TimeSum = number

class UndergroundSystem {
  private checkInInfo: Map<PassengerId, [StartStationName, StartTime]>
  private checkOutInfo: Map<TwoStationNames, [TimeSum, Count]>

  constructor() {
    this.checkInInfo = new Map()
    this.checkOutInfo = new Map()
  }

  // 编号为 id 的乘客在 t 时刻进入地铁站 stationName
  checkIn(id: number, stationName: string, t: number): void {
    this.checkInInfo.set(id, [stationName, t])
  }

  // 编号为 id 的乘客在 t 时刻离开地铁站 stationName
  checkOut(id: number, stationName: string, t: number): void {
    const [startStationName, startTime] = this.checkInInfo.get(id)!
    const key = `${startStationName}#${stationName}`
    const [sum, count] = this.checkOutInfo.get(key) || [0, 0]
    this.checkOutInfo.set(key, [sum + t - startTime, count + 1])
  }

  // 返回从地铁站 startStation 到地铁站 endStation 的平均花费时间。
  getAverageTime(startStation: string, endStation: string): number {
    const key = `${startStation}#${endStation}`
    const [sum, count] = this.checkOutInfo.get(key)!
    return sum / count
  }
}

const system = new UndergroundSystem()
system.checkIn(45, 'Leyton', 3)
system.checkIn(32, 'Paradise', 8)
system.checkIn(27, 'Leyton', 10)
system.checkOut(45, 'Waterloo', 15)
system.checkOut(27, 'Waterloo', 20)
system.checkOut(32, 'Cambridge', 22)
console.log(system.getAverageTime('Paradise', 'Cambridge')) // 14
console.log(system.getAverageTime('Leyton', 'Waterloo')) // 11
system.checkIn(10, 'Leyton', 24)
console.log(system.getAverageTime('Leyton', 'Waterloo')) // 11
system.checkOut(10, 'Waterloo', 38)
console.log(system.getAverageTime('Leyton', 'Waterloo')) // 12

// 符合自然的思路：
// 1.上车记录乘客id 与[起点车站，起点时间] 的关系
// 2. 下车将(起点车站,结束车站)作为key 记录 [总次数，总时间和]

export {}
