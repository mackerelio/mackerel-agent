// auto-generated file

package main

import "github.com/motemen/go-cli"

func init() {
	cli.Use(
		&cli.Command{
			Name:   "",
			Action: doMain,
			Short:  "mackerel-agent",
			Long:   "mackerel-agent [options]\n\nmain process of mackerel-agent",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "version",
			Action: doVersion,
			Short:  "display version of mackerel-agent",
			Long:   "version\n\ndisplay the version of mackerel-agent",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "configtest",
			Action: doConfigtest,
			Short:  "configtest",
			Long:   "configtest\n\ndo configtest",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "retire",
			Action: doRetire,
			Short:  "retire the host",
			Long:   "retire [-force]\n\nretire the host",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "once",
			Action: doOnce,
			Short:  "output onetime",
			Long:   "once\n\noutput metrics and meta data of the host one time.\nThese data are only displayed and not posted to Mackerel.",
		},
	)
}
