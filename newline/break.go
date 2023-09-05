package newline

import (
	"auto-NewLine/data"
	"errors"
	"math"
	"strconv"
	"strings"
)

func Break(infos []WordInfo, stg data.Setting, min int, aim float64, max int) (string, error) {
	breaked := [][]WordInfo{}
	allline := infos

	for len(allline) > 0 && getWordLength(allline) > max {
		maxline, residue := splitSliceWordCount(allline, max, false)

		if len(maxline) == 0 { //エラー処理
			maxlen := 0
			for _, r := range residue {
				maxlen = int(math.Max(float64(maxlen), float64(r.End-r.Start)))
			}
			return "", errors.New("words longer than -max value found. make that value larger than " + strconv.Itoa(maxlen))
		}

		line, space := splitSliceWordCount(maxline, min, true)
		line, space = moveItemLeftToReght(line, space) // 単語の後の改行位置を探るので、最小行から一単語を候補へ移動

		offset := 0
		if len(line) > 0 {
			offset = line[0].Start
		}

		brkvals := getBreakValues(space, stg.PriorityWait, offset, aim)
		brkidx := getMaxIndex(brkvals)

		allline = append(space[brkidx+1:], residue...)
		breaked = append(breaked, append(line, space[:brkidx+1]...))
	}

	breaked = append(breaked, allline)
	return merge(breaked, stg), nil
}

func moveItemLeftToReght(from, to []WordInfo) ([]WordInfo, []WordInfo) {
	return from[:len(from)-1], append([]WordInfo{from[len(from)-1]}, to...)
}

func getWordLength(infos []WordInfo) int {
	return infos[len(infos)-1].End - infos[0].Start
}

func splitSliceWordCount(slice []WordInfo, count int, isover bool) ([]WordInfo, []WordInfo) {
	offset := slice[0].Start
	splidx := 0

	for i, v := range slice {
		if v.End-offset >= count {
			if isover || v.End-offset == count {
				splidx = i
			}

			break
		}

		splidx = i
	}

	return slice[:splidx+1], slice[splidx+1:]
}

func getBreakValues(space []WordInfo, wait float64, offset int, aim float64) []float64 {
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

func merge(candidates [][]WordInfo, stg data.Setting) string {
	lines := []string{}

	for _, brks := range candidates {
		line := ""

		for _, v := range brks {
			line += v.Word
		}

		if len(line) > 0 {
			lines = append(lines, deleteBothEnds(line, stg.DeleteFixs))
		}
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

	if sidx == len(str) {
		return str[sidx:]
	}

	for _, delstr := range delstrs {
		if strings.HasSuffix(str, delstr) {
			eidx = len(str) - len(delstr)
			break
		}
	}

	return str[sidx:eidx]
}
