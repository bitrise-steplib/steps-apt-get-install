package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/cache"
	"github.com/kballard/go-shellquote"
)

const (
	cacheInputNone = "none"
	cacheInputAll  = "all"
)

// ConfigsModel ...
type ConfigsModel struct {
	Packages   string `env:"packages,required"`
	Options    string `env:"options"`
	Upgrade    string `env:"upgrade"`
	CacheLevel string `env:"cache_level,opt[all,none]"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Packages: %s", configs.Packages)
	log.Printf("- Options: %s", configs.Options)
	log.Printf("- Upgrade: %s", configs.Upgrade)
	log.Printf("- Cache Level: %s", configs.CacheLevel)
}

func (configs ConfigsModel) validate() error {
	if configs.Packages == "" {
		return errors.New("no Packages parameter specified")
	}
	if configs.Upgrade != "" && configs.Upgrade != "yes" && configs.Upgrade != "no" {
		return fmt.Errorf("invalid 'Upgrade' specified (%s), valid options: [yes no]", configs.Upgrade)
	}
	return nil
}

func main() {
	var configs ConfigsModel
	if err := stepconf.Parse(&configs); err != nil {
		fail("Issue with input: %s", err)
	}

	fmt.Println()
	configs.print()
	fmt.Println()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	if err := applyCacheConfig(configs.CacheLevel); err != nil {
		fail("Could not apply caching: %s", err)
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

func applyCacheConfig(cacheConfig string) error {
	switch cacheConfig {
	case cacheInputNone:
		return applyNoCache()
	case cacheInputAll:
		return applyAllCache()
	default:
		return fmt.Errorf("invalid cache level, no such configuration: %s", cacheConfig)
	}
}

func applyNoCache() error {
	if err := disablePackageCache(); err != nil {
		return fmt.Errorf("could not disable package cache: %s", err)
	}
	if err := removeDownloadedPackages(); err != nil {
		return fmt.Errorf("could not remove downloaded package cache: %s", err)
	}
	if err := removeLookupFiles(); err != nil {
		return fmt.Errorf("could not remove lookup files: %s", err)
	}
	return nil
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

func disablePackageCache() error {
	cacheFilePth := "/etc/apt/apt.conf.d/00_disable-cache"
	cacheFileCnt := `Dir::Cache::pkgcache "";` + "\n" + `Dir::Cache::srcpkgcache "";` + "\n" + `Dir::Cache "";` + "\n" + `Dir::Cache::archives "";`
	return fileutil.AppendStringToFile(cacheFilePth, cacheFileCnt)
}

func removeDownloadedPackages() error {
	return command.RemoveAll("/var/cache/apt/archives/")
}

func removeLookupFiles() error {
	return command.RemoveAll("/var/cache/apt/pkgcache.bin", "/var/cache/apt/pkgcache.bin")
}

func removeDockerCleanFile() error {
	return command.RemoveAll("/etc/apt/apt.conf.d/docker-clean")
}
