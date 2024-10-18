package service

type PollRepository interface{

}

type PollSrv struct{
	PollRepository
}

func NewPollSrv(pollRepo PollRepository) *PollSrv{
	return &PollSrv{PollRepository: pollRepo}
}