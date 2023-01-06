package surge

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/openapi"
)

func WalletAddress() string {
	wallet, _ := nkn.ClientAddrToWalletAddr(GetAccountAddress())
	return wallet
}

func WalletTransfer(address string, amount string, fee string) (bool, string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = fee

	result, err := client.Transfer(address, amount, config)
	if err != nil {
		pushError("Transfer failed", err.Error())
		return false, ""
	}
	log.Println("Transfered " + amount + " nkn to " + address + " txHash: " + result)
	return true, result
}

func WalletBalance() string {
	amount, err := client.Balance()
	if err != nil {
		pushError("Transfer failed", err.Error())
		return "-1"
	}
	return amount.String()
}

func CalculateFee(Fee string) (string, error) {
	avgFee, err := openapi.GetAvgFee()
	if err != nil {
		return "0.0", err
	}

	avgFeeFloat, _ := strconv.ParseFloat(avgFee, 64)

	feePercent := 0.2
	lowFee := avgFeeFloat - avgFeeFloat*feePercent
	highFee := avgFeeFloat + avgFeeFloat*feePercent

	switch Fee {
	case "0":
		return "0", nil
	case "33":
		return fmt.Sprintf("%f", lowFee), nil
	case "66":
		return avgFee, nil
	case "100":
		return fmt.Sprintf("%f", highFee), nil
	}

	return "0", nil
}

//ValidateBalanceForTransaction returns a boolean for whether there is enough balance to make a transation
func ValidateBalanceForTransaction(Amount float64, Fee float64, UtilTransaction bool) (bool, error) {
	if !UtilTransaction && Amount < 0.00000001 {
		return false, errors.New("minimum tip amount is 0.00000001")
	}

	balance := WalletBalance()
	balanceFloat, _ := strconv.ParseFloat(balance, 64)

	if Amount+Fee > balanceFloat {
		return false, errors.New("not enough nkn available required: " + fmt.Sprintf("%f", Amount+Fee) + " available: " + balance)
	}

	return true, nil
}

func IsSubscriptionActive(TopicEncoded string) (bool, error) {
	subs, err := client.GetSubscribers(TopicEncoded, 0, 1000, false, true)
	if err != nil {
		return false, err
	}

	for sub := range subs.Subscribers.Map {
		if sub == GetAccountAddress() {
			return true, nil
		}
	}

	for sub := range subs.SubscribersInTxPool.Map {
		if sub == GetAccountAddress() {
			return true, nil
		}
	}
	return false, nil
}

/*Wallet Features
- Import wallet from private key
- Export wallet private key
- ✔ Send transaction with (amount, fee, toAddress)
- ✔ WalletInfo (retrieve personal wallet address + wallet balance)
- ✔ Set/Get transaction fee default


- Optional:
seed + password wallet files instead of private keys?
get network average fee for last x amount of blocks
*/
