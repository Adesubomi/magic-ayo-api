package entity

type BitcoinWallet struct {
	ID       string
	Mnemonic string
	Address  string
}

func (bW BitcoinWallet) GenerateAddress(index int) (*BitcoinAddress, error) {
	return nil, nil
}
