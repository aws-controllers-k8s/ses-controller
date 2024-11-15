package configuration_set

import (
	"context"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"

	"github.com/aws-controllers-k8s/ses-controller/pkg/util"
)

func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	_ *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	return util.ValidateImmutableResource(ctx, rm.getImmutableFieldChanges(delta), desired)
}
