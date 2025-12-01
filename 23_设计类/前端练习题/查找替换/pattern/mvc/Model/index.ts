import { ISearchResult, ISearchState, ISearchStore } from '../types'

const INITIAL_STATE: ISearchState = {
  keyword: '',
  results: new Map(),
  activeResult: undefined
}

/**
 * 实际项目中可替换为 Mobx/Zustand.
 */
export class SearchStore implements ISearchStore {
  private _state: ISearchState = INITIAL_STATE
  private _listeners: ((state: ISearchState, prevState: ISearchState) => void)[] = []

  dispose(): void {
    this._listeners = []
    this._state = INITIAL_STATE
  }

  setState(newState: Partial<ISearchState>) {
    const prevState = this._state
    this._state = { ...this._state, ...newState }
    this._listeners.forEach(listener => listener(this._state, prevState))
  }

  subscribe(listener: (state: ISearchState, prevState: ISearchState) => void): () => void {
    this._listeners.push(listener)
    return () => {
      this._listeners = this._listeners.filter(l => l !== listener)
    }
  }

  setModuleResults(moduleId: string, results: ISearchResult[]): void {
    const newResults = new Map(this._state.results)
    if (results.length === 0) {
      newResults.delete(moduleId)
    } else {
      newResults.set(moduleId, results)
    }
    this.setState({ results: newResults })
  }

  getModuleResults(moduleId: string): ISearchResult[] {
    return this._state.results.get(moduleId) || []
  }

  getAllModuleResults(): Map<string, ISearchResult[]> {
    return this._state.results
  }

  setActiveResult(result: ISearchResult | undefined): void {
    this.setState({ activeResult: result })
  }

  getActiveResult(): ISearchResult | undefined {
    return this._state.activeResult
  }
}

export {}
