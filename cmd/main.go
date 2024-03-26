package main

import (
	"log"
	"os"

	as "github.com/osmosis-labs/autenticator-test/cmd/seed/auth"
	cls "github.com/osmosis-labs/autenticator-test/cmd/seed/cl"
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

	EdgeChainID = "edgenet"
	EdgeAddress = "161.35.19.190:9090"

	LocalChainID = "edgenet"
	LocalAddress = "0.0.0.0:9090"

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

	localCmd := &cobra.Command{
		Use:   "local",
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

	// Localnet Config SetUp
	localConfig := config.SetUp(
		LocalChainID,
		LocalAddress,
		[]string{
			TestKeyValidator,
			TestKeyUser1,
			TestKeyUser2,
			TestKeyUser3,
			TestKeyUser4,
		},
		DefaultDenoms,
	)

	// LocalCommands
	authenticatorCmd.AddCommand(
		as.SeedCreateOneClickTradingAccount(localConfig),
		as.SeedSwapCmd(localConfig),
		as.SeedRemoveAllAuthenticators(localConfig),
		as.SeedCreateCosigner(localConfig),
	)

	clCmd.AddCommand(
		cls.StartClIncentiveFlow(localConfig),
		cls.StartClSwapAndTransferPositionFlow(localConfig),
	)

	stakingCmd.AddCommand(
		vs.StartValidatorFlow(localConfig),
	)

	localCmd.AddCommand(
		authenticatorCmd,
		clCmd,
		stakingCmd,
	)

	// Edgenet Config SetUp
	edgeConfig := config.SetUp(
		EdgeChainID,
		EdgeAddress,
		[]string{
			TestKeyValidator,
			TestKeyUser1,
			TestKeyUser2,
			TestKeyUser3,
			TestKeyUser4,
		},
		DefaultDenoms,
	)

	edgeCmd := &cobra.Command{
		Use:   "edge",
		Short: "the edge command interacts with an edgenet deployed here: " + EdgeAddress,
	}

	edgeAuthenticatorCmd := &cobra.Command{
		Use:   "auth",
		Short: "auth has a seeds that run to interact with authenticators",
	}

	edgeCLCmd := &cobra.Command{
		Use:   "cl",
		Short: "cl has a seeds that run to interact with the concentrated-liquidity",
	}

	edgeStakingCmd := &cobra.Command{
		Use:   "stake",
		Short: "stake has a seeds that run to interact with the staking module",
	}

	edgeAuthenticatorCmd.AddCommand(
		as.SeedCreateOneClickTradingAccount(edgeConfig),
		as.SeedSwapCmd(edgeConfig),
		as.SeedRemoveAllAuthenticators(edgeConfig),
		as.SeedCreateCosigner(edgeConfig),
	)

	edgeCLCmd.AddCommand(
		cls.StartClIncentiveFlow(edgeConfig),
		cls.StartClSwapAndTransferPositionFlow(edgeConfig),
	)

	edgeStakingCmd.AddCommand(
		vs.StartValidatorFlow(edgeConfig),
	)

	edgeCmd.AddCommand(
		edgeAuthenticatorCmd,
		edgeCLCmd,
		edgeStakingCmd,
	)

	// ROOT command
	rootCmd.AddCommand(
		localCmd,
		edgeCmd,
	)

	return rootCmd
}
