import type { ICapsuleInfo } from '@/components/ChatProvider';
import type { ICommandItem } from '../Commands';
import { findCommandItem, findCommandItemByName } from '../Commands';

export const CAPSULE_BOUNDARY = '\u200B';

const MENTION_NAME_PATTERN = `[^@\\s${CAPSULE_BOUNDARY}]+`;

/** 胶囊匹配正则：@name 后跟零宽空格 */
export const CAPSULE_REGEX = new RegExp(`@(${MENTION_NAME_PATTERN})${CAPSULE_BOUNDARY}`, 'g');

/** 自动补全胶囊正则：@name 后跟换行或行尾（用于粘贴时自动补充零宽空格） */
export const AUTO_COMPLETE_CAPSULE_REGEX = new RegExp(`@(${MENTION_NAME_PATTERN})(?=\\n|$)`, 'g');

export const extractMentionQuery = (
  text: string,
  cursorPos: number
): { query: string; atIndex: number } | null => {
  let atIndex = -1;
  for (let i = cursorPos - 1; i >= 0; i--) {
    if (text[i] === '@') {
      atIndex = i;
      break;
    }
    if (text[i] === ' ' || text[i] === '\n' || text[i] === CAPSULE_BOUNDARY) {
      break;
    }
  }
  if (atIndex === -1) return null;
  return { query: text.slice(atIndex + 1, cursorPos), atIndex };
};

export const matchCapsules = (
  text: string,
  callback: (name: string, from: number, to: number) => boolean | undefined | void
) => {
  const regex = new RegExp(CAPSULE_REGEX.source, 'g');
  let match;
  while ((match = regex.exec(text)) !== null) {
    const name = match[1];
    if (!name || name.length === 0) continue;
    const from = match.index;
    const to = match.index + match[0].length;
    if (callback(name, from, to) === false) {
      break;
    }
  }
};

export const extractAllCapsules = (text: string, commands: ICommandItem[]): ICapsuleInfo[] => {
  const capsules: ICapsuleInfo[] = [];
  matchCapsules(text, (name) => {
    const item = findCommandItemByName(commands, name) || findCommandItem(commands, name);
    capsules.push({ id: item?.id || name, item });
  });
  return capsules;
};

export const DEFAULT_CAPSULE_ICON_SVG =
  '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/></svg>';

export const DISABLED_CAPSULE_ICON_SVG =
  '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>';
