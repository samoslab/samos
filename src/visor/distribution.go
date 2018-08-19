package visor

import (
	"github.com/samoslab/samos/src/coin"
)

const (
	// MaxCoinSupply is the maximum supply of skycoins
	MaxCoinSupply uint64 = 3e8 // 300,000,000 million

	// DistributionAddressesTotal is the number of distribution addresses
	DistributionAddressesTotal uint64 = 300

	// DistributionAddressInitialBalance is the initial balance of each distribution address
	DistributionAddressInitialBalance uint64 = MaxCoinSupply / DistributionAddressesTotal

	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64 = 30

	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 = 5

	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 = 60 * 60 * 24 * 365 // 1 year
)

func init() {
	if MaxCoinSupply%DistributionAddressesTotal != 0 {
		panic("MaxCoinSupply should be perfectly divisible by DistributionAddressesTotal")
	}
}

// GetDistributionAddresses returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 300 addresses.
func GetDistributionAddresses() []string {
	addrs := make([]string, len(distributionAddresses))
	for i := range distributionAddresses {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetUnlockedDistributionAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddresses() []string {
	// The first InitialUnlockedCount (30) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (30) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (30) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, InitialUnlockedCount)
	for i := range distributionAddresses[:InitialUnlockedCount] {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetLockedDistributionAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddresses[InitialUnlockedCount:] {
		addrs[i] = distributionAddresses[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}

// TransactionIsLocked returns true if the transaction spends locked outputs
func TransactionIsLocked(inUxs coin.UxArray) bool {
	lockedAddrs := GetLockedDistributionAddresses()
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true
		}
	}

	return false
}

var distributionAddresses = [DistributionAddressesTotal]string{
	"bQcvGbswSUtTBUNg1VRqWaPcJxWUmFMBp5",
	"Eawai1bY2UwdXkf5JEmMBTNxBcLc519JeJ",
	"z7qaCXUwTUrRFbaUehTrYCM5RAJyxWMJue",
	"2PBcFVpoHTkjaABdTjZEcAhvAtY8fPEvtGJ",
	"PWcvJxdnVqSvqJVPX7LRJib4Wym6JEzfzd",
	"QyL5ry5xbGjMNC8SeQ8Dhjs3XYQhXS2yco",
	"NVnUzah5NSEVvcizzKGpPouP55UXVrBhMa",
	"2Lp7SwycXHs9kqCvYc2ewx4kFAukL15RB7j",
	"2AQh6i6pjkMJLAeKxrFxLQSaogUvVZz4Yno",
	"wKpp5vUH6NoCHV7TczgfCwuocD9Y9CpNi2",
	"jf1TDFdQi9H65WAT74njoZ4gL7TqsZaED3",
	"2Mc91tWJvYYqMXY7kuy92PZzKoJyKufH5MK",
	"2dmqyymsbbKZZykPNVkGMCDNYDe4qTbjiCY",
	"ZtdXaB4ehb8kzz9FhaQdnfMZvG9Sy2tiqi",
	"3DLb9Q9SdSdq7t7Hrot5MW2zGRLSeHviLg",
	"YEdwpypjuZNEEnyh27X6gSiQ4cqcPRq3jW",
	"JBoD1pqUoiWCdPte93dKTbdEcE33zn493S",
	"2QZphW5UUzqaR1UfXWb8UvomCN3SXQRPqhM",
	"2JytUt5SqnTScf8d6HDhkSJ488bVKGTMQYd",
	"FQc4K8oAjxmyxU36vrupm8opoGrPrSBLH1",
	"294ThuQukcphF3X8FPJigoLwL4nhruKQFQj",
	"9wVXYPe4zJeZx2kxiLyjDm7M9ZA3kzTBxC",
	"VALdSUJNAaJqFhii1dYGcCL6xQh17zcsSo",
	"2e2koT7rXwF9hnvXywbafgwih57sqA1nkaY",
	"2XKxooUN24r6mNbzFUKDq1UgqCf9jSZauFt",
	"yNq9QpAYnY7TfTSKvtCqbkaytQBvUmtaDD",
	"AmjnGiVR4MZL8qBB3bRusibLt4AdRLiMkS",
	"24pHydBTAfahezspeUEFUwF4FEYcgprongH",
	"81XAcUoC4ZX9ZXCR8fXmfDpuSHpbBu1ppe",
	"2KrEVRkeqJsZvLuZjNw8YFmbakGqh6yoA39",
	"2M8gqHXwi33Jbd1QHmjHt859UvviXxZTkRJ",
	"7yiMcmgmRU64jFvKbqnhToSER38FHPbXNt",
	"4viJnDB4aUxfxuYNEtciHya3sEnH46zxsj",
	"G75VkqvsGqiVEkhkLKJdNUDgbXQUMyNSgk",
	"2BTqXpfcLbs7KNnudDPr3zCTM3Zy8Qxi58R",
	"X1rRqTU7z46snoRTtcjpjGy2XeduA7Km5J",
	"25Zp5MgwRRxeABzPr21znb9s11buWdLtF5x",
	"2SPs6cq9SoshCjVDeYsv1oc5qHpio1tVqt1",
	"2SRBdHDV7Esadew5K7zN875KjdwT2DjLAWc",
	"NcPwCaXLEfw4LxbUZEz3Xmf8m9ZyCDYepa",
	"W3xufLkNkXEnEy4MdtjMgRo2HB9k7hv5EP",
	"VpXSXhvjE8S2HcxWrtCrbUan444JqBg98e",
	"gHSTiwo93E7uA3NsptgGder3SST6LN7byd",
	"2TGZdF9vhUXHUh3zMm5eJz4L7b43KCKFYKw",
	"2FKkW4oHC6ds26T449LdfEAja6B2GkMe25m",
	"2Sqkh5hWM5bAZrSPbZkD3pNxy8kxh8jdrr4",
	"2iKqwe5UvQkd1JoMqoHqjz9pWbZyEhSGQR4",
	"27hswtovejQuEG1WeoqkpQ5WEzBTkWZgTba",
	"iUEjqxYPBRUEgivi8h6Nqrnrae3CPy9eWu",
	"mYGBZtNV63B1A2X38UWbBvWCh8pLSA6vZm",
	"2DmGPuUgESyBFsPTjjE8R745tq3zi12mKtV",
	"2EJpw48fE4UunX3r3AsPvAHx52G8QkQbrFZ",
	"D8JscwuokksNX9n6SCJuyZNRp47johxRPn",
	"REtxJ11iYP2CsHzNYefu1wiKwiZ4Zi39TL",
	"avfFxtcNCGNVjDsJcLvSNVycDm1ee6THvi",
	"rwWntiq57sKd4kFtfR4xaJhgyetmMzwVBV",
	"2Ss4jKrE8PSsZ28L4Eghqt52xvVnUv4zuah",
	"PMVx5Qhh3pZCknke1XfMyYtazc61nTtGha",
	"2aDiqSeNGEyiCRHGVdvMTRTT86doJ5pMp7D",
	"2RtmXqUJMFLk8ci3iU7pozPW4XhX1rtJYxz",
	"HrApr6a79QGhdL24m8MoDDnHwbVpbyEzyk",
	"2HKPvVZw12vfE1V1yoSm4YcYgjEuUxUdTNX",
	"P2fh7HurKvCRvHvw5bSRFvKcvmA5QxCGL5",
	"Y26x3GPHGjgvVJSRWotgeRahRkwQ6bY7a8",
	"cfaaXZA4fgexmQuYq1pAey21mbaa1sufr6",
	"gJt1ZqGgsZzz6rc3qZbX1mM1QNG38ZrGLZ",
	"FVcXybyHdEXtCJShHiCzW9NbCGpAQqev94",
	"6wfHSyEBdHzPPSkbNvXuBpDNBrxZ3k6Eb7",
	"pvTc9n2VRV97Ch9BT6qD4DJvDukhLdVi32",
	"2YW5xhwMwUA8nTDz6U9SgFY8CPsRxCTtqCP",
	"kC1i1cL1yLMLFZHKwVoBPQU2Z3d6EZfv6e",
	"2emfaof8DjfmUw6tdhbqfNVNRp1AkUM2iQh",
	"2Ni8NYtUmCzL2ck1JcufdxJK9wKCNEDTLBP",
	"3JbQX3kn8DeMjCZpPtZLbgmgVv2qBdbqya",
	"2HZuSogp2pAmi67w2kVBdHiRPvzoNUNKSTX",
	"ZXMbf6vBzZKUxbmFQyiUzDxQJidxH9H9Kn",
	"Mj2wtEyXA8zNDpF8ADQPYqfjCuQvLNS6iM",
	"trrV8GQqM7JagWdiAYwCEf6oQab6aKYgNz",
	"vW14mvW2SCBuhmEnzicxYoEzTrd8fJnSYy",
	"nrv8pT2nmKPaQfZYgbY5aPVFvBfiF93Jf5",
	"J9B5haDDtFXwEdHShUGT8hNE6bWbThzmTw",
	"cpaL3mkFE4crqA9b3RHYVCpQomsMqKq5sW",
	"Br31GN6X5wAK7mHnCrAV5JoB1LyNLNoZTi",
	"2YoMBdXkrvCJwJi4No3yMiDStehmDfXXGDE",
	"2mwFMPtjrhJGRCJDEtAjNzA6Br7vJe3xBWs",
	"89rd1VkngUdbwKDqapLD8vexWwLfYPjZNm",
	"2M5PcRUm9NMUwmy28bCaEYop4pwEQWirgHQ",
	"2kpDUgxsF61fxaDVjbUykZfDhAzYgrPhpKW",
	"2Z3FJ8b11WZCL6CJ3FpwrZvskBEjL7Q1qjk",
	"o96VUmp1eoPQAUB1U8FAzhAKea8nn8Wc6q",
	"wnWYmXBqJ6A3FoHcFtHr9bq5KA6wnTcTt2",
	"22CnEk2bxg9SM5DXjfhAiAXEdUJ62vDJ7pk",
	"HTPRrvfPrt3Ki3cbwq9my6jzD1v9b7nQpM",
	"Boyikw2M6aMm557CzbuBUfnYtLyCKBjYXr",
	"2ZeTyvB2dzGRdjhbhVuJavtFqwKFtTkN3Bw",
	"eCWF1dGx8wAAh5KtMh4aggoUNGWHFZ85MA",
	"27ZPHRi39iM28LwDBUMMuuimyb1Y5xLUtrA",
	"XHKudj6n5x98qPjd44AfeQJAm52VqMD228",
	"2XH6j58XVZa489qrNG5qAaF7UQMWJbxVJbh",
	"MhERK1Kda2HaNQY4JwmbHtvpkdUT9gavov",
	"zDHRgxqpo3meDTn8xktyZBCyfbjqipfZuq",
	"2E283J2uqeP6y2fFLTLs3FGfKdJBraPFhu6",
	"dxVWumd875FmF7dHvNsBzJidfXdCk8Ya1A",
	"2Fp5RAE8F35jSikHvPUCf762xKvf3jnVhSs",
	"hw5jxKwSxHXereEB53J53rsRNCz1xYcCs7",
	"fKBN5gXwR267X45GDZxZQXb1WPY3QT9qrc",
	"iiRwcCZiepD796STcmYuCw8XF54ZxvT5SS",
	"29mWjpBuKubDk5yf6JkA3469MK3sFmytkDZ",
	"WUm8r15oZ96uwCNv12AMbsDuKeBoH6tes7",
	"v26z72eKiJ8f1nMAH8zWP2UEp7nMvzNMBP",
	"22dMQPdZd18wZyPkCHeVqFAQhpxiwTVi5rJ",
	"2ajgnvGqYPZv2iPFXTQFPBt3L1LuxgJWe5k",
	"vLZizutLpML1BuQ1CyK92QXkAiHq8NoLTJ",
	"2L83y8pZvTcRAH66sSTtoZ3eAk2npvAAm5j",
	"GbpJN3XGs6CJBSLyBC9qsNxhKac9Z2GpBT",
	"2YVRS7pBtHW8WJjxenpHrQrWXjgq8UiDaWP",
	"2a7bitkkV1qW1uPUdjajRw35w35hsP4hoTz",
	"YbeL3BhZ2atdKWSriuP9o8mpBsSbrtMJLj",
	"2RVDUUH1xUVSP2EBPHADNXtYXoxKJGjrXww",
	"Tt27DTYmZYusZso57DJBNs74z5rArFNwbB",
	"4SWpyVUrTTLy22fbTxCNxKJeUpj5LRaYC7",
	"2QiCcte9kfafF9xD7YAKS4hjEoSMxs4krse",
	"MBbQvfeP19KqM84UbzMmcgSKRUufyvLVCS",
	"jZPuaX4xfkfwUnLcsJYdmeFoBPSxH2x96r",
	"yjygjimDYbwPLV91vQGBuSDBP8jETVY8XY",
	"L4X85AzDGwA7wudb1LEFudZXq8ogdEg8ts",
	"2C2Vjzi8reK5AmqgSaB2AH3MS1UTN93PqcT",
	"2BxQQ5XCCTtv3ZyCHXLVNmVzE3eztbcqQgo",
	"28YFHWdHpQF7ShvXNSHfkuwrPhTqJZFkTAv",
	"2Qn5DKwvCPxg26wg9Eq22pNd3YoZLYpLXSX",
	"21uGAQgbPc9vDXC5dfJaTV99QypTSt2qgRx",
	"bGTQwpCSR2xeMSQFCqcM97mDRmjPTseEdy",
	"24JC9mP8p34oV2YPzX3dPg6sP1GKoHNNyas",
	"215ULXCRpwz6MeCFjsTVmN3obiswxShqhiZ",
	"p4RKiTU1EnxZjbTKYoU3V6nwj59uJu5GCJ",
	"jLfKYN1Kmbm3eLTwUjarZ6MN7xE2N6TVNW",
	"qgRx9ocghmgkP9uezMCC9x5C4F5giCgDW7",
	"tQK8VvcTYeHbScwSiQpYZg6ZkBhiXFtD9U",
	"kfJ5Z2aHQdPMYF859egn5taoXZr6QQL9Ht",
	"cTrp6MTTKYkTT5xR6iQ1mDeiy6npt6SxYe",
	"2YUEN18Wg3hVbJjQGXwKrGfnHd18sAPqxZ5",
	"2YydmShr1361V5nmSGb91nXJFRv2qH3CA3Z",
	"vvCTD2VNvVMQBUR45VXBfmrQ6vs7Qrkh9e",
	"225ndLmWnZh9fNRLjxUaHdNcB3AitUaNLZx",
	"2B5coDWoZRfcs5Z3Fqtnj4MhGUajph2zFrQ",
	"2VN6R1hpJAnDhtqMjCjSpzNkr2UJVDPdGX4",
	"tdT4dBXQW7rXFwCDv1QJ8XV3dfyQmVkpfW",
	"2Fz9MDVMqMVyrPV2cZ1gUBQc9SZqCeNC7V6",
	"2Bs6otwm4TeyKgEGyPLR6M23JsnXCCxTYx1",
	"23Q1EuSfLkb1D5ryh6NByUrLkR8WZhES3Xw",
	"2WdzzbRrnw9rnqMaJKCUXQTUqoKCaS3tqJ3",
	"pSd77tuaA39FzQ36V2BJ3Bgig3couWXjNq",
	"gPQEMJq9zJNNiKQdQ6mBPqNYdBSVhaJWp6",
	"2jh7pcXyUgNBWXNy6fa37xhaevvGeEV7fSV",
	"2N8Pp432kgTxGRNYprmB4sg6w7Gizd1x1Cg",
	"zxygRkgXWL43F13unrWEiF9vDYGCcXnc6p",
	"2kx6wmiy6PrUeG6kEgaJJVdLL2LWabC8kKa",
	"2DPtgJBhFkZ8NDqkEeZK13RTSNVHxVf4dRR",
	"AMFAKV3CdTTiMCH1Tou9Csyyz84iXFEZhV",
	"2dkzykxRRyQKWZ35B65q9ZpZnxy9xcpqsqe",
	"2VceeGTS2daH4esD3uKJoM5NpwcrctnFsij",
	"28TnzmwSGJc876qJJoSbkCUftF9m3C47GQH",
	"9qBd6KkCDYRnXXNs7uqMq48JpBwbnmNsGb",
	"27LPXNquQUJkFUXjHoNBb62hERL2J2HXPv6",
	"QZxfxDwrRWLKtMHsLiHDXKTYCbxmj6mmYC",
	"2fBiNZWcyxXkkgC5S3zSQ1npvkKPfQxAJd1",
	"VKHXedxzeqHtAxpnMwhdxK3RgqcFX56HJX",
	"2YPyNgFzi9Y1xzjEVUEn2agDV2Y8FRdD27d",
	"2fXw8nGiLHjgEzfgBLY89YzEaeFQmLTU1pS",
	"22EKziZsC8K9q1NsRqVYSNJMdwo6wBqCER2",
	"aPzGykEZUJSuqdYFtUYhAFVRkdDCDpWVbm",
	"zZNK8R5YNunXns8UciVmSCjpZrpMTHopuF",
	"24W5K9Q9QP69xW1qaQPG87hHKDUbYeYZudo",
	"DdBACnaRjbXxHxg3vQbJv4rjUsxvk1MymD",
	"tUhqGAYaByWKVMEhBPUDrJMUAzCD8S7oY6",
	"XbPT1zboEzVTxvnCE7xmMFedEG6bdqHtB7",
	"sK1C38ZEA2FZZde93RfJYuxfzWPfZxnmAX",
	"2MJgQMtVXpXvWad6nYHwWzit3hQhh45YyB1",
	"2t7xGgZnHWtdSSqj7aPt8uFSZYSUYD26yd",
	"2Uu1Nso9RETH8vYuaVyLdSAX8cmsys4f78Z",
	"2UGS3pdnGFbjeQeNCfdPByuPSEb2epJJaq1",
	"36uvK8tvozvLtVih93DsuAwFwLqJBK5f5B",
	"2UKNLPVogsEoEcP9UiUQYKkdtzNCEQeg4J8",
	"2k2Zg4TNDmaz15RPwRpCoXdFk5oZVMtJDW",
	"2PM1S6uibybNdY7VhqaD9hBnJbH8FDPsrXL",
	"2bvrRxmgm7oNv2kPsQKjdkFDoD6PavzB6bD",
	"UbzMzyTUUpQnnc48rvhhstwFNDQZ95pMgr",
	"VT8QJVNdNsWAbH65sLm2fYaLKhkk6VeEnr",
	"2jyeQGgkHNVo3roLEJTPZhmmLD4eW5mSAZy",
	"26Wz1NMkNYUSMM2Hz9jbV4VriACv3QbXUmA",
	"bdjvoDVAEi78fkV8i3ZnmhtsSnrKq7UF1m",
	"aq1uV8kZgoikyu4HR2pTq2yzyjePkzEMrz",
	"DonejqcKQuSVQrJw8a1xpVqqq8NuSpjcKY",
	"2X8jRfh82kQgtN8rZBWop1Qv1BdeGAVcUXn",
	"VEAerQ3rkS9dU5oPcQkfY7wmGteqvCQqj7",
	"y5hdeXUncjVxoLwUa5BJSPyciEABjPQqQw",
	"SJKprJry28zs4wn5LoS7cunoiWtgZRbu2w",
	"DocWD4m65iVrYxb156Jpw5pgP1HWwBNR45",
	"2SE7yfRQzXkjzccR7SyGputh11f4W8w1EWz",
	"2Twmvom3q7krXGkEbKrcubu7tfVAzgqf2Hm",
	"v33Tx3cd1Sf4gssvU9t77YMHeaW5RY27iQ",
	"LpQuDcfZzywm333jNzVXjE7LvhCrF8WuXH",
	"28xEe17B1HzFuCEYrEEk4kycDMMLL91veGS",
	"EJGvBvCjy77ivTRk4YwuFzz6TevW8NiBYS",
	"2gvrYZBvU1XNo7KmRvAAEXqWrg18UAzphKJ",
	"sJn7BYvYJKRYKC2DXgATPiqbF1qTe2ewj",
	"2VZrNSiMXJHJ2YsCFCo712yNzbXhzwprxAy",
	"2WU38RhC3BGjeqZHRL2ipfoeP5jZEWfyUYs",
	"nHpAJjtPPdaErraabve4pY1uyuf95p1g8z",
	"22TA7Y8t993gCDx2CeywwYSejnj7fT6ur8h",
	"xFx89hXyBmKr4iaFMHQ6dXVWdvQn4gdi5q",
	"p2Utmk9xvuCuoW4U31ppjTMg3mVxsh1iXT",
	"2iR5SteGvg56Eec2ZCgSErBugz3ar1Wyz7K",
	"3F184SSmcs57qoeqYXekniHEjMVEj5iVE5",
	"2DNpVTnxZXdzvz1jyQ8t9EzXN8KhyeMZjGv",
	"94E31JiSiHg3sG3zqYkZU8rkrdYi4LN6YA",
	"mbLuMe54RCRsjeguEdeFRYipyu6CFTtnxw",
	"2eDb3hvwzdMHAvN2i7tvWzeKHM554ssEq3v",
	"2gL9RHa6HkPFArxCEZ4EvyV6rGmiFeMwC3Y",
	"24WhaG95msroj8PYzCcLQ2rdQgZAwsKQCjn",
	"2kxF7RMEo4R9DVGZDnEHPR34pTbRG5GooSA",
	"v8pQvgVFhW6JAvbzf9pC9yQrmEXAnojBGf",
	"2BSHqSTrEwMiLBHM39FBkAbozXVJeUSuUx4",
	"7UkP6r3W3LVJ97MCioEcBdXeyS7XgmoXoo",
	"kbQ3QTQXucD3WXCCKvR4ppdGgtr885oDDL",
	"2SPrK4xRkryRgwyispQCrDCgLzwDNUrhf21",
	"yKULXtXRZTF8uh5qo7GPfbsyP5Qzyo3yXx",
	"2DEmQyDghEK4ko2U1zUb95VeT15Y4gVqnfN",
	"umtZwgY4uZxquM3xZp8jcSqgUvycpeLK5k",
	"SobYAc8vDUYNXU7edanb33JCAJLPYdhHZm",
	"fid71H9w2MhJjfqDNhxZWaVXQjKpLFvhsA",
	"L94pCcGH4wyQkYUm7oJUJEtLyQ3Z65VLm3",
	"2jAQq1pCS2ygoLKc5CKogSqStVfRCRCK4Vu",
	"nctih8tKUsF4QL71whSR5nP97PyRR1FrHm",
	"dUJw4egMU7fku8ZwsS78AMpBYwCjXpo4Hb",
	"vggjFbBzReq52KLeZh3gVZrT9rVqnvmU5h",
	"2DRFS2hztFnzfRSvXrguxefCMTyMBV98WhW",
	"2CLDXK75qASC99Hd2jE5xPH1RJkzrdsLEJw",
	"245QNvLR71Wg4Dm4Wbgm667Btg4avh8qNWz",
	"wGfPyGujTZ6CjrQLRFshy2xxAoEkZCDL8n",
	"2Y95NYMiaidiNEZ3qiLCWiDbVhmUy9zKN14",
	"2D2HpMU2daYxeCS6tg9hfaK9cgykUyA51Rt",
	"24uGt257QxYVfKTHhM1PCuNEsrfoyykoGzp",
	"uUraLYiLFf7sRR7UPdyDAXbb4HebXphW8T",
	"T5NK4s4Nk9BFCH7iDwwMn371B8nZwRyYLJ",
	"21L6E4vREqrzH5wBE2vCZi1T491n5KU4sQD",
	"LGw33n4iWnsNmWcC7quK1aL3G2aVuaFCC3",
	"91iT8egc5zxerrc8ea3VMstpeEH2YNBPQH",
	"PY1soViFqMBxpG8jFyz6bi6GMGpxdP8Lau",
	"FrRsZHvfA2VQPLz5E8sEdSav6Lti7U3ubK",
	"YjMotcSHusH4gdXWLrB9Y5R5gNCLtZscE6",
	"2gt1TRMSWDjNRKyGFnVkCurvjq2Binx8mg2",
	"qKbj1DJBdQTpW19XTuKAiMSzPqyUn9M2hf",
	"2LbkT2eHBFXtSZ8E4khfm5mSP9McJNZtkyi",
	"KDPpqXP8DcsfmH2P1CppS4czFgK1tfoyAW",
	"XmBfpC2kTDjsWR2iaHeYjKGaJX63YKYrGx",
	"2TV17eBzkiDz1hYt5FPPfdaPxnT41WFx6x7",
	"nAcd6pq9pBpEhDdcZzLh7wtPZEcMitXoXs",
	"2WN9pQ8BU4PzpVdpeYmcFBkWooHEJMjY3jU",
	"2W5DovWCjaG3M4JUJN1fCxSnNZFLshnGNYD",
	"2FtnzSp9KqeddZshnXaeagShVFRpq9RMv4Y",
	"2kNoRFnBfzxufn7BEtb4uA6kCh8KFtSbZD6",
	"Ug76ak76n5ecvUBZD5ZxdNNqPgE7WPqMUc",
	"o1K9QSukgDznRrjvHvWLoLit2uz8cfkjXu",
	"L9ncyMRYC4MDyEQhiacR7e2mGD7GXAzK6T",
	"Cx1TQTrST7oRB7wdUKhJVRKH4QwiA1zxqH",
	"2SdAMr7d5u3V4VhnSkzsokEWAM7uwvNfGxK",
	"2ZdW6HjEvKkXnziNf5SZhMpH6kBCesgCxu1",
	"2B7NZfoaNXU47jZjkgbzzQkVEgJCGP86hbt",
	"ukLJ3AdS2wr5BYsush2yta5AVRrz1Z3Zd2",
	"P8TZZakUxUDTmJafu1hEqbgaeyu5mpy5ZS",
	"2NcaMRPvQdJYtANaNnraDS6niG3BtMR5ZVT",
	"24hvwkZexTLygRbKdjHxED2QZn4WY7osjDE",
	"mHUuUbCkQPwoXikzt2sPq36bmihH11VdeL",
	"Q8jPYScREy2v11yJktjufNq8MbMKDy5xaE",
	"NGkktTZbi25cp8sTWqerchRemeHfGR7Z9B",
	"3Vu9RzgCYuL86prnhBSr6NDokdkXyLuxfF",
	"2e55vdTjtJSJwXoRSEvDGw9AJSLjyAaVidS",
	"e1tMTijpmGAEUHMu8X71rPGZ2vgDWyJbLk",
	"224mTv4W5rK7FdTTbdjVE2cdebWHMX76CtS",
	"2TMbaCoqXDktuTAHehdqXykaYm86ZDqdufY",
	"Eh3CeYsQWTosm6t7JagVXZ7LT5it7zYvwx",
	"iBSXZBCDSLroZp61bbFoFqLUjJzULJkhNu",
	"v2VVbZLhVERK578GrjUuG7W7ZhKY9n5M4G",
	"HgX6NJxDuebtf9CVgWxmBhBwK5SxiqdnU2",
	"JujXaNEnpSXVvXRTxTuDFCYY9FZSAp1CXA",
	"6SCmRdmfyXdRWUs5UXoxESmGDX6oQn6ZRc",
	"WP4Hpbn3eDP8rujh4tCKY3Um7zC9ByfFr5",
	"G3He6Wudtrcrc9wEdjpd3u6SW36zYcXNHZ",
	"LvcF8mR5HToXw87hfhFuormLAM6mAKu551",
	"S4fYiLgAmv8fX9C313KLZVFwdqRK3LeJVw",
	"26irRk4EDEcS8GZAviEMi13HHLVz99M2xUs",
	"rzQUGJpjZKW9cxxjrRrbss3w5sPFqVwKLa",
	"KT6NaUTakswMzPUMrp4QsfB8tZKw7JBQZ6",
	"2C4A2Hy6d4ux3FHrMMxamyLk1hvqLfS5eNt",
	"h83fcrewvgbHos9LAty4C1ZbhamfwMe2Kp",
	"czz6ZJKKnBEZcwop14ro9dkxhSeG5W5ky7",
	"Z7M6VJZucWM72tqwQ4RoMHUMnRuSLJ9vSy",
	"RM8gwixF9uiKN2zUEHYCuHjaGah2tZbu3N",
	"MrnECsYbBcp8VBBGZ55zGTEXQcvT9izM1K",
}
