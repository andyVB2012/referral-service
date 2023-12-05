package kafka

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"strconv"

	"github.com/andyVB2012/referral-service/block"
	"github.com/andyVB2012/referral-service/logger"
	"github.com/andyVB2012/referral-service/repository"
	"github.com/joho/godotenv"
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
var (
	isConsumerHealthy int32 = 1 // Use an atomic variable for thread-safe updates
)

func IsConsumerHealthy() bool {
	return atomic.LoadInt32(&isConsumerHealthy) == 1
}

func RunnConsumers(repo repository.Repository) {
	logger.Log.Info("Running consumers")
	err := godotenv.Load() // Loads from .env by default, or specify .env file
	if err != nil {
		logger.Log.Fatalf("Error loading .env file: %v", err)
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092" // Default value if the environment variable is not set
	}

	logger.Log.Info("Running consumers")
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval: time.Second,
		Brokers:        []string{kafkaBrokers},
		GroupID:        "referral-service",
		GroupTopics:    []string{"stfx.stream.block.arb.stfx.vault.v2", "stfx.stream.block.eth.stfx.vault.v2", "stfx.stream.block.arb.stfx.q", "stfx.stream.block.eth.stfx.q", "stfx.stream.block.eth.stfx.stake", "stfx.stream.block.arb.stfx.stake", "stfx.stream.block.arb.stfx.subscription.v1"},
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Log.Error("Failed to read: ", err.Error())
			atomic.StoreInt32(&isConsumerHealthy, 0)
			break
		}
		atomic.StoreInt32(&isConsumerHealthy, 1)
		processMessage(m, repo)
	}
}

func processMessage(m kafka.Message, repo repository.Repository) {
	var event block.BlockEvent // Use the correct struct from your protobuf

	if err := proto.Unmarshal(m.Value, &event); err != nil {
		logger.Log.Error("Error unmarshalling message: ", err)
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
		logger.Log.Info("UnTracked Contract", event.GetBlockIndex().GetContractAddr())
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
			logger.Log.Info("Deposit Q")
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				logger.Log.Error("Error converting amount to int: ", err.Error())
				return
			} else {
				err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
				if err != nil {
					logger.Log.Error("Error adding trader account: ", err.Error())
					return
				}

				err1 := repo.AddVariable(context.Background(), traderAddr, "deposited", amountTo)
				if err1 != nil {
					logger.Log.Error("Error adding variable to deposit: ", err.Error())
					return
				}
				logger.Log.Infof("Deposit to Q Added %v to %v", amountTo, traderAddr)
			}
		}
	case "Withdraw":
		// tokenAddr := block.UnmarshalString(event.GetBlockIndex(), "tokenAddr")
		traderAcc := block.UnmarshalString(event.GetBlockIndex(), "traderAcc")
		amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
		traderAddr := block.UnmarshalString(event.GetBlockIndex(), "traderAddr")
		if repo.IsTraderAddrInDb(context.Background(), traderAddr) {
			logger.Log.Info("Withdraw Q")
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				logger.Log.Error("Error converting amount to int: ", err.Error())
				return
			} else {
				err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
				if err != nil {
					logger.Log.Error("Error adding trader account: ", err.Error())
					return
				}
				err1 := repo.AddVariable(context.Background(), traderAddr, "withdrawn", amountTo)
				if err1 != nil {
					logger.Log.Error("Error adding variable to withdraw: ", err.Error())
					return
				}
				logger.Log.Infof("Withdraw from Q Added %v to %v", amountTo, traderAddr)
			}
		}
	case "CreateTraderAccount":
		traderAcc := block.UnmarshalString(event.GetBlockIndex(), "traderAcc")
		traderAddr := block.UnmarshalString(event.GetBlockIndex(), "traderAddr")
		if repo.IsTraderAddrInDb(context.Background(), traderAddr) {
			logger.Log.Info("CreateTraderAccount Q")
			err := repo.AddTraderAcctIfNotExists(context.Background(), traderAddr, traderAcc)
			if err != nil {
				logger.Log.Error("Error adding trader account: ", err.Error())
				return
			}
			logger.Log.Infof("CreateTraderAccount Q Added %v to %v", traderAcc, traderAddr)
		}
	default:
		logger.Log.Info("Untracked Event Q", event.GetBlockIndex().GetEvent())
	}
}

