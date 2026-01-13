package email

type Reciever struct {
	Email string
}

type Params struct {
	Subject  string
	Message  string
	Reciever Reciever
}
