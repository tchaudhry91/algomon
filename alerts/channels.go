package alerts

import "context"

type Alerter interface {
	Alert(context.Context) error
}
