package kafka

import (
	"context"
	"fmt"
	"time"

	"strconv"

	"github.com/andyVB2012/referral-service/block"
	"github.com/andyVB2012/referral-service/repository"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

var (
	subscriptionContractAddrArb = "0x6cA7dB3BbAa82699202d547E28b8d122e217107a"
	qContractAddrArb            = "0xddfDb14Ac8107b004Fb3A1fB5960572433A45aCf"
	vaultContractAddrArb        = "0x56AFA0771B6E1E0A33bf050f029ffe49c00F1BB9"
	stakeContractAddrArb        = "0x961C24c0eD7285A96B804F177bB38dc4a89BADc1"
	uSDCe                       = "0xff970a61a04b1ca14834a43f5de4533ebddb5cc8"
)
var (
	qContractAddrEth     = "0xF9AA0BBe42EBF0BeBc5a44BAcE12baB5F00A342B"
	vaultContractAddrEth = "0x8910A8AA84d5a30F7603BA4B5Ec263cFB8C4D49c"
	stakeContractAddrEth = "0xa16218fbada83d2d299b57a88362abd91f11ab56"
)

func RunnConsumers(repo repository.Repository) {
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval: time.Second,
		Brokers:        []string{"localhost:9092"},                                                                                                                                                                                                                                                   // Replace with your Kafka broker addresses
		GroupID:        "referral-service",                                                                                                                                                                                                                                                           // Replace with your consumer group ID
		GroupTopics:    []string{"stfx.stream.block.arb.stfx.vault.v2", "stfx.stream.block.eth.stfx.vault.v2", "stfx.stream.block.arb.stfx.q", "stfx.stream.block.eth.stfx.q", "stfx.stream.block.eth.stfx.stake", "stfx.stream.block.arb.stfx.stake", "stfx.stream.block.arb.stfx.subscription.v1"}, // Replace with your topic name
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Failed to read", err)
			break
		}
		processMessage(m, repo)
	}
}

func processMessage(m kafka.Message, repo repository.Repository) {
	var event block.BlockEvent // Use the correct struct from your protobuf

	if err := proto.Unmarshal(m.Value, &event); err != nil {
		fmt.Printf("Error unmarshalling message: %v", err)
	}

	// case for event.GetBlockIndex().GetContractAddr()
	switch event.GetBlockIndex().GetContractAddr() {
	case qContractAddrArb, qContractAddrEth:
		digestQEvent(&event, repo)
	case vaultContractAddrArb, vaultContractAddrEth:
		digestVaultEvent(&event, repo)
	case stakeContractAddrArb, stakeContractAddrEth:
		digestStakeEvent(&event, repo)
	case subscriptionContractAddrArb:
		digestSubscriptionEvent(&event, repo)
	default:
		fmt.Println("Unknown Contract")
	}

}

