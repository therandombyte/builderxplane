package logging

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

// setting klog as the backend
func SetFilteredKlogLogger(log logr.Logger) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(fs)
	fs.Parse([]string{"--v=3"})

	klogr := logr.New(&requestThrottlingFilter{log.GetSink()})
	klog.SetLogger(klogr)
}

type requestThrottlingFilter struct {
	logr.LogSink
}
