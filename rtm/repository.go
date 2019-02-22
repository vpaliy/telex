package rtm

import (
	"github.com/vpaliy/telex/model"
	"github.com/vpaliy/telex/store"
	"github.com/vpaliy/telex/utils"
)

type Content interface{}

type TokenizedChannelRequest struct {
	Token string
	ChannelRequest
}

type TokenizedPostRequest struct {
	Token string
	CreateMessageRequest
}

type Repository interface {
	FetchChannel(string) (*model.Channel, error)
	JoinChannel(*TokenizedChannelRequest) (Content, error)
	LeaveChannel(*TokenizedChannelRequest) (Content, error)
	PostMessage(*TokenizedPostRequest) (Content, error)
}

type repository struct {
	userStore         store.UserStore
	messageStore      store.MessageStore
	subscriptionStore store.MessageStore
	channelStore      store.ChannelStore
}

func (r *repository) FetchChannel(ch string) (*model.Channel, error) {
	return r.channelStore.Fetch(ch)
}

func (r *repository) PostMessage(request *TokenizedPostRequest) (Content, error) {
	channel, err := r.channelStore.Fetch(request.Channel)
	// if there is an error
	if channel == nil || err != nil {
		return nil, err
	}
	// check if the user is subscribed to the chat
	currentUser := utils.GetUserFromToken(request.Token)
	if !channel.HasUser(currentUser.ID) {
		// TODO: return an error that the user is not here
		return nil, nil
	}
	// submit a message
	message := request.ToMessage(
		currentUser.ID, channel.ID,
	)
	// create message and send an error message if fails
	if err := r.messageStore.Create(message); err != nil {
		// TODO: notify that the message has failed to be written to the db
		return nil, err
	}
	return message, nil
}

func (r *repository) JoinChannel(request *TokenizedChannelRequest) (Content, error) {
	channel, err := r.channelStore.Fetch(request.Channel)
	// notify about the error
	if channel == nil || err != nil {
		return nil, err
	}
	// get the user from JWT token
	user := utils.GetUserFromToken(request.Token)
	// create a subscription
	subscription := channel.CreateSubscription(user.ID)
	// create it
	if err := r.subscriptionStore.Create(channel, subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (r *repository) LeaveChannel(request *TokenizedChannelRequest) (Content, error) {
	return nil, nil
}
