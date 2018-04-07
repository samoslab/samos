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
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

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
	"b9JKtor5PyDJTjSogY7rqeDhMqtwAwjXuq",
	"2Z494SmHNKrhiQpBw3QtgaNovpWMG1jCBcZ",
	"2ezdii8WLXJRduvyjCqrCBANKKqeJqaJBgc",
	"2PErMyZYYtQUEjRmaasrqPPGePpVD2QBtCd",
	"YtYHHiWkHHZhqACoPRn4xjvm6YMRUUHSmA",
	"4kRYRsV8cuTvrV1dJzqshVYBRtjCh3p4Qh",
	"2MjLV4uggqLCkpC1uoKBogwU67SXTosnYkQ",
	"9gTZWDtVq7aVy22V9iESnKJGW2k9ma4NXb",
	"d1qgY5zvY8nKJJYuzQWTRV89vHyEnKY9GF",
	"25pcR28iA7mXZWjqMu12N7gwtdtLULvfgEJ",
	"24CDF6psd2WsLdyVz53zorNvu1iEsrbvexL",
	"2hVXZzpRe82p7GfgxJW4TPZn9mZAuSvyxhy",
	"2TBm78uxazExLyvzmuVKi9soKEuRvPZwn5c",
	"N8i7kSkDnTUbmm2Bn65nBCwjFfucBF7hUW",
	"QJe9H8VMphMmUQ65nkFRufVTf7jnpa7r8w",
	"26mKcaBpJfnPq1PaPBsRC1fGc9WC7H412og",
	"nxSFQwPaj6HWqCnQNwGEMGZ1dsVpxcMpbt",
	"2KxYQaLzHKiT1WDq7XoYKsQEmEMQK3iCEJi",
	"8bK8bQwwp7RvaLWM2BoAZ3CmuzUZsqXMWN",
	"ikpqw5zesZEVyDA1qydDReUg9uogKMRPzi",
	"pFLmBdF1Z1Bze7Rg81C4mVawog8GqzD8JT",
	"2NoP3SRXJkLuRMGmN9vdBpcJygeqrxxhse1",
	"VpEAPkuHUxrjVrKf5VRttd2Vis3xp6viBt",
	"GUrCX38aut9PdQZhCzF99Ej2ahPvdqKtrc",
	"2SP5FYVZpEhF5hZgvLj64tXT7AAtRY87WvA",
	"2TAiGppge2dx864NDUqL5me9mZCcAcRK8oW",
	"TLc44UT9XvsRvb73rviD3bFvG75R2pB647",
	"2eZLY5SMNLAtBbTgEPB6XMaD6WRCKphFbL8",
	"zrazwfvPCJcdEc2gELuoRRHJ154fUBnTHg",
	"277G8go25xZDPh7N2vCeT3fCfMjfaVix8Ft",
	"VR1KVU8ouUjTD4FVeyW2C6xCDMYbbP3vi5",
	"2hTXDCHTqeKd88hjE5gtMijv9DLp6SYnhFk",
	"2AcrEtHH2m2mq2pbb1Gctq6avrSg4tDetfb",
	"WtE3Ux2nrG8Y7DzczHgDjMTqPpr5TQkAuW",
	"2EVTfoQToAEq95ePru4gKaDuSYaqpnUGWxf",
	"td3ooSVJT2XSVKvsW9THwhDefhNNRQs5td",
	"eGN6uRyHf8VAWQLpyqegSsRHLmgapk1SfL",
	"27rxXbsBAur9vyLtp5rpt5JSH595Q59JCqN",
	"24WCqNxiE5tJpDbNeWcX58yrhdm5fnvvBKo",
	"g7Wf9hePuU9pFZfsZm1J3trDGt9nDF4hPQ",
	"mUj1ea5sXCZzM8RkJbESJMJc2E4W413RAh",
	"rWF2qjuynBriBoXAxeDg32q3bqUT2YH196",
	"4rvpc2C8b59JcD8dFdSM2PnBc6Xp22scPS",
	"LdfsVejcmEHGXEwT5HB3yEPmf3jcgpiPnm",
	"duo2hLqhQCbT52TUBfjvhE7uUwDGpb4Kma",
	"c7eBsHA9hnyFN5auxoaosKTmQereM6SK3j",
	"JNtgjHwgvnXLyDzyzDmijFcVpVeF7PuGwu",
	"2m6phmLb3wSbsVreH1XsrsoXx4TKSrEyZ1N",
	"2HE2bHbuMBkRJuzXnhhCYZUREim8gaWC8e",
	"wKGGWJctfsowQgyr4f87i2XT6j8wDMQYNX",
	"iYVweDJMPe4SaCXr6BormkVA2prAji347t",
	"AW7DGGBf4QLdUYsdhxgxJ9LKWTaqt6T6CA",
	"2NnZe8ERD8vakTszLfYZs3QEbgp3cgQPLJ4",
	"sfmb4TrTeMHjUU3J4J3McmTTN3Myx4g7x3",
	"2Xg5TAdyGHmCZbHZMtXCxYhbaaDcyNeWwRr",
	"2YoRqsvhiM6v5fQB9PayqbBBH4ri1LR6LKM",
	"2HkeyzGNRjarXaVnCSbsz3W3JkypC3N8vt1",
	"2W9yWvKjkyd9ormmVgL7ZAkXrEkFtUUmMNx",
	"Mnx5fna8DsgdVzfVGvEoAr7N6kJqjXEAW5",
	"2Scxg9iR1YVGA4TptCfE2NT3rgjCieSPCAY",
	"aBJdAFK6t2Qts35cMeFWY1m7xvWqk9z1ec",
	"owCf7yAUCkmtXZF1CJBpwb8LwVFXkgyaYi",
	"SzadZtmUSfmh7W6TzMg9Q6bgD2aqaEDm95",
	"2hP9PuMCj6yAtpudZrm5Wh6PQYGe5z5d3BT",
	"2jbvynpTzKL85JNi5g29eiJsW5vNpmvKMak",
	"2RK14nCUYiw74c7zZjnS7bxQPw1VTiCuDqg",
	"yWQBBrMZxnoewe2Qn5pMyPwPDWqJZwztYd",
	"2QS4rYxprieQQ8M2pXQNRyYnJPG4rxEcZs1",
	"2ZcqgrXpTVPtPaQccuynD4wcegJu11n4wCv",
	"gj2UzYpC7JikqPaEHqd6zAthSLsfncCUjY",
	"21PbxrXMFzXW6HsGjPKUSSAZkwP4xr5wCTW",
	"xyKvmzgLT27Ub5e3kQWk7PsUTCiVUntSuZ",
	"2HHYWPq1GENBXZfPY3RCsvsLBPrqBntwJ2S",
	"2MHJut5ktT1iDegTvaT7GQZnCmCrnh8Szef",
	"2FC6Uyh2R2brGcnKYQTkd4HqqtnGLdj7ztF",
	"hfRWcecnZjFHucZ1RntTUHJkdRPSYTMCA2",
	"2WrgUWhjtzX9Cw5YnanTkPtxz1r2pyRboks",
	"2Bhg2rvmm4mpL1Aup575JDwMfy5yNHCjJbG",
	"28NgehMnUHhv6vPLpmx8XHHw9Mt5DLemJqn",
	"79orEhmUwePwfZGKqenei3bCzricjtv1r2",
	"b51TrAcWBPowZjKjj3eNENEpTFoomfYKyZ",
	"yi8KUNChQ9kkTZEjbZZU6UNuXEx6NKGkFo",
	"G6dsYaqh3AUUQHRpwGcnYdcA9oeVzwPpRC",
	"2kc65nKwB6gs5rrybZDR3932C2DiRoYR1Kp",
	"mPgWBJ3JaMpbYwH8kwhheLabZrkXGXECfk",
	"2JxhYmcft4LK1zeFaDBTZSgC9dK2rrm6BGN",
	"2kM8oPA7zPE7yDE4XBYQ8CuQrPagJtC9Be3",
	"nSQXCFEpzwxQN1mhchU7rBuUWL9BiQGGa1",
	"MaMDjTDPZo9LFfkw2bSrdJt76PFBciURam",
	"pHi77kaFiMxMPkdtiPGyPArRjhNf7yVvxP",
	"Sdvy5vuLV8CAWUKJ36xxb4AN7QAibakhTB",
	"2CXSmr9KkupsFxTw4yfW535uPq56VAbuD5F",
	"UZdgZaFjwTjKs673shqdQsbR5GbVL86sC3",
	"2ZVMK6QEdcLoyZpHhfLX9apSbJURh7ZVbKw",
	"2AXQKjXMRxpWc4ibz2Wg4VCU62EG4QRDPw",
	"2SFMvd3jLNiQTqk9TVHnsodibx5YEe8Trzv",
	"2dzXfMPZzs88HeQcxRaypg3dQDoNUq7sntx",
	"23bbk8wW5diMU23CcqcB7tZFd7FXYsfhx35",
	"2UdwauW9NptjjWt9VCbpAt2foekbECFEKYb",
	"2Q8j6Wffndryb3SLuqDuaaSBVrHxB2RvTJr",
	"x2MYdsUcKeF2zNYRizSpmSnsfrar9YJkL8",
	"2E48seih7ZmmmR4bKPYw6E4GzzLhGTQBpCM",
	"kMmh2uHEzM54Y69NJqpEcNrJKXVXW958wS",
	"2aMQ4nLxxpgqm6S2o35ZyBMjafTULSSaAYE",
	"AxH75hcaMwq9P2YR8YMrheUdzzshRJ9N1T",
	"2dm1p7Kw7DC9uBPPHTem6e7GejqekEDPceK",
	"jGJWojYZzKRcrAbFRMfNdoTPa54ogAemwY",
	"2MthG5A4W1kPBE2zmEW4juaxu5F7xPK4NPa",
	"GfjzQWRbApKzKKTnBpmQnXEmnXjY8Us4ag",
	"i7nBpUVota7cf5SxVbN8fn7vMuvqz5qEB6",
	"roWrfLgnXE2tbyMw9ks36MtsKxfpDJyu8R",
	"b9YqfDnUD1ncSbGsGi3D9aWDbv9jmLotpn",
	"2cJbgNqzmAsyDgnRMEeTaL6sLqmyx1bNe2i",
	"uvo7goXXMhbbNSzJtWoMbmZ4VGBqecYUB6",
	"LEQNxkYiiC2tDt5BQSYYGm9ofohAsLL6sH",
	"2Pd1GE2Lg8TmbwRsoBYyV6zcSXcBbBrud19",
	"2gTCJAff4d7DeYFgLRVQ3xnmwXmdmUTX4Dt",
	"2gsu34JMwXE6jA4M6NLta1jbuPUvB1rQVWA",
	"2f4b5HMtJmsNmS7D92yNTJX9Tu4ENCXKGLo",
	"9MYnJVUukuXuTdLGYkt9NpcXTNxA4GS6NX",
	"evwoyyPKhaBB4T78Dpu2q1eHK7g5LPciTn",
	"c3RxQYvid9z2Dezg7yiDJxRAwDS3Fin3WK",
	"mdzhi6zLSwQi73pdcFL45dwrsKu7GZp6U3",
	"LiK5D5Stb28qGMFVDjEBaWYn9iU4vxaPEy",
	"2Ap9XP3phGVfmnh9QyrdiTxH8CJSGh2y9Wv",
	"2Jt5GAAdJL8yQE2uychUiy5AMBPM7kkVuAn",
	"2YzqptvELE4byUo54bgR1czoXhej3kobFn",
	"23efA9MRF9npe99KYd5UBZmdzoZeCVzUBkm",
	"QyLniqffazEGzo1XkLwpHApeMUTp6hgB5q",
	"ytEWqD4xazD3SnsH4yVzXcnrDox5rMeQQu",
	"2QJK4JRLKPZWzTHitq3bYCJqDTNWrYqCcZu",
	"2kQ2fw424PKBGGzLGVkcXnZu5nsxs9bTtda",
	"DsHo7iaezba3KwN3Umd8pyXNTq8MuoA1cW",
	"eizzV6tZmYk3atgF5f5bgYHay8jWmX6SSB",
	"2B3fDyC7i7R2jQ4h3xjtrmW6rkHApWWapS1",
	"23jDqCadRJrFHKmudfsHPzS2Kn6PwGJCFUm",
	"29PkokcEK4ZRJJu2DJw1a9WnuSNnrZfnvEk",
	"2AE7WXtGsiz8NqiPiB3HQF9nDPm6T95qPsX",
	"2EVZTG5uRbMTCzn32KxEkNeMUYCH4yK2fxq",
	"28WYGPqWQfmTEeAVfY6c17WvXekiFmgJVoj",
	"J9N1PXmW3YgD8tfv4dwt3j96cBtMnjMiGK",
	"DWAfVFzKoBYzxWX7L2SBefHBAMqzZpyRSN",
	"8ZHnVvLNFNEC2odX1hXmq5MhRVLHz4RiQq",
	"7zT9gBu5ovjNpNzFSyFoLBSrKC8mmrCFdp",
	"FXUJ9zFiMVtG4VGHyJuG3XeeE9ciHpbSeX",
	"6bnJqa9jeP5Q2MjjVGTfZcawLd7mfGvwTG",
	"2R1NuAtHjqaqNSftDyCct49oKcFHDd3MgjN",
	"whRQ8DS19QotEt5Uu2psudq6Mok9rfE9c9",
	"fbsRBNaCJxdn29xhB7Re6ZwBSk4QDDNVzV",
	"G2z3kAQpzVCRnEConLtjjX4xJ2MjmufXRx",
	"2kSTJU6fyx8ShxGguzJ1UogPkkgcoQPxiyE",
	"aZa8ikyF4aB92xoUJ5cu2r77GSd8vvk3c",
	"2MEpfwsBnwe4FwDpe7eXPtkGjJ5yFx2br7t",
	"24JKAWxr5tt2RsUxQnqTLJGEFHFgFAcGZwA",
	"BXru8c2yqiu9c7bh5sSsMYc7TNLKtD12Ao",
	"KT8QpNJ4VakjLf18rbu2dtcFTqCQzjVW3n",
	"Sq74LD2FPzKLKTdScE1vCwNSvrNX5sRTxj",
	"4Fvoq5JPTLBpvJMPU7uaMMHwHZicYe2gs2",
	"taAHTQMy6KjbvywYqNgMJ69LxxNb5b9Pez",
	"Pdggejt4UZXaz7UxJM79dYNwSDK68xCoGd",
	"6VcdMA1EDqeZXw3CvNA2stqLjVA1eMabDF",
	"Awu9UisDwc8WDs4uzFKpb7NGxU2xtkVHKC",
	"kGYyGEcxMcdzNedG2oHkv77mu61EKMHXZy",
	"26M8xTr8X7SdyYtDrJPGQkqdDv7h3kwacyo",
	"VDS4Es8E51CL61d4AqJZTEMBq56KFDpgmX",
	"2epmGinBwfXPhmcQqaoVHPVt9osDW2iUz5W",
	"2ADC8qXQ3VqxbDEUKkGNtqDXz2UBXVQan3x",
	"2Z6rGnVth7ufU3apNvDB3mLLXjRjV2kz6FX",
	"2Lt11nFb4qTwyRtg8sgGrL4dZVtvYPNdL5h",
	"2LKArV5tzdTEe3hZWRdhM1kXTw4JK4bPVDH",
	"2FpFE9GHdbB7FWKDvRrK7gaVnsrhj2iK6PX",
	"Dm89zT7v27AWCTCnijNJbE5Hd3paUZjFEW",
	"UcHSSQ83fYPNco9wxzX6aW44B8Q4WdJgNz",
	"jkkHGTKA8eEpfRpnDHXTERTjpnnxD9h7vf",
	"K7R2smMZowGgQAhuF1S36g9ouau3fqSo1F",
	"T7nxpqAUCnF2GDfbeo73ogVP9cad2Fs4b9",
	"rQFKL4tfkiPbhY1vc2byWPqgK1W67VMRr9",
	"2E9CC4ugc4TnuZeZoWqcZJbdKrsXERq3Kiz",
	"4uV71cMsibHSHxFkgnFDiqaHccPnAmsspJ",
	"2SDZFWnQ3XkoppTegTCkqoETFasBYuB8cfr",
	"2B5AjF2YSrepVaSfnQygSDJFN8SZ7bUEdGV",
	"28F6RmS3a5amTsERE3BphgZrrw3MdUBQYav",
	"26QF1q7YBbG8uAw8jmLwVasYcGjvULX6XsY",
	"2hnHhxHzkKFXJZwxDvtrLYCqnTNxtK3xszP",
	"bfZTCJo4VEH2ZgJqnSYz8jwKYDqCVBDruZ",
	"26pehDC8rXteyoUht6KRyTswq8mE7Hdm4V",
	"2KeSz9cmT2GKQ7ttjNRd2RKXPwzVvk3dmyW",
	"KVbuPjaus4QXkNvGBynQy9ViFriVNbUHjX",
	"2WoVMRmEEHC6YrYeKV7FR8wecTgB8mKCcbt",
	"262d5hYH4AZarR7yvzwDavrk84X9xmeRLco",
	"nuZVreiKAQo6hTGMhnNvw7KK9F5M4w2QKc",
	"274HRGCX1bXX5mL9rfChcg7C6FUX5Xy2DyV",
	"sikZ2wnoepbrGWxR1jXejAEwH5j9vH6sp3",
	"nGZnnV2AeTRtpCM2EziKaP9ggjGTAfA1Dx",
	"jivLNQYkeUXUUrpD6uuXoXbwrrmRZoKq8c",
	"2QN5A1qqpTPEpxnZAFWSTMND5ttmzuC7Fqb",
	"2gzEBHkJD2Tef2R2JhRV64Ncy9RXxaN6sXR",
	"TNLFB9emzA8g5jBw2AL43E582NZPoC4VeJ",
	"2Tmjksf1CyJngD37hH74DaG2eusTHrAzwcd",
	"uC1bKeomxQNETv3tLJ5GwynwMQhH8aZHcj",
	"2n1DP64k88UttcHZH3DePN5u5ivgfraeRyg",
	"2RLsctJvE2H6rsePM4UPrsBjdaPSz2Xthtg",
	"2CisrwzNj1MmuafN2xyNyCqX9fo3tnZySRJ",
	"wHLyWKBAMcvs4hhaPZxMy6DGQVzBQyyzAm",
	"2koRxzcw2YVehkdNFe1Y6nd9vSMDeANuDUD",
	"2EfByhGnQLCTPkgJgZLt5rjk3TnUtDKRgku",
	"Jxuo8MakYZDnpu56DV5bL1t92qvZg1VTny",
	"JmaqUne28mEyweMH8vQejVkJ8BE47Bu4su",
	"2ySG8aeKWEwAoyoqfEJXJAhYrzJ7C83wFC",
	"sF2tnctXe6SknBjyNUMg6TwWUAgDd3wh3e",
	"i4L9hw8YdV3FReBsmw55jPVSSVy3j6P77W",
	"2H5M42jvDCxYkBv363HeA6FLbPYSVp8ezYy",
	"2YL2imtRhooHYsoQ5paCgC4w4pzmq39JPTy",
	"zbx2zANgVoVECuaoTq2TttDVN3FAzhenYF",
	"2GJnyHjqYwRrbkZ5patb4j3Q3w7iCAqWXot",
	"2JEMngafTGR9hp4rtVvi9Ruc98HHTNV8oh3",
	"2EUWdaUG5EBjtFekenPT6JcTRnn4E48QaWS",
	"f8ZBmY2mVsciAQ8bcLx5P4hmGsbw59SpfL",
	"FtVgWsMzUSLe3sXukhez6G6XoPdd3QvQ3M",
	"2QcmEbgVQ9bEFAxhjCzUYu9cwpLKbeHQ9bZ",
	"2Nr6v9faM9m6eoDXSjSkcxAGejDeNPfTXoY",
	"t78Gfi6nrXRUvYkurxPpsESnQq6zEicAp4",
	"24RUPgfZ1prmDnFzZjndQ6BDenRwd3itS7q",
	"R9ytcV1Xi9K2anpaNwBcRJsVnaa7FQoi3s",
	"dZP6CJwCYNbCG4Fe45C3vFLYt8kBhEoaw2",
	"qNeEKDXSNwWgUW3qzZyuEr8J8Qw8kyu9C",
	"G4FXaPpF9i41vT9HcqcqvmZ4vPboxC8uzV",
	"2kKX35v6s6Hj39jE2b3aDwaaQWqbUjDQyjH",
	"2dB2t7smnaM4UhzHXuucxMJLCqPRDiHPMBK",
	"du8zmyKP92GzzHXvULCrNXAAU1BDbSUH2c",
	"kVW2kPoBSeYLzCji35GFg3jFJjMYXqzrNq",
	"s38EGGfSoMBuVsyaAgrA7qZNy7qBTJ3MRu",
	"q12jXpSzCorsRrqgHR6Q8euTrUmKXP7ri2",
	"2j1YRCsWgZGpp3voMWyXgWhjvcfCQoAUrYE",
	"2CDDiAUN1gxHTnehD298tgaxakSAncKn4RH",
	"SB52MtJ28r4DLBVoqHXymJXD2UXDuhF3Vj",
	"DtLaPqvb1nSvqiPsaNvQXBGbtg6zfZyfPC",
	"hdgiopiLuTeLJEaNZygZjuzd13CPXX5HNE",
	"HdrsJvddxsvFTovZeuDWXxqqb8aetStjLx",
	"MBuQtxwmhTUVytMfP7fx23ZyXt2TDh49SP",
	"n5X2kuVn4q8ZaMzSJCDBkEibR6KQ2cLGjP",
	"2PEdggFJDStMq7pYvcvUvgcvYTDYQLo86Fc",
	"2MpeMfcqbFFwR2tkMah5qvU6MkTGgcMHgQ",
	"HB1W6h16eAevbjbs9SCdMZeQsDkZjVFiNG",
	"DqaRUEtVHqLaaspKfYe4uDEsmcMzFZo9CY",
	"2PkSPGQBDbdp7W888oB1Ccts3QtxjfoWqZM",
	"wRx8u15Cn1gAv1rtu1o26z6FXcVi3xyuoE",
	"mMR2AE6RGHQxMFC4XFQirCXDRgqueN18PD",
	"2QCwkkYRHBr6i56UULRBgXP8iJiZPwueXDW",
	"27R7Gk5eEYrH9nMmnpzwkv4M4pWkWusVqoX",
	"2mGZzv2LM5tthBM7cCnfUU5HVGiVcUzPcuy",
	"23L6v9qtX15cN8fWWaDX57n4Kjr5on5UgAE",
	"2Z8bVQ6Z1CpmVUwyJFWmmPjfymGCqyU4cPc",
	"xjckK8YVwxGTUkVEUfyrUtzG5P57okDerA",
	"2V9afCkHBprDAeLeTToCto7RRmc3xHFnjKR",
	"1AvG9jvQJk2tPW1FAqKvqFBDRA1oXGctii",
	"BCpCAvrConRAwVinbS87KTCZ1ikvdyJPgh",
	"2exxgBRqeYeC9KAhQhRorNmNvzzaZUb3m4N",
	"ozsZcCzfSp5Ar7ArH1BPnxUuFmXNB1v49p",
	"2Y2UP1qG2cF52e4bHrfywkNbDdQag62BYW",
	"2kPryPw7GFiW2sFSfapkrGQAbA9Co58G6pG",
	"25v5sd6m4nbns3NUaTyVBmXMk96NCe8pRcr",
	"2d2z3Q2cU8SoYUPmCAGfdXwxz3fFD3515ft",
	"5QEvMMGgm97kepMcp8Yb69RRiJ4xSUvHF3",
	"dJG2cUVr7tjzPE1sR3yz8Qx2vr82EkRzNn",
	"yuCpxGt6htRZ6fYrjp3WkX9aGqNax2Q6es",
	"2SU2DfowSvtivsgjGdki43HwB1piMf3CWj1",
	"H2G9ddpdQ6LdyApF4aLkMywY3sEsBpd29h",
	"v87ceXuHu74ZCyyoGi9Jh9KvkP7yUWfMs1",
	"2mqAat8KsEb7jbTuCaitbWCknoFJu8HHazd",
	"2Yv1h6T43fCzBpYACtJLoiggGnbqEwPaChg",
	"hfLvaKte4MCjD6q2EWzJmZjiwDGvQ7Pxvm",
	"2JmTFdoz9y5jqfD2w72fXNbC8oAoVvsW18s",
	"22KHpkycMaqidYdDKz1QyVxjAyix8gidX1Y",
	"2ixr7n2zog4zzmDH5sUVdkuCd6dFitiMVEw",
	"23rErDvAFeSJA4VuPF3EwuZV3EJnDHVS4HB",
	"2N9KwGSzWjkoXHuKm1tGAnfZkPchYQCzgqD",
	"29FXKjg6Lqa8j1B9ceb67NsVE3WVsgp69V5",
	"257E2A1pCMv1TbyfssiYRS7g1zN7jguLjLH",
	"2Po4zAiGYU6Zjar2vZCPPa72vm6ymRpYNFn",
	"2mZXxFQzunCwwcDJbCoKPJaVrf6PeyTXH5A",
	"2MQ4LVmaqm8XYN4FjohUbhHx72KbBksccjP",
	"26uRPvr5A65U5PQEY9rV1icDnV5hskZmXXK",
	"2T7DAY9EwgDNo72vvRVkVNKDUpWpQ6NJ5RT",
	"p4GAV417P5sPmHF94qWXRLKJ6e7AHBoAMr",
	"2UqDr1Py5LocokHjnzj9drsh3domBjjX37X",
	"HxzzJWMzaXTJmumkeZELRypke7rXfAUZQF",
	"22piejdepGeqyjRHVBjX1nukvcKRs1q7DZ9",
	"2NhEUHur6VHti68nhD5dmgkZgMpvc581DZ9",
	"26ayeSAuPvicNvYaAVL7zevF9YMkkVSUHk4",
	"nCU4XpXxoGdsmzFdT4JzMJcgFgBRZdkH2d",
	"o3XMSd6EnAHXLZ5zKYzW9NY13Bkm7byWr6",
	"id6Y3dEYR6zX6JE4u2ApiMqpVh14vfQSgQ",
	"2mmMjUpQ6osYKZvYhrXwpTSANijN2xh35N9",
	"2KcWbqa8JuV4JjGeaZSCTUBk9hFuKnbip4S",
	"2GSNtxoQbVxNSX3Cum3AMtE4fXytd88P42b",
	"w7YLRL87zA5AhMe9EkHQnQmRcaoR63qsL6",
	"FD3FkifbVSLBktKirBRam66opkCY4FmQcN",
	"2YiW4G9yUnwjCGaV7Nb33FM4Eu9KzkUVmvE",
	"2SsZo8cxfGDcR5hYzDreJYEXM6vAQTkDfex",
}
