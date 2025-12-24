import { EditorView } from '@codemirror/view';

type Extension = any;

/**
 * 编辑器主题配置
 */
export const editorTheme: Extension = EditorView.theme({
  '&': {},
  '.cm-scroller': { minHeight: '50px', maxHeight: '200px' },
  '&.cm-focused': { outline: 'none' },
  '.cm-line': { lineHeight: '30px', padding: '0 10px' },
  '.cm-selectionBackground': { backgroundColor: '#c7d2fe !important' },
  '&.cm-focused .cm-selectionBackground': { backgroundColor: '#c7d2fe !important' },
  '.cm-content': { caretColor: '#4338ca' },
  '& .cm-line ::selection': { backgroundColor: '#c7d2fe', color: '#000' },
  '& .cm-content ::selection': { backgroundColor: '#c7d2fe', color: '#000' }
});
