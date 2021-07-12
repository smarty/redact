package redact

//type FullDOB struct { //FIXME
//	dobRedaction *dobRedaction
//}
//
//func (this *FullDOB) findMatch(input []byte) {
//
//}
//func(this *dobRedaction) isMonth(input byte){
//
//}
//
//func isMonth(first, last byte, length int) bool {
//	candidates, found := months[first]
//	if !found {
//		return false
//	}
//	candidate, found := candidates[last]
//	if !found {
//		return false
//	}
//	for _, number := range candidate {
//		if number == length {
//			return true
//		}
//	}
//	return false
//}
//
//var (
//	months = map[byte]map[byte][]int{
//		'J': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}, 'N': []int{3}, 'Y': []int{7, 4}, 'E': []int{4}, 'L': []int{3}},
//		'F': {'b': []int{3}, 'y': []int{8}, 'B': []int{3}, 'Y': []int{8}},
//		'M': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}, 'H': []int{5}, 'R': []int{3}, 'Y': []int{3}},
//		'A': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}, 'G': []int{3}, 'T': []int{6}, 'L': []int{5}, 'R': []int{3}},
//		'S': {'r': []int{9}, 'p': []int{3}, 'R': []int{9}, 'P': []int{3}},
//		'O': {'t': []int{3}, 'r': []int{7}, 'T': []int{3}, 'R': []int{7}},
//		'N': {'v': []int{3}, 'r': []int{9}, 'V': []int{3}, 'R': []int{9}},
//		'D': {'r': []int{8}, 'c': []int{3}, 'R': []int{8}, 'C': []int{3}},
//	}
//	validFirst = map[byte]bool{
//		'J': true,
//		'F': true,
//		'M': true,
//		'A': true,
//		'S': true,
//		'O': true,
//		'N': true,
//		'D': true,
//	}
//)