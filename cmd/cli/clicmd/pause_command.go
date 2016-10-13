package clicmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fagongzi/netproxy/pkg/proxy"
	"github.com/spf13/cobra"
)

// NewPauseCommand returns the cobra command for "pause".
func NewPauseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause <addr>",
		Short: "Pause proxy on addr",
		Run:   pauseCommandFunc,
	}

	return cmd
}

// pauseCommandFunc executes the "pause" command.
func pauseCommandFunc(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("http://%s%s", Global.Endpoints, proxy.APIProxies)
	cli := &http.Client{}
	fmt.Printf("<%s> send to server\n", args[0])
	request, _ := http.NewRequest("DELETE", url, bytes.NewReader([]byte(args[0])))
	rsp, err := cli.Do(request)

	if err != nil {
		fmt.Println(err)
	} else {
		defer rsp.Body.Close()
		data, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("server return %s\n", data)
		}
	}
}
