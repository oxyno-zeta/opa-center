package partitions

import (
	"time"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

// Limit of element search.
const ListLimit = 2

type RetentionCleanTask struct {
	s          *service
	logger     log.Logger
	inProgress bool
}

func (r *RetentionCleanTask) Description() string { return "Retention clean processing task" }
func (r *RetentionCleanTask) Key() int            { return 1 }

func (r *RetentionCleanTask) endCurrentTask() {
	r.inProgress = false
}

func (r *RetentionCleanTask) buildCurrentLogger() log.Logger {
	return r.logger.WithField("task-id", time.Now().Unix())
}

func (r *RetentionCleanTask) Execute() {
	// Build logger
	logger := r.buildCurrentLogger()
	// Check if another run isn't already in progress
	if r.inProgress {
		logger.Info("Another retention clean is already in progress => Skipping this run")

		return
	}

	// Store fact that task is in progress
	r.inProgress = true
	// Defer end current task
	defer r.endCurrentTask()

	logger.Info("Starting retention clean processing task")

	err := r.runTask(logger)
	// Check error
	if err != nil {
		logger.Error(err)

		return
	}

	logger.Info("Retention clean processing task ended")
}

func (r *RetentionCleanTask) runTask(logger log.Logger) error {
	// Initialize page input
	pageIn := &pagination.PageInput{Limit: ListLimit}

	// While page input exists, loop to run the task
	for pageIn != nil {
		// Get all partition paginated
		list, pOut, err := r.s.dao.GetAllPaginated(
			pageIn,
			nil,
			nil,
			&models.Projection{ID: true, DecisionLogRetention: true, StatusDataRetention: true},
		)
		// Check error
		if err != nil {
			return err
		}

		// Loop over the list
		for _, item := range list {
			// For each partition try to parse durations if they exists
			// Decision log duration case
			if item.DecisionLogRetention != "" {
				retentionDuration, err := time.ParseDuration(item.DecisionLogRetention)
				// Check error
				if err != nil {
					return err
				}
				// Start retention clean process on decision logs
				err = r.s.decisionLogsSvc.ManageRetention(logger, retentionDuration, item.ID)
				// Check error
				if err != nil {
					return err
				}
			}
			// Decision log duration case
			if item.StatusDataRetention != "" {
				retentionDuration, err := time.ParseDuration(item.StatusDataRetention)
				// Check error
				if err != nil {
					return err
				}
				// Start retention clean process on decision logs
				err = r.s.statusesSvc.ManageRetention(logger, retentionDuration, item.ID)
				// Check error
				if err != nil {
					return err
				}
			}
		}

		// Clean pagination input
		pageIn = nil
		// Calculate next pagination input if there is a next page
		if pOut.HasNext {
			pageIn = &pagination.PageInput{Limit: ListLimit, Skip: pOut.Skip + ListLimit}
		}
	}

	return nil
}
