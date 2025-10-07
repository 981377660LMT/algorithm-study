import { TLKeyboardEventInfo } from '../types/event-types'
import { StateNode } from './StateNode'

export class RootState extends StateNode {
	static override id = 'root'
	static override initial = ''
	static override children() {
		return []
	}

	override onKeyDown(info: TLKeyboardEventInfo) {
		// todo: move this logic up to the tldraw library, as the "zoom" tool only exists there
		switch (info.code) {
			case 'KeyZ': {
				if (!(info.shiftKey || info.ctrlKey)) {
					const currentTool = this.getCurrent()
					if (currentTool && currentTool.getCurrent()?.id === 'idle' && this.children!['zoom']) {
						this.editor.setCurrentTool('zoom', { ...info, onInteractionEnd: currentTool.id })
					}
				}
				break
			}
		}
	}
}
