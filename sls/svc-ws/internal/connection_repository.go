package ws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Jimeux/sls-chat-app/sls/lib/awsiface"
)

// Connection stores a reference to an active WebSocket connection.
type Connection struct {
	ConnectionID string `dynamodbav:"connection_id" json:"connectionId"`
	Username     string `dynamodbav:"username,omitempty" json:"username,omitempty"` // not unique
}

// Repository implements CRUD operations on the Connection table.
type Repository struct {
	ddb       awsiface.DynamoDB
	tableName *string
}

func NewRepository(ddb awsiface.DynamoDB, tableName string) *Repository {
	return &Repository{
		ddb:       ddb,
		tableName: &tableName,
	}
}

// NewRepositoryFromConfig is a convenience function for creating
// a Repository from an aws.Config instance.
func NewRepositoryFromConfig(cfg aws.Config, tableName string) *Repository {
	ddb := dynamodb.NewFromConfig(cfg)
	return NewRepository(ddb, tableName)
}

// SaveConnection creates or updates a Connection record with the given values.
func (r *Repository) SaveConnection(ctx context.Context, connID, username string) error {
	conn := &Connection{
		ConnectionID: connID,
		Username:     username,
	}
	av, err := attributevalue.MarshalMap(conn)
	if err != nil {
		return fmt.Errorf("failed to marshal connection: %w", err)
	}

	if _, err := r.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: r.tableName,
	}); err != nil {
		return fmt.Errorf("failed PutItem for connection: %w", err)
	}
	return nil
}

func (r *Repository) RemoveConnection(ctx context.Context, connID string) error {
	conn := &Connection{ConnectionID: connID}
	key, err := attributevalue.MarshalMap(conn)
	if err != nil {
		return fmt.Errorf("failed to marshal Connection: %w", err)
	}

	if _, err := r.ddb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key:       key,
	}); err != nil {
		return fmt.Errorf("failed PutItem for connection: %w", err)
	}
	return nil
}

func (r *Repository) GetConnection(ctx context.Context, connID string) (*Connection, error) {
	keyConn := &Connection{ConnectionID: connID}
	key, err := attributevalue.MarshalMap(keyConn)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Connection: %w", err)
	}

	res, err := r.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var conn *Connection
	if err := attributevalue.UnmarshalMap(res.Item, &conn); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conns: %w", err)
	}
	return conn, nil
}

func (r *Repository) GetConnections(ctx context.Context) ([]*Connection, error) {
	res, err := r.ddb.Scan(ctx, &dynamodb.ScanInput{
		TableName: r.tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	var conns []*Connection
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &conns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conns: %w", err)
	}
	return conns, nil
}
