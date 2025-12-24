import { EditorView, WidgetType } from '@codemirror/view';

import { findCommandItem, findCommandItemByName } from '../../Commands';
import { DEFAULT_CAPSULE_ICON_SVG, DISABLED_CAPSULE_ICON_SVG, CAPSULE_BOUNDARY } from '../utils';
import type { CapsuleGlobalState } from '../types';
import styles from '../index.module.less';

/**
 * 将 @mention 渲染为胶囊.
 */
export class MentionCapsuleWidget extends WidgetType {
  private name: string;
  private globalState: CapsuleGlobalState;

  private deleteMousedownHandler: ((e: MouseEvent) => void) | null = null;
  private deleteClickHandler: ((e: MouseEvent) => void) | null = null;
  private deleteBtnElement: HTMLElement | null = null;

  constructor(name: string, globalState: CapsuleGlobalState) {
    super();
    this.name = name;
    this.globalState = globalState;
  }

  toDOM(view: EditorView): HTMLElement {
    const commands = this.globalState.commandsGetter?.() || [];
    const commandItem = findCommandItemByName(commands, this.name) || findCommandItem(commands, this.name);

    const isDisabled = !commandItem;

    const displayName = commandItem?.name || this.name;
    const typeLabel = commandItem?.typeLabel || '';
    const capsuleIconSvg = isDisabled
      ? DISABLED_CAPSULE_ICON_SVG
      : commandItem?.capsuleIconSvg || DEFAULT_CAPSULE_ICON_SVG;
    const typeLabelStyle = commandItem?.typeLabelStyle;

    const capsule = document.createElement('span');
    capsule.className = isDisabled
      ? `${styles.mentionCapsule} ${styles.mentionCapsuleDisabled}`
      : styles.mentionCapsule;
    capsule.setAttribute('data-mention', this.name);

    if (isDisabled) {
      capsule.setAttribute('title', `${displayName}（已失效）`);
    } else {
      capsule.setAttribute('title', typeLabel ? `${typeLabel}: ${displayName}` : displayName);
    }

    if (!isDisabled && typeLabelStyle) {
      if (typeLabelStyle.background) capsule.style.background = typeLabelStyle.background as string;
      if (typeLabelStyle.color) capsule.style.color = typeLabelStyle.color as string;
      if (typeLabelStyle.borderColor) capsule.style.borderColor = typeLabelStyle.borderColor as string;
    }

    // 阻止点击胶囊时改变光标位置（删除按钮除外）
    capsule.onmousedown = (e) => {
      const target = e.target as HTMLElement;
      if (target.closest(`.${styles.capsuleDelete}`)) {
        e.stopPropagation();
        return;
      }
      e.preventDefault();
    };

    capsule.onclick = (e) => {
      if ((e.target as HTMLElement).closest(`.${styles.capsuleDelete}`)) return;
      e.preventDefault();
      e.stopPropagation();

      if (isDisabled) return;

      if (commandItem?.onClickCapsule) {
        commandItem.onClickCapsule(commandItem);
      } else {
        this.globalState.capsuleClickHandler?.(this.name, commandItem);
      }
    };

    if (capsuleIconSvg) {
      const icon = document.createElement('span');
      icon.className = styles.capsuleIcon;
      icon.innerHTML = capsuleIconSvg;
      capsule.appendChild(icon);
    }

    const nameSpan = document.createElement('span');
    nameSpan.className = styles.capsuleName;
    nameSpan.textContent = displayName;
    capsule.appendChild(nameSpan);

    const deleteBtn = document.createElement('span');
    deleteBtn.className = styles.capsuleDelete;
    deleteBtn.innerHTML =
      '<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>';

    this.deleteMousedownHandler = (e: MouseEvent) => {
      e.stopPropagation();
      e.stopImmediatePropagation();
      e.preventDefault();
    };

    this.deleteClickHandler = (e: MouseEvent) => {
      e.stopPropagation();
      e.stopImmediatePropagation();
      e.preventDefault();

      const pos = view.posAtDOM(capsule);
      if (pos === null || pos === undefined) return;

      const capsuleLength = 1 + this.name.length + CAPSULE_BOUNDARY.length;
      view.dispatch({
        changes: { from: pos, to: pos + capsuleLength, insert: '' }
      });
    };

    deleteBtn.addEventListener('mousedown', this.deleteMousedownHandler, true);
    deleteBtn.addEventListener('click', this.deleteClickHandler, true);

    deleteBtn.style.pointerEvents = 'auto';
    deleteBtn.style.position = 'relative';
    deleteBtn.style.zIndex = '10';

    this.deleteBtnElement = deleteBtn;

    capsule.appendChild(deleteBtn);

    return capsule;
  }

  destroy(): void {
    if (this.deleteBtnElement) {
      if (this.deleteMousedownHandler) {
        this.deleteBtnElement.removeEventListener('mousedown', this.deleteMousedownHandler, true);
      }
      if (this.deleteClickHandler) {
        this.deleteBtnElement.removeEventListener('click', this.deleteClickHandler, true);
      }
    }
    this.deleteBtnElement = null;
    this.deleteMousedownHandler = null;
    this.deleteClickHandler = null;
  }

  eq(other: MentionCapsuleWidget): boolean {
    return this.name === other.name;
  }

  ignoreEvent(event: Event): boolean {
    return event.type === 'mousedown' || event.type === 'click';
  }
}
