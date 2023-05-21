package main

import (
	builtin "ice/Builtin"
	com "ice/Common"
	ice "ice/Lang"
	shell "iceShell/Shell"
	"os"
	"runtime/pprof"
)

func exitOnError(err ice.Err) {
	com.PrintAndLog(err.Error())
	os.Exit(-1)
}

func main() {
	com.InitLog()
	defer com.CloseLog()

	var filename string
	var globalMemFile string = "GlobalMem.ice"

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-t":
			com.Trace = true
		default:
			filename = arg
		}
	}

	if len(filename) != 0 {
		// Init Profiling
		if profilingFile, err := os.Create("profile.pprof"); err == nil {
			defer profilingFile.Close()
			if err = pprof.StartCPUProfile(profilingFile); err != nil {
				com.Warning("Failed to start CPU Profiling: ", err)
			} else {
				defer pprof.StopCPUProfile()
			}
		} else {
			com.Warning("Failed to create profile.pprof")
		}

		//Run from file
		ctx := builtin.NewContext(globalMemFile)
		defer builtin.SaveGlobalContext(globalMemFile, ctx)

		obj := ctx.RunFile(filename)
		if ice.IsErr(obj) {
			exitOnError(ice.NewErrorCtx(ctx, "Run from file", "Startup failed", obj))
		}
	} else {
		//Run interactive shell
		shell.Initialize()
		if err := shell.RunShell(globalMemFile); err != nil {
			exitOnError(ice.NewError("Run Shell failed", err))
		}
	}

	com.Nl()
}
