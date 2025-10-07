import { StateNode } from 'tldraw'

// There's a guide at the bottom of this file!

export class ScreenshotIdle extends StateNode {
	static override id = 'idle'

	// [1]
	override onPointerDown() {
		this.parent.transition('pointing')
	}
}

/*
[1]
When we the user makes a pointer down event, we transition to the pointing state.
*/
