package clicmd

import (
	"bytes"
	"fmt"

	"net/http"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// NewResumeCommand returns the cobra command for "resume".
func NewResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume <addr>",
		Short: "Resume proxy on addr",
		Run:   resumeCommandFunc,
	}

	return cmd
}

// resumeCommandFunc executes the "resume" command.
func resumeCommandFunc(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("http://%s/api/proxy/resume", Global.Endpoints)
	cli := &http.Client{}
	fmt.Printf("<%s> send to server\n", args[0])
	request, _ := http.NewRequest("PUT", url, bytes.NewReader([]byte(args[0])))
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
