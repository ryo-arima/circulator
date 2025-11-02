package repository

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// MySQLRepository handles direct MySQL database operations
type MySQLRepository struct {
	config config.BaseConfig
	db     *gorm.DB
}

// NewMySQLRepository creates a new MySQLRepository instance
func NewMySQLRepository(cfg config.BaseConfig, dsn string) (*MySQLRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		cfg.Logger.ERROR(config.SRMCONN, "Failed to connect to MySQL", map[string]interface{}{
			"error": err.Error(),
			"dsn":   dsn,
		})
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	repo := &MySQLRepository{
		config: cfg,
		db:     db,
	}

	// Auto-migrate tables
	if err := repo.autoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	cfg.Logger.INFO(config.SRMINIT, "MySQL repository initialized", map[string]interface{}{
		"dsn": dsn,
	})

	return repo, nil
}

// autoMigrate performs database migrations
func (r *MySQLRepository) autoMigrate() error {
	r.config.Logger.DEBUG(config.SRMMIG, "Performing database migrations", nil)

	models := []interface{}{
		&model.Agent{},
		&model.AgentInfo{},
		&model.AgentProcessingConfig{},
	}

	for _, model := range models {
		if err := r.db.AutoMigrate(model); err != nil {
			r.config.Logger.ERROR(config.SRMERR, "Failed to migrate model", map[string]interface{}{
				"error": err.Error(),
				"model": fmt.Sprintf("%T", model),
			})
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	r.config.Logger.INFO(config.SRMSUCC, "Database migrations completed", map[string]interface{}{
		"model_count": len(models),
	})

	return nil
}

// GetAllAgents retrieves all agents from database
func (r *MySQLRepository) GetAllAgents(ctx context.Context) ([]*model.Agent, error) {
	r.config.Logger.DEBUG(config.SRMGA, "Getting all agents from database", nil)

	var agents []*model.Agent
	if err := r.db.WithContext(ctx).Find(&agents).Error; err != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to retrieve agents", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to retrieve agents: %w", err)
	}

	r.config.Logger.INFO(config.SRMSUCC, "Successfully retrieved agents", map[string]interface{}{
		"agent_count": len(agents),
	})

	return agents, nil
}

// GetAgentByUUID retrieves an agent by UUID
func (r *MySQLRepository) GetAgentByUUID(ctx context.Context, uuid string) (*model.Agent, error) {
	r.config.Logger.DEBUG(config.SRMGAU, "Getting agent by UUID", map[string]interface{}{
		"uuid": uuid,
	})

	var agent model.Agent
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.config.Logger.WARN(config.SRMNOTFOUND, "Agent not found", map[string]interface{}{
				"uuid": uuid,
			})
			return nil, fmt.Errorf("agent not found: %s", uuid)
		}
		r.config.Logger.ERROR(config.SRMERR, "Failed to retrieve agent by UUID", map[string]interface{}{
			"error": err.Error(),
			"uuid":  uuid,
		})
		return nil, fmt.Errorf("failed to retrieve agent by UUID: %w", err)
	}

	r.config.Logger.DEBUG(config.SRMSUCC, "Successfully retrieved agent by UUID", map[string]interface{}{
		"uuid": uuid,
		"id":   agent.ID,
	})

	return &agent, nil
}

// CreateAgent creates a new agent in database
func (r *MySQLRepository) CreateAgent(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	r.config.Logger.INFO(config.SRMCA, "Creating agent in database", map[string]interface{}{
		"uuid":     agent.UUID,
		"hostname": agent.Hostname,
	})

	if err := r.db.WithContext(ctx).Create(agent).Error; err != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to create agent", map[string]interface{}{
			"error":    err.Error(),
			"uuid":     agent.UUID,
			"hostname": agent.Hostname,
		})
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	r.config.Logger.INFO(config.SRMSUCC, "Successfully created agent", map[string]interface{}{
		"uuid": agent.UUID,
		"id":   agent.ID,
	})

	return agent, nil
}

// UpdateAgent updates an existing agent in database
func (r *MySQLRepository) UpdateAgent(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	r.config.Logger.INFO(config.SRMUA, "Updating agent in database", map[string]interface{}{
		"uuid": agent.UUID,
		"id":   agent.ID,
	})

	if err := r.db.WithContext(ctx).Save(agent).Error; err != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to update agent", map[string]interface{}{
			"error": err.Error(),
			"uuid":  agent.UUID,
			"id":    agent.ID,
		})
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}

	r.config.Logger.INFO(config.SRMSUCC, "Successfully updated agent", map[string]interface{}{
		"uuid": agent.UUID,
		"id":   agent.ID,
	})

	return agent, nil
}

// DeleteAgent deletes an agent from database
func (r *MySQLRepository) DeleteAgent(ctx context.Context, uuid string) error {
	r.config.Logger.INFO(config.SRMDA, "Deleting agent from database", map[string]interface{}{
		"uuid": uuid,
	})

	result := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&model.Agent{})
	if result.Error != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to delete agent", map[string]interface{}{
			"error": result.Error.Error(),
			"uuid":  uuid,
		})
		return fmt.Errorf("failed to delete agent: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.config.Logger.WARN(config.SRMNOTFOUND, "Agent not found for deletion", map[string]interface{}{
			"uuid": uuid,
		})
		return fmt.Errorf("agent not found for deletion: %s", uuid)
	}

	r.config.Logger.INFO(config.SRMSUCC, "Successfully deleted agent", map[string]interface{}{
		"uuid": uuid,
	})

	return nil
}

// CountAgents returns the total number of agents
func (r *MySQLRepository) CountAgents(ctx context.Context) (int64, error) {
	r.config.Logger.DEBUG(config.SRMCOUNT, "Counting agents in database", nil)

	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Agent{}).Count(&count).Error; err != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to count agents", map[string]interface{}{
			"error": err.Error(),
		})
		return 0, fmt.Errorf("failed to count agents: %w", err)
	}

	r.config.Logger.DEBUG(config.SRMSUCC, "Successfully counted agents", map[string]interface{}{
		"count": count,
	})

	return count, nil
}

// Close closes the database connection
func (r *MySQLRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	
	if err := sqlDB.Close(); err != nil {
		r.config.Logger.ERROR(config.SRMERR, "Failed to close database connection", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	r.config.Logger.INFO(config.SRMCLOSE, "MySQL repository closed", nil)
	return nil
}