package newline

import (
	"auto-NewLine/data"
	"math"
	"regexp"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type WordInfo struct {
	Priority float64
	Word     string
	Start    int
	End      int
	Features []string
}

func GetInfo(text string, bptns []data.BreakPattern, userDictPath string) ([]WordInfo, error) {
	udic, err := dict.NewUserDict(userDictPath)
	if err != nil {
		return []WordInfo{}, err
	}

	ter, err := tokenizer.New(ipa.DictShrink(), tokenizer.OmitBosEos(), tokenizer.UserDict(udic))
	if err != nil {
		return []WordInfo{}, err
	}

	return Analyze(ter.Tokenize(text), bptns)
}

func Analyze(tokens []tokenizer.Token, bptns []data.BreakPattern) ([]WordInfo, error) {
	result := []WordInfo{}
	var ltoken *tokenizer.Token = nil
	var rtoken *tokenizer.Token = nil

	for i := 0; i < len(tokens); i++ {
		ltoken = rtoken
		rtoken = &tokens[i]
		if ltoken == nil {
			continue
		}

		p, err := getPriority(ltoken, rtoken, bptns)
		if err != nil {
			return []WordInfo{}, err
		}

		result = append(
			result,
			WordInfo{
				Word:     ltoken.Surface,
				Priority: p,
				Start:    ltoken.Start,
				End:      ltoken.End,
				Features: ltoken.Features(),
			})
	}
	last := tokens[len(tokens)-1]
	result = append(
		result,
		WordInfo{
			Word:     last.Surface,
			Priority: math.MaxInt32,
			Start:    last.Start,
			End:      last.End,
			Features: last.Features(),
		})

	return result, nil
}

func getPriority(ltoken *tokenizer.Token, rtoken *tokenizer.Token, bptns []data.BreakPattern) (float64, error) {
	def := 0.0

	for _, bptn := range bptns {
		if len(bptn.Patterns) == 0 {
			def = bptn.Priority
			continue
		}

		valid, err := checkPriority(bptn, ltoken, rtoken)
		if err != nil {
			return 0, err
		}

		if valid {
			return bptn.Priority, nil
		}
	}

	return def, nil
}

func checkPriority(bptn data.BreakPattern, ltoken *tokenizer.Token, rtoken *tokenizer.Token) (bool, error) {
	for _, ptn := range bptn.Patterns {
		if len(ptn) == 0 {
			break
		}

		lptn, rptn := getKeyValue(ptn)

		lok, err := match('$', ltoken, lptn)
		if err != nil {
			return false, err
		}

		rok, err := match('$', rtoken, rptn)
		if err != nil {
			return false, err
		}

		if lok && rok {
			return true, nil
		}
	}

	return false, nil
}

func match(prefix rune, token *tokenizer.Token, ptn string) (bool, error) {
	if []rune(ptn)[0] == prefix {
		return regexp.MatchString(string([]rune(ptn)[1:]), token.Surface)
	}

	for _, f := range token.Features() {
		matched, err := regexp.MatchString(ptn, f)

		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func getKeyValue(m map[string]string) (key string, value string) {
	for k, v := range m {
		return k, v
	}

	return "", ""
}
