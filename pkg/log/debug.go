package log

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-isatty"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	"os"
	"reflect"
	"runtime"
	"strings"
	"text/tabwriter"
)

type FmtColor string

const (
	CBlack   FmtColor = "\u001b[90m"
	CRed     FmtColor = "\u001b[91m"
	CCyan    FmtColor = "\u001b[96m"
	CGreen   FmtColor = "\u001b[92m"
	CYellow  FmtColor = "\u001b[93m"
	CBlue    FmtColor = "\u001b[94m"
	CMagenta FmtColor = "\u001b[95m"
	CWhite   FmtColor = "\u001b[97m"
	CReset   FmtColor = "\u001b[0m"
)

func FiberRequestDebug(ctx *fiber.Ctx) error {
	handlerName := ""
	currentRoute := ctx.Route()
	currentPath := string(ctx.Request().URI().Path())

	for _, routes := range ctx.App().Stack() {
		for _, r := range routes {

			if r.Path == currentPath {
				for _, handler := range r.Handlers {
					hName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

					splits := strings.Split(hName, "/")
					if len(splits) > 0 {
						lastHandler := splits[len(splits)-1]

						pkgSplits := strings.Split(lastHandler, ".")

						if len(pkgSplits) > 0 {
							handlerName = pkgSplits[len(pkgSplits)-1]
						}
					}
				}
			}
		}
	}

	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") == "1" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}

	w := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)

	if handlerName != " - - ?? - - " {
		_, _ = fmt.Fprintf(w, "%s[%s]\t%s | %s%s\t%s | %s%s\t%s | %s%s%s\n", CBlue, currentRoute.Method, CWhite, CGreen, currentPath, CWhite, CYellow, handlerName, CWhite, CCyan, currentRoute.Name, CReset)
	} else {
		_, _ = fmt.Fprintf(w, "%s[%s]\t%s | %s%s\t%s | %s%s\t%s | %s%s%s\n", CBlue, currentRoute.Method, CWhite, CGreen, currentPath, CWhite, CRed, " -- NO HANDLER FOUND --", CWhite, CCyan, currentRoute.Name, CReset)
	}

	_ = w.Flush()

	return ctx.Next()
}

func coloredPrint(txt string, color FmtColor) {
	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") == "1" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}

	w := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)
	_, _ = fmt.Fprintf(w, "%s%s%s\n", color, txt, CReset)

	_ = w.Flush()
}

func PrintlnRed(txt string) {
	coloredPrint(txt, CRed)
}

func PrintlnGreen(txt string) {
	coloredPrint(txt, CGreen)
}
