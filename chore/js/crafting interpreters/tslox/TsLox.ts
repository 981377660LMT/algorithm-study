import { Scanner } from './Scanner'

class TsLox {
  private _hadError = false

  run(source: string): void {
    const scanner = new Scanner(source, { reportError: this.error }) // !TOOD: dependency injection
  }

  error(line: number, message: string): void {
    this._report(line, '', message)
  }

  private _report(line: number, where: string, message: string): void {
    console.log(`[line ${line}] Error ${where}: ${message}`)
    this._hadError = true
  }
}
