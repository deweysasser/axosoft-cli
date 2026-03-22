package program

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestOptions_Version(t *testing.T) {
	var program Options

	exitValue := -1
	exitFunc = func(x int) {
		exitValue = x
	}
	defer func() { exitFunc = os.Exit }()

	out := capturer.CaptureStdout(func() {
		_, err := program.Parse([]string{"--version", "serve"})
		assert.NoError(t, err)
	})

	assert.Equal(t, 0, exitValue)
	assert.Equal(t, "unknown\n", out)
}

func TestOptions_ParseServe(t *testing.T) {
	var program Options

	ctx, err := program.Parse([]string{"serve"})
	assert.NoError(t, err)
	assert.Equal(t, "serve", ctx.Command())
}

func TestOptions_ParseLogin(t *testing.T) {
	var program Options

	ctx, err := program.Parse([]string{"login"})
	assert.NoError(t, err)
	assert.Equal(t, "login", ctx.Command())
}
