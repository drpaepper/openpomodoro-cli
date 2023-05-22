package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/drpaepper/go-openpomodoro"
	"github.com/drpaepper/openpomodoro-cli/hook"
	"github.com/justincampbell/go-countdown"
	"github.com/justincampbell/go-countdown/format"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "break [duration]",
		Short: "Take a break",
		RunE:  breakCmd,
	}

	RootCmd.AddCommand(command)
}

func breakCmd(cmd *cobra.Command, args []string) error {
	d := settings.DefaultBreakDuration.Minutes()

	p := openpomodoro.NewPomodoro()
	if len(args) > 0 {
		x, err := strconv.Atoi(args[0])
		if err != nil {
			p.Duration = time.Duration(x) * time.Minute
		}
	} else {
		p.Duration = time.Duration(int(d)) * time.Minute
	}

	p.Description = "BREAK"
	p.StartTime = time.Now().Add(-agoFlag)
	var tags []string
	tags = make([]string, 1)
	tags[0] = "BREAK"
	p.Tags = tags

	if err := client.Start(p); err != nil {
		return err
	}

	if err := hook.Run(client, "break"); err != nil {
		return err
	}

	return statusCmd(cmd, args)
}

func wait(d time.Duration) error {
	err := countdown.For(d, time.Second).Do(func(c *countdown.Countdown) error {
		fmt.Printf("\r%s", format.MinSec(c.Remaining()))
		return nil
	})

	fmt.Println()

	return err
}
