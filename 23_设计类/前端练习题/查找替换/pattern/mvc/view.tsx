import React, { useEffect, useRef, useState } from 'react'
import Draggable from 'react-draggable'

import { SearchController } from './model'
import './view.css'

enum SearchMode {
  Search = 'search',
  Replace = 'replace'
}

const SearchPanel: React.FC<{ controller: SearchController }> = ({ controller }) => {
  // 1. 业务状态 (来自 Controller)
  const [state, setState] = useState(controller.store.state)

  // 2. UI 状态 (View 自治)
  const [searchMode, setSearchMode] = useState<SearchMode>(SearchMode.Search)
  const [isReplaceVisible, setIsReplaceVisible] = useState(true) // 比如可以通过 props 控制是否允许替换
  const [isExpanded, setIsExpanded] = useState(true) // 面板展开/折叠

  // 3. 输入状态
  const [inputValue, setInputValue] = useState(state.keyword)
  const [replaceValue, setReplaceValue] = useState('')

  const isCompositionRef = useRef(false)
  const nodeRef = useRef(null)

  // 1. 仅负责订阅状态更新 (只在 controller 实例变化时执行一次)
  useEffect(() => {
    const dispose = controller.store.subscribe(setState)
    return () => dispose()
  }, [controller])

  // 2. 负责同步外部状态到本地 input (当 state.keyword 变化时)
  useEffect(() => {
    if (state.keyword !== inputValue) {
      setInputValue(state.keyword)
    }
  }, [state.keyword]) // 移除 inputValue 依赖，避免循环

  useEffect(() => {
    const timer = setTimeout(() => {
      // 只有当输入值变化，且不在输入法组合过程中时，才触发搜索
      if (inputValue !== controller.store.state.keyword && !isCompositionRef.current) {
        controller.search(inputValue)
      }
    }, 300)
    return () => clearTimeout(timer)
  }, [inputValue, controller])

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'f') {
        e.preventDefault()
        // 聚焦输入框逻辑...
      }
      if (e.key === 'Enter') {
        if (e.shiftKey) controller.prev()
        else controller.next()
      }
      if (e.key === 'Escape') {
        // 关闭搜索框或清空
        controller.search('')
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [controller])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value)
  }

  const handleCompositionStart = () => {
    isCompositionRef.current = true
  }

  const handleCompositionEnd = (e: React.CompositionEvent<HTMLInputElement>) => {
    isCompositionRef.current = false
    setInputValue(e.currentTarget.value)
  }

  const toggleMode = (mode: SearchMode) => {
    setSearchMode(mode)
  }

  const handleReplaceAll = () => {
    const count = state.results.filter(r => r.canReplace).length
    if (count > 1000) {
      alert('暂不支持全部替换超过1000个可替换结果') // 简单模拟 Toast
      return
    }
    controller.replaceAll(replaceValue)
  }

  return (
    <Draggable handle=".drag-handle" nodeRef={nodeRef}>
      <div ref={nodeRef} className="search-modal">
        {/* 拖拽手柄 + 折叠按钮 */}
        <div className="header">
          <div className="drag-handle" style={{ cursor: 'move', flex: 1 }}>
            :::
          </div>
          <button onClick={() => setIsExpanded(!isExpanded)}>{isExpanded ? '_' : '[]'}</button>
        </div>

        {isExpanded && (
          <>
            {/* 纯 UI 状态控制的 Tabs */}
            <div className="tabs">
              <button
                className={searchMode === SearchMode.Search ? 'active' : ''}
                onClick={() => toggleMode(SearchMode.Search)}
              >
                查找
              </button>
              {isReplaceVisible && (
                <button
                  className={searchMode === SearchMode.Replace ? 'active' : ''}
                  onClick={() => toggleMode(SearchMode.Replace)}
                >
                  替换
                </button>
              )}
            </div>

            {/* 查找输入框 */}
            <div className="input-row">
              <input
                value={inputValue}
                onChange={e => setInputValue(e.target.value)}
                onCompositionStart={() => (isCompositionRef.current = true)}
                onCompositionEnd={e => {
                  isCompositionRef.current = false
                  setInputValue(e.currentTarget.value)
                }}
                placeholder="查找..."
              />
              <span className="count">
                {state.results.length > 0 ? state.currentIndex + 1 : 0} / {state.results.length}
              </span>
            </div>

            {/* 替换输入框 (根据本地 searchMode 渲染) */}
            {searchMode === SearchMode.Replace && (
              <div className="input-row">
                <input
                  value={replaceValue}
                  onChange={e => setReplaceValue(e.target.value)}
                  placeholder="替换为..."
                />
              </div>
            )}

            {/* 操作按钮 */}
            <div className="actions">
              <button onClick={() => controller.prev()}>Prev</button>
              <button onClick={() => controller.next()}>Next</button>

              {searchMode === SearchMode.Replace && (
                <>
                  <button onClick={() => controller.replace(replaceValue)}>替换</button>
                  <button onClick={handleReplaceAll}>全部替换</button>
                </>
              )}
            </div>
          </>
        )}
      </div>
    </Draggable>
  )
}