func digestQEvent(event *block.BlockEvent, repo repository.Repository) {
	switch event.GetBlockIndex().GetEvent() {
	case "Deposit":
		traderAddr := block.UnmarshalString(event.GetBlockIndex(), "traderAddr")
		if repo.IsTraderAddrInDb(context.Background(), traderAddr) {
			tokenAddr := block.UnmarshalString(event.GetBlockIndex(), "tokenAddr")
			traderAcc := block.UnmarshalString(event.GetBlockIndex(), "traderAcc")
			amount := ""
			if tokenAddr == uSDCe {
				amount = block.UnmarshalString(event.GetBlockIndex(), "amount")
			} else {
				amount = block.UnmarshalString(event.GetBlockIndex(), "returnAmount")
			}
			fmt.Println("Deposit Q")
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				fmt.Println("Error converting amount to int: ", err)
			} else {

				fmt.Println("traderAcc: ", traderAcc)
				fmt.Println("traderAddr: ", traderAddr)
				fmt.Println("amount: ", amount)
				err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
				if err != nil {
					fmt.Println("Error adding trader account: ", err)
				}

				err1 := repo.AddVariable(context.Background(), traderAddr, "deposited", amountTo)
				if err1 != nil {
					fmt.Println("Error adding variable to deposit: ", err)
				}
			}
		}
	case "Withdraw":
		// tokenAddr := block.UnmarshalString(event.GetBlockIndex(), "tokenAddr")
		traderAcc := block.UnmarshalString(event.GetBlockIndex(), "traderAcc")
		amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
		traderAddr := block.UnmarshalString(event.GetBlockIndex(), "traderAddr")
		if repo.IsTraderAddrInDb(context.Background(), traderAddr) {
			fmt.Println("Withdraw Q")
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				fmt.Println("Error converting amount to int: ", err)
			} else {
				err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
				if err != nil {
					fmt.Println("Error adding trader account: ", err)
				}
				err1 := repo.AddVariable(context.Background(), traderAddr, "withdrawn", amountTo)
				if err1 != nil {
					fmt.Println("Error adding variable to withdraw: ", err)
				}
			}
		}
	case "CreateTraderAccount":
		traderAcc := block.UnmarshalString(event.GetBlockIndex(), "traderAcc")
		traderAddr := block.UnmarshalString(event.GetBlockIndex(), "traderAddr")
		if repo.IsTraderAddrInDb(context.Background(), traderAddr) {
			fmt.Println("CreateTraderAccount Q")
			fmt.Println("traderAcc: ", traderAcc)
			fmt.Println("traderAddr: ", traderAddr)
			err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
			if err != nil {
				fmt.Println("Error adding trader account: ", err)
			}
		}
	default:
		fmt.Println("Unknown Event Q", event.GetBlockIndex().GetEvent())
	}
}

func digestVaultEvent(event *block.BlockEvent, repo repository.Repository) {
	switch event.GetBlockIndex().GetEvent() {
	case "CreateStv":
		stv := block.UnmarshalString(event.GetBlockIndex(), "stvId")
		managerAddr := block.UnmarshalString(event.GetBlockIndex(), "managerAddr")
		if repo.IsTraderAddrInDb(context.Background(), managerAddr) {
			fmt.Println("CreateStv")
			fmt.Println("managerAddr: ", managerAddr)
			fmt.Println("stv: ", stv)
			repo.AddVariable(context.Background(), managerAddr, "vaultscreated", 1)
		}
	case "Deposit":
		stvId := block.UnmarshalString(event.GetBlockIndex(), "stvId")
		investorAddr := block.UnmarshalString(event.GetBlockIndex(), "investorAddr")
		amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
		if repo.IsTraderAddrInDb(context.Background(), investorAddr) {
			fmt.Println("Deposit vault")
			fmt.Println("stvId: ", stvId)
			fmt.Println("investorAddr: ", investorAddr)
			fmt.Println("amount: ", amount)
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				fmt.Println("Error converting amount to int: ", err)
			} else {
				err := repo.AddVariable(context.Background(), investorAddr, "vaultsinvestedin", 1)
				if err != nil {
					fmt.Println("Error adding variable to vaultsinvestedin: ", err)
				}
				err1 := repo.AddVariable(context.Background(), investorAddr, "manualinvted", amountTo)
				if err1 != nil {
					fmt.Println("Error adding variable to manualinvted: ", err)
				}
				err2 := repo.AddVariable(context.Background(), investorAddr, "totalinvested", amountTo)
				if err2 != nil {
					fmt.Println("Error adding variable to totalinvested: ", err)
				}
			}
		}
	case "Cancel":
		fmt.Println("Cancel vault ")
		stvId := block.UnmarshalString(event.GetBlockIndex(), "stvId")
		fmt.Println("stvId: ", stvId)
	case "DepositWithSubscription":
		investorAddr := block.UnmarshalString(event.GetBlockIndex(), "investorAddr")
		if investorAddr == "0x07c5aa477c140cc13c31b5c508f2f4f0ff9f47d8" {
			fmt.Println("DepositWithSubscription vault ")
		}
		if repo.IsTraderAccountInDb(context.Background(), investorAddr) {
			amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
			traderAddr, err := repo.GetTraderAddrFromTraderAcct(context.Background(), investorAddr)
			if err != nil {
				fmt.Println("Error getting traderAddr from traderAcct: ", err)
			} else {
				amountTo, err := strconv.Atoi(amount)
				if err != nil {
					fmt.Println("Error converting amount to int: ", err)
				} else {
					repo.AddVariable(context.Background(), traderAddr, "vaultsinvestedin", 1)
					repo.AddVariable(context.Background(), traderAddr, "subscriptioninvested", amountTo)
					repo.AddVariable(context.Background(), traderAddr, "totalinvested", amountTo)
				}
			}
		}
	default:
		// fmt.Println("Unknown Event", event.GetBlockIndex().GetEvent())
	}
}

