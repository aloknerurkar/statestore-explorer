package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethersphere/bee/pkg/statestore/leveldb"
)

var usage string = `Statestore explorer
Usage:
	ssexplorer <Path to LevelDB statestore>
`

var actionUsage string = `
Actions:
	get		Get value for key
	count		Get count for prefix
	list		List values for prefix
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		return
	}

	if st, err := os.Stat(os.Args[1]); err != nil || !st.Mode().IsDir() {
		fmt.Println("LevelDB path incorrect\n")
		fmt.Println(usage)
		return
	}

	st, err := leveldb.NewStateStore(os.Args[1], nil)
	if err != nil {
		fmt.Println("Path not a statestore. Error opening:", err.Error())
		fmt.Println(usage)
		return
	}

	defer st.Close()

	fmt.Println("Statestore explorer")
	fmt.Println("----------------------")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n failed reading input: %s", err.Error())
			return
		}
		in = strings.Replace(in, "\n", "", -1)

		tokens := strings.Split(in, " ")
		switch tokens[0] {
		case "get":
			if tokens[1] == "" || len(tokens) > 2 {
				fmt.Println("usage: get <key>")
				continue
			}
			data, err := st.DB().Get([]byte(tokens[1]), nil)
			if err != nil {
				fmt.Printf("failed to get value: %v\n", err)
				continue
			}
			fmt.Println(string(data))
		case "count":
			if len(tokens) > 2 {
				fmt.Println("usage: count <prefix>")
				continue
			}
			prefix := ""
			if len(tokens) > 1 {
				prefix = tokens[1]
			}
			count := 0
			st.Iterate(prefix, func(_, _ []byte) (bool, error) {
				count++
				return false, nil
			})
			// For entries which use shed, there is a prefix added before key
			if count == 0 {
				np := append([]byte{1}, []byte(prefix)...)
				st.Iterate(string(np), func(_, _ []byte) (bool, error) {
					count++
					return false, nil
				})
			}
			fmt.Println("count: ", count)
		case "list":
			if len(tokens) > 3 {
				fmt.Println("usage: list <prefix>")
				continue
			}
			prefix := ""
			if len(tokens) > 1 && tokens[1] != "all" {
				prefix = tokens[1]
			}
			skip, limit := 0, 50
			if len(tokens) > 2 {
				skip, _ = strconv.Atoi(tokens[2])
			}
			for {
				skipLoop := 0
				current := 0
				st.Iterate(prefix, func(k, _ []byte) (bool, error) {
					if skipLoop < skip {
						skipLoop++
						return false, nil
					}
					if current < limit {
						current++
						skip++
						fmt.Println(string(k))
						return false, nil
					}
					return true, nil
				})
				if current == 0 && skipLoop == 0 {
					np := append([]byte{1}, []byte(prefix)...)
					prefix = string(np)
					continue
				}
				if current == 0 {
					fmt.Println("Done")
					return
				}
				fmt.Println("---- Press 'Enter' to load more, any other key to exit ----", skipLoop, current)
				c, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("\n failed reading input: %s", err.Error())
					return
				}
				c = strings.Replace(c, "\n", "", -1)
				if string(c) != "" {
					fmt.Println("exiting", string(c))
					break
				}
			}
		case "exit":
			return
		default:
			fmt.Println(actionUsage)
			continue
		}
	}
}
