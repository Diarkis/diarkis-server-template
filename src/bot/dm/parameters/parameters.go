// © 2019-2024 Diarkis Inc. All rights reserved.

package parameters

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Diarkis/diarkis-server-template/bot/dm/loadtest"
)

func ParseParams(args []string) *loadtest.Params {
	p := &loadtest.Params{}
	for _, v := range args {
		list := strings.Split(v, "=")
		// invalid format:
		// valid format is $(name)=$(value)
		if len(list) != 2 {
			continue
		}
		name := list[0]
		value := list[1]
		switch name {
		case "host":
			p.Host = value
		case "bots":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Howmany = n
		case "protocol":
			p.Protocol = strings.ToUpper(value)
		case "size":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Size = n
		case "interval":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Interval = int64(n)
		}
	}
	return p
}
