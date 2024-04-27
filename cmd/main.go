package main

import (
	"log"
	"os"

	auctionSeed "github.com/osmosis-labs/autenticator-test/cmd/seed/auction"
	as "github.com/osmosis-labs/autenticator-test/cmd/seed/auth"
	cls "github.com/osmosis-labs/autenticator-test/cmd/seed/cl"
	gov "github.com/osmosis-labs/autenticator-test/cmd/seed/gov"
	vs "github.com/osmosis-labs/autenticator-test/cmd/seed/staking"
	"github.com/osmosis-labs/autenticator-test/pkg/config"

	"github.com/spf13/cobra"
)

// Test data for the seeds to run
const (
	GrpcConnectionTimeoutSeconds = 10
	TestKeyValidator             = "9ff80c31b47c7f2946654f569a6b1530db78d7fa5b3ea16db82570cdfd6d43f6"
	TestKeyUser1                 = "48d23cc417a30674e907a2403f109f082d92e197823d02e6a423c6aeb8e41204"
	TestKeyUser2                 = "6e67cda92a2ffa21242e8a01e03f93d13b8b3b3094e75e58fee480f16f98855a"
	TestKeyUser3                 = "40fc464087a28a93e697615f9585af7d763c8bd4b9cd50412c19c74fa501af41"

	// TestUser4 is not in the auth store
	TestKeyUser4         = "3d23af3840f0535863518fa8bbb8b98a231aa0bd2eb181911bfd8930f0ada7f9"
	AccountAddressPrefix = "osmo"

	EdgeChainID = "smartaccount"
	EdgeAddress = "164.92.247.225:9090"

	LocalChainID = "edgenet"
	LocalAddress = "161.35.19.190:9090"

	//	MainChainID = "osmosis-1"
	//	LocalAddress = ":9090"

	MainnetAddress = "grpc.osmosis.zone:9090"
)

var DefaultDenoms = map[string]string{
	"OsmoDenom":     "uosmo",
	"IonDenom":      "uion",
	"StakeDenom":    "stake",
	"AtomDenom":     "uatom",
	"DaIBCiDenom":   "ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7",
	"OsmoIBCDenom":  "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
	"StakeIBCDenom": "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787",
	"UstIBCDenom":   "ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC",
	"LuncIBCDenom":  "ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0",
	"AtomIBCDenom":  "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"UsdcIBCDenom":  "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4",
}

const (
	appName = "osmosis-test"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCmd returns the root command for parser.
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "osmosis-test has a variety of seeds that run against localnet, testnet, and mainnet",
	}

	// localnet commands
	localCmd := SetUpCmds("local", LocalChainID, LocalAddress)

	// edgenet commands
	edgeCmd := SetUpCmds("edge", EdgeChainID, EdgeAddress)

	// ROOT command
	rootCmd.AddCommand(
		localCmd,
		edgeCmd,
	)

	return rootCmd
}

func SetUpCmds(cmdName, chainID, address string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: "the local command interacts with a local node deployed here: " + LocalAddress,
	}

	authenticatorCmd := &cobra.Command{
		Use:   "auth",
		Short: "auth has a seeds that run to interact with authenticators",
	}

	clCmd := &cobra.Command{
		Use:   "cl",
		Short: "cl has a seeds that run to interact with the concentrated-liquidity",
	}

	stakingCmd := &cobra.Command{
		Use:   "stake",
		Short: "stake has a seeds that run to interact with the staking module",
	}

	govCmd := &cobra.Command{
		Use:   "gov",
		Short: "gov has a seeds that run to interact with the gov module",
	}

	auctionCmd := &cobra.Command{
		Use:   "auction",
		Short: "auction has a seeds that run to interact with the auction module",
	}

	conf := config.SetUp(
		chainID,
		address,
		[]string{
			TestKeyValidator,
			TestKeyUser1,
			TestKeyUser2,
			TestKeyUser3,
			TestKeyUser4,
		},
		DefaultDenoms,
	)

	authenticatorCmd.AddCommand(
		as.SeedCreateOneClickTradingAccount(conf),
		as.SeedSwapCmd(conf),
		as.SeedRemoveAllAuthenticators(conf),
		as.SeedCreateCosigner(conf),
	)

	auctionCmd.AddCommand(
		auctionSeed.StartAuctionFlow(conf),
	)

	clCmd.AddCommand(
		cls.StartClIncentiveFlow(conf),
		cls.StartClSwapAndTransferPositionFlow(conf),
	)

	stakingCmd.AddCommand(
		vs.StartValidatorFlow(conf),
	)

	govCmd.AddCommand(
		gov.StartGovernanceFlow(conf),
	)

	cmd.AddCommand(
		authenticatorCmd,
		auctionCmd,
		clCmd,
		stakingCmd,
		govCmd,
	)

	return cmd
}
