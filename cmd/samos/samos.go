package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/samoslab/samos/src/api/webrpc"
	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
	"github.com/samoslab/samos/src/daemon"
	"github.com/samoslab/samos/src/gui"
	"github.com/samoslab/samos/src/util/browser"
	"github.com/samoslab/samos/src/util/cert"
	"github.com/samoslab/samos/src/util/file"
	"github.com/samoslab/samos/src/util/logging"
	"github.com/samoslab/samos/src/visor"
	"github.com/samoslab/samos/src/wallet"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "1.2.0"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	help = false

	logger = logging.MustGetLogger("main")

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "73874031bb09f93f7a7acfb23cd6b823abaa33cdeb07806ea57c998c715b99d07ee7b599df1c35526cfee8d91c7204194bbf8880f6646714e72dcca1b80d09de00"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2PzndHacXbmM8GNjMsA5dDTiyQFiKzjpFzX"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "02aecd90febe163da3c4ac5bb711d9a87b2950d11413541acc9bda17fbda47954e"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	TrustAddressStr          = "EX8omhDyjKtc8zHGp1KZwn7usCndaoJxSe"
	BlockchainTrustPubkeyStr = "02aecd90febe163da3c4ac5bb711d9a87b2950d11413541acc9bda17fbda47954e"
	BlockchainTrustSeckeyStr = ""

	TrustPubkeyListStr = "02aecd90febe163da3c4ac5bb711d9a87b2950d11413541acc9bda17fbda47954e,02d15bf28c4ed2c39b35b2be2f8bcde1318e2b3b65fe2a676db39b520bee9bfe86,02e99a1338841e8b1f192337d2c6157045faa0cfe3b8a02210283aed7f5ad6880d"

	// BlockchainSeckeyFile encrypted seckey file
	BlockchainSeckeyFile = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1426562704
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 300e12

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"47.52.211.167:8858",
		"47.74.7.161:8858",
		"47.254.130.80:8858",
		"47.52.222.166:8858",
	}
)

// Config records the node's configuration
type Config struct {
	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Disables wallet API
	EnableWalletAPI bool
	// Disable CSRF check in the wallet api
	DisableCSRF bool
	// Enable /wallet/seed api endpoint
	EnableSeedAPI bool

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max outgoing connections to maintain
	MaxOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface      bool
	WebInterfacePort  int
	WebInterfaceAddr  string
	WebInterfaceCert  string
	WebInterfaceKey   string
	WebInterfaceHTTPS bool

	RPCInterface     bool
	RPCInterfacePort int
	RPCInterfaceAddr string

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.samos
	DataDirectory string
	// GUI directory contains assets for the html gui
	GUIDirectory string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	// Wallet crypto type
	WalletCryptoType string

	RunMaster bool

	GenesisSignature cipher.Sig
	GenesisTimestamp uint64
	GenesisAddress   cipher.Address

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

	BlockchainTrustPubkey cipher.PubKey
	BlockchainTrustSeckey cipher.SecKey
	TrustPubkeyList       []cipher.PubKey

	AgreeNum int
	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo string

	DBPath       string
	DBReadOnly   bool
	Arbitrating  bool
	RPCThreadNum uint // rpc number
	LogToFile    bool
}

func (c *Config) register() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.EnableWalletAPI, "enable-wallet-api", c.EnableWalletAPI, "Enable the wallet API")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable csrf check")
	flag.BoolVar(&c.EnableSeedAPI, "enable-seed-api", c.EnableSeedAPI, "enable /wallet/seed api")
	flag.StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "cert.pem file for web interface HTTPS. If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "key.pem file for web interface HTTPS. If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface, "enable the rpc interface")
	flag.IntVar(&c.RPCInterfacePort, "rpc-interface-port", c.RPCInterfacePort, "port to serve rpc interface on")
	flag.StringVar(&c.RPCInterfaceAddr, "rpc-interface-addr", c.RPCInterfaceAddr, "addr to serve rpc interface on")
	flag.UintVar(&c.RPCThreadNum, "rpc-thread-num", c.RPCThreadNum, "rpc thread number")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address", c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.samos)")
	flag.StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file (defaults to ~/.samos/data.db)")
	flag.BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db read-only")
	flag.StringVar(&c.ConnectTo, "connect-to", c.ConnectTo, "connect to this ip only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "Run the http profiling interface")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the html gui")

	// Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster, "run the daemon as blockchain master server")

	flag.StringVar(&BlockchainPubkeyStr, "master-public-key", BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&BlockchainSeckeyStr, "master-secret-key", BlockchainSeckeyStr, "secret key, set for master")
	flag.StringVar(&BlockchainSeckeyFile, "master-secret-file", BlockchainSeckeyFile, "encrypted secret key file, set for master")

	flag.StringVar(&GenesisAddressStr, "genesis-address", GenesisAddressStr, "genesis address")
	flag.StringVar(&GenesisSignatureStr, "genesis-signature", GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&BlockchainTrustPubkeyStr, "trust-public-key", BlockchainTrustPubkeyStr, "public key of the trust node")
	flag.StringVar(&BlockchainTrustSeckeyStr, "trust-secret-key", BlockchainTrustSeckeyStr, "secret key, set for trust node")
	flag.StringVar(&TrustAddressStr, "trust-address", TrustAddressStr, "trust node address")
	flag.StringVar(&TrustPubkeyListStr, "trust-pubkey-list", TrustPubkeyListStr, "trust pubkey list")
	flag.IntVar(&c.AgreeNum, "agreeNum", c.AgreeNum, "agree num for pbft")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.samos/wallet/")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "The maximum outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "The peer list size")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	flag.StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
}

