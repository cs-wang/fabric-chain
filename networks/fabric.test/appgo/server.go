package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hyperledger/fabric/common/flogging"
)

const (
	cmdRoot = "core"
)

var (
	// Logging
	logger = logging.MustGetLogger("DelayInsurance.app")
	versionFlag bool
	ContextMap map[string]*CurrentContext
	appStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the app.",
		Long:  `Starts a app that interacts with the network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			initNetWork()
			return serve(args)
		},
	}
	resetStartCmd = &cobra.Command{
		Use:   "update",
		Short: "update the app",
		Long:  "update a app that interacts with the network",
		RunE: func(cmd *cobra.Command, args []string) error {
			initNetWorkDonotSetChannelAndChaincode()
			return serve(args)
		},
	}
)


// The main command describes the service and
// defaults to printing the help message.
var mainCmd = &cobra.Command{
	Use: "app",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loggingSpec := viper.GetString("logging_level")
		if loggingSpec == "" {
			loggingSpec = "DEBUG"
		}
		flogging.InitFromSpec(loggingSpec)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			VersionPrint()
		} else {
			cmd.HelpFunc()(cmd, args)
		}
	},
}

func main() {
	// Logging
	var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	// viper init
	viper.AddConfigPath("../fixtures/config")
	viper.SetConfigName(cmdRoot)
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Define command-line flags that are valid for all peer commands and
	// subcommands.
	mainFlags := mainCmd.PersistentFlags()
	mainFlags.BoolVarP(&versionFlag, "version", "v", false, "Display current version of fabric peer server")
	mainCmd.AddCommand(VersionCmd())
	mainCmd.AddCommand(appStartCmd)
	mainCmd.AddCommand(resetStartCmd)
	runtime.GOMAXPROCS(viper.GetInt("core.gomaxprocs"))

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
	logger.Info("Exiting.....")
}

func initNetWorkDonotSetChannelAndChaincode() {
	setNetWork(false)
}

func initNetWork() {
	setNetWork(true)
}

func setNetWork(doInit bool) {
	APPClient := APPClient{
		ConfigFile:      "../config/config_app.yaml",
		ConnectEventHub: true,
	}
	APPClient.InitAPPClient()
	channels := viper.GetStringMap("channels")
	ContextMap = make(map[string]*CurrentContext)
	for currentChannelK, _ := range channels {
		currentChannelConfig := viper.GetString("channels." + currentChannelK + ".config")
		currentChannelPeers := viper.GetStringSlice("channels." + currentChannelK + ".peers")
		chaincodeid := viper.GetString("channels." + currentChannelK + ".chaincodeid")
		chaincodepath := viper.GetString("channels." + currentChannelK + ".chaincodepath")
		chaincodever := viper.GetString("channels." + currentChannelK + ".chaincodever")
		currentChainCode := ChainCode{
			ChainCodeId:chaincodeid,
			ChainCodeVersion:chaincodever,
			ChainCodePath:chaincodepath,
		}
		currentChannel := Channel{
			AppClient:&APPClient,
			ChannelId:currentChannelK,
			ChannelOrgIDs:   currentChannelPeers,
			ChannelConfig :currentChannelConfig,
			ChainCode:currentChainCode,
		}
		currentChannel.Organizations = make([]*Organization, len(currentChannelPeers))
		for i, currentPeerOrg := range currentChannelPeers {

			currentOrg := Organization{
				OrgID:           currentPeerOrg,
				OrgPath:         currentPeerOrg[4:8],
			}
			currentAdaptor := CurrentContext{
				AppClient:&APPClient,
				Channel:&currentChannel,
				Organization:&currentOrg,
			}
			err := currentAdaptor.SetOrgUser()
			if err != nil {
				logger.Error("currentAdaptor.SetOrgUser error:", err)
			}
			currentAdaptor.SetPeers()
			currentAdaptor.ConnectEventHub()
			currentChannel.Organizations[i] = currentAdaptor.Organization
			ContextMap[currentChannelK + currentPeerOrg] = &currentAdaptor

		}
		currentChannel.SetChannel()
		if doInit {
			currentChannel.CreateAndJoinChannel()
			currentChannel.InstallAndInstantiateCC(nil)
		}
	}
}
