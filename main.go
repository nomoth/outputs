package main

import (
	"fmt"
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"outputs/ipc"
)

var (
	verbose = kingpin.Flag("verbose", "display commands executed").Short('v').Bool()
	outputs = kingpin.Arg("outputs", "output order from left to right").Strings()
)

func main() {
	kingpin.Parse()
	conn := ipc.NewConnection()
	defer conn.Close()
	list, err := conn.GetOutputs()
	if err != nil {
		log.Println(err)
	}

	if len(*outputs) > 0 {
		setOutputs(conn, list)
	} else {
		listOutputs(conn, list)
	}
}

func setOutputs(conn *ipc.Connection, available []*ipc.Output) {
	x := 0
	for _, n := range *outputs {
		for _, o := range available {
			if o.Name == n && o.Active {
				cmd := fmt.Sprintf("output %s pos %d 0", o.Name, x)
				if *verbose {
					fmt.Println(cmd)
				}
				err := conn.Run(cmd)
				if err != nil {
					fmt.Printf("error: %s\n", err)
				}
				x += int(float32(o.CurrentMode.Width) / o.Scale)
				break
			}
		}
	}
}

func listOutputs(conn *ipc.Connection, available []*ipc.Output) {
	for _, o := range available {
		if o.Active {
			fmt.Printf("%s scale:%1.2f %dx%d %dHz\n", o.Name, o.Scale, o.CurrentMode.Width, o.CurrentMode.Height, o.CurrentMode.Refresh/1000)
		} else {
			fmt.Printf("%s disabled\n", o.Name)
		}
	}
}
