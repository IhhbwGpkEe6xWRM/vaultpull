// Package prompt provides a lightweight interactive confirmation helper used
// by vaultpull before performing potentially destructive actions, such as
// overwriting an existing .env file when secrets have changed.
//
// Basic usage:
//
//	c := prompt.New()
//	ok, err := c.Confirm("Overwrite .env with updated secrets?")
//	if err != nil {
//		return err
//	}
//	if !ok {
//		fmt.Println("aborted")
//		return nil
//	}
//
// For testing, inject custom readers and writers via NewWithReadWriter.
package prompt
