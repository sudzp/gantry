package storage

import (
	"context"
	"fmt"
	"time"

	"gantry/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorage implements MongoDB-based storage
type MongoStorage struct {
	client       *mongo.Client
	database     *mongo.Database
	workflows    *mongo.Collection
	workflowRuns *mongo.Collection
}

// NewMongoStorage creates a new MongoDB storage instance
func NewMongoStorage(uri, database string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(database)

	return &MongoStorage{
		client:       client,
		database:     db,
		workflows:    db.Collection("workflows"),
		workflowRuns: db.Collection("workflow_runs"),
	}, nil
}

// SaveWorkflow saves a workflow to MongoDB
func (s *MongoStorage) SaveWorkflow(wf *models.Workflow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": wf.Name}
	update := bson.M{"$set": wf}
	opts := options.Update().SetUpsert(true)

	_, err := s.workflows.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save workflow: %w", err)
	}

	return nil
}

// GetWorkflow retrieves a workflow by name
func (s *MongoStorage) GetWorkflow(name string) (*models.Workflow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wf models.Workflow
	err := s.workflows.FindOne(ctx, bson.M{"name": name}).Decode(&wf)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("workflow '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return &wf, nil
}

// ListWorkflows returns all workflows
func (s *MongoStorage) ListWorkflows() ([]*models.Workflow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.workflows.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	defer func() { _ = cursor.Close(ctx) }()

	var workflows []*models.Workflow
	if err := cursor.All(ctx, &workflows); err != nil {
		return nil, fmt.Errorf("failed to decode workflows: %w", err)
	}

	return workflows, nil
}

// DeleteWorkflow deletes a workflow
func (s *MongoStorage) DeleteWorkflow(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.workflows.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("workflow '%s' not found", name)
	}

	return nil
}

// SaveRun saves a workflow run
func (s *MongoStorage) SaveRun(run *models.WorkflowRun) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clone to avoid mutex issues
	clone := run.Clone()

	filter := bson.M{"id": run.ID}
	update := bson.M{"$set": clone}
	opts := options.Update().SetUpsert(true)

	_, err := s.workflowRuns.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save run: %w", err)
	}

	return nil
}

// DeleteRunsByWorkflow deletes all runs for a workflow
func (s *MongoStorage) DeleteRunsByWorkflow(workflowName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.workflowRuns.DeleteMany(ctx, bson.M{"workflow_name": workflowName})
	if err != nil {
		return fmt.Errorf("failed to delete runs: %w", err)
	}

	if result.DeletedCount > 0 {
		fmt.Printf("Deleted %d runs for workflow '%s'\n", result.DeletedCount, workflowName)
	}

	return nil
}

// GetRun retrieves a run by ID
func (s *MongoStorage) GetRun(id string) (*models.WorkflowRun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var run models.WorkflowRun
	err := s.workflowRuns.FindOne(ctx, bson.M{"id": id}).Decode(&run)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("run '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to get run: %w", err)
	}

	return &run, nil
}

// ListRuns returns all runs, sorted by start time (newest first)
func (s *MongoStorage) ListRuns() ([]*models.WorkflowRun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "started_at", Value: -1}})
	cursor, err := s.workflowRuns.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list runs: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var runs []*models.WorkflowRun
	if err := cursor.All(ctx, &runs); err != nil {
		return nil, fmt.Errorf("failed to decode runs: %w", err)
	}

	return runs, nil
}

// UpdateRun updates an existing run
func (s *MongoStorage) UpdateRun(run *models.WorkflowRun) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clone to avoid mutex issues
	clone := run.Clone()

	filter := bson.M{"id": run.ID}
	update := bson.M{"$set": clone}

	result, err := s.workflowRuns.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update run: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("run '%s' not found", run.ID)
	}

	return nil
}

// Close closes the MongoDB connection
func (s *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}