var home = file.UserHome()

var devConfig = Config{
	// Disable peer exchange
	DisablePEX: false,
	// Don't make any outgoing connections
	DisableOutgoingConnections: false,
	// Don't allowing incoming connections
	DisableIncomingConnections: false,
	// Disables networking altogether
	DisableNetworking: false,
	// Enable wallet API
	EnableWalletAPI: false,
	// Enable seed API
	EnableSeedAPI: false,
	// Disable CSRF check in the wallet api
	DisableCSRF: false,
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly: false,
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address: "",
	//gnet uses this for TCP incoming and outgoing
	Port: 8858,
	// MaxOutgoingConnections is the maximum outgoing connections allowed.
	MaxOutgoingConnections: 16,
	DownloadPeerList:       false,
	PeerListURL:            "https://www.samos.io/blockchain/peers.txt",
	// How often to make outgoing connections, in seconds
	OutgoingConnectionsRate: time.Second * 5,
	PeerlistSize:            65535,
	// Wallet Address Version
	//AddressVersion: "test",
	// Remote web interface
	WebInterface:             true,
	WebInterfacePort:         8640,
	WebInterfaceAddr:         "127.0.0.1",
	WebInterfaceCert:         "",
	WebInterfaceKey:          "",
	WebInterfaceHTTPS:        false,
	PrintWebInterfaceAddress: false,

	RPCInterface:     true,
	RPCInterfacePort: 8650,
	RPCInterfaceAddr: "127.0.0.1",
	RPCThreadNum:     5,

	LaunchBrowser: false,
	// Data directory holds app data -- defaults to ~/.samos
	DataDirectory: filepath.Join(home, ".samos"),
	// Web GUI static resources
	GUIDirectory: "./src/gui/static/",
	// Logging
	ColorLog:        true,
	LogLevel:        "INFO",
	LogToFile:       false,
	DisablePingPong: false,

	// Wallets
	WalletDirectory:  "",
	WalletCryptoType: string(wallet.CryptoTypeScryptChacha20poly1305),

	// Timeout settings for http.Server
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 60 * time.Second,
	IdleTimeout:  120 * time.Second,

	// Centralized network configuration
	RunMaster:        false,
	BlockchainPubkey: cipher.PubKey{},
	BlockchainSeckey: cipher.SecKey{},

	BlockchainTrustPubkey: cipher.PubKey{},
	BlockchainTrustSeckey: cipher.SecKey{},
	TrustPubkeyList:       []cipher.PubKey{},
	AgreeNum:              0,

	GenesisAddress:   cipher.Address{},
	GenesisTimestamp: GenesisTimestamp,
	GenesisSignature: cipher.Sig{},

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU: false,
	// Where the file is written to
	ProfileCPUFile: "samos.prof",
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf: false,
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo: "",
}

func init() {
	applyConfigMode()
}

func applyConfigMode() {
	switch ConfigMode {
	case "":
	case "STANDALONE_CLIENT":
		devConfig.EnableWalletAPI = true
		devConfig.EnableSeedAPI = true
		devConfig.LaunchBrowser = true
		devConfig.DisableCSRF = false
		devConfig.DownloadPeerList = true
		devConfig.RPCInterface = false
		devConfig.WebInterface = true
		devConfig.LogToFile = false
		devConfig.ColorLog = true
	default:
		panic("Invalid ConfigMode")
	}
}

