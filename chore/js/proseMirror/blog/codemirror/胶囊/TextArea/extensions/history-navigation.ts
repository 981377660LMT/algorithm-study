import { EditorView } from '@codemirror/view';
import type { HistoryNavigationRefs } from '../types';

/**
 * 历史导航 keymap.
 */
export const createHistoryNavigationKeymap = (refs: HistoryNavigationRefs) => {
  const navigateUp = (view: EditorView): boolean => {
    const inputHistoryList = refs.inputHistory.current;
    if (inputHistoryList.length === 0) return false;

    const { state } = view;
    const cursorPos = state.selection.main.head;
    if (cursorPos !== 0) {
      view.dispatch({ selection: { anchor: 0 } });
      return true;
    }

    const currentIndex = refs.historyIndex.current;
    const nextIndex = currentIndex + 1;
    if (nextIndex >= inputHistoryList.length) return false;

    if (currentIndex === -1) {
      refs.draft.current = refs.valueRef.current;
    }

    refs.isNavigating.current = true;
    refs.setHistoryIndex(nextIndex);

    const historyValue = String(inputHistoryList[nextIndex] || '');
    refs.onChange(historyValue);

    refs.currentHistoryValue.current = historyValue;

    setTimeout(() => {
      const editorView = refs.editorRef.current?.$view;
      if (editorView) {
        editorView.dispatch({ selection: { anchor: 0 } });
      }
    }, 0);
    return true;
  };

  const navigateDown = (view: EditorView): boolean => {
    const inputHistoryList = refs.inputHistory.current;
    const currentIndex = refs.historyIndex.current;
    if (currentIndex === -1) return false;

    const { state } = view;
    const cursorPos = state.selection.main.head;
    const docLength = state.doc.length;
    if (cursorPos !== docLength) {
      view.dispatch({ selection: { anchor: docLength } });
      return true;
    }

    const currentValue = refs.valueRef.current;
    const originalHistoryValue = refs.currentHistoryValue.current;
    if (currentValue !== originalHistoryValue) {
      // 内容已修改，视为新的编辑记录，退出历史浏览模式
      refs.isNavigating.current = true;
      refs.setHistoryIndex(-1);
      refs.draft.current = currentValue;
      refs.currentHistoryValue.current = '';
      return true;
    }

    const nextIndex = currentIndex - 1;
    refs.isNavigating.current = true;
    if (nextIndex === -1) {
      // 回到草稿
      refs.setHistoryIndex(-1);
      refs.onChange(refs.draft.current);
      refs.currentHistoryValue.current = '';
    } else {
      refs.setHistoryIndex(nextIndex);
      const historyValue = String(inputHistoryList[nextIndex] || '');
      refs.onChange(historyValue);
      refs.currentHistoryValue.current = historyValue;
    }

    setTimeout(() => {
      const editorView = refs.editorRef.current?.$view;
      if (editorView) {
        editorView.dispatch({ selection: { anchor: editorView.state.doc.length } });
      }
    }, 0);
    return true;
  };

  return [
    {
      key: 'ArrowUp',
      run: (view: EditorView): boolean => {
        const { state } = view;
        const line = state.doc.lineAt(state.selection.main.head);
        if (line.number === 1) {
          return navigateUp(view);
        }
        return false;
      }
    },
    {
      key: 'ArrowDown',
      run: (view: EditorView): boolean => {
        const { state } = view;
        const line = state.doc.lineAt(state.selection.main.head);
        if (line.number === state.doc.lines) {
          return navigateDown(view);
        }
        return false;
      }
    }
  ];
};
