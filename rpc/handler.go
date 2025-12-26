package rpc
import(
	"context"
)

type Handler func(ctx context.Context, params map[string]any) (any, error)

type Server struct {
	handlers map[string]Handler
}
