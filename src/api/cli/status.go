package cli

import (
	gcli "github.com/urfave/cli"

	"github.com/samoslab/samos/src/api/webrpc"
)

func statusCmd() gcli.Command {
	name := "status"
	return gcli.Command{
		Name:         name,
		Usage:        "Check the status of current samos node",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rpcClient := RPCClientFromContext(c)
			status, err := rpcClient.GetStatus()
			if err != nil {
				return err
			}

			return printJSON(struct {
				webrpc.StatusResult
				RPCAddress string `json:"webrpc_address"`
			}{
				StatusResult: *status,
				RPCAddress:   rpcClient.Addr,
			})
		},
	}
}
