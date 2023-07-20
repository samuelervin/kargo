package handler

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubev1alpha1 "github.com/akuity/kargo/api/v1alpha1"
	typesv1alpha1 "github.com/akuity/kargo/internal/api/types/v1alpha1"
	svcv1alpha1 "github.com/akuity/kargo/pkg/api/service/v1alpha1"
)

type GetStageV1Alpha1Func func(
	context.Context,
	*connect.Request[svcv1alpha1.GetStageRequest],
) (*connect.Response[svcv1alpha1.GetStageResponse], error)

func GetStageV1Alpha1(
	kc client.Client,
) GetStageV1Alpha1Func {
	return func(
		ctx context.Context,
		req *connect.Request[svcv1alpha1.GetStageRequest],
	) (*connect.Response[svcv1alpha1.GetStageResponse], error) {
		if req.Msg.GetProject() == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("project should not be empty"))
		}
		if req.Msg.GetName() == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name should not be empty"))
		}

		if err := kc.Get(ctx, client.ObjectKey{Name: req.Msg.GetProject()}, &corev1.Namespace{}); err != nil {
			if kubeerr.IsNotFound(err) {
				return nil, connect.NewError(connect.CodeNotFound,
					fmt.Errorf("project %q not found", req.Msg.GetProject()))
			}
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		var stage kubev1alpha1.Stage
		if err := kc.Get(ctx, client.ObjectKey{
			Namespace: req.Msg.GetProject(),
			Name:      req.Msg.GetName(),
		}, &stage); err != nil {
			if kubeerr.IsNotFound(err) {
				return nil, connect.NewError(connect.CodeNotFound, err)
			}
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return connect.NewResponse(&svcv1alpha1.GetStageResponse{
			Stage: typesv1alpha1.ToStageProto(stage),
		}), nil
	}
}