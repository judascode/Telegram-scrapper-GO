package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"telegram.downloloader.com/config"
	"telegram.downloloader.com/prompt"
	"telegram.downloloader.com/service"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/query/channels/participants"
	"github.com/gotd/td/tg"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh/terminal"
)

type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

type termAuth struct {
	noSignUp
	phone string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func getNewFlow() auth.Flow {
	promptUserNumber := promptui.Prompt{
		Label: "Enter your telegram number (5511987654321)",
	}

	userNumber, err := promptUserNumber.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	flow := auth.NewFlow(
		termAuth{phone: userNumber},
		auth.SendCodeOptions{},
	)
	return flow
}

func callBackParticipants(size int) func(ctx context.Context, participant participants.Elem) error {
	iterator := 0

	callBackParticipants := func(ctx context.Context, participant participants.Elem) error {
		iterator++
		user, ok := participant.User()
		if !ok {
			return nil
		}
		fmt.Println(iterator, " of ", size, "-", user.ID)
		time.Sleep(1 * time.Second)
		return nil
	}

	return callBackParticipants
}

func main() {

	opts, err := telegram.OptionsFromEnvironment(telegram.Options{})
	if err != nil {
		panic(err)
	}

	client := telegram.NewClient(config.GetConfig().APP_ID, config.GetConfig().APP_HASH, opts)
	if err = client.Run(context.Background(), func(ctx context.Context) error {

		authStatus, err := client.Auth().Status(ctx)
		if err != nil {
			panic(err)
		}
		if !authStatus.Authorized {
			if err := client.Auth().IfNecessary(ctx, getNewFlow()); err != nil {
				return err
			}
		}
		rawClient := tg.NewClient(client)
		sender := message.NewSender(rawClient)

		channelName, err := prompt.NameOfGroup()
		if err != nil {
			panic(err)
		}

		service.SaveUsers(ctx, sender, rawClient, channelName)
		return nil
	}); err != nil {
		panic(err)
	}
}
