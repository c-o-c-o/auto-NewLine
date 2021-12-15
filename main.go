package main

import (
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	app := &cli.App{
		Name:            "auto-NewLine",
		Usage:           Version,
		Description:     "指定したテキストファイルを最大文字数に合わせていい感じに改行します",
		Version:         Version,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "text",
				Aliases:  []string{"t"},
				Usage:    "改行するテキストファイルのパス",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "maxlen",
				Aliases: []string{"l"},
				Value:   30,
				Usage:   "最大文字数",
			},
			&cli.StringFlag{
				Name:    "encode",
				Aliases: []string{"e"},
				Value:   "shift-jis",
				Usage:   "エンコード",
			},
		},
		Action: appfunc,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appfunc(c *cli.Context) error {
	maxlen := c.Int("maxlen")
	textp := c.Path("text")
	enc := map[string]encoding.Encoding{
		"utf-8":     nil,
		"shift-jis": japanese.ShiftJIS,
	}

	text, err := os.ReadFile(textp)
	if err != nil {
		return err
	}

	e := enc[c.String("encode")]
	if e != nil {
		text, _, err = transform.Bytes(e.NewDecoder(), text)
		if err != nil {
			return err
		}
	}

	if len([]rune(string(text))) < maxlen {
		return nil
	}

	anlys, err := AnalyzeLanguage(string(text))
	if err != nil {
		return err
	}
	minstrs := JoinSimilar(anlys)
	nlstrs := JoinStrings(minstrs, maxlen)
	wbytes := []byte(strings.Join(nlstrs, "\r\n"))

	if e != nil {
		wbytes, _, err = transform.Bytes(e.NewEncoder(), wbytes)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(textp, wbytes, 0777)
}

func JoinStrings(minstrs []string, maxlen int) []string {
	lines := []string{}
	line := ""

	for _, str := range minstrs {
		if len([]rune(line))+len([]rune(str)) <= maxlen {
			line += str
		} else {
			lines = append(lines, line)
			line = str
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func JoinSimilar(strs []string) []string {
	rslt := []string{}

	for i := 0; i < len(strs); i++ {
		if i >= len(strs)-1 {
			rslt = append(rslt, strs[i])
			break
		}

		now := []rune(strs[i])
		next := []rune(strs[i+1])

		if IsSimilar(now, next, unicode.Han) {
			rslt = append(rslt, strs[i]+strs[i+1])
			i++
			continue
		}

		if IsSimilar(now, next, unicode.Katakana) {
			rslt = append(rslt, strs[i]+strs[i+1])
			i++
			continue
		}

		if IsSimilar(now, next, unicode.Digit) {
			rslt = append(rslt, strs[i]+strs[i+1])
			i++
			continue
		}

		rslt = append(rslt, strs[i])
	}

	return rslt
}

func IsSimilar(l []rune, r []rune, rt *unicode.RangeTable) bool {
	return unicode.In(l[len(l)-1], rt) && unicode.In(r[0], rt)
}

func AnalyzeLanguage(text string) ([]string, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	tokens := t.Tokenize(text)
	anlys := []string{}

	for _, token := range tokens {
		f := token.Features()
		if len(anlys) != 0 && (Included(f, "接", "助", "記号", "非自立", "サ変")) {
			anlys[len(anlys)-1] += token.Surface
		} else {
			anlys = append(anlys, token.Surface)
		}
	}
	return anlys, nil
}

func Included(strs []string, substrs ...string) bool {
	for _, s := range strs {
		for _, ss := range substrs {
			if strings.Contains(s, ss) {
				return true
			}
		}
	}
	return false
}
