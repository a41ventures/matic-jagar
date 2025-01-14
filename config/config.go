package config

import (
	"os/user"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Telegram bot details struct
	Telegram struct {
		// BotToken is the token of your telegram bot
		BotToken string `mapstructure:"tg_bot_token"`
		// ChatID is the id of telegarm chat which will be used to get alerts
		ChatID int64 `mapstructure:"tg_chat_id"`
	}

	// SendGrid stores sendgrid API credentials
	SendGrid struct {
		// Token of sendgrid account
		Token string `mapstructure:"sendgrid_token"`
		// ToEmailAddress is the email to which all the alerts will be sent
		ReceiverEmailAddress string `mapstructure:"receiver_email_address"`
		// SendgridEmail is the email of sendgrid account which will be used to send mail alerts
		SendgridEmail string `mapstructure:"account_email"`
		// SendgridName is the name of sendgrid account which will be used to send mail alerts
		SendgridName string `mapstructure:"sendgrid_account_name"`
	}

	// Scraper defines the time intervals for multiple scrapers to fetch the data
	Scraper struct {
		// Rate is to call and get the data for specified targets on that particular time interval
		Rate string `mapstructure:"rate"`
		// ValidatorRate is to call and fetch the data from validatorStatus target on that time interval
		ValidatorRate string `mapstructure:"validator_rate"`
		// ContractRate is to call and fetch the data from smart contract realted targets on that time interval
		ContractRate string `mapstructure:"contract_rate"`
		// CommandsRate is to check the for telegram commands from telegram chat and returns the data
		CommandsRate string `mapstructure:"tg_commnads_rate"`
	}

	// InfluxDB stores influxDB credntials
	InfluxDB struct {
		// Port on which influxdb is running
		Port string `mapstructure:"port"`
		// IP to connect to influxdb where it is running
		IP string `mapstructure:"ip"`
		// Database is the name of the influxdb database to store the data
		Database string `mapstructure:"database"`
		// Username is the name of the user of influxdb
		Username string `mapstructure:"username"`
		// Password of influxdb
		Password string `mapstructure:"password"`
	}

	// Endpoints defines multiple API base-urls to fetch the data
	Endpoints struct {
		EthRPCEndpoint         string `mapstructure:"eth_rpc_endpoint"`
		BorRPCEndpoint         string `mapstructure:"bor_rpc_end_point"`
		BorExternalRPC         string `mapstructure:"bor_external_rpc"`
		HeimdallRPCEndpoint    string `mapstructure:"heimdall_rpc_endpoint"`
		HeimdallLCDEndpoint    string `mapstructure:"heimdall_lcd_endpoint"`
		HeimdallExternalRPC    string `mapstructure:"heimdall_external_rpc"`
		PolygonStakingEndpoint string `mapstructure:"polygon_staking_endpoint"`
	}

	// ValDetails stores the validator meta details
	ValDetails struct {
		// ValidatorHexAddress will be used to get account balances, proposals and used to check missed blocks, etc
		ValidatorHexAddress string `mapstructure:"validator_hex_addr"`
		// SignerAddress will be used to get latest block, current poposer etc
		SignerAddress string `mapstructure:"signer_address"`
		// ValidatorName is the moniker of your validator which will be used to display in alerts messages
		ValidatorName string `mapstructure:"validator_name"`
		// StakeManagerContract is the address of stake manager contract which will be used to get vaidator share contract etc
		StakeManagerContract string `mapstructure:"stake_manager_contract"`
		// ValidatorNumber is the number of your validator at staking.polygon.technology
		ValidatorNumber string `mapstructure:"validator_number"`
	}

	// EnableAlerts struct which holds options to enalbe/disable alerts
	EnableAlerts struct {
		EnableTelegramAlerts bool `mapstructure:"enable_telegram_alerts"`
		EnableEmailAlerts    bool `mapstructure:"enable_email_alerts"`
	}

	// RegularStatusAlerts defines time-slots to receive validator status alerts
	RegularStatusAlerts struct {
		// AlertTimings is the array of time slots to send validator status alerts
		AlertTimings []string `mapstructure:"alert_timings"`
	}

	// AlerterPreferences stores individual alert settings to enable/disable particular alert
	AlerterPreferences struct {
		BalanceChangeAlerts string `mapstructure:"balance_change_alerts"`
		VotingPowerAlerts   string `mapstructure:"voting_power_alerts"`
		ProposalAlerts      string `mapstructure:"proposal_alerts"`
		BlockDiffAlerts     string `mapstructure:"block_diff_alerts"`
		MissedBlockAlerts   string `mapstructure:"missed_block_alerts"`
		NumPeersAlerts      string `mapstructure:"num_peers_alerts"`
		NodeSyncAlert       string `mapstructure:"node_sync_alert"`
		NodeStatusAlert     string `mapstructure:"node_status_alert"`
		EthLowBalanceAlert  string `mapstructure:"eth_low_balance_alert"`
	}

	//  AlertingThreshold defines threshold condition for different alert-cases.
	// `Alerter` will send alerts if the condition reaches the threshold
	AlertingThreshold struct {
		// NumPeersThreshold is to alert when the connected peers falls below this threshold
		NumPeersThreshold int64 `mapstructure:"num_peers_threshold"`
		// MissedBlocksThreshold is to alert when validator misses continuous missed bocks
		// Alerter will send alerts when the missed blocks count reaches the configured threshold
		MissedBlocksThreshold int64 `mapstructure:"missed_blocks_threshold"`
		// BlockDiffThreshold is to send alerts when the difference b/w network and validator
		// block height reaches the given threshold
		BlockDiffThreshold int64 `mapstructure:"block_diff_threshold"`
		// EthBalanceThreshold is to send alerts when the etherium balance falls below the configured threshold
		EthBalanceThreshold float64 `mapstructure:"eth_balance_threshold"`
	}

	// Config defines all the configurations required for the app
	Config struct {
		Endpoints           Endpoints           `mapstructure:"rpc_and_lcd_endpoints"`
		ValDetails          ValDetails          `mapstructure:"validator_details"`
		EnableAlerts        EnableAlerts        `mapstructure:"enable_alerts"`
		RegularStatusAlerts RegularStatusAlerts `mapstructure:"regular_status_alerts"`
		AlerterPreferences  AlerterPreferences  `mapstructure:"alerter_preferences"`
		AlertingThresholds  AlertingThreshold   `mapstructure:"alerting_threholds"`
		Scraper             Scraper             `mapstructure:"scraper"`
		Telegram            Telegram            `mapstructure:"telegram"`
		SendGrid            SendGrid            `mapstructure:"sendgrid"`
		InfluxDB            InfluxDB            `mapstructure:"influxdb"`
	}
)

// ReadFromFile to read config details using viper
func ReadFromFile() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configPath := path.Join(usr.HomeDir, `.matic-jagar/config/`)
	log.Printf("Config Path : %s", configPath)

	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

// Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
