package zerolog

import (
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type (
	Logger struct {
		standard *log.Logger
		zerolog  *zerolog.Logger
		Out, Err *os.File
	}
)

var (
	logID = strconv.FormatInt(time.Now().UnixMicro(), 36)
	sOUT  = os.Stdout
	sERR  = os.Stderr
)

func NewZeroLog(ctx context.Context, c ...io.Writer) *Logger {
	ws, z := make([]io.Writer, 0), new(zerolog.Logger)

	for j := 0; j < len(c); j++ {
		if c[j] != nil {
			ws = append(ws, c[j])
		}
	}

	if len(c) > 0 {
		switch len(ws) {
		case 0:
			*z = zerolog.Nop()
		case 1:
			*z = zerolog.New(ws[0]).With().Timestamp().Stack().Logger()
		default:
			*z = zerolog.New(zerolog.MultiLevelWriter(ws...))
		}
	} else if zz := zerolog.Ctx(ctx); zz != nil {
		*z = *zz
	}

	dir := os.TempDir()
	tempOUT, _ := os.Create(dir + "/golang-clean-architecture-" + logID + "-out.log")
	tempERR, _ := os.Create(dir + "/golang-clean-architecture-" + logID + "-err.log")

	return &Logger{log.New(z, "", 0), z, tempOUT, tempERR}
}

func (x *Logger) S() *log.Logger     { return x.standard }
func (x *Logger) Z() *zerolog.Logger { return x.zerolog }
func (x *Logger) Level(level string) *Logger {
	lv, err := zerolog.ParseLevel(strings.ToLower(level))
	if err == nil {
		*x.zerolog = x.zerolog.Level(lv)
		x.standard.SetOutput(x.zerolog)
	}

	return x
}
func (x *Logger) Unswap() {
	os.Stdout, os.Stderr = sOUT, sERR

	log.SetOutput(os.Stderr)
}
