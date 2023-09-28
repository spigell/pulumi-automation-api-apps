package pulumi

import (
	"context"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
)

func (p *Pulumi) Destroy(ctx context.Context) error {
	opts := []optdestroy.Option{}
	opts = append(opts, optdestroy.ProgressStreams(os.Stdout))

	_, err := p.stack.Destroy(ctx, opts...)

	return err
}
