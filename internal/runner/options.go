package runner

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/projectdiscovery/cloudlist/pkg/schema"
	"github.com/projectdiscovery/gologger"
	"gopkg.in/yaml.v2"
)

// Options contains the configuration options for cloudlist.
type Options struct {
	JSON      bool   // JSON returns JSON output
	Silent    bool   // Silent Display results only
	Version   bool   // Version returns the version of the tool.
	Verbose   bool   // Verbose prints verbose output.
	Hosts     bool   // Hosts specifies to fetch only DNS Names
	IPAddress bool   // IPAddress specifes to fetch only IP Addresses
	Config    string // Config is the location of the config file.
	Output    string // Output is the file to write found results too.
	Provider  string // Provider specifies what providers to fetch assets for.
}

var defaultConfigLocation = path.Join(userHomeDir(), "/.config/cloudlist/config.yaml")

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}

	flag.BoolVar(&options.JSON, "json", false, "Show json output")
	flag.BoolVar(&options.Silent, "silent", false, "Show only results in output")
	flag.BoolVar(&options.Version, "version", false, "Show version of cloudlist")
	flag.BoolVar(&options.Verbose, "v", false, "Show Verbose output")
	flag.BoolVar(&options.Hosts, "host", false, "Show only hosts in output")
	flag.BoolVar(&options.IPAddress, "ip", false, "Show only IP addresses in output")
	flag.StringVar(&options.Config, "config", defaultConfigLocation, "Configuration file to use for enumeration")
	flag.StringVar(&options.Output, "o", "", "File to write output to (optional)")
	flag.StringVar(&options.Provider, "provider", "", "Provider to fetch assets from (optional)")
	flag.Parse()

	options.configureOutput()
	showBanner()

	if options.Version {
		gologger.Infof("Current Version: %s\n", Version)
		os.Exit(0)
	}
	checkAndCreateConfigFile(options)
	return options
}

// configureOutput configures the output on the screen
func (options *Options) configureOutput() {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.MaxLevel = gologger.Verbose
	}
	if options.Silent {
		gologger.MaxLevel = gologger.Silent
	}
}

// readConfig reads the config file from the options
func readConfig(configFile string) (schema.Options, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := schema.Options{}
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		if err == io.EOF {
			return nil, errors.New("invalid configuration file provided")
		}
		return nil, err
	}
	return config, nil
}

// checkAndCreateConfigFile checks if a config file exists,
// if not creates a default.
func checkAndCreateConfigFile(options *Options) {
	if options.Config == defaultConfigLocation {
		os.MkdirAll(path.Dir(options.Config), os.ModePerm)
		if _, err := os.Stat(defaultConfigLocation); os.IsNotExist(err) {
			if writeErr := ioutil.WriteFile(defaultConfigLocation, []byte(defaultConfigFile), os.ModePerm); writeErr != nil {
				gologger.Warningf("Could not write default output to %s: %s\n", defaultConfigLocation, writeErr)
			}
		}
	}
}

func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		gologger.Fatalf("Could not get user home directory: %s\n", err)
	}
	return usr.HomeDir
}

const defaultConfigFile = `# Configuration file for cloudlist enumeration agent
#- # provider is the name of the provider
#  provider: do
#  # profile is the name of the provider profile
#  profile: xxxx
#  # digitalocean_token is the API key for digitalocean cloud platform
#  digitalocean_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
#
#- # provider is the name of the provider
#  provider: scw
#  # scaleway_access_key is the access key for scaleway API
#  scaleway_access_key: SCWXXXXXXXXXXXXXX
#  # scaleway_access_token is the access token for scaleway API
#  scaleway_access_token: xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx
#
#- # provider is the name of the provider
#  provider: aws
#  # profile is the name of the provider profile
#  profile: staging
#  # aws_access_key is the access key for AWS account
#  aws_access_key: AKIAXXXXXXXXXXXXXX
#  # aws_secret_key is the secret key for AWS account
#  aws_secret_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
#  # aws_session_token session token for temporary security credentials retrieved via STS (optional)
#  aws_session_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
