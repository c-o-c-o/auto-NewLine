package newline

import (
	"auto-NewLine/data"
	"math"
	"strings"
)

func Break(infos []Info, stg data.Setting, min int, aim float64, max int) string {
	breaked := [][]Info{}

	for len(infos) > 0 && getWordLength(infos) > max {
		space, residue := splitSliceWordCount(infos, max, false)
		line, space := splitSliceWordCount(space, min, true)

		brkvals := getBreakValues(space, stg.PriorityWait, line[0].Start, aim)
		brkidx := getMaxIndex(brkvals)

		line = append(line, space[:brkidx+1]...)
		infos = append(space[brkidx+1:], residue...)

		breaked = append(breaked, line)
	}

	breaked = append(breaked, infos)
	return merge(breaked, stg)
}

func getWordLength(infos []Info) int {
	return infos[len(infos)-1].End - infos[0].Start
}

func splitSliceWordCount(slice []Info, count int, isover bool) ([]Info, []Info) {
	offset := slice[0].Start
	splidx := 0

	for i, v := range slice {
		if v.End-offset >= count {
			splidx = i
			if isover && v.End-offset != count {
				splidx += 1
			}

			break
		} else { //最小分割ワードの長さがcountを超えていた場合
			splidx = 1
		}
	}

	return slice[:splidx], slice[splidx:]
}

func getBreakValues(space []Info, wait float64, offset int, aim float64) []float64 {
	rslt := []float64{}

	for _, v := range space {
		distance := math.Abs(aim - float64(v.End-offset))
		rslt = append(rslt, v.Priority*wait-distance)
	}

	return rslt
}

func getMaxIndex(slice []float64) int {
	idx := 0
	max := 0.0

	for i, v := range slice {
		if max < v {
			max = v
			idx = i
		}
	}

	return idx
}

func merge(infos [][]Info, stg data.Setting) string {
	lines := []string{}

	for _, infos := range infos {
		line := ""

		for _, v := range infos {
			line += v.Word
		}

		lines = append(lines, deleteBothEnds(line, stg.DeleteFixs))
	}

	return strings.Join(lines, stg.BreakStr)
}

func deleteBothEnds(str string, delstrs []string) string {
	sidx := 0
	eidx := len(str)

	for _, delstr := range delstrs {
		if strings.HasPrefix(str, delstr) {
			sidx = len(delstr)
			break
		}
	}

	for _, delstr := range delstrs {
		if strings.HasSuffix(str, delstr) {
			eidx = len(str) - len(delstr)
			break
		}
	}

	return str[sidx:eidx]
}
