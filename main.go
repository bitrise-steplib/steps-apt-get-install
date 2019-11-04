package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/cache"
	"github.com/kballard/go-shellquote"
)

type config struct {
	Packages   string `env:"packages,required"`
	Options    string `env:"options"`
	Upgrade    string `env:"upgrade,opt[yes,no]"`
	CacheLevel string `env:"cache_level,opt[all,none]"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	var configs config
	if err := stepconf.Parse(&configs); err != nil {
		fail("Issue with input: %s", err)
	}

	fmt.Println()
	stepconf.Print(configs)
	fmt.Println()

	if configs.CacheLevel == "all" {
		if err := applyAllCache(); err != nil {
			fail("Could not apply caching: %s", err)
		}
	}

	log.Infof("$ apt-get %s", command.PrintableCommandArgs(false, []string{"update"}))
	if err := command.RunCommand("apt-get", "update"); err != nil {
		fail("Can't perform apt-get update: %s", err)
	}

	var cmdArgs []string
	if configs.Upgrade == "yes" {
		cmdArgs = append(cmdArgs, "upgrade", "-y")
	} else {
		cmdArgs = append(cmdArgs, "install", "-y")
	}
	if configs.Options != "" {
		args, err := shellquote.Split(configs.Options)
		if err != nil {
			fail("Can't split options: %s", err)
		}
		cmdArgs = append(cmdArgs, args...)
	}
	packages := strings.Split(configs.Packages, " ")
	cmdArgs = append(cmdArgs, packages...)

	fmt.Println()
	log.Infof("$ apt-get %s", command.PrintableCommandArgs(false, cmdArgs))
	if err := command.RunCommand("apt-get", cmdArgs...); err != nil {
		fail("Can't install packages:  %s", err)
	}
}

func applyAllCache() error {
	if err := removeDockerCleanFile(); err != nil {
		return fmt.Errorf("could not remove docker clean file: %s", err)
	}
	c := cache.New()
	c.IncludePath("/var/cache/apt/archives")
	if err := c.Commit(); err != nil {
		return fmt.Errorf("could not add packages to cache: %s", err)
	}
	return nil
}

func removeDockerCleanFile() error {
	return command.RemoveAll("/etc/apt/apt.conf.d/docker-clean")
}
