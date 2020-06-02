package gcp10

import (
	"testing"
	"time"

	"github.com/zeepin/ZeepinChain-Crypto/keypair"
	"github.com/zeepin/ZeepinChain/core/types"
	"github.com/zeepin/Zeepin-Go-Sdk"
	"github.com/zeepin/Zeepin-Go-Sdk/account"
)

const scriptHash = "0f27a43a74c963e07c0b633aff49ebb269e6d727"

func TestGcp10(t *testing.T) {
	contractAddr := scriptHash

	zptSdk := Zeepin_Go_Sdk.NewZeepinSdk()
	zptSdk.NewRpcClient().SetAddress("http://localhost:20336")
	Gcp10 := NewGcp10(contractAddr, zptSdk)

	wallet, err := zptSdk.OpenWallet("../../wallet.json")
	if err != nil {
		t.Fatal(err)
	}
	if wallet.GetAccountCount() < 2 {
		t.Fatal("account not enough")
	}
	acc, err := wallet.GetDefaultAccount([]byte("passwordtest"))
	if err != nil {
		t.Fatal(err)
	}
	name, err := Gcp10.Name(acc)
	if err != nil {
		t.Fatal(err)
	}
	symbol, err := Gcp10.Symbol(acc)
	if err != nil {
		t.Fatal(err)
	}
	decimals, err := Gcp10.Decimals(acc)
	if err != nil {
		t.Fatal(err)
	}
	totalSupply, err := Gcp10.TotalSupply(acc)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := Gcp10.BalanceOf(acc, acc.Address.ToBase58())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("name %s, symbol %s, decimals %d, totalSupply %d, balanceOf %s is %d",
		name, symbol, decimals, totalSupply, acc.Address.ToBase58(), balance)

	anotherAccount, err := wallet.GetAccountByIndex(2, []byte("passwordtest"))
	if err != nil {
		t.Fatal(err)
	}
	m := 2
	multiSignAddr, err := types.AddressFromMultiPubKeys([]keypair.PublicKey{acc.PublicKey, anotherAccount.PublicKey}, m)
	if err != nil {
		t.Fatal(err)
	}
	amount := "1000"
	gasPrice := uint64(500)
	gasLimit := uint64(500000)
	transferTx, err := Gcp10.Transfer(acc, multiSignAddr.ToBase58(), amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferTx %s: from %s to multi-sign addr %s, amount %d", transferTx.ToHexString(),
		acc.Address.ToBase58(), multiSignAddr.ToBase58(), amount)
	accounts := []*account.Account{acc, anotherAccount}
	transferMultiSignTx, err := Gcp10.MultiSignTransfer(accounts, m, acc.Address.ToBase58(), amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferMultiSignTx %s: from %s to multi-sign addr %s, amount %d", transferMultiSignTx.ToHexString(),
		multiSignAddr.ToBase58(), acc.Address.ToBase58(), amount)
	approveTx, err := Gcp10.Approve(acc, multiSignAddr.ToBase58(), amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("approveTx %s: owner %s approve to multi-sign spender addr %s, amount %d", approveTx.ToHexString(),
		acc.Address.ToBase58(), multiSignAddr.ToBase58(), amount)
	multiSignTransferFromTx, err := Gcp10.MultiSignTransferFrom(accounts, m, acc.Address.ToBase58(), multiSignAddr.ToBase58(), amount,
		gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multiSignTransferFromTx %s: owner %s, multi-sign spender addr %s, to %s, amount %d",
		multiSignTransferFromTx.ToHexString(), acc.Address.ToBase58(), multiSignAddr.ToBase58(), multiSignAddr.ToBase58(),
		amount)
	multiSignApproveTx, err := Gcp10.MultiSignApprove(accounts, m, acc.Address.ToBase58(), amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multiSignApproveTx %s: multi-sign owner %s approve to spender addr %s, amount %d",
		multiSignApproveTx.ToHexString(), multiSignAddr.ToBase58(), acc.Address.ToBase58(), amount)
	transferFromTx, err := Gcp10.TransferFrom(acc, multiSignAddr.ToBase58(), acc.Address.ToBase58(), amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferFromTx %s: multi-sign owner %s, spender addr %s, to %s, amount %d",
		transferFromTx.ToHexString(), multiSignAddr.ToBase58(), acc.Address.ToBase58(), acc.Address.ToBase58(), amount)
	_, _ = zptSdk.WaitForGenerateBlock(30 * time.Second)

	eventsFromTx, err := Gcp10.FetchTxTransferEvent(transferTx.ToHexString())
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range eventsFromTx {
		t.Logf("tx %s transfer event: %s", transferTx.ToHexString(), evt.String())
	}

	height := uint32(1791727)
	eventsFromBlock, err := Gcp10.FetchBlockTransferEvent(height)
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range eventsFromBlock {
		t.Logf("block %d transfer event: %s", height, evt.String())
	}
}
