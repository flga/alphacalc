package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"github.com/flga/lba2/internal/format"
	"github.com/flga/lba2/internal/lba"
)

func process(_ js.Value, args []js.Value) interface{} {
	data := args[0]
	req := lba.Request{
		Rsn:      data.Get("rsn").String(),
		Ironman:  data.Get("ironman").Bool(),
		Progress: lba.Progress(data.Get("progress").Int()),
		Charges:  data.Get("charges").Int(),
		Pts: lba.PtsCounter{
			Attacker:  data.Get("pts").Get("attacker").Int(),
			Defender:  data.Get("pts").Get("defender").Int(),
			Healer:    data.Get("pts").Get("healer").Int(),
			Collector: data.Get("pts").Get("collector").Int(),
		},
		Bxp: lba.BxpCounter{
			Agility:    data.Get("bxp").Get("agility").Int(),
			Firemaking: data.Get("bxp").Get("firemaking").Int(),
			Mining:     data.Get("bxp").Get("mining").Int(),
		},
		King:        data.Get("king").Int(),
		Queen:       data.Get("queen").Int(),
		Hm10Tickets: data.Get("hm10tickets").Int(),
	}

	order, err := lba.Process(req)
	if err != nil {
		return js.ValueOf(err.Error())
	}

	output := &bytes.Buffer{}

	fmt.Fprintln(output, "== Request ==============================")
	fmt.Fprintln(output, req.Print())

	fmt.Fprintln(output)
	fmt.Fprintln(output, "== Net ==================================")
	fmt.Fprintln(output, order.Print())

	fmt.Fprintln(output)
	fmt.Fprintln(output, "== Breakdown ============================")
	var cost int
	var dur time.Duration
	for i, s := range order.Breakdown {
		cost += s.Cost()
		dur += s.Duration()

		if i > 0 {
			fmt.Fprintln(output, "------------------------------------------")
		}
		fmt.Fprintln(output, s.Print())
	}
	fmt.Fprintln(output)
	fmt.Fprintln(output, "== Total ================================")
	fmt.Fprintf(output, "Cost: %s\n", format.Int(cost))
	fmt.Fprintf(output, "Duration: %s\n", dur)

	return js.ValueOf(output.String())
}

func registerCallbacks() {
	js.Global().Set("process", js.FuncOf(process))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
