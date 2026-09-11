package bootstrap

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"golang.org/x/time/rate"
)

type blockBurstLimiter struct {
	*rate.Limiter
}

func (l blockBurstLimiter) WaitN(ctx context.Context, total int) error {
	for total > 0 {
		n := l.Burst()
		if l.Limiter.Limit() == rate.Inf || n > total {
			n = total
		}
		err := l.Limiter.WaitN(ctx, n)
		if err != nil {
			return err
		}
		total -= n
	}
	return nil
}

// InitStreamLimit installs the stream limiters used by the download and upload
// paths. The per-direction limits are no longer configurable, so every limiter
// is installed unlimited and the wrappers stay in place for the callers that
// already route through them.
func InitStreamLimit() {
	unlimited := func(limiter *stream.Limiter) {
		*limiter = blockBurstLimiter{Limiter: rate.NewLimiter(rate.Inf, 0)}
	}
	unlimited(&stream.ClientDownloadLimit)
	unlimited(&stream.ClientUploadLimit)
	unlimited(&stream.ServerDownloadLimit)
	unlimited(&stream.ServerUploadLimit)
}
