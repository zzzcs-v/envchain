package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/envchain/internal/encrypt"
)

var (
	encryptPassphrase string
	encryptDecrypt    bool
)

func init() {
	encryptCmd := &cobra.Command{
		Use:   "encrypt [value]",
		Short: "Encrypt or decrypt a single env var value",
		Args:  cobra.ExactArgs(1),
		RunE:  runEncrypt,
	}

	encryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "passphrase for encryption/decryption (required)")
	encryptCmd.Flags().BoolVarP(&encryptDecrypt, "decrypt", "d", false, "decrypt the given value instead of encrypting")
	_ = encryptCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	value := args[0]
	e := encrypt.NewEncryptor(encryptPassphrase)

	if encryptDecrypt {
		plaintext, err := e.Decrypt(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}
		fmt.Println(plaintext)
		return nil
	}

	ciphertext, err := e.Encrypt(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	fmt.Println(ciphertext)
	return nil
}
