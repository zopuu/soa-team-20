package api

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zopuu/soa-team-20/Backend/services/shopping_service/internal/model"
	"github.com/zopuu/soa-team-20/Backend/services/shopping_service/internal/repo"
	shoppingpb "github.com/zopuu/soa-team-20/Backend/services/shopping_service/proto/shoppingpb"
)

type Server struct {
	shoppingpb.UnimplementedShoppingServiceServer
	repo *repo.ShoppingRepo
}

func NewServer(db *mongo.Database) *Server {
	r := repo.NewShoppingRepo(db)
	if err := r.EnsureIndexes(context.Background()); err != nil {
		log.Printf("index create error: %v", err)
	}
	return &Server{repo: r}
}

func (s *Server) AddToCart(ctx context.Context, req *shoppingpb.AddToCartRequest) (*shoppingpb.AddToCartResponse, error) {
	cart, err := s.repo.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	item := model.OrderItem{
		TourID: req.Item.TourId,
		Name:   req.Item.Name,
		Price:  req.Item.Price,
	}
	cart.Items = append(cart.Items, item)
	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}
	return &shoppingpb.AddToCartResponse{Cart: cartToProto(cart)}, nil
}

func (s *Server) RemoveFromCart(ctx context.Context, req *shoppingpb.RemoveFromCartRequest) (*shoppingpb.RemoveFromCartResponse, error) {
	cart, err := s.repo.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// remove first matching tour_id
	newItems := make([]model.OrderItem, 0, len(cart.Items))
	removed := false
	for _, it := range cart.Items {
		if !removed && it.TourID == req.TourId {
			removed = true
			continue
		}
		newItems = append(newItems, it)
	}
	cart.Items = newItems
	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}
	return &shoppingpb.RemoveFromCartResponse{Cart: cartToProto(cart)}, nil
}

func (s *Server) GetCart(ctx context.Context, req *shoppingpb.GetCartRequest) (*shoppingpb.GetCartResponse, error) {
	cart, err := s.repo.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &shoppingpb.GetCartResponse{Cart: cartToProto(cart)}, nil
}

func (s *Server) Checkout(ctx context.Context, req *shoppingpb.CheckoutRequest) (*shoppingpb.CheckoutResponse, error) {
	cart, err := s.repo.GetCart(ctx, req.UserId)
	println("Checkout cart items:", len(cart.Items), " for user:", req.UserId)
	if err != nil {
		return nil, err
	}
	tokens := make([]model.TourPurchaseToken, 0, len(cart.Items))
	protoTokens := make([]*shoppingpb.TourPurchaseToken, 0, len(cart.Items))
	now := time.Now()
	for _, it := range cart.Items {
		t := model.TourPurchaseToken{
			ID:          uuid.NewString(),
			UserID:      req.UserId,
			TourID:      it.TourID,
			TourName:    it.Name,
			Price:       it.Price,
			PurchasedAt: now,
		}
		tokens = append(tokens, t)
		protoTokens = append(protoTokens, &shoppingpb.TourPurchaseToken{
			Id:          t.ID,
			UserId:      t.UserID,
			TourId:      t.TourID,
			TourName:    t.TourName,
			Price:       t.Price,
			PurchasedAt: timestamppb.New(t.PurchasedAt),
		})
	}

	// store tokens and clear cart in a simple transactional order (no real transaction)
	if err := s.repo.SaveTokens(ctx, tokens); err != nil {
		return nil, err
	}
	if err := s.repo.ClearCart(ctx, req.UserId); err != nil {
		return nil, err
	}

	return &shoppingpb.CheckoutResponse{Tokens: protoTokens}, nil
}

func (s *Server) GetTokens(ctx context.Context, req *shoppingpb.GetTokensRequest) (*shoppingpb.GetTokensResponse, error) {
	toks, err := s.repo.GetTokens(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	out := make([]*shoppingpb.TourPurchaseToken, 0, len(toks))
	for _, t := range toks {
		out = append(out, &shoppingpb.TourPurchaseToken{
			Id:          t.ID,
			UserId:      t.UserID,
			TourId:      t.TourID,
			TourName:    t.TourName,
			Price:       t.Price,
			PurchasedAt: timestamppb.New(t.PurchasedAt),
		})
	}
	return &shoppingpb.GetTokensResponse{Tokens: out}, nil
}

func cartToProto(c *model.Cart) *shoppingpb.Cart {
	total := 0.0
	items := make([]*shoppingpb.OrderItem, 0, len(c.Items))
	for _, it := range c.Items {
		total += it.Price
		items = append(items, &shoppingpb.OrderItem{
			TourId: it.TourID, Name: it.Name, Price: it.Price,
		})
	}
	return &shoppingpb.Cart{UserId: c.UserID, Items: items, Total: total}
}
