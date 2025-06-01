package tableconfig

type filterModel struct {
	NodeStr rune                  //内容
	Subli   map[rune]*filterModel //屏蔽子集合
	IsEnd   bool                  //是否为结束
}

type configFilter struct {
	FilterList map[rune]*filterModel //屏蔽字树
}

func filterFor(li map[rune]*filterModel, rowr []rune, index int) bool {
	if len(rowr) <= index {
		return true
	}
	fmd, ok := li[rowr[index]]
	if !ok {
		fmd = new(filterModel)
		fmd.NodeStr = rowr[index]
		fmd.Subli = make(map[rune]*filterModel)
		li[rowr[index]] = fmd
	}
	index++
	fmd.IsEnd = filterFor(fmd.Subli, rowr, index)
	return false
}

type SenWordConfig struct {
	ID      int32  `json:"id"`
	SenWord string `json:"sensitivewords"`
}
type SenWordConfigCol struct {
	SenWordConfigList []SenWordConfig
	// configFilter      *configFilter
	SenWordConfigMap map[string]bool
}

func (s *SenWordConfigCol) InitMap() {
	// s.configFilter = &configFilter{}
	s.SenWordConfigMap = make(map[string]bool)

	// li := make(map[rune]*filterModel)
	for _, row := range s.SenWordConfigList {
		if row.SenWord == "" {
			continue
		}
		// rowr := []rune(row.SenWord)
		// fmd, ok := li[rowr[0]]
		// if !ok {
		// 	fmd = new(filterModel)
		// 	fmd.NodeStr = rowr[0]
		// 	fmd.Subli = make(map[rune]*filterModel)
		// 	li[rowr[0]] = fmd
		// }
		// fmd.IsEnd = filterFor(fmd.Subli, rowr, 1)

		s.SenWordConfigMap[row.SenWord] = true
	}
	// s.configFilter.FilterList = li

	return
}

func (s *SenWordConfigCol) ReplaceSenWord(data string) string {
	// filterli := s.configFilter.FilterList
	arr := []rune(data)
	indexList := make([][2]int, 0)
	for i := 0; i < len(arr); i++ {
		for j := i; j < len(arr); j++ {
			word := ""
			if i == j {
				word = string(arr[i : i+1])
			} else {
				if i == len(arr)-1 {
					word = string(arr[i:])
				} else {
					word = string(arr[i : j+1])
				}
			}
			if _, ok := s.SenWordConfigMap[word]; ok {
				indexList = append(indexList, [2]int{i, j})
				// for n := i; n <= j; n++ {
				// 	arr[n] = rune('*')
				// }
			}
		}
		// fmd, ok := filterli[arr[i]]
		// if !ok {
		// 	continue
		// }
		// if ok, index := FilterCheckFor(arr, i+1, fmd.Subli); ok {
		// 	arr[i] = rune('*')
		// 	i = index
		// }
	}

	for i := 0; i < len(indexList); i++ {
		for n := indexList[i][0]; n <= indexList[i][1]; n++ {
			arr[n] = rune('*')
		}
	}
	return string(arr)
}

//屏蔽字操作，这个方法就是外部调用的入口方法
func (s *SenWordConfigCol) CheckSenWord(data string) bool {
	arr := []rune(data)
	for i := 0; i < len(arr); i++ {
		for j := i; j < len(arr); j++ {
			word := ""
			if i == j {
				word = string(arr[i : i+1])
			} else {
				if i == len(arr)-1 {
					word = string(arr[i:])
				} else {
					word = string(arr[i : j+1])
				}
			}
			if _, ok := s.SenWordConfigMap[word]; ok {
				return true
			}
		}
	}
	return false
}

//递归调用检查屏蔽字
func filterCheckFor(arr []rune, index int, filterli map[rune]*filterModel) (bool, int) {
	if len(arr) <= index {
		return false, index
	}
	if arr[index] == rune(' ') {
		if ok, i := filterCheckFor(arr, index+1, filterli); ok {
			arr[index] = rune('*')
			return true, i
		}
	}
	fmd, ok := filterli[arr[index]]
	if !ok {
		return false, index
	}
	if fmd.IsEnd {
		arr[index] = rune('*')
		ok, i := filterCheckFor(arr, index+1, fmd.Subli)
		if ok {
			return true, i
		}
		return true, index
	} else if ok, i := filterCheckFor(arr, index+1, fmd.Subli); ok {
		arr[index] = rune('*')
		return true, i
	}
	return false, index
}

var SenWordConfigs *SenWordConfigCol
