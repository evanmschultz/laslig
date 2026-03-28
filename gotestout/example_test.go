package gotestout_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/evanmschultz/laslig"
	"github.com/evanmschultz/laslig/gotestout"
)

// ExampleRender shows one small plain-text gotestout render.
func ExampleRender() {
	stream := strings.NewReader(`{"Action":"pass","Package":"example/pkg","Test":"TestRender","Elapsed":0.01}
{"Action":"pass","Package":"example/pkg","Elapsed":0.01}
`)

	summary, _ := gotestout.Render(os.Stdout, stream, gotestout.Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
	})

	fmt.Printf("tests=%d packages=%d failures=%t\n", summary.TotalTests(), summary.TotalPackages(), summary.HasFailures())

	// Output:
	// [PKG PASS] example/pkg (0.01s)
	//
	// Test summary
	//   tests: 1
	//   passed: 1
	//   failed: 0
	//   skipped: 0
	//   packages: 1
	//   pkg passed: 1
	//   pkg failed: 0
	//   pkg skipped: 0
	//
	// [SUCCESS] All tests passed
	//   1 test passed across 1 package.
	// tests=1 packages=1 failures=false
}
