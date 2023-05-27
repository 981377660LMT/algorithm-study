import { RangeModeQuery } from '../离线查询/根号分治/RangeModeQuery'

class MajorityChecker {
  private readonly _rmq: RangeModeQuery

  constructor(arr: number[]) {
    this._rmq = new RangeModeQuery(arr)
  }

  query(left: number, right: number, threshold: number): number {
    const [mode, freq] = this._rmq.query(left, right + 1)
    return freq >= threshold ? mode : -1
  }
}
