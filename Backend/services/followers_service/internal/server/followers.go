package server

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zopuu/soa-team-20/Backend/services/followers_service/proto/followerspb"
)

type FollowersServer struct {
	followerspb.UnimplementedFollowersServiceServer
	Driver neo4j.DriverWithContext
}

// Follow a user
func (s *FollowersServer) Follow(ctx context.Context, req *followerspb.FollowRequest) (*followerspb.FollowResponse, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	fmt.Println("GOT REQUEST", req.FollowerId, req.FolloweeId)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MERGE (a:User {id: $followerId})
			 MERGE (b:User {id: $followeeId})
			 MERGE (a)-[:FOLLOWS]->(b)`,
			map[string]any{
				"followerId": req.FollowerId,
				"followeeId": req.FolloweeId,
			})
		return nil, err
	})
	fmt.Println("RESPONSE FROM DB", err)
	if err != nil {
		return &followerspb.FollowResponse{Success: false}, err
	}
	return &followerspb.FollowResponse{Success: true}, nil
}

// Unfollow a user
func (s *FollowersServer) Unfollow(ctx context.Context, req *followerspb.FollowRequest) (*followerspb.FollowResponse, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MATCH (a:User {id: $followerId})-[r:FOLLOWS]->(b:User {id: $followeeId})
			 DELETE r`,
			map[string]any{
				"followerId": req.FollowerId,
				"followeeId": req.FolloweeId,
			})
		return nil, err
	})
	if err != nil {
		return &followerspb.FollowResponse{Success: false}, err
	}
	return &followerspb.FollowResponse{Success: true}, nil
}

// Get users this user is following
func (s *FollowersServer) GetFollowing(ctx context.Context, req *followerspb.UserRequest) (*followerspb.UsersResponse, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (:User {id: $id})-[:FOLLOWS]->(u:User) RETURN u.id`,
			map[string]any{"id": req.UserId})
		if err != nil {
			return nil, err
		}
		var ids []string
		for records.Next(ctx) {
			ids = append(ids, records.Record().Values[0].(string))
		}
		return ids, nil
	})
	if err != nil {
		return nil, err
	}
	return &followerspb.UsersResponse{UserIds: result.([]string)}, nil
}

func (s *FollowersServer) GetRecommendations(ctx context.Context, req *followerspb.UserRequest) (*followerspb.UsersResponse, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (:User {id: $id})-[:FOLLOWS]->(:User)-[:FOLLOWS]->(rec:User)
             WHERE rec.id <> $id AND NOT (:User {id: $id})-[:FOLLOWS]->(rec)
             RETURN DISTINCT rec.id`,
			map[string]any{"id": req.UserId})
		if err != nil {
			return nil, err
		}

		var ids []string
		for records.Next(ctx) {
			ids = append(ids, records.Record().Values[0].(string))
		}
		return ids, records.Err()
	})
	if err != nil {
		return nil, err
	}

	return &followerspb.UsersResponse{UserIds: result.([]string)}, nil
}

// Get followers of a user
func (s *FollowersServer) GetFollowers(ctx context.Context, req *followerspb.UserRequest) (*followerspb.UsersResponse, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (u:User)-[:FOLLOWS]->(:User {id: $id}) RETURN u.id`,
			map[string]any{"id": req.UserId})
		if err != nil {
			return nil, err
		}
		var ids []string
		for records.Next(ctx) {
			ids = append(ids, records.Record().Values[0].(string))
		}
		return ids, nil
	})
	if err != nil {
		return nil, err
	}
	return &followerspb.UsersResponse{UserIds: result.([]string)}, nil
}