// Parse prepare the config
func (c *Config) Parse() {
	c.register()
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	//if c.RunMaster == true && BlockchainSeckeyStr == "" {
	//	if BlockchainSeckeyFile == "" {
	//		logger.Error("master-secret-file must not empty")
	//		os.Exit(1)
	//	}
	//	if _, err := os.Stat(BlockchainSeckeyFile); os.IsNotExist(err) {
	//		logger.Error("%s not exists", BlockchainSeckeyFile)
	//		os.Exit(1)
	//	}
	//	fmt.Printf("input password\n")
	//	key, err := gopass.GetPasswd()
	//	if err != nil {
	//		logger.Error("password input error")
	//		os.Exit(1)
	//	}
	//	encryptMsg, err := ioutil.ReadFile(BlockchainSeckeyFile)
	//	if err != nil {
	//		logger.Error("read secret file failed, %+v", err)
	//		os.Exit(1)
	//	}

	//	msg, err := encrypt.Decrypt(key, string(encryptMsg))
	//	if err != nil {
	//		logger.Error("decrypt failed, please input corrent password: error %+v", err)
	//		os.Exit(1)
	//	}

	//	BlockchainSeckeyStr = msg
	//}

	c.postProcess()
}

func (c *Config) postProcess() {
	var err error
	if GenesisSignatureStr != "" {
		c.GenesisSignature, err = cipher.SigFromHex(GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}
	if GenesisAddressStr != "" {
		c.GenesisAddress, err = cipher.DecodeBase58Address(GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}
	if BlockchainPubkeyStr != "" {
		c.BlockchainPubkey, err = cipher.PubKeyFromHex(BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey, err = cipher.SecKeyFromHex(BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		BlockchainSeckeyStr = ""
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey = cipher.SecKey{}
	}
	if BlockchainTrustPubkeyStr != "" {
		c.BlockchainTrustPubkey, err = cipher.PubKeyFromHex(BlockchainTrustPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if BlockchainTrustSeckeyStr != "" {
		c.BlockchainTrustSeckey, err = cipher.SecKeyFromHex(BlockchainTrustSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		BlockchainTrustPubkeyStr = ""
	}
	if BlockchainTrustPubkeyStr != "" {
		c.BlockchainTrustSeckey = cipher.SecKey{}
	}

	if TrustPubkeyListStr != "" {
		pubkeys := strings.Split(TrustPubkeyListStr, ",")
		sort.Strings(pubkeys)
		for _, pubkey := range pubkeys {
			trustPubkeyStr := cipher.MustPubKeyFromHex(pubkey)
			c.TrustPubkeyList = append(c.TrustPubkeyList, trustPubkeyStr)
		}
	}

	c.DataDirectory, err = file.InitDataDir(c.DataDirectory)
	panicIfError(err, "Invalid DataDirectory")

	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	if c.WalletDirectory == "" {
		c.WalletDirectory = filepath.Join(c.DataDirectory, "wallets")
	}

	if c.DBPath == "" {
		c.DBPath = filepath.Join(c.DataDirectory, "data.db")
	}

	if c.RunMaster {
		// Run in arbitrating mode if the node is master
		c.Arbitrating = true
	}

	// Don't open browser to load wallets if wallet apis are disabled.
	if c.EnableWalletAPI {
		c.GUIDirectory = file.ResolveResourceDirectory(c.GUIDirectory)
	} else {
		c.LaunchBrowser = false
	}
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func printProgramStatus() {
	p := pprof.Lookup("goroutine")
	if err := p.WriteTo(os.Stdout, 2); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)
	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func catchInterrupt(quit chan<- struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	close(quit)

	// If ctrl-c is called again, panic so that the program state can be examined.
	// Ctrl-c would be called again if program shutdown was stuck.
	go catchInterruptPanic()
}

// catchInterruptPanic catches os.Interrupt and panics
func catchInterruptPanic() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	printProgramStatus()
	panic("SIGINT")
}

func createGUI(c *Config, d *daemon.Daemon, host string, quit chan struct{}) (*gui.Server, error) {
	var s *gui.Server
	var err error

	config := gui.Config{
		StaticDir:       c.GUIDirectory,
		DisableCSRF:     c.DisableCSRF,
		EnableWalletAPI: c.EnableWalletAPI,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		IdleTimeout:     c.IdleTimeout,
	}

	if c.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		if err := cert.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Samosd"); err != nil {
			logger.Errorf("gui.CreateCertIfNotExists failure: %v", err)
			return nil, err
		}

		s, err = gui.CreateHTTPS(host, config, d, c.WebInterfaceCert, c.WebInterfaceKey)
	} else {
		s, err = gui.Create(host, config, d)
	}
	if err != nil {
		logger.Errorf("Failed to start web GUI: %v", err)
		return nil, err
	}

	return s, nil
}

func initLogFile(dataDir string) (*os.File, error) {
	logDir := filepath.Join(dataDir, "logs")
	if err := createDirIfNotExist(logDir); err != nil {
		logger.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
		return nil, fmt.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
	}

	// open log file
	tf := "2006-01-02-030405"
	logfile := filepath.Join(logDir, fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), Version))

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		logger.Errorf("os.OpenFile(%s) failed: %v", logfile, err)
		return nil, err
	}

	hook := logging.NewWriteHook(f)
	logging.AddHook(hook)

	return f, nil
}

func initProfiling(httpProf, profileCPU bool, profileCPUFile string) {
	if profileCPU {
		f, err := os.Create(profileCPUFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if httpProf {
		go func() {
			log.Println(http.ListenAndServe("localhost:8989", nil))
		}()
	}
}

func configureDaemon(c *Config) daemon.Config {
	//cipher.SetAddressVersion(c.AddressVersion)
	dc := daemon.NewConfig()
	dc.Pex.DataDirectory = c.DataDirectory
	dc.Pex.Disabled = c.DisablePEX
	dc.Pex.Max = c.PeerlistSize
	dc.Pex.DownloadPeerList = c.DownloadPeerList
	dc.Pex.PeerListURL = c.PeerListURL
	dc.Daemon.DisableOutgoingConnections = c.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.DisableNetworking
	dc.Daemon.Port = c.Port
	dc.Daemon.Address = c.Address
	dc.Daemon.LocalhostOnly = c.LocalhostOnly
	dc.Daemon.OutgoingMax = c.MaxOutgoingConnections
	dc.Daemon.DataDirectory = c.DataDirectory
	dc.Daemon.LogPings = !c.DisablePingPong

	if c.OutgoingConnectionsRate == 0 {
		c.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.OutgoingConnectionsRate
	dc.Visor.Config.IsMaster = c.RunMaster

	dc.Visor.Config.BlockchainPubkey = c.BlockchainPubkey
	dc.Visor.Config.BlockchainSeckey = c.BlockchainSeckey

	dc.Visor.Config.BlockchainTrustPubkey = c.BlockchainTrustPubkey
	dc.Visor.Config.BlockchainTrustSeckey = c.BlockchainTrustSeckey
	dc.Visor.Config.TrustPubkeyList = c.TrustPubkeyList
	dc.Visor.Config.AgreeNum = c.AgreeNum

	dc.Visor.Config.GenesisAddress = c.GenesisAddress
	dc.Visor.Config.GenesisSignature = c.GenesisSignature
	dc.Visor.Config.GenesisTimestamp = c.GenesisTimestamp
	dc.Visor.Config.GenesisCoinVolume = GenesisCoinVolume
	dc.Visor.Config.DBPath = c.DBPath
	dc.Visor.Config.DBReadOnly = c.DBReadOnly
	dc.Visor.Config.Arbitrating = c.Arbitrating
	dc.Visor.Config.EnableWalletAPI = c.EnableWalletAPI
	dc.Visor.Config.WalletDirectory = c.WalletDirectory
	dc.Visor.Config.BuildInfo = visor.BuildInfo{
		Version: Version,
		Commit:  Commit,
		Branch:  Branch,
	}
	dc.Visor.Config.EnableSeedAPI = c.EnableSeedAPI

	dc.Gateway.EnableWalletAPI = c.EnableWalletAPI

	// Initialize wallet default crypto type
	cryptoType, err := wallet.CryptoTypeFromString(c.WalletCryptoType)
	if err != nil {
		log.Panic(err)
	}

	dc.Visor.Config.WalletCryptoType = cryptoType

	return dc
}

// Run starts the samos node
func Run(c *Config) {
	defer func() {
		// try catch panic in main thread
		if r := recover(); r != nil {
			logger.Errorf("recover: %v\nstack:%v", r, string(debug.Stack()))
		}
	}()

	logLevel, err := logging.LevelFromString(c.LogLevel)
	if err != nil {
		logger.Error("Invalid -log-level:", err)
		return
	}

	logging.SetLevel(logLevel)

	if c.ColorLog {
		logging.EnableColors()
	} else {
		logging.DisableColors()
	}

	var logFile *os.File
	if c.LogToFile {
		var err error
		logFile, err = initLogFile(c.DataDirectory)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	scheme := "http"
	if c.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical().Infof("Full address: %s", fullAddress)
	if c.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
	}

	initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)

	var wg sync.WaitGroup

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		catchInterrupt(quit)
	}()

	// Watch for SIGUSR1
	wg.Add(1)
	func() {
		defer wg.Done()
		go catchDebug()
	}()

	// creates blockchain instance
	dconf := configureDaemon(c)

	logger.Infof("Opening database %s", dconf.Visor.Config.DBPath)
	db, err := visor.OpenDB(dconf.Visor.Config.DBPath, dconf.Visor.Config.DBReadOnly)
	if err != nil {
		logger.Errorf("Database failed to open: %v. Is another samos instance running?", err)
		return
	}

	d, err := daemon.NewDaemon(dconf, db, DefaultConnections)
	if err != nil {
		logger.Error(err)
		return
	}

	var rpc *webrpc.WebRPC
	if c.RPCInterface {
		rpcAddr := fmt.Sprintf("%v:%v", c.RPCInterfaceAddr, c.RPCInterfacePort)
		rpc, err = webrpc.New(rpcAddr, webrpc.Config{
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			IdleTimeout:  c.IdleTimeout,
			ChanBuffSize: 1000,
			WorkerNum:    c.RPCThreadNum,
		}, d.Gateway)
		if err != nil {
			logger.Error(err)
			return
		}
		rpc.ChanBuffSize = 1000
		rpc.WorkerNum = c.RPCThreadNum
	}

	var webInterface *gui.Server
	if c.WebInterface {
		webInterface, err = createGUI(c, d, host, quit)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	// Debug only - forces connection on start.  Violates thread safety.
	if c.ConnectTo != "" {
		if err := d.Pool.Pool.Connect(c.ConnectTo); err != nil {
			logger.Errorf("Force connect %s failed, %v", c.ConnectTo, err)
			return
		}
	}

	errC := make(chan error, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := d.Run(); err != nil {
			logger.Error(err)
			errC <- err
		}
	}()

	// start the webrpc
	if c.RPCInterface {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := rpc.Run(); err != nil {
				logger.Error(err)
				errC <- err
			}
		}()
	}

	if c.WebInterface {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := webInterface.Serve(); err != nil {
				logger.Error(err)
				errC <- err
			}
		}()

		if c.LaunchBrowser {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Wait a moment just to make sure the http interface is up
				time.Sleep(time.Millisecond * 100)

				logger.Infof("Launching System Browser with %s", fullAddress)
				if err := browser.Open(fullAddress); err != nil {
					logger.Error(err)
					return
				}
			}()
		}
	}

	/*
		time.Sleep(5)
		tx := InitTransaction()
		_, _, err = d.Visor.V.InjectTransaction(tx)
		if err != nil {
			log.Panic(err)
		}
	*/

	select {
	case <-quit:
	case err := <-errC:
		logger.Error(err)
	}

	logger.Info("Shutting down...")
	if rpc != nil {
		rpc.Shutdown()
	}
	if webInterface != nil {
		webInterface.Shutdown()
	}
	d.Shutdown()
	wg.Wait()

	logger.Info("Goodbye")

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file")
		}
	}
}

func main() {
	devConfig.Parse()
	Run(&devConfig)
}

// InitTransaction creates the initialize transaction
func InitTransaction() coin.Transaction {
	var tx coin.Transaction

	output := cipher.MustSHA256FromHex("c8d02e832715e862f9c6ac7e763151483bb990e916cf5a75defa0d8029c23399")
	tx.PushInput(output)

	addrs := visor.GetDistributionAddresses()

	if len(addrs) != 300 {
		log.Panic("Should have 300 distribution addresses")
	}

	// 1 million per address, measured in droplets
	if visor.DistributionAddressInitialBalance != 1e6 {
		log.Panic("visor.DistributionAddressInitialBalance expected to be 1e6*1e6")
	}

	for i := range addrs {
		addr := cipher.MustDecodeBase58Address(addrs[i])
		tx.PushOutput(addr, visor.DistributionAddressInitialBalance*1e6, 1)
	}
	seckeys := make([]cipher.SecKey, 1)
	seckey := ""
	seckeys[0] = cipher.MustSecKeyFromHex(seckey)
	tx.SignInputs(seckeys)
	/*
	 */

	tx.UpdateHeader()

	err := tx.Verify()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("signature= %s", tx.Sigs[0].Hex())
	return tx
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0777)
}
