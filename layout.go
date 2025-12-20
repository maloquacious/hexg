// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

//// BottomRightHex returns (0,0,0) if the list of hexes is empty
//func BottomRightHex(l Screen_i, hexes ...Hex) Hex {
//	var maxHex Hex
//	var maxOffset OffsetCoord
//
//	for n, h := range hexes {
//		if n == 0 {
//			maxHex, maxOffset = h, l.HexToOffsetCoord(h)
//			continue
//		}
//		oc := l.HexToOffsetCoord(h)
//		if oc.Row > maxOffset.Row || (oc.Row == maxOffset.Row && oc.Col > maxOffset.Col) {
//			maxHex, maxOffset = h, l.HexToOffsetCoord(h)
//		}
//	}
//
//	return maxHex
//}
//
//// TopLeftHex returns (0,0,0) if the list of hexes is empty
//func TopLeftHex(l Screen_i, hexes ...Hex) Hex {
//	var minHex Hex
//	var minOffset OffsetCoord
//
//	for n, h := range hexes {
//		if n == 0 {
//			minHex, minOffset = h, l.HexToOffsetCoord(h)
//			continue
//		}
//		oc := l.HexToOffsetCoord(h)
//		if oc.Row < minOffset.Row || (oc.Row == minOffset.Row && oc.Col < minOffset.Col) {
//			minHex, minOffset = h, oc
//		}
//	}
//
//	return minHex
//}
