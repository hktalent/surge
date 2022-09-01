// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all blockchain related functions
*/

package surge

import (
	"log"
	"strconv"

	"github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
)

var TransactionFee string

func subscribeToPubSub(topic string) bool {
	config := &nkn.DefaultTransactionConfig

	calculatedFee, err := CalculateFee(TransactionFee)
	if err != nil {
		pushError("Error on subscribe to topic", err.Error())
		return false
	}
	config.Fee = calculatedFee

	feeFloat, _ := strconv.ParseFloat(config.Fee, 64)
	hasBalance, _ := ValidateBalanceForTransaction(0, feeFloat, true)
	if !hasBalance {
		pushError("Error on subscribe to topic", "Not enough fee in wallet, consider depositing NKN or if possible lower transaction fees in the wallet settings.")
		updateTopicSubscriptionState(topic, 0)
		return false
	}

	updateTopicSubscriptionState(topic, 1)
	txnHash, err := client.Subscribe("", topic, constants.SubscriptionDuration, constants.TransactionMeta, config)
	if err != nil {
		log.Println("Subsription transaction failed for topic:", topic, "error:", err)
		updateTopicSubscriptionState(topic, 1)
		return false
	} else {
		log.Println("Subscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
	updateTopicSubscriptionState(topic, 2)
	return true
}

func unsubscribeToPubSub(topic string) bool {
	config := &nkn.DefaultTransactionConfig

	var err error = nil
	config.Fee, err = CalculateFee(TransactionFee)
	if err != nil {
		pushError("Error on unsubscribe to topic", err.Error())
		return false
	}

	initialState := 0
	currentState, exists := topicEncodedSubcribeStateMap[topic]
	if exists {
		initialState = currentState
	}

	feeFloat, _ := strconv.ParseFloat(config.Fee, 64)
	hasBalance, _ := ValidateBalanceForTransaction(0, feeFloat, true)
	if !hasBalance {
		pushError("Error on unsubscribe to topic", "Not enough fee in wallet, consider depositing NKN or if possible lower transaction fees in the wallet settings.")
		updateTopicSubscriptionState(topic, initialState)
		return false
	}
	updateTopicSubscriptionState(topic, 1)

	txnHash, err := client.Unsubscribe("", topic, config)
	if err != nil {
		log.Println("Probably not subscribed to:", topic, "error:", err)
	} else {
		log.Println("Unsubscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
	updateTopicSubscriptionState(topic, 0)
	return true
}
