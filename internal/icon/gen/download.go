package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {

	svgUrl := "https://fonts.gstatic.com/s/i/short-term/release/materialsymbolsrounded/%s/wght200fill1/%s.svg"

	icons := []string{
		"skip_previous", "skip_next", "play_arrow", "pause",
		"expand_more", "expand_less", "search", "close",
		"playlist_add", "navigate_before", "navigate_next",
		"volume_up", "volume_down", "volume_mute",
	}
	sizes := []string{"48px"}

	jf := jen.NewFilePath("github.com/mlctrez/goapp-audioplayer/internal/icon")

	for _, icon := range icons {
		for _, size := range sizes {
			svgPath := fmt.Sprintf("internal/icon/svg/%s_%s.svg", icon, size)
			if _, err := os.Stat(svgPath); os.IsNotExist(err) {
				fmt.Println("downloading", svgPath)

				url := fmt.Sprintf(svgUrl, icon, size)

				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}
				create, err := os.Create(svgPath)
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(create, resp.Body)
				if err != nil {
					panic(err)
				}
				_ = resp.Body.Close()
			}
			svgBytes, err := os.ReadFile(svgPath)
			if err != nil {
				panic(err)

			}

			jf.Line()
			jf.Func().Id(IconName(icon, size)).Params().Params(jen.String()).Block(
				jen.Return(jen.Lit(string(svgBytes))),
			)
		}
	}

	create, err := os.Create("internal/icon/icons.go")
	if err != nil {
		panic(err)
	}
	err = jf.Render(create)
	if err != nil {
		panic(err)
	}
	_ = create.Close()
}

func IconName(icon, size string) string {
	var result string
	for _, s := range strings.Split(icon, "_") {
		result += Capitalize(s)
	}
	result += strings.TrimSuffix(size, "px")
	return result

}

func Capitalize(in string) string {
	if len(in) > 0 {
		return strings.ToUpper(in[0:1]) + in[1:]
	}
	return in
}
