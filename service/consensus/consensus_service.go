package consensus

import "it-chain/service/blockchain"

type ConsensusService interface{

	//Consensus 시작
	StartConsensus(transaction blockchain.Transaction)

	StopConsensus()

	//consensus메세지는 모두 이쪽으로 받는다.
	ReceiveConsensusMessage(consensusMsg *ConsensusMessage)
}