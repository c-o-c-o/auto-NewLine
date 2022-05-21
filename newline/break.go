package newline

import (
	"auto-NewLine/data"
	"errors"
	"math"
	"strconv"
	"strings"
)

func Break(infos []Info, stg data.Setting, min int, aim float64, max int) (string, error) {
	breaked := [][]Info{}

	for len(infos) > 0 && getWordLength(infos) > max {
		space, residue := splitSliceWordCount(infos, max, false)

		if len(space) == 0 { //エラー処理
			maxlen := 0
			for _, r := range residue {
				maxlen = int(math.Max(float64(maxlen), float64(r.End-r.Start)))
			}
			return "", errors.New("words longer than -max value found. make that value larger than " + strconv.Itoa(maxlen))
		}

		line, space := splitSliceWordCount(space, min, true)

		if len(space) != 0 {
			offset := 0
			if len(line) > 0 {
				offset = line[0].Start
			}

			brkvals := getBreakValues(space, stg.PriorityWait, offset, aim)
			brkidx := getMaxIndex(brkvals)

			line = append(line, space[:brkidx+1]...)
			infos = append(space[brkidx+1:], residue...)
		} else {
			infos = residue
		}

		breaked = append(breaked, line)
	}

	breaked = append(breaked, infos)
	return merge(breaked, stg), nil
}

func getWordLength(infos []Info) int {
	return infos[len(infos)-1].End - infos[0].Start
}

func splitSliceWordCount(slice []Info, count int, isover bool) ([]Info, []Info) {
	offset := slice[0].Start
	splidx := 1

	for i, v := range slice {
		if v.End-offset > count {
			splidx = i
			if isover && v.End-offset != count {
				splidx += 1
			}

			break
		}
	}

	if splidx >= len(slice) {
		return slice, []Info{}
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
	max := -math.MaxFloat64

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
