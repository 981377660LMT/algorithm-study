import {
  Decoration,
  DecorationSet,
  EditorView,
  MatchDecorator,
  ViewPlugin,
  ViewUpdate
} from '@codemirror/view';
import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget';
import type { CapsuleGlobalState } from '../types';
import { CAPSULE_REGEX } from '../utils';

export const createCapsulePlugin = (globalState: CapsuleGlobalState) => {
  const matcher = new MatchDecorator({
    regexp: CAPSULE_REGEX,
    decoration: (match, view, pos) => {
      const name = match[1];
      if (!name || name.length === 0) return null;
      return Decoration.replace({
        widget: new MentionCapsuleWidget(name, globalState),
        inclusive: false
      });
    }
  });

  return ViewPlugin.fromClass(
    class CapsuleViewPlugin {
      decorations: DecorationSet;

      constructor(view: EditorView) {
        this.decorations = matcher.createDeco(view);
      }

      update(update: ViewUpdate) {
        alert(1);

        this.decorations = matcher.updateDeco(update, this.decorations);
      }
    },
    {
      decorations: (plugin) => plugin.decorations,
      provide(plugin) {
        return EditorView.atomicRanges.of((view) => view.plugin(plugin)?.decorations ?? Decoration.none);
      }
    }
  );
};
