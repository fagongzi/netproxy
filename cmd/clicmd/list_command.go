package clicmd

import (
	"fmt"

	"net/http"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// NewListCommand returns the cobra command for "list".
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [options]",
		Short: "List the clients",
		Run:   listCommandFunc,
	}

	return cmd
}

// listCommandFunc executes the "list" command.
func listCommandFunc(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("http://%s/api/clients", Global.Endpoints)
	cli := &http.Client{}
	rsp, err := cli.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer rsp.Body.Close()
		data, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n", data)
		}
	}
}
