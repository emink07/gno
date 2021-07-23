package client

import (
	"crypto/sha256"
	"fmt"

	"github.com/gnolang/gno/pkgs/command"
	"github.com/gnolang/gno/pkgs/crypto/bip39"
)

type GenerateOptions struct {
	CustomEntropy bool
}

func runGenerateCmd(cmd *command.Command) error {
	opts := cmd.Options.(GenerateOptions)
	customEntropy := opts.CustomEntropy

	var entropySeed []byte

	if customEntropy {
		// prompt the user to enter some entropy
		inputEntropy, err := cmd.GetString("> WARNING: Generate at least 256-bits of entropy and enter the results here:")
		if err != nil {
			return err
		}
		if len(inputEntropy) < 43 {
			return fmt.Errorf("256-bits is 43 characters in Base-64, and 100 in Base-6. You entered %v, and probably want more", len(inputEntropy))
		}
		conf, err := cmd.GetConfirmation(fmt.Sprintf("> Input length: %d", len(inputEntropy)))
		if err != nil {
			return err
		}
		if !conf {
			return nil
		}

		// hash input entropy to get entropy seed
		hashedEntropy := sha256.Sum256([]byte(inputEntropy))
		entropySeed = hashedEntropy[:]
	} else {
		// read entropy seed straight from crypto.Rand
		var err error
		entropySeed, err = bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		return err
	}
	cmd.Println(mnemonic)

	return nil
}