func digestSubscriptionEvent(event *block.BlockEvent, repo repository.Repository) {
	switch event.GetBlockIndex().GetEvent() {
	case "Subscribe":
		managerAddr := block.UnmarshalString(event.GetBlockIndex(), "managerAddr")
		subscriberAddr := block.UnmarshalString(event.GetBlockIndex(), "subscriberAddr")
		subscriberAccountAddr := block.UnmarshalString(event.GetBlockIndex(), "subscriberAccountAddr")
		maxLimit := block.UnmarshalString(event.GetBlockIndex(), "maxLimit")
		if repo.IsTraderAddrInDb(context.Background(), subscriberAddr) {
			fmt.Println("Subscription Contract")
			fmt.Println("Subscribe")
			fmt.Println("managerAddr: ", managerAddr)
			fmt.Println("subscriberAddr: ", subscriberAddr)
			fmt.Println("subscriberAccountAddr: ", subscriberAccountAddr)
			fmt.Println("maxLimit: ", maxLimit)
			repo.AddVariable(context.Background(), subscriberAddr, "subscriptions", 1)
		}

	case "Unsubscribe":
		managerAddr := block.UnmarshalString(event.GetBlockIndex(), "managerAddr")
		subscriberAddr := block.UnmarshalString(event.GetBlockIndex(), "subscriberAddr")
		subscriberAccountAddr := block.UnmarshalString(event.GetBlockIndex(), "subscriberAccountAddr")
		if repo.IsTraderAddrInDb(context.Background(), subscriberAddr) {
			fmt.Println("Unsubscribe")
			fmt.Println("managerAddr: ", managerAddr)
			fmt.Println("subscriberAddr: ", subscriberAddr)
			fmt.Println("subscriberAccountAddr: ", subscriberAccountAddr)
			repo.Unsubscribe(context.Background(), subscriberAddr)
		}
	default:
		// fmt.Println("Unknown Event", event.GetBlockIndex().GetEvent())
	}
}

func digestStakeEvent(event *block.BlockEvent, repo repository.Repository) {
	fmt.Println("Stake Contract")
	switch event.GetBlockIndex().GetEvent() {
	case "AddStake":
		stkrAddr := block.UnmarshalString(event.GetBlockIndex(), "stkrAddr")
		amount := block.UnmarshalString(event.GetBlockIndex(), "stkAmount")
		if repo.IsTraderAddrInDb(context.Background(), stkrAddr) {
			fmt.Println("Deposit Stake")
			fmt.Println("stkrAddr: ", stkrAddr)
			fmt.Println("amount: ", amount)
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				fmt.Println("Error converting amount to int: ", err)
			} else {
				repo.AddVariable(context.Background(), stkrAddr, "stakedamount", amountTo)
			}
		}
	default:
		fmt.Println("Unknown Event Stake", event.GetBlockIndex().GetEvent())
	}
}