func digestVaultEvent(event *block.BlockEvent, repo repository.Repository) {
	switch event.GetBlockIndex().GetEvent() {
	case "CreateStv":
		stv := block.UnmarshalString(event.GetBlockIndex(), "stvId")
		managerAddr := block.UnmarshalString(event.GetBlockIndex(), "managerAddr")
		if repo.IsTraderAddrInDb(context.Background(), managerAddr) {

			logger.Log.Info("CreateStv")
			err := repo.AddVariable(context.Background(), managerAddr, "vaultscreated", 1)
			if err != nil {
				logger.Log.Error("Error adding variable to vaultscreated: ", err.Error())
				return
			}
			logger.Log.Infof("CreateStv Added %v to %v", stv, managerAddr)
		}
	case "Deposit":
		stvId := block.UnmarshalString(event.GetBlockIndex(), "stvId")
		investorAddr := block.UnmarshalString(event.GetBlockIndex(), "investorAddr")
		amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
		if repo.IsTraderAddrInDb(context.Background(), investorAddr) {
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				logger.Log.Error("Error converting amount to int: ", err.Error())
			} else {
				err := repo.AddVariable(context.Background(), investorAddr, "vaultsinvestedin", 1)
				if err != nil {
					logger.Log.Error("Error adding variable to vaultsinvestedin: ", err.Error())
					return
				}
				err1 := repo.AddVariable(context.Background(), investorAddr, "manualinvted", amountTo)
				if err1 != nil {
					logger.Log.Error("Error adding variable to manualinvted: ", err.Error())
					return
				}
				err2 := repo.AddVariable(context.Background(), investorAddr, "totalinvested", amountTo)
				if err2 != nil {
					logger.Log.Error("Error adding variable to totalinvested: ", err.Error())
					return
				}
				logger.Log.Infof("Deposit to Vault %v Added %v to %v", stvId, amountTo, investorAddr)
			}
		}
	case "DepositWithSubscription":
		investorAddr := block.UnmarshalString(event.GetBlockIndex(), "investorAddr")
		if repo.IsTraderAccountInDb(context.Background(), investorAddr) {
			amount := block.UnmarshalString(event.GetBlockIndex(), "amount")
			traderAddr, err := repo.GetTraderAddrFromTraderAcct(context.Background(), investorAddr)
			if err != nil {
				logger.Log.Error("Error getting traderAddr from traderAcct: ", err.Error())
			} else {
				amountTo, err := strconv.Atoi(amount)
				if err != nil {
					logger.Log.Error("Error converting amount to int: ", err.Error())
					return
				} else {
					err1 := repo.AddVariable(context.Background(), traderAddr, "vaultsinvestedin", 1)
					if err1 != nil {
						logger.Log.Error("Error adding variable to vaultsinvestedin: ", err.Error())
						return
					}
					err2 := repo.AddVariable(context.Background(), traderAddr, "subscriptioninvested", amountTo)
					if err2 != nil {
						logger.Log.Error("Error adding variable to subscriptioninvested: ", err.Error())
						return
					}
					err3 := repo.AddVariable(context.Background(), traderAddr, "totalinvested", amountTo)
					if err3 != nil {
						logger.Log.Error("Error adding variable to totalinvested: ", err.Error())
						return
					}
					logger.Log.Infof("DepositWithSubscription to Vault Added %v to %v", amountTo, traderAddr)
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
		if repo.IsTraderAddrInDb(context.Background(), subscriberAddr) {
			err := repo.AddVariable(context.Background(), subscriberAddr, "subscriptions", 1)
			if err != nil {
				logger.Log.Error("Error adding variable to subscriptions: ", err.Error())
				return
			}
			logger.Log.Infof("Subscribe Added %v to %v", subscriberAddr, managerAddr)
		}

	case "Unsubscribe":
		managerAddr := block.UnmarshalString(event.GetBlockIndex(), "managerAddr")
		subscriberAddr := block.UnmarshalString(event.GetBlockIndex(), "subscriberAddr")
		if repo.IsTraderAddrInDb(context.Background(), subscriberAddr) {
			err := repo.Unsubscribe(context.Background(), subscriberAddr)
			if err != nil {
				logger.Log.Error("Error unsubscribing: ", err.Error())
				return
			}
			logger.Log.Infof("Unsubscribe Added %v to %v", subscriberAddr, managerAddr)
		}
	default:
		// fmt.Println("Unknown Event", event.GetBlockIndex().GetEvent())
	}
}

func digestStakeEvent(event *block.BlockEvent, repo repository.Repository) {
	switch event.GetBlockIndex().GetEvent() {
	case "AddStake":
		stkrAddr := block.UnmarshalString(event.GetBlockIndex(), "stkrAddr")
		amount := block.UnmarshalString(event.GetBlockIndex(), "stkAmount")
		if repo.IsTraderAddrInDb(context.Background(), stkrAddr) {
			amountTo, err := strconv.Atoi(amount)
			if err != nil {
				logger.Log.Error("Error converting amount to int: ", err.Error())
				return
			} else {
				err := repo.AddVariable(context.Background(), stkrAddr, "stakedamount", amountTo)
				if err != nil {
					logger.Log.Error("Error adding variable to stakedamount: ", err.Error())
					return
				}
				logger.Log.Infof("Deposit Stake Added %v to %v", amountTo, stkrAddr)
			}
		}
	default:
		logger.Log.Info("UnTracked Event Stake", event.GetBlockIndex().GetEvent())
	}
}
