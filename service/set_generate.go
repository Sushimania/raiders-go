package service

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/dchest/blake2b"
	"github.com/gosuri/uilive"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"raiders-go/model"
	"raiders-go/properties"
	"strconv"
	"strings"
	"time"
)

var (
	// Information of Bitcoin addresses and balance
	mChainState map[string]uint64

	// result of hashing private key
	blob []byte
	// the number of 0
	difficulty int
	difficultyPrefix string

	searchingSpeed int

	optionLogFlag bool

	authToken string
)

func SetGenerate(eosAccountName string, machineId string) {
	// get JWT
	authToken = getAuthToken(eosAccountName, machineId)

	count := 0
	optionLogFlag = false

	// get difficulty from server
	// [code]

	difficulty = 5
	applyDifficulty()

	fmt.Println("Loading balance data for all Bitcoin addresses")
	fmt.Println("Please wait a few minutes...")

	// Ctrl + c : Show total count
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		// Program finished
		fmt.Println(" Total Search Address: ", count)
		os.Exit(0)
	}()

	// load chainstate file
	loadChainStates()

	// for showing on terminal
	writer := uilive.New()
	writer.Start()

	fmt.Println("Start a raid...")
	startTime := time.Now()

	// connect Bitcoin mainnet
	chainParams := &chaincfg.MainNetParams
	for {
		t := time.Now()

		// create a private key
		privKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			log.Fatalf("Failed to create private key, err: %v", err)
		}

		addrPubKey, err := btcutil.NewAddressPubKey(
			privKey.PubKey().SerializeUncompressed(), chainParams)
		if err != nil {
			log.Fatalf("Failed to calculate public key, err: %v", err)
		}

		rcvAddr := addrPubKey.AddressPubKeyHash().EncodeAddress()
		wif, err := btcutil.NewWIF(privKey, chainParams, false)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		// option for showing logs
		if optionLogFlag {
			fmt.Println("--------------------------------------------------------------------------------------")
			fmt.Println("[" + strconv.Itoa(count)+ "] " + t.String())
			fmt.Println("[" + strconv.Itoa(count)+ "] " + "Private Key: " + wif.String() + "	address: " + rcvAddr)
		}

		// Proof of work
		hashWork(wif.String())

		// for test when I found something
		//if count == 500000 {
		//	// 292929
		//	rcvAddr = "3PhV7nQziDpxaty6P6gQqDKpyUa3pNsW6S"
		//}

		if matchAddress(rcvAddr) {
			// I found Bitcoin!
			// send a private key to server
			// [code]
		}

		// Calculate speed of searching addresses
		if count % 100000 == 0 && count >= 100000 {
			elapsedTime := time.Since(startTime)

			// searching speed per second
			speed := float64(100000 * 1000000000) / float64(elapsedTime)

			//fmt.Println("speed: ", float64(100000 * 1000000000) / float64(elapsedTime), " Keys/s")
			//fmt.Printf("%.0f Keys/s \n", speed)

			fmt.Fprintf(writer, "Raiding... %.0f Keys/s\n", speed)
			searchingSpeed, err = strconv.Atoi(fmt.Sprintf("%.0f", speed))
			if err != nil {
				log.Fatalf("err: %v", err)
			}

			startTime = time.Now()

			applyDifficulty()
		}

		count++
	}

	fmt.Println("Ends the search.")
}

func loadChainStates() {
	// balances file has 100 addresses for fast loading. When deploy for live production, use full data.
	bytes, err := ioutil.ReadFile("chainstate/balances") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(bytes) // convert content to a 'string'

	chainStateRawArr := strings.Fields(str)
	mChainState = make(map[string]uint64)

	// parsing from text to map variable
	length := len(chainStateRawArr)
	var tempArr []string
	for i := 0; i < length; i++ {
		// [Bitcoin address];[satoshi]
		tempArr = strings.Split(chainStateRawArr[i], ";")
		u64, _ := strconv.ParseUint(tempArr[1], 10, 64)

		// map is faster than slice(https://www.darkcoding.net/software/go-slice-search-vs-map-lookup/)
		mChainState[tempArr[0]] = uint64(u64)
	}

	//balanceData := mChainState["3Cbq7aT1tY8kMxWLbitaG7yT6bPbKChq64"]
	//fmt.Println("prs: ", balanceData)

	fmt.Println("Loading is complete")
}

func applyDifficulty() {
	difficultyPrefix = ""
	for i := 0; i < difficulty; i++ {
		difficultyPrefix += "0"
	}
}

func matchAddress(address string) bool {
	matchFlag := false

	if mChainState[address] > 0 {
		matchFlag = true
	}

	return matchFlag
}

func hashWork(privateKey string) {
	// blake2b PoW
	h := blake2b.New256()
	h.Write([]byte(privateKey + properties.BLAKE2B_SALT))
	hashResult := fmt.Sprintf("%x", h.Sum(nil))

	if strings.HasPrefix(hashResult, difficultyPrefix) {
		// found a proper hash!

		// Send nonce data to server with JWT
		// After validation on server, user will receive rewards.
		url := properties.RAIDSPLATFORM_URL + "/api/submitpow"
		payload := strings.NewReader("{\n  \"nonce\" : \"" + privateKey + "\"\n}")
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("Authorization", "Basic " + properties.BASIC_AUTH_KEY)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Auth-Token", authToken)
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var block model.Block
		_ = json.Unmarshal(body, &block)

		fmt.Println("block: ", block)

		if block.BlockReward <= 0 {
			log.Fatalf("err: submitpow error")
		}

		fmt.Println("BlockReward: %v", block.BlockReward)
	}
}

func getAuthToken(eosAccountName string, machineId string) string {
	url := properties.RAIDSPLATFORM_URL + "/api/getauthtoken"
	payload := strings.NewReader("{\n  \"eosAccountName\" : \"" + eosAccountName + "\",\n  \"machineId\" : \"" + machineId + "\"\n}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Authorization", "Basic " + properties.BASIC_AUTH_KEY)
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var getAuthToken model.GetAuthToken
	_ = json.Unmarshal(body, &getAuthToken)

	fmt.Println("getauthtoken response: ", string(body))

	return getAuthToken.XAuthToken
}