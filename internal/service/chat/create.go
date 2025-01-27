package chat

import "context"

func (s *serv) Create(ctx context.Context, usernames []string) (int64, error) {
	id, err := s.chatRepository.Create(ctx, usernames)
	if err != nil {
		return 0, err
	}

	return id, nil
}
