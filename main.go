package main

import (
	// standard
	"bufio"
	"fmt"
	"os"

	// external
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/mitchellh/go-homedir"
	coinApi "github.com/miguelmota/go-coinmarketcap"
)

var (
	coinpath = kingpin.Arg("coins", "File containing list of coins, default '~/.coins'").Default("~/.coins").String()
	coinsrc string
)

func init() {
	kingpin.Parse()
	var err error
	if *coinpath == "~/.coins" {
		coinsrc, err = homedir.Expand(*coinpath)
		if err != nil {
			panic(err)
		}
	} else {
		coinsrc = *coinpath
	}
}

func getCoins(path string) ([]string, error) {
	coins := []string{}

	file, err := os.Open(path)
	if err != nil {
		return coins, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		coins = append(coins, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return coins, err
	}

	return coins, nil
}

func main() {
	coins, err := getCoins(coinsrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, coinname := range(coins) {
		coin, err := coinApi.GetCoinData(coinname)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %.2f\n", coin.Name, coin.PriceUsd)
	}
}
