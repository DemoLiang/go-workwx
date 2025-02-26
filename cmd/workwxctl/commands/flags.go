package commands

import (
	"fmt"
	"net/http"
	"os"
)

const (
	flagCorpID            = "corpid"
	flagCorpSecret        = "corpsecret"
	flagAgentID           = "agentid"
	flagWebhookKey        = "webhook-key"
	flagQyapiHostOverride = "qyapi-host-override"
	flagTLSKeyLogFile     = "tls-key-logfile"

	flagName   = "name"
	flagOwner  = "owner"
	flagUser   = "user"
	flagChatID = "chatid"

	flagMessageType  = "message-type"
	flagSafe         = "safe"
	flagToUser       = "to-user"
	flagToUserShort  = "u"
	flagToParty      = "to-party"
	flagToPartyShort = "p"
	flagToTag        = "to-tag"
	flagToTagShort   = "t"
	flagToChat       = "to-chat"
	flagToChatShort  = "c"

	flagMediaID          = "media-id"
	flagThumbMediaID     = "thumb-media-id"
	flagDescription      = "desc"
	flagTitle            = "title"
	flagAuthor           = "author"
	flagURL              = "url"
	flagPicURL           = "pic-url"
	flagButtonText       = "button-text"
	flagSourceContentURL = "source-content-url"
	flagDigest           = "digest"

	flagMediaType = "media-type"

	flagMentionUser        = "mention-user"
	flagMentionMobile      = "mention-mobile"
	flagMentionMobileShort = "m"
)

type cliOptions struct {
	CorpID     string
	CorpSecret string
	AgentID    int64
	WebhookKey string

	QYAPIHostOverride string
	TLSKeyLogFile     string
}

func mustGetConfig(c *cli.Context) *cliOptions {

	return &cliOptions{
		CorpID:     c.String(flagCorpID),
		CorpSecret: c.String(flagCorpSecret),
		AgentID:    c.Int64(flagAgentID),
		WebhookKey: c.String(flagWebhookKey),

		QYAPIHostOverride: c.String(flagQyapiHostOverride),
		TLSKeyLogFile:     c.String(flagTLSKeyLogFile),
	}
}

//
// impl cliOptions
//

func (c *cliOptions) makeHTTPClient() *http.Client {
	if c.TLSKeyLogFile == "" {
		return http.DefaultClient
	}

	f, err := os.OpenFile(c.TLSKeyLogFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("can't open TLS key log file for writing: %+v\n", err)
		panic(err)
	}

	fmt.Fprintf(f, "# SSL/TLS secrets log file, generated by go\n")

	return &http.Client{
		Transport: newTransportWithKeyLog(f),
	}
}

func (c *cliOptions) makeWorkwxClient() *workwx.Workwx {
	if c.CorpID == "" {
		panic("corpid must be set")
	}

	if c.CorpSecret == "" {
		panic("corpsecret must be set")
	}

	if c.AgentID == 0 {
		panic("agentid must be set (for now; may later lift the restriction)")
	}

	httpClient := c.makeHTTPClient()
	if c.QYAPIHostOverride != "" {
		// wtf think of a way to change this
		return workwx.New(c.CorpID,
			workwx.WithQYAPIHost(c.QYAPIHostOverride),
			workwx.WithHTTPClient(httpClient),
		)
	}
	return workwx.New(c.CorpID, workwx.WithHTTPClient(httpClient))
}

func (c *cliOptions) MakeWorkwxApp() *workwx.WorkwxApp {
	return c.makeWorkwxClient().WithApp(c.CorpSecret, c.AgentID)
}

func (c *cliOptions) makeWebhookClient() *workwx.WebhookClient {
	if c.WebhookKey == "" {
		panic("webhook key must be set")
	}

	httpClient := c.makeHTTPClient()
	if c.QYAPIHostOverride != "" {
		// wtf think of a way to change this
		return workwx.NewWebhookClient(c.WebhookKey,
			workwx.WithQYAPIHost(c.QYAPIHostOverride),
			workwx.WithHTTPClient(httpClient),
		)
	}
	return workwx.NewWebhookClient(c.WebhookKey, workwx.WithHTTPClient(httpClient))

}
