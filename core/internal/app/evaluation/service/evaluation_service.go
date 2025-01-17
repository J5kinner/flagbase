package service

import (
	"context"
	evaluationmodel "core/internal/app/evaluation/model"
	flagmodel "core/internal/app/flag/model"
	flagrepo "core/internal/app/flag/repository"
	segmentrepo "core/internal/app/segment/repository"
	segmentrulemodel "core/internal/app/segmentrule/model"
	segmentrulerepo "core/internal/app/segmentrule/repository"
	targetingmodel "core/internal/app/targeting/model"
	targetingrepo "core/internal/app/targeting/repository"
	targetingrulemodel "core/internal/app/targetingrule/model"
	targetingrulerepo "core/internal/app/targetingrule/repository"
	cons "core/internal/pkg/constants"
	rsc "core/internal/pkg/resource"
	"core/internal/pkg/srvenv"
	"core/pkg/evaluator"
	"core/pkg/hashutil"
	"core/pkg/model"
	res "core/pkg/response"
)

type Service struct {
	Senv              *srvenv.Env
	FlagRepo          *flagrepo.Repo
	SegmentRepo       *segmentrepo.Repo
	SegmentRuleRepo   *segmentrulerepo.Repo
	TargetingRepo     *targetingrepo.Repo
	TargetingRuleRepo *targetingrulerepo.Repo
}

func NewService(senv *srvenv.Env) *Service {
	return &Service{
		Senv:              senv,
		FlagRepo:          flagrepo.NewRepo(senv),
		SegmentRepo:       segmentrepo.NewRepo(senv),
		SegmentRuleRepo:   segmentrulerepo.NewRepo(senv),
		TargetingRepo:     targetingrepo.NewRepo(senv),
		TargetingRuleRepo: targetingrulerepo.NewRepo(senv),
	}
}

// Get returns a set raw (non-evaluated) flagsets
// (*) atk: access_type <= service
func (s *Service) Get(
	atk rsc.Token,
	a evaluationmodel.RootArgs,
) (*model.Flagset, *res.Errors) {
	var e res.Errors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var o model.Flagset

	fl, _e := s.FlagRepo.List(context.Background(), flagmodel.RootArgs{
		WorkspaceKey: a.WorkspaceKey,
		ProjectKey:   a.ProjectKey,
	})
	if _e != nil {
		e.Append(cons.ErrorInternal, _e.Error())
	}

	for _, f := range fl {
		t, err := s.TargetingRepo.Get(ctx, targetingmodel.RootArgs{
			WorkspaceKey:   a.WorkspaceKey,
			ProjectKey:     a.ProjectKey,
			FlagKey:        f.Key,
			EnvironmentKey: a.EnvironmentKey,
		})
		if err != nil {
			e.Append(cons.ErrorInternal, err.Error())
		}

		tr, _err := s.TargetingRuleRepo.List(ctx, targetingrulemodel.RootArgs{
			WorkspaceKey:   a.WorkspaceKey,
			ProjectKey:     a.ProjectKey,
			FlagKey:        f.Key,
			EnvironmentKey: a.EnvironmentKey,
		})
		if _err != nil {
			e.Append(cons.ErrorInternal, _err.Error())
		}

		var _r []*model.Rule
		for _, _tr := range tr {
			switch _tr.Type {
			case string(rsc.Trait):
				_r = append(_r, &model.Rule{
					ID:             _tr.ID,
					RuleType:       _tr.Type,
					TraitKey:       _tr.TraitKey,
					TraitValue:     _tr.TraitValue,
					Operator:       _tr.Operator,
					Negate:         _tr.Negate,
					RuleVariations: _tr.RuleVariations,
				})
			case string(rsc.Segment):
				if string(_tr.SegmentKey) != "" {
					sr, err := s.SegmentRuleRepo.List(ctx, segmentrulemodel.RootArgs{
						WorkspaceKey:   a.WorkspaceKey,
						ProjectKey:     a.ProjectKey,
						EnvironmentKey: a.EnvironmentKey,
						SegmentKey:     _tr.SegmentKey,
					})
					if err != nil {
						e.Append(cons.ErrorInternal, err.Error())
					}

					for _, _sr := range sr {
						_r = append(_r, &model.Rule{
							ID:             _tr.ID,
							RuleType:       _tr.Type,
							TraitKey:       _sr.TraitKey,
							TraitValue:     _sr.TraitValue,
							Operator:       _sr.Operator,
							Negate:         _sr.Negate,
							RuleVariations: _tr.RuleVariations,
						})
					}
				}
				// TODO match identity ~ case string(rsc.Identity):
			}
		}

		o = append(o, &model.Flag{
			ID:                    t.ID,
			FlagKey:               string(f.Key),
			UseFallthrough:        !t.Enabled,
			FallthroughVariations: t.FallthroughVariations,
			Rules:                 _r,
		})
	}

	return &o, &e
}

// Evaluate returns an evaluated flagset given the user context
// (*) atk: access_type <= service
func (s *Service) Evaluate(
	atk rsc.Token,
	ectx model.Context,
	a evaluationmodel.RootArgs,
) (*model.Evaluations, *res.Errors) {
	var e res.Errors

	r, err := s.Get(atk, a)
	if !err.IsEmpty() {
		e.Extend(err)
	}

	var o model.Evaluations

	for _, flag := range *r {
		salt := hashutil.HashKeys(
			string(a.WorkspaceKey),
			string(a.ProjectKey),
			string(a.EnvironmentKey),
			flag.FlagKey,
			ectx.Identifier,
		)
		o = append(o, evaluator.Evaluate(*flag, salt, ectx))
	}

	return &o, &e
}
