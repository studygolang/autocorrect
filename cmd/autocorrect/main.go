package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/studygolang/autocorrect"

	"github.com/urfave/cli"
)

var (
	version = "1.0.0"
	build   = "100"
)

func main() {
	app := cli.NewApp()
	app.Name = "autocorrect"
	app.Usage = "自动给中英文之间加入合理的空格并纠正专用名词大小写"
	app.Version = fmt.Sprintf("%s (%s)", version, build)
	app.Copyright = "(c) 2018 studygolang.com"
	app.Authors = []cli.Author{
		{
			Name:  "Xu Xinhua",
			Email: "polaris@studygolang.com",
		},
	}
	app.EnableBashCompletion = true
	app.Flags = globalFlags
	app.Action = action
	app.Commands = commands
	app.Run(os.Args)
}

var globalFlags = []cli.Flag{}

var cmdCommonflags = []cli.Flag{
	cli.StringFlag{
		Name:  "outfile, f",
		Value: "",
		Usage: "输出到哪个文件 `FILE`",
	},
	cli.BoolFlag{
		Name:  "w",
		Usage: "将结果直接写入源文件，而不是输出到标准输出，如果提供了 outfile 选项，忽略此选项",
	},
}

var action = func(c *cli.Context) error {
	return cli.ShowAppHelp(c)
}

var bashComplete = func(c *cli.Context) {
	// This will complete if no args are passed
	if c.NArg() > 0 {
		return
	}
	for _, name := range commandNames {
		fmt.Println(name)
	}
}

var commandNames = []string{"space", "correct", "convert"}

var commands = []cli.Command{
	{
		Name:    "space",
		Aliases: []string{"s"},
		// Category: "",
		Usage:        "给中英文之间加入合理的空格",
		Action:       processSpace,
		Flags:        cmdCommonflags,
		BashComplete: bashComplete,
	},
	{
		Name:    "correct",
		Aliases: []string{"c"},
		// Category: "",
		Usage:        "纠正专用名词大小写",
		Action:       processCorrect,
		Flags:        cmdCommonflags,
		BashComplete: bashComplete,
	},
	{
		Name:    "convert",
		Aliases: []string{"a"},
		// Category: "",
		Usage:        "同时做 space 和 correct 的事情",
		Action:       processConvert,
		Flags:        cmdCommonflags,
		BashComplete: bashComplete,
	},
}

func processStd(ctx *cli.Context) (string, error) {
	content := ""

	if ctx.NArg() > 0 {
		arg := ctx.Args().First()
		if !exits(arg) {
			content = arg
		} else {
			b, err := ioutil.ReadFile(arg)
			if err != nil {
				return "", err
			}
			content = string(b)
		}

		return content, nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
	}

	return content, nil
}

func processSpace(ctx *cli.Context) error {
	content, err := processStd(ctx)
	if err != nil {
		return err
	}

	content = autocorrect.AutoSpace(content)
	return output(ctx, content)
}

func processCorrect(ctx *cli.Context) error {
	content, err := processStd(ctx)
	if err != nil {
		return err
	}

	content = autocorrect.AutoCorrect(content)
	return output(ctx, content)
}

func processConvert(ctx *cli.Context) error {
	content, err := processStd(ctx)
	if err != nil {
		return err
	}

	content = autocorrect.Convert(content)
	return output(ctx, content)
}

func output(ctx *cli.Context, content string) error {
	if ctx.IsSet("outfile") {
		outfile := ctx.String("outfile")
		file, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(content)

		return err
	}

	if ctx.IsSet("w") {
		if ctx.NArg() > 0 {
			arg := ctx.Args().First()
			if !exits(arg) {
				fmt.Println(content)
				return nil
			}

			ioutil.WriteFile(arg, []byte(content), 0666)

			return nil
		}
	}

	fmt.Println(content)

	return nil
}

func exits(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
