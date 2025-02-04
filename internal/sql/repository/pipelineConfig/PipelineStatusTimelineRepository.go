package pipelineConfig

import (
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type TimelineStatus string

const (
	TIMELINE_STATUS_GIT_COMMIT            TimelineStatus = "GIT COMMIT"
	TIMELINE_STATUS_KUBECTL_APPLY_STARTED TimelineStatus = "KUBECTL APPLY STARTED"
	TIMELINE_STATUS_KUBECTL_APPLY_SYNCED  TimelineStatus = "KUBECTL APPLY SYNCED"
	TIMELINE_STATUS_APP_HEALTHY           TimelineStatus = "HEALTHY"
	TIMELINE_STATUS_APP_DEGRADED          TimelineStatus = "DEGRADED"
)

type PipelineStatusTimelineRepository interface {
	SaveTimeline(timeline *PipelineStatusTimeline) error
	UpdateTimeline(timeline *PipelineStatusTimeline) error
	FetchTimelinesByPipelineId(pipelineId int) ([]*PipelineStatusTimeline, error)
	FetchTimelinesByWfrId(wfrId int) ([]*PipelineStatusTimeline, error)
	FetchTimelineOfLatestWfByCdWorkflowIdAndStatus(pipelineId int, status TimelineStatus) (*PipelineStatusTimeline, error)
	CheckTimelineExistsOfLatestWfByCdWorkflowIdAndStatus(cdWorkflowId int, status TimelineStatus) (bool, error)
}

type PipelineStatusTimelineRepositoryImpl struct {
	dbConnection *pg.DB
	logger       *zap.SugaredLogger
}

func NewPipelineStatusTimelineRepositoryImpl(dbConnection *pg.DB,
	logger *zap.SugaredLogger) *PipelineStatusTimelineRepositoryImpl {
	return &PipelineStatusTimelineRepositoryImpl{
		dbConnection: dbConnection,
		logger:       logger,
	}
}

type PipelineStatusTimeline struct {
	tableName                    struct{}       `sql:"pipeline_status_timeline" pg:",discard_unknown_columns"`
	Id                           int            `sql:"id,pk"`
	InstalledAppVersionHistoryId int            `sql:"installed_app_version_history_id"`
	CdWorkflowRunnerId           int            `sql:"cd_workflow_runner_id"`
	Status                       TimelineStatus `sql:"status"`
	StatusDetail                 string         `sql:"status_detail"`
	StatusTime                   time.Time      `sql:"status_time"`
	sql.AuditLog
}

func (impl *PipelineStatusTimelineRepositoryImpl) SaveTimeline(timeline *PipelineStatusTimeline) error {
	err := impl.dbConnection.Insert(timeline)
	if err != nil {
		impl.logger.Errorw("error in saving timeline of cd pipeline status", "err", err, "timeline", timeline)
		return err
	}
	return nil
}

func (impl *PipelineStatusTimelineRepositoryImpl) UpdateTimeline(timeline *PipelineStatusTimeline) error {
	err := impl.dbConnection.Update(timeline)
	if err != nil {
		impl.logger.Errorw("error in updating timeline of cd pipeline status", "err", err, "timeline", timeline)
		return err
	}
	return nil
}

func (impl *PipelineStatusTimelineRepositoryImpl) FetchTimelinesByPipelineId(pipelineId int) ([]*PipelineStatusTimeline, error) {
	var timelines []*PipelineStatusTimeline
	err := impl.dbConnection.Model(&timelines).
		Join("INNER JOIN cd_workflow_runner wfr ON wfr.id = pipeline_status_timeline.cd_workflow_runner_id").
		Join("INNER JOIN cd_workflow cw ON cw.id=wfr.cd_workflow_id").
		Where("cw.pipelineId = ?", pipelineId).Select()
	if err != nil {
		impl.logger.Errorw("error in getting timelines by pipelineId", "err", err, "pipelineId", pipelineId)
		return nil, err
	}
	return timelines, nil
}

func (impl *PipelineStatusTimelineRepositoryImpl) FetchTimelinesByWfrId(wfrId int) ([]*PipelineStatusTimeline, error) {
	var timelines []*PipelineStatusTimeline
	err := impl.dbConnection.Model(&timelines).
		Where("cd_workflow_runner_id = ?", wfrId).Select()
	if err != nil {
		impl.logger.Errorw("error in getting timelines by wfrId", "err", err, "wfrId", wfrId)
		return nil, err
	}
	return timelines, nil
}

func (impl *PipelineStatusTimelineRepositoryImpl) FetchTimelineOfLatestWfByCdWorkflowIdAndStatus(cdWorkflowId int, status TimelineStatus) (*PipelineStatusTimeline, error) {
	timeline := &PipelineStatusTimeline{}
	err := impl.dbConnection.Model(timeline).
		Join("INNER JOIN cd_workflow_runner wfr ON wfr.id = pipeline_status_timeline.cd_workflow_runner_id").
		Join("INNER JOIN cd_workflow cw ON cw.id=wfr.cd_workflow_id").
		Where("cw.id = ?", cdWorkflowId).
		Where("pipeline_status_timeline.status = ?", status).
		Order("cw.id DESC").Limit(1).Select()
	if err != nil {
		impl.logger.Errorw("error in getting timeline of latest wf by pipelineId and status", "err", err, "cdWorkflowId", cdWorkflowId)
		return nil, err
	}
	return timeline, nil
}

func (impl *PipelineStatusTimelineRepositoryImpl) CheckTimelineExistsOfLatestWfByCdWorkflowIdAndStatus(cdWorkflowId int, status TimelineStatus) (bool, error) {
	timeline := &PipelineStatusTimeline{}
	exists, err := impl.dbConnection.Model(timeline).
		Join("INNER JOIN cd_workflow_runner wfr ON wfr.id = pipeline_status_timeline.cd_workflow_runner_id").
		Join("INNER JOIN cd_workflow cw ON cw.id=wfr.cd_workflow_id").
		Where("cw.id = ?", cdWorkflowId).
		Where("pipeline_status_timeline.status = ?", status).
		Order("cw.id DESC").Limit(1).Exists()
	if err != nil {
		impl.logger.Errorw("error in getting timeline of latest wf by pipelineId and status", "err", err, "cdWorkflowId", cdWorkflowId)
		return false, err
	}
	return exists, nil
}
