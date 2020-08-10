package gombal

type Messaging struct {
	Sender    UserWithRef `json:"sender,omitempty"`
	Recipient User        `json:"recipient,omitempty"`
	Timestamp int         `json:"timestamp,omitempty"`
	Message   Message     `json:"message,omitempty"`
	Postback  Postback    `json:"postback,omitempty"`
}

type UserWithRef struct {
	ID      string `json:"id,omitempty"`
	UserRef string `json:"user_ref,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Postback struct {
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type Message struct {
	MID         string        `json:"mid,omitempty"`
	Text        string        `json:"text,omitempty"`
	QuickReply  *QuickReply   `json:"quick_reply,omitempty"`
	ReplyTo     *ReplyTo      `json:"reply_to,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
}

type ReplyTo struct {
	MID string `json:"mid,omitempty"`
}

type QuickReply struct {
	Payload string `json:"payload,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Payload struct {
	URL       string `json:"url,omitempty"`
	Title     string `json:"title,omitempty"`
	StickerId string `json:"sticker_id,omitempty"`
}

type Entry struct {
	ID         string      `json:"id,omitempty"`
	Time       int         `json:"time,omitempty"`
	Messagings []Messaging `json:"messaging,omitempty"`
}

type Callback struct {
	Object  string  `json:"object,omitempty"`
	Entries []Entry `json:"entry,omitempty"`
}

type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}
