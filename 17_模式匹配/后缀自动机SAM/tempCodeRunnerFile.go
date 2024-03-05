char]; ok {
	// 	lastNode, nextNode := sam.Nodes[lastPos], sam.Nodes[tmp]
	// 	if lastNode.MaxLen+1 == nextNode.MaxLen {
	// 		return tmp
	// 	} else {
	// 		newQ := int32(len(sam.Nodes))
	// 		sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1))
	// 		for k, v := range nextNode.Next {
	// 			sam.Nodes[newQ].Next[k] = v
	// 		}
	// 		sam.Nodes[tmp].Link = newQ
	// 		for lastPos != -1 && sam.Nodes[lastPos].Next[char] == tmp {
	// 			sam.Nodes[lastPos].Next[char] = newQ
	// 			lastPos = sam.Nodes[lastPos].Link
	// 		}
	// 		return newQ
	// 	}
	// }